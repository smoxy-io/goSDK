package db

import (
	"context"
	"fmt"
	"github.com/smoxy-io/goSDK/util/db/dgraph/migrator"
	"github.com/smoxy-io/goSDK/util/errors"
)

func MigrateSchema(ctx context.Context, fromVersion string, toVersion string) error {
	if fromVersion == "" {
		return errors.ErrInvalid.WithVars("fromVersion")
	}

	if toVersion == "" {
		return errors.ErrInvalid.WithVars("toVersion")
	}

	if toVersion == "latest" {
		// perform all database migrations starting from `fromVersion`
		return migrator.MigrateFrom(ctx, fromVersion)
	}

	newMigratorFn := migrator.GetMigratorFactory(fromVersion, toVersion)

	if newMigratorFn == nil {
		return fmt.Errorf("no migration defined for %s -> %s", fromVersion, toVersion)
	}

	m := newMigratorFn(ctx)

	return m.Apply()
}
