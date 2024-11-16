package models

import (
	"context"
	"encoding/json"
	"github.com/shurcooL/graphql"
	proto "github.com/smoxy-io/goSDK/pkg/proto/smoxy/util/db/dgraph"
	db "github.com/smoxy-io/goSDK/util/db/dgraph"
	"github.com/smoxy-io/goSDK/util/db/dgraph/queries"
	"github.com/smoxy-io/goSDK/util/errors"
)

type MigrationStatus graphql.String

const (
	MigrationStatusSuccess    MigrationStatus = "SUCCESS"
	MigrationStatusRolledBack MigrationStatus = "ROLLED_BACK"
	MigrationStatusFailed     MigrationStatus = "FAILED"
	MigrationStatusInProgress MigrationStatus = "IN_PROGRESS"
)

func GetSchemaMetaData(ctx context.Context) (*proto.SchemaMetaData, error) {
	dgraph, cErr, ctx := db.GetClient(ctx)

	if cErr != nil {
		return nil, cErr
	}

	txn := db.GetTxn(ctx)
	extTxn := txn != nil

	if !extTxn {
		txn = dgraph.NewReadOnlyTxn()
		defer txn.Discard(ctx)
	}

	qResp, qErr := txn.Query(ctx, queries.GetSchemaMetaData)

	if qErr != nil {
		return nil, qErr
	}

	var r queries.QuerySchemaMetaData

	if err := json.Unmarshal(qResp.Json, &r); err != nil {
		return nil, err
	}

	if len(r.SchemaMetaData) < 1 {
		return nil, errors.ErrNotFound
	}

	return r.SchemaMetaData[0], nil
}
