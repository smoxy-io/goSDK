package queries

import proto "github.com/smoxy-io/goSDK/pkg/proto/smoxy/util/db/dgraph"

type QuerySchemaMetaData struct {
	SchemaMetaData []*proto.SchemaMetaData `json:"schemaMetaData"`
}

const (
	// there should only ever be 1 SchemaMetaData node
	GetSchemaMetaData = `
query {
  schemaMetaData(func: type(SchemaMetaData), first: 1) {
    ` + SchemaMetaDataFragment + `
  }
}
`
)

// Fragments
const (
	SchemaMetaDataFragment = `
id: uid
version: SchemaMetaData.version
lastUpdateTime: SchemaMetaData.lastUpdateTime
lastMigrationTime: SchemaMetaData.lastMigrationTime
lastMigrationStatus: SchemaMetaData.lastMigrationStatus
isMigrating: SchemaMetaData.isMigrating
migrations: SchemaMetaData.migrations {
  ` + MigrationFragment + `
}
`
)

const (
	MigrationFragment = `
id: uid
fromVersion: Migration.fromVersion
toVersion: Migration.toVersion
lastStatus: Migration.lastStatus
results: Migration.results {
  ` + MigrationResultFragment + `
}
`
)

const (
	MigrationResultFragment = `
id: uid
startTime: MigrationResult.startTime
endTime: MigrationResult.endTime
status: MigrationResult.status
`
)
