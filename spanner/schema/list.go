package schema

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/Jumpaku/gotaface/schema"
	gf_spanner "github.com/Jumpaku/gotaface/spanner"
	"github.com/samber/lo"
)

type lister struct {
	queryer gf_spanner.Queryer
}

func NewLister(queryer gf_spanner.Queryer) schema.Lister {
	return lister{queryer: queryer}
}

var _ schema.Lister = lister{}

func (lister lister) List(ctx context.Context) (tables []string, err error) {
	sql := `--sql query table name and parent information
SELECT TABLE_NAME AS Name,
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_SCHEMA = ''
ORDER BY TABLE_NAME ASC`
	type Table struct {
		Name string
	}
	list, err := gf_spanner.ScanRowsStruct[Table](lister.queryer.Query(ctx, spanner.Statement{SQL: sql}))
	if err != nil {
		return nil, fmt.Errorf(`fail to list tables: %w`, err)
	}

	return lo.Map(list, func(item Table, _ int) string { return item.Name }), nil
}
