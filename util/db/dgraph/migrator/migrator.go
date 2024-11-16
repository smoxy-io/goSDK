package migrator

import (
	"context"
	"fmt"
)

type Migrator interface {
	Apply() error
	CanMigrate() bool
	FromVersion() string
	ToVersion() string
}

type MigrationOption func() Migrator

type NewMigrator func(ctx context.Context) Migrator

var (
	migrators = map[string]NewMigrator{}
)

func RegisterMigrator(fromVersion string, toVersion string, migrator NewMigrator) {
	migrators[genMigratorKey(fromVersion, toVersion)] = migrator
}

func GetMigratorFactory(fromVersion string, toVersion string) NewMigrator {
	if m, ok := migrators[genMigratorKey(fromVersion, toVersion)]; ok {
		return m
	}

	return nil
}

func GetMigratorByFromVersion(ctx context.Context, fromVersion string) Migrator {
	for _, migrator := range migrators {
		m := migrator(ctx)

		if m.FromVersion() == fromVersion {
			return m
		}
	}

	return nil
}

func MigrateFrom(ctx context.Context, fromVersion string) error {
	if fromVersion == "" {
		fromVersion = "v0.0.0"
	}

	for {
		migrator := GetMigratorByFromVersion(ctx, fromVersion)

		if migrator == nil {
			break
		}

		if !migrator.CanMigrate() {
			return fmt.Errorf("migration %s is not allowed", genMigratorKey(migrator.FromVersion(), migrator.ToVersion()))
		}

		fmt.Printf("Performing migration %s\n", genMigratorKey(migrator.FromVersion(), migrator.ToVersion()))

		if err := migrator.Apply(); err != nil {
			return err
		}

		fromVersion = migrator.ToVersion()
	}

	return nil
}

func genMigratorKey(fromVersion string, toVersion string) string {
	return fromVersion + " -> " + toVersion
}
