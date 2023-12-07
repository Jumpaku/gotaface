package schema

import (
	"context"
	"fmt"

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
	Name            string   `json:"name"`
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

	schemaTable, err := getTable(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

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

func getTable(ctx context.Context, tx gf_sqlite3.Queryer, table string) (SchemaTable, error) {
	sql := `--sql query table name and parent information
SELECT
	TABLE_NAME AS Name,
	IFNULL(PARENT_TABLE_NAME, "") AS Parent,
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_NAME = ?`
	rows, err := tx.QueryxContext(ctx, sql, table)
	if err != nil {
		return SchemaTable{}, fmt.Errorf(`fail to get table %s: %w`, table, err)
	}
	found, err := gf_sqlite3.ScanRowsStruct[SchemaTable](rows)
	if err != nil {
		return SchemaTable{}, fmt.Errorf(`fail to get table %s: %w`, table, err)
	}
	if len(found) == 0 {
		return SchemaTable{}, fmt.Errorf("table %q not found", table)
	}
	return found[0], nil
}

func queryColumns(ctx context.Context, tx gf_sqlite3.Queryer, table string) ([]SchemaColumn, error) {
	sql := `--sql query column information
SELECT
	COLUMN_NAME AS Name,
	SPANNER_TYPE AS Type,
	(IS_NULLABLE = 'YES') AS Nullable,
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_NAME = ?
ORDER BY ORDINAL_POSITION`
	rows, err := tx.QueryxContext(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get columns of %s: %w`, table, err)
	}
	columns, err := gf_sqlite3.ScanRowsStruct[SchemaColumn](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get columns of %s: %w`, table, err)
	}

	return columns, nil
}

func queryPrimaryKey(ctx context.Context, tx gf_sqlite3.Queryer, table string) ([]string, error) {
	sql := `--sql query primary key information
SELECT
	kcu.COLUMN_NAME AS Name
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu
	JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS tc
	ON kcu.CONSTRAINT_NAME = tc.CONSTRAINT_NAME 
        AND kcu.TABLE_NAME = tc.TABLE_NAME
WHERE kcu.TABLE_NAME = ? AND tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
ORDER BY kcu.ORDINAL_POSITION`
	type PrimaryKey struct{ Name string }
	rows, err := tx.QueryxContext(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get primary key of %s: %w`, table, err)
	}
	primaryKey, err := gf_sqlite3.ScanRowsStruct[PrimaryKey](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get primary key of %s: %w`, table, err)
	}
	return lo.Map(primaryKey, func(it PrimaryKey, i int) string { return it.Name }), nil
}

func queryForeignKeys(ctx context.Context, tx gf_sqlite3.Queryer, table string) ([]SchemaForeignKey, error) {
	sql := `--sql query foreign key information
SELECT
	tc.CONSTRAINT_NAME AS Name,
	ctu.TABLE_NAME AS ReferencedTable,
	ARRAY(
		SELECT kcu.COLUMN_NAME
		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu 
		WHERE kcu.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
		ORDER BY kcu.ORDINAL_POSITION
	) AS ReferencingKey,
	ARRAY(
		SELECT kcu.COLUMN_NAME
		FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu 
		WHERE kcu.CONSTRAINT_NAME = rc.UNIQUE_CONSTRAINT_NAME
		ORDER BY kcu.ORDINAL_POSITION
	) AS ReferencedKey
FROM
	INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
	JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc ON rc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
	JOIN INFORMATION_SCHEMA.CONSTRAINT_TABLE_USAGE ctu ON ctu.CONSTRAINT_NAME = rc.UNIQUE_CONSTRAINT_NAME
WHERE tc.CONSTRAINT_TYPE = 'FOREIGN KEY' AND tc.TABLE_NAME = ?
ORDER BY Name`
	rows, err := tx.QueryxContext(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get foreign keys of %s: %w`, table, err)
	}
	foreignKeys, err := gf_sqlite3.ScanRowsStruct[SchemaForeignKey](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get foreign keys of %s: %w`, table, err)
	}
	return foreignKeys, nil
}

func queryUniqueKeys(ctx context.Context, tx gf_sqlite3.Queryer, table string) ([]SchemaUniqueKey, error) {
	sql := `--sql query unique key information
WITH
	EXCLUDE_FK_BACKING AS (
		SELECT rc.UNIQUE_CONSTRAINT_NAME AS Name
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
		JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc ON rc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
		JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc2 ON tc2.CONSTRAINT_NAME = rc.UNIQUE_CONSTRAINT_NAME AND tc2.CONSTRAINT_TYPE = 'UNIQUE'
		WHERE tc.CONSTRAINT_TYPE = 'FOREIGN KEY'
	)
SELECT
	idx.INDEX_NAME AS Name,
	ARRAY(
		SELECT idxc.COLUMN_NAME
		FROM INFORMATION_SCHEMA.INDEX_COLUMNS idxc
		WHERE idx.INDEX_NAME = idxc.INDEX_NAME
		ORDER BY idxc.ORDINAL_POSITION
	) AS Key
FROM INFORMATION_SCHEMA.INDEXES idx
WHERE
	idx.TABLE_NAME = ?
	AND idx.IS_UNIQUE
	AND INDEX_TYPE = "INDEX"
	AND idx.INDEX_NAME NOT IN (SELECT Name FROM EXCLUDE_FK_BACKING)
ORDER BY Name`
	rows, err := tx.QueryxContext(ctx, sql, table)
	if err != nil {
		return nil, fmt.Errorf(`fail to get unique keys of %s: %w`, table, err)
	}
	uniqueKeys, err := gf_sqlite3.ScanRowsStruct[SchemaUniqueKey](rows)
	if err != nil {
		return nil, fmt.Errorf(`fail to get unique keys of %s: %w`, table, err)
	}
	return uniqueKeys, nil
}
