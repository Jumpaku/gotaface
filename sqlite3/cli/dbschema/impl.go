package dbschema

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Jumpaku/gotaface/cli/dbschema"
	sqlite3_schema "github.com/Jumpaku/gotaface/sqlite3/ddl/schema"
	_ "github.com/mattn/go-sqlite3"
)

type DBSchema struct{}

func (DBSchema) Exec(ctx context.Context, driver string, dataSource string) (dbschema.DBSchemaOutput, error) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, fmt.Errorf(`fail to create sqlite3 client: %w`, err)
	}
	defer db.Close()

	schema, err := sqlite3_schema.NewFetcher(db).Fetch(ctx)
	if err != nil {
		return nil, fmt.Errorf(`fail to fetch table schema: %w`, err)
	}

	return schema.(*sqlite3_schema.Schema), nil
}
