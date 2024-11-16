package migrator

import (
	"context"
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	proto "github.com/smoxy-io/goSDK/pkg/proto/smoxy/util/db/dgraph"
	db "github.com/smoxy-io/goSDK/util/db/dgraph"
	"github.com/smoxy-io/goSDK/util/db/dgraph/models"
	"strconv"
	"time"
)

type AppyFn func(ctx context.Context) error

// Migration base migrator struct
type Migration struct {
	fromVersion string
	toVersion   string
	ctx         context.Context
	PBar        *progressbar.ProgressBar
	schemaMeta  *proto.SchemaMetaData
	migration   *proto.Migration
	result      *proto.MigrationResult
}

func (m *Migration) Apply() error {
	return errors.New("not implemented")
}

// Migrate wrapper structs call this function to reduce boilerplate code that's needed for a migration
func (m *Migration) Migrate(applyFn AppyFn) error {
	// check if migration can be performed
	if !m.CanMigrate() {
		return fmt.Errorf("cannot perform %s –> %s migration", m.fromVersion, m.toVersion)
	}

	// start the migration transaction
	nCtx, cErr := db.StartTxn(m.ctx)

	if cErr != nil {
		return cErr
	}

	defer db.Rollback(nCtx) // ensure the setup transaction is rolled back. safe to call after successful commit

	if m.migration == nil {
		// create a new migration
		m.migration = &proto.Migration{
			FromVersion: m.fromVersion,
			ToVersion:   m.toVersion,
			LastStatus:  proto.MigrationStatus_SUCCESS,
			Results:     make([]*proto.MigrationResult, 0),
		}

		if err := models.Set(nCtx, m.migration); err != nil {
			return err
		}

		m.schemaMeta.Migrations = append(m.schemaMeta.Migrations, m.migration)

		if err := models.Set(nCtx, m.schemaMeta); err != nil {
			return err
		}
	}

	// create a new migration result
	m.result = &proto.MigrationResult{
		Migration: m.migration,
		StartTime: strconv.FormatInt(time.Now().UTC().Unix(), 10),
		EndTime:   nil,
		Status:    proto.MigrationStatus_IN_PROGRESS,
	}

	if err := models.Set(nCtx, m.result); err != nil {
		return err
	}

	m.migration.LastStatus = m.result.Status
	m.migration.Results = append(m.migration.Results, m.result)

	if err := models.Set(nCtx, m.migration); err != nil {
		return err
	}

	isMigrating := true
	m.schemaMeta.IsMigrating = &isMigrating
	lastMigStatus := m.result.Status
	m.schemaMeta.LastMigrationStatus = &lastMigStatus

	if err := models.Set(nCtx, m.schemaMeta); err != nil {
		return err
	}

	if err := db.Commit(nCtx); err != nil {
		return err
	}

	// run the migration
	mCtx, mcErr := db.StartTxn(m.ctx) // build the migration context off the original context

	if mcErr != nil {
		return mcErr
	}

	defer db.Rollback(mCtx) // ensure the migration transaction is rolled back. safe to call after successful commit

	aErr := applyFn(mCtx)

	updateMetaOnErr := func() {
		m.result.Status = proto.MigrationStatus_FAILED
		endTime := strconv.FormatInt(time.Now().UTC().Unix(), 10)
		m.result.EndTime = &endTime

		m.migration.LastStatus = m.result.Status

		*m.schemaMeta.IsMigrating = false
		*m.schemaMeta.LastMigrationStatus = m.result.Status
		m.schemaMeta.LastMigrationTime = &endTime
		m.schemaMeta.Version = m.fromVersion

		// attempt to update the metadata to reflect migration failure
		// - won't work if failure is due to db being offline
		eCtx, ecErr := db.StartTxn(m.ctx)

		if ecErr != nil {
			return
		}

		defer db.Rollback(eCtx)

		if err := models.Set(eCtx, m.result); err != nil {
			return
		}

		if err := models.Set(eCtx, m.migration); err != nil {
			return
		}

		if err := models.Set(eCtx, m.schemaMeta); err != nil {
			return
		}

		_ = db.Commit(eCtx)
	}

	if aErr != nil {
		updateMetaOnErr()
		return aErr
	}

	// add the update to the schema metadata version and the migration results to the migration txn
	m.result.Status = proto.MigrationStatus_SUCCESS
	endTime := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	m.result.EndTime = &endTime

	m.migration.LastStatus = m.result.Status

	m.schemaMeta.Version = m.toVersion
	*m.schemaMeta.IsMigrating = false
	*m.schemaMeta.LastMigrationStatus = m.result.Status
	m.schemaMeta.LastMigrationTime = m.result.EndTime

	if err := models.Set(mCtx, m.result); err != nil {
		updateMetaOnErr()
		return err
	}

	if err := models.Set(mCtx, m.migration); err != nil {
		updateMetaOnErr()
		return err
	}

	if err := models.Set(mCtx, m.schemaMeta); err != nil {
		updateMetaOnErr()
		return err
	}

	// check that schema metadata is updated
	cm, cmErr := models.GetSchemaMetaData(mCtx)

	if cmErr != nil {
		updateMetaOnErr()
		return cmErr
	}

	if cm.Version != m.toVersion {
		updateMetaOnErr()
		return fmt.Errorf("schema version (%s) not updated to %s", m.toVersion)
	}

	if cm.LastMigrationTime == nil || *cm.LastMigrationTime != *m.schemaMeta.LastMigrationTime {
		updateMetaOnErr()
		return fmt.Errorf("schema last migration time not updated")
	}

	// apply the changes to the database
	if err := db.Commit(mCtx); err != nil {
		updateMetaOnErr()
		return err
	}

	return nil
}

func (m *Migration) CanMigrate() bool {
	md, err := models.GetSchemaMetaData(m.ctx)

	if err != nil {
		return false
	}

	if md == nil {
		return false
	}

	m.schemaMeta = md

	if md.Version != m.fromVersion {
		return false
	}

	if md.IsMigrating != nil && *md.IsMigrating {
		return false
	}

	if md.LastMigrationStatus != nil && md.LastMigrationStatus.Name() == string(models.MigrationStatusInProgress) {
		return false
	}

	for _, mi := range md.Migrations {
		if mi.FromVersion == m.fromVersion || mi.ToVersion == m.toVersion {
			// this migrator has been run before
			m.migration = mi

			status := mi.LastStatus.Name()

			if status == string(models.MigrationStatusSuccess) || status == string(models.MigrationStatusInProgress) {
				return false
			}

			return true
		}
	}

	return true
}

func (m *Migration) FromVersion() string {
	return m.fromVersion
}

func (m *Migration) ToVersion() string {
	return m.toVersion
}

func NewMigration(fromVersion string, toVersion string, ctx context.Context) *Migration {
	return &Migration{
		fromVersion: fromVersion,
		toVersion:   toVersion,
		ctx:         ctx,
		PBar:        nil,
	}
}
