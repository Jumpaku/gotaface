package test

import (
	"context"
	"flag"
	"fmt"
	"github.com/Jumpaku/sqanner/tokenize"
	"strings"
	"testing"

	"cloud.google.com/go/spanner"

	spanner_admin "cloud.google.com/go/spanner/admin/database/apiv1"
	spanner_adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"github.com/samber/lo"
)

var DataSource = flag.String("data-source", "", "data source name")

func SkipWithoutDataSource(t *testing.T) {
	if *DataSource == "" {
		t.Skipf("flag %q is not specified", "data-source")
	}
}
func Setup(t *testing.T) (adminClient *spanner_admin.DatabaseAdminClient, client *spanner.Client) {
	t.Helper()

	if *DataSource == "" {
		t.Fatal("if flag -data-source is not set, it will be skipped")
	}

	ctx := context.Background()
	adminClient, err := spanner_admin.NewDatabaseAdminClient(ctx)
	if err != nil {
		t.Fatalf(`fail to create spanner admin client: %v`, err)
	}
	t.Cleanup(func() { adminClient.Close() })

	err = adminClient.DropDatabase(ctx, &spanner_adminpb.DropDatabaseRequest{Database: *DataSource})
	if err != nil {
		t.Logf(`fail to drop spanner database %q: %v`, *DataSource, err)
	}

	splitDataSource := strings.Split(*DataSource, "/")
	parent := strings.Join(splitDataSource[:4], "/")
	database := splitDataSource[5]

	op, err := adminClient.CreateDatabase(ctx, &spanner_adminpb.CreateDatabaseRequest{
		Parent:          parent,
		CreateStatement: fmt.Sprintf("CREATE DATABASE %s", database),
	})
	if err != nil {
		t.Fatalf(`fail to create spanner database in %q: %v`, parent, err)
	}

	if _, err := op.Wait(ctx); err != nil {
		t.Fatalf(`fail to create spanner database in %q: %v`, parent, err)
	}
	t.Cleanup(func() { adminClient.DropDatabase(ctx, &spanner_adminpb.DropDatabaseRequest{Database: *DataSource}) })

	client, err = spanner.NewClient(ctx, *DataSource)
	if err != nil {
		t.Fatalf(`fail to create spanner client with %q: %v`, *DataSource, err)
	}
	t.Cleanup(func() { client.Close() })

	return adminClient, client
}

func InitDDLs(t *testing.T, adminClient *spanner_admin.DatabaseAdminClient, database string, stmts []string) {
	t.Helper()

	removeComment := func(stmt string, _ int) string {
		tokens, err := tokenize.Tokenize([]rune(stmt))
		if err != nil {
			panic(err)
		}
		tokens = lo.Filter(tokens, func(token tokenize.Token, _ int) bool { return token.Kind != tokenize.TokenComment })
		return lo.Reduce(tokens, func(agg string, token tokenize.Token, _ int) string {
			return agg + string(token.Content)
		}, "")
	}
	ctx := context.Background()
	ddl := &spanner_adminpb.UpdateDatabaseDdlRequest{
		Database:   database,
		Statements: lo.Map(stmts, removeComment),
	}

	op, err := adminClient.UpdateDatabaseDdl(ctx, ddl)
	if err != nil {
		t.Fatalf(`fail to execute ddl: %v`, err)
	}
	if err := op.Wait(ctx); err != nil {
		t.Fatalf(`fail to wait create tables: %v`, err)
	}
}
