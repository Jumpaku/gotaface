package dbschema

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/Jumpaku/gotaface/cli/dbschema"
	spanner_schema "github.com/Jumpaku/gotaface/spanner/ddl/schema"
)

type DBSchema struct{}

func (DBSchema) Exec(ctx context.Context, driver string, dataSource string) (dbschema.DBSchemaOutput, error) {
	client, err := spanner.NewClient(ctx, dataSource)
	if err != nil {
		return nil, fmt.Errorf(`fail to create spanner client: %w`, err)
	}
	defer client.Close()

	tx := client.ReadOnlyTransaction()
	defer tx.Close()

	schema, err := spanner_schema.NewFetcher(tx).Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf(`fail to fetch table schema: %w`, err)
	}

	return schema.(*spanner_schema.Schema), nil
}
