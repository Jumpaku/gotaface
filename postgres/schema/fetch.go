package schema

import (
	"context"
	"fmt"
	"slices"

	"github.com/Jumpaku/go-assert"
	gf_postgres "github.com/Jumpaku/gotaface/postgres"
	"github.com/Jumpaku/gotaface/schema"
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
	queryer gf_postgres.Queryer
}

func NewFetcher(queryer gf_postgres.Queryer) fetcher {
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

func queryColumns(ctx context.Context, tx gf_postgres.Queryer, table string) ([]SchemaColumn, error) {
	sql := `--sql query column information
SELECT 
	column_name AS "Name",
	data_type AS "Type",
	is_nullable = 'YES' AS "Nullable"
FROM information_schema.columns
WHERE table_name = $1
ORDER BY ordinal_position`
	rows, err := tx.Query(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get columns of %s: %w`, table, err)
	}
	type column struct {
		Name     string `db:"Name"`
		Type     string `db:"Type"`
		Nullable bool   `db:"Nullable"`
	}
	columns, err := gf_postgres.ScanRowsStruct[column](rows)
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

func queryPrimaryKey(ctx context.Context, tx gf_postgres.Queryer, table string) ([]string, error) {
	sql := `--sql query primary key information
SELECT
    kcu.column_name AS "Name"
FROM information_schema.table_constraints AS tc
     JOIN information_schema.key_column_usage AS kcu
          ON kcu.constraint_name = tc.constraint_name
WHERE kcu.table_name = $1 AND tc.constraint_type = 'PRIMARY KEY'
ORDER BY kcu.ordinal_position;`
	type key struct {
		Name string `db:"Name"`
	}
	rows, err := tx.Query(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get primary key of %s: %w`, table, err)
	}
	primaryKey, err := gf_postgres.ScanRowsStruct[key](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get primary key of %s: %w`, table, err)
	}
	return lo.Map(primaryKey, func(it key, i int) string { return it.Name }), nil
}

func queryForeignKeys(ctx context.Context, tx gf_postgres.Queryer, table string) ([]SchemaForeignKey, error) {
	sql := `--sql query foreign key information
SELECT
    tc.constraint_name AS "Name",
    ctu.table_name AS "ReferencedTable",
    kcu1.column_name AS "ReferencingKey",
    kcu2.column_name AS "ReferencedKey"
FROM
    information_schema.table_constraints tc
        JOIN information_schema.referential_constraints rc
            ON rc.constraint_name = tc.constraint_name
        JOIN information_schema.constraint_table_usage ctu
            ON ctu.constraint_name = rc.unique_constraint_name
        JOIN information_schema.key_column_usage kcu1
            ON kcu1.constraint_name = rc.constraint_name
        JOIN information_schema.key_column_usage kcu2
            ON kcu2.constraint_name = rc.unique_constraint_name
                AND kcu2.ordinal_position = kcu1.ordinal_position
WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = $1
ORDER BY "Name", kcu1.ordinal_position;

`
	rows, err := tx.Query(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get foreign keys of %s: %w`, table, err)
	}
	type fkRow struct {
		Name            string `db:"Name"`
		ReferencedTable string `db:"ReferencedTable"`
		ReferencingKey  string `db:"ReferencingKey"`
		ReferencedKey   string `db:"ReferencedKey"`
	}
	fkRows, err := gf_postgres.ScanRowsStruct[fkRow](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get foreign keys of %s: %w`, table, err)
	}

	group := lo.GroupBy(fkRows, func(fkRow fkRow) string { return fkRow.Name })
	groupNames := lo.Keys(group)
	slices.Sort(groupNames)

	var foreignKeys []SchemaForeignKey
	for _, id := range groupNames {
		g := group[id]
		foreignKeys = append(foreignKeys, SchemaForeignKey{
			ReferencedTable: g[0].ReferencedTable,
			ReferencedKey:   lo.Map(g, func(fkRow fkRow, _ int) string { return fkRow.ReferencedKey }),
			ReferencingKey:  lo.Map(g, func(fkRow fkRow, _ int) string { return fkRow.ReferencingKey }),
		})
	}

	return foreignKeys, nil
}

func queryUniqueKeys(ctx context.Context, tx gf_postgres.Queryer, table string) ([]SchemaUniqueKey, error) {
	sql := `--sql query primary key information
SELECT
    tc.constraint_name AS "Name",
    kcu.column_name AS "ColumnName"
FROM information_schema.table_constraints AS tc
	 JOIN information_schema.key_column_usage AS kcu
		  ON kcu.constraint_name = tc.constraint_name
WHERE kcu.table_name = $1 AND tc.constraint_type = 'UNIQUE'
ORDER BY tc.constraint_name, kcu.ordinal_position;`
	rows, err := tx.Query(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get unique keys of %s: %w`, table, err)
	}
	type ukRow struct {
		Name       string `db:"Name"`
		ColumnName string `db:"ColumnName"`
	}
	ukRows, err := gf_postgres.ScanRowsStruct[ukRow](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get unique keys of %s: %w`, table, err)
	}
	group := lo.GroupBy(ukRows, func(ukRow ukRow) string { return ukRow.Name })
	groupNames := lo.Keys(group)
	slices.Sort(groupNames)

	var uniqueKeys []SchemaUniqueKey
	for _, name := range groupNames {
		g := group[name]
		uk := SchemaUniqueKey{
			Key: lo.Map(g, func(ukRow ukRow, _ int) string { return ukRow.ColumnName }),
		}

		uniqueKeys = append(uniqueKeys, uk)
	}

	return uniqueKeys, nil
}
