package dbschema

import (
	"context"
	"encoding/json"

	"github.com/Jumpaku/gotaface/ddl/schema"
)

type DBSchemaOutput interface {
	json.Marshaler
	schema.Schema
}

type DBSchema interface {
	Exec(ctx context.Context, driver string, dataSource string) (DBSchemaOutput, error)
}
