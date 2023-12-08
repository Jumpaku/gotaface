package schema

import (
	"context"
	"fmt"
	"slices"

	"github.com/Jumpaku/go-assert"
	"github.com/Jumpaku/gotaface/schema"
	gf_sqlite3 "github.com/Jumpaku/gotaface/sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type SchemaColumn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}
type SchemaForeignKey struct {
	ReferencedTable string   `json:"referenced_table"`
	ReferencedKey   []string `json:"referenced_key"`
	ReferencingKey  []string `json:"referencing_key"`
}
type SchemaUniqueKey struct {
	Name string   `json:"name"`
	Key  []string `json:"key"`
}
type SchemaTable struct {
	Name        string             `json:"name"`
	Columns     []SchemaColumn     `json:"columns"`
	PrimaryKey  []string           `json:"primary_key"`
	ForeignKeys []SchemaForeignKey `json:"foreign_key"`
	UniqueKeys  []SchemaUniqueKey  `json:"unique_key"`
}

type fetcher struct {
	queryer sqlx.QueryerContext
}

func NewFetcher(queryer gf_sqlite3.Queryer) fetcher {
	return fetcher{queryer: queryer}
}

var _ schema.Fetcher[SchemaTable] = fetcher{}

func (fetcher fetcher) Fetch(ctx context.Context, table string) (SchemaTable, error) {
	wrapError := func(err error) (SchemaTable, error) {
		assert.Params(err != nil, "wrapped error must be not nil")
		return SchemaTable{}, fmt.Errorf(`fail to fetch schema of %s: %w`, table, err)
	}

	schemaTable := SchemaTable{Name: table}

	var err error
	schemaTable.Columns, err = queryColumns(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

	schemaTable.PrimaryKey, err = queryPrimaryKey(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

	schemaTable.ForeignKeys, err = queryForeignKeys(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

	schemaTable.UniqueKeys, err = queryUniqueKeys(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

	return schemaTable, nil
}

func queryColumns(ctx context.Context, tx gf_sqlite3.Queryer, table string) ([]SchemaColumn, error) {
	sql := `--sql query column information
SELECT 
	"name" AS Name,
	"type" AS Type,
	"notnull" = 0 AS Nullable
FROM pragma_table_info(?)
ORDER BY "cid"`
	rows, err := tx.QueryxContext(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get columns of %s: %w`, table, err)
	}
	type column struct {
		Name     string `db:"Name"`
		Type     string `db:"Type"`
		Nullable bool   `db:"Nullable"`
	}
	columns, err := gf_sqlite3.ScanRowsStruct[column](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get columns of %s: %w`, table, err)
	}

	return lo.Map(columns, func(column column, index int) SchemaColumn {
		return SchemaColumn{
			Name:     column.Name,
			Type:     column.Type,
			Nullable: column.Nullable,
		}
	}), nil
}

func queryPrimaryKey(ctx context.Context, tx gf_sqlite3.Queryer, table string) ([]string, error) {
	sql := `--sql query primary key information
SELECT
	"name" AS Name
FROM pragma_table_info(?)
WHERE "pk" > 0
ORDER BY "pk"`
	type key struct {
		Name string `db:"Name"`
	}
	rows, err := tx.QueryxContext(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get primary key of %s: %w`, table, err)
	}
	primaryKey, err := gf_sqlite3.ScanRowsStruct[key](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get primary key of %s: %w`, table, err)
	}
	return lo.Map(primaryKey, func(it key, i int) string { return it.Name }), nil
}

func queryForeignKeys(ctx context.Context, tx gf_sqlite3.Queryer, table string) ([]SchemaForeignKey, error) {
	sql := `--sql query foreign key information
SELECT
	"id" AS Id,
	"seq" AS Seq,
	"table" AS ReferencedTable,
	"from" AS ReferencingKey,
	"to" AS ReferencedKey
FROM pragma_foreign_key_list(?)
ORDER BY "id", "seq"`
	rows, err := tx.QueryxContext(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get foreign keys of %s: %w`, table, err)
	}
	type fkRow struct {
		Id              int64  `db:"Id"`
		Seq             int64  `db:"Seq"`
		ReferencedTable string `db:"ReferencedTable"`
		ReferencingKey  string `db:"ReferencingKey"`
		ReferencedKey   string `db:"ReferencedKey"`
	}
	fkRows, err := gf_sqlite3.ScanRowsStruct[fkRow](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get foreign keys of %s: %w`, table, err)
	}

	group := lo.GroupBy(fkRows, func(fkRow fkRow) int64 { return fkRow.Id })
	groupIDs := lo.MapToSlice(group, func(id int64, _ []fkRow) int64 { return id })
	slices.Sort(groupIDs)

	var foreignKeys []SchemaForeignKey
	for _, id := range groupIDs {
		g := group[id]
		foreignKeys = append(foreignKeys, SchemaForeignKey{
			ReferencedTable: g[0].ReferencedTable,
			ReferencedKey:   lo.Map(g, func(fkRow fkRow, _ int) string { return fkRow.ReferencedKey }),
			ReferencingKey:  lo.Map(g, func(fkRow fkRow, _ int) string { return fkRow.ReferencingKey }),
		})
	}

	return foreignKeys, nil
}

func queryUniqueKeys(ctx context.Context, tx gf_sqlite3.Queryer, table string) ([]SchemaUniqueKey, error) {
	sql := `--sql query unique key information
SELECT 
    pil."seq" AS Seq,
    pil."name" AS Name,
    pil."origin" = "c" AS Named,
    pii."name" AS ColName
FROM pragma_index_info(pil.name) AS pii
    JOIN pragma_index_list(?) AS pil
WHERE pil."unique" AND (pil."origin" = "c" OR pil."origin" = "u")
ORDER BY pil."seq", pii."seqno"`
	rows, err := tx.QueryxContext(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get unique keys of %s: %w`, table, err)
	}
	type ukRow struct {
		Seq     int64  `db:"Seq"`
		Name    string `db:"Name"`
		Named   bool   `db:"Named"`
		ColName string `db:"ColName"`
	}
	ukRows, err := gf_sqlite3.ScanRowsStruct[ukRow](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get unique keys of %s: %w`, table, err)
	}
	group := lo.GroupBy(ukRows, func(ukRow ukRow) int64 { return ukRow.Seq })
	groupIDs := lo.MapToSlice(group, func(id int64, _ []ukRow) int64 { return id })
	slices.Sort(groupIDs)

	var uniqueKeys []SchemaUniqueKey
	for _, id := range groupIDs {
		g := group[id]
		uk := SchemaUniqueKey{
			Key: lo.Map(g, func(ukRow ukRow, _ int) string { return ukRow.ColName }),
		}
		if g[0].Named {
			uk.Name = g[0].Name
		}

		uniqueKeys = append(uniqueKeys, uk)
	}

	return uniqueKeys, nil
}
