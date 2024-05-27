package schema

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/Jumpaku/go-assert"
	"github.com/Jumpaku/gotaface/schema"
	gf_spanner "github.com/Jumpaku/gotaface/spanner"
	"github.com/samber/lo"
)

type fetcher struct {
	queryer gf_spanner.Queryer
}

func NewFetcher(queryer gf_spanner.Queryer) schema.Fetcher[Table] {
	return fetcher{queryer: queryer}
}

var _ schema.Fetcher[Table] = fetcher{}

func (fetcher fetcher) Fetch(ctx context.Context, table string) (Table, error) {
	wrapError := func(err error) (Table, error) {
		assert.Params(err != nil, "wrapped error must be not nil")
		return Table{}, fmt.Errorf(`fail to fetch schema of %s: %w`, table, err)
	}

	schemaTable, err := getTable(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

	schemaTable.Columns, err = queryColumns(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

	if schemaTable.View {
		return schemaTable, nil
	}

	schemaTable.PrimaryKey, err = queryPrimaryKey(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

	schemaTable.ForeignKeys, err = queryForeignKeys(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

	schemaTable.Indexes, err = queryIndexes(ctx, fetcher.queryer, table)
	if err != nil {
		return wrapError(err)
	}

	return schemaTable, nil
}

func getTable(ctx context.Context, tx gf_spanner.Queryer, table string) (Table, error) {
	sql := `--sql query table name and parent information
SELECT
	TABLE_NAME AS Name,
	(TABLE_TYPE = 'VIEW') AS IsView,
	IFNULL(PARENT_TABLE_NAME, "") AS Parent,
FROM INFORMATION_SCHEMA.TABLES
WHERE TABLE_NAME = @Table`
	found, err := gf_spanner.ScanRowsStruct[Table](tx.Query(ctx, spanner.Statement{
		SQL:    sql,
		Params: map[string]interface{}{"Table": table},
	}))
	if err != nil {
		return Table{}, fmt.Errorf(`fail to get table %s: %w`, table, err)
	}
	if len(found) == 0 {
		return Table{}, fmt.Errorf("table %q not found", table)
	}
	return found[0], nil
}

func queryColumns(ctx context.Context, tx gf_spanner.Queryer, table string) ([]Column, error) {
	sql := `--sql query column information
SELECT
	COLUMN_NAME AS Name,
	SPANNER_TYPE AS Type,
	(IS_NULLABLE = 'YES') AS Nullable,
	(IS_GENERATED = 'ALWAYS') AS Generated,
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_NAME = @Table
ORDER BY ORDINAL_POSITION`
	columns, err := gf_spanner.ScanRowsStruct[Column](tx.Query(ctx, spanner.Statement{
		SQL:    sql,
		Params: map[string]interface{}{"Table": table},
	}))
	if err != nil {
		return nil, fmt.Errorf(`fail to get columns of %s: %w`, table, err)
	}

	return columns, nil
}

func queryPrimaryKey(ctx context.Context, tx gf_spanner.Queryer, table string) ([]string, error) {
	sql := `--sql query primary key information
SELECT
	kcu.COLUMN_NAME AS Name
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu
	JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS tc
	ON kcu.CONSTRAINT_NAME = tc.CONSTRAINT_NAME 
        AND kcu.TABLE_NAME = tc.TABLE_NAME
WHERE kcu.TABLE_NAME = @Table AND tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
ORDER BY kcu.ORDINAL_POSITION`
	type PrimaryKey struct{ Name string }
	primaryKey, err := gf_spanner.ScanRowsStruct[PrimaryKey](tx.Query(ctx, spanner.Statement{
		SQL:    sql,
		Params: map[string]interface{}{"Table": table},
	}))
	if err != nil {
		return nil, fmt.Errorf(`fail to get primary key of %s: %w`, table, err)
	}
	return lo.Map(primaryKey, func(it PrimaryKey, i int) string { return it.Name }), nil
}

func queryForeignKeys(ctx context.Context, tx gf_spanner.Queryer, table string) ([]ForeignKey, error) {
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
WHERE tc.CONSTRAINT_TYPE = 'FOREIGN KEY' AND tc.TABLE_NAME = @Table
ORDER BY Name`
	foreignKeys, err := gf_spanner.ScanRowsStruct[ForeignKey](tx.Query(ctx, spanner.Statement{
		SQL:    sql,
		Params: map[string]interface{}{"Table": table},
	}))
	if err != nil {
		return nil, fmt.Errorf(`fail to get foreign keys of %s: %w`, table, err)
	}
	return foreignKeys, nil
}

func queryIndexes(ctx context.Context, tx gf_spanner.Queryer, table string) ([]Index, error) {
	sql := `--sql query unique key information
WITH EXCLUDE_FK_BACKING AS (
	SELECT rc.UNIQUE_CONSTRAINT_NAME AS Name
	FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
		JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc ON rc.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
		JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc2 ON tc2.CONSTRAINT_NAME = rc.UNIQUE_CONSTRAINT_NAME AND tc2.CONSTRAINT_TYPE = 'UNIQUE'
	WHERE tc.CONSTRAINT_TYPE = 'FOREIGN KEY'
)
SELECT
	idx.INDEX_NAME AS Name,
	idx.IS_UNIQUE AS IsUnique,
	ARRAY(
		SELECT AS STRUCT
			idxc.COLUMN_NAME AS Name,
			idxc.COLUMN_ORDERING = 'DESC' AS IsDesc,
		FROM INFORMATION_SCHEMA.INDEX_COLUMNS idxc
		WHERE idx.INDEX_NAME = idxc.INDEX_NAME AND idxc.ORDINAL_POSITION IS NOT NULL
		ORDER BY idxc.ORDINAL_POSITION
	) AS Key
FROM INFORMATION_SCHEMA.INDEXES idx
WHERE
	idx.TABLE_NAME = @Table
	AND INDEX_TYPE = 'INDEX'
	AND idx.INDEX_NAME NOT IN (SELECT Name FROM EXCLUDE_FK_BACKING)
ORDER BY Name`
	uniqueKeys, err := gf_spanner.ScanRowsStruct[Index](tx.Query(ctx, spanner.Statement{
		SQL:    sql,
		Params: map[string]interface{}{"Table": table},
	}))
	if err != nil {
		return nil, fmt.Errorf(`fail to get unique keys of %s: %w`, table, err)
	}
	return uniqueKeys, nil
}
