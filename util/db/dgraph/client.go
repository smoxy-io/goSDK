package dgraph

import (
	"context"
	"github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/graphql"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math"
	"os"
)

const (
	ENV_HOST         = "DGRAPH_DB_HOST"
	ENV_GRPC_PORT    = "DGRAPH_DB_GRPC_PORT"
	ENV_GRAPHQL_PORT = "DGRAPH_DB_GRAPHQL_PORT"
)

const (
	clientCtxKey    = "db-client"
	gqlClientCtxKey = clientCtxKey + "-graphql"
	txnCtxKey       = "db-txn"
	txnDepthKey     = "db-txn-depth"

	DefaultHost        = "localhost"
	DefaultGrpcPort    = "9080"
	DefaultGraphQLPort = "8080"

	DgraphAdminTokenHeader = "X-Dgraph-AuthToken"

	SchemaVersion = "v0.0.1"
)

func NewClient() (*dgo.Dgraph, error) {
	//TODO: setup TLS
	connOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithInitialWindowSize(math.MaxInt32),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}

	conn, err := grpc.Dial(GetEndpoint(), connOpts...)

	if err != nil {
		return nil, err
	}

	//TODO: login to namespace

	return dgo.NewDgraphClient(api.NewDgraphClient(conn)), nil
}

func GetClient(ctx context.Context) (*dgo.Dgraph, error, context.Context) {
	// get the client from the context
	dbClient, ok := ctx.Value(clientCtxKey).(*dgo.Dgraph)

	if !ok {
		// ctx[clientCtxKey] is unset or not a Dgraph client
		// create it
		dbc, err := NewClient()

		if err != nil {
			return nil, err, ctx
		}

		switch c := ctx.(type) {
		case *gin.Context:
			c.Set(clientCtxKey, dbc)
		default:
			ctx = context.WithValue(ctx, clientCtxKey, dbc)
		}

		dbClient = dbc
	}

	return dbClient, nil, ctx
}

func NewGraphQLClient() (*graphql.Client, error) {
	//TODO: use custom httpClient for granular production settings and
	//      authentication
	client := graphql.NewClient(GetGraphQLEndpoint(), nil)

	return client, nil
}

func GetGraphQLClient(ctx context.Context) (*graphql.Client, error, context.Context) {
	// get the client from the context
	dbClient, ok := ctx.Value(gqlClientCtxKey).(*graphql.Client)

	if !ok {
		// ctx[clientCtxKey] is unset or not a Dgraph client
		// create it
		dbc, err := NewGraphQLClient()

		if err != nil {
			return nil, err, ctx
		}

		switch c := ctx.(type) {
		case *gin.Context:
			c.Set(gqlClientCtxKey, dbc)
		default:
			ctx = context.WithValue(ctx, gqlClientCtxKey, dbc)
		}

		dbClient = dbc
	}

	return dbClient, nil, ctx
}

func BackgroundContext(c context.Context) (context.Context, error) {
	if c == nil {
		c = context.Background()
	}

	_, cErr, ctx := GetClient(c)

	if cErr != nil {
		return ctx, cErr
	}

	_, cgErr, ctx := GetGraphQLClient(ctx)

	if cgErr != nil {
		return ctx, cgErr
	}

	return ctx, nil
}

func GetTxn(ctx context.Context) *dgo.Txn {
	txn, ok := ctx.Value(txnCtxKey).(*dgo.Txn)

	if !ok {
		return nil
	}

	return txn
}

func getTxnDepth(ctx context.Context) int {
	depth, ok := ctx.Value(txnDepthKey).(int)

	if !ok {
		return 0
	}

	return depth
}

func incrementTxnDepth(ctx context.Context, delta int) context.Context {
	depth := getTxnDepth(ctx)

	return context.WithValue(ctx, txnDepthKey, depth+delta)
}

func StartTxn(ctx context.Context, readOnly ...bool) (context.Context, error) {
	if t := GetTxn(ctx); t != nil {
		// context already has an active transaction. increment transaction depth
		return incrementTxnDepth(ctx, 1), nil
	}

	dgraph, err, nCtx := GetClient(ctx)

	if err != nil {
		return ctx, err
	}

	var txn *dgo.Txn

	if len(readOnly) > 0 && readOnly[0] {
		txn = dgraph.NewReadOnlyTxn()
	} else {
		txn = dgraph.NewTxn()
	}

	return context.WithValue(nCtx, txnCtxKey, txn), nil
}

func Commit(ctx context.Context) error {
	if depth := getTxnDepth(ctx); depth != 0 {
		// this is a nested transaction. the actual commit will be done at the top level context
		return nil
	}

	txn := GetTxn(ctx)

	if txn == nil {
		// no txn to commit
		return nil
	}

	return txn.Commit(ctx)
}

func Rollback(ctx context.Context) error {
	// not checking txn depth here as rollback is safe to call anytime
	// and if rollback is called anywhere in the nested txn tree, the entire txn tree should be rolled back
	txn := GetTxn(ctx)

	if txn == nil {
		// no txn to rollback
		return nil
	}

	return txn.Discard(ctx)
}

func GetHost() string {
	host := os.Getenv(ENV_HOST)

	if host == "" {
		return DefaultHost
	}

	return host
}

func GetGrpcPort() string {
	port := os.Getenv(ENV_GRPC_PORT)

	if port == "" {
		return DefaultGrpcPort
	}

	return port
}

func GetGraphQLPort() string {
	port := os.Getenv(ENV_GRAPHQL_PORT)

	if port == "" {
		return DefaultGraphQLPort
	}

	return port
}

func GetEndpoint() string {
	return GetHost() + ":" + GetGrpcPort()
}

func GetGraphQLEndpoint() string {
	return "http://" + GetHost() + ":" + GetGraphQLPort() + "/graphql"
}
