package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	gf_sqlite3 "github.com/Jumpaku/gotaface/sqlite3"
	"github.com/jmoiron/sqlx"
)

func SkipIfNoEnv(t *testing.T) {
	t.Helper()

	if GetEnvSQLite3() == "" {
		t.Skipf(`environment variable %s are required`, EnvTestSQLite3)
	}
}

func Setup(t *testing.T, database string) (db *sqlx.DB, teardown func()) {
	t.Helper()

	SkipIfNoEnv(t)

	db, err := sqlx.Open("sqlite3", database)
	if err != nil {
		t.Fatalf(`fail to create spanner admin client: %v`, err)
	}

	teardown = func() {
		db.Close()
		os.Remove(database)
	}

	return db, teardown
}

func InitDDLs(t *testing.T, db *sqlx.DB, stmts []string) {
	t.Helper()

	ctx := context.Background()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		t.Fatalf(`fail to begin transaction: %v`, err)
	}
	defer tx.Rollback()

	for _, stmt := range stmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			t.Fatalf(`fail to execute ddl: %v`, err)
		}
	}

	tx.Commit()
}

func InitDMLs(t *testing.T, db *sqlx.DB, stmts []string) {
	t.Helper()

	ctx := context.Background()
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		t.Fatalf(`fail to begin transaction: %v`, err)
	}
	defer tx.Rollback()

	for _, stmt := range stmts {
		_, err := tx.ExecContext(ctx, stmt)
		t.Fatalf(`fail to execute ddl: %v`, err)
	}

	tx.Commit()
}

func ListRows[Row any](t *testing.T, tx gf_sqlite3.Queryer, from string) []*Row {
	t.Helper()

	itr, err := tx.QueryxContext(context.Background(), fmt.Sprintf(`SELECT * FROM %s`, from))
	if err != nil {
		t.Fatalf(`fail to query row: %v`, err)
	}

	rowsStruct := []*Row{}
	for itr.Next() {
		var row Row
		if err := itr.StructScan(&row); err != nil {
			t.Fatalf(`fail to scan row: %v`, err)
		}

		rowsStruct = append(rowsStruct, &row)
	}

	return rowsStruct
}

func FindRow[Row any](t *testing.T, tx gf_sqlite3.Queryer, from string, where map[string]any) *Row {
	t.Helper()

	cond := " TRUE"
	args := []any{}
	for key, val := range where {
		cond += ` AND ` + key + ` = ?`
		args = append(args, val)
	}

	itr, err := tx.QueryxContext(context.Background(), fmt.Sprintf(`SELECT * FROM %s WHERE %s`, from, cond), args...)
	if err != nil {
		t.Fatalf(`fail to query row: %v`, err)
	}

	rowsStruct := []*Row{}
	for itr.Next() {
		var row Row
		if err := itr.StructScan(&row); err != nil {
			t.Fatalf(`fail to scan row: %v`, err)
		}

		rowsStruct = append(rowsStruct, &row)
	}
	if len(rowsStruct) == 0 {
		return nil
	}

	return rowsStruct[0]
}
