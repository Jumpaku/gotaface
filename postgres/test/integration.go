package test

import (
	"context"
	"flag"
	"testing"

	"github.com/jackc/pgx/v5"
)

var DataSource = flag.String("data-source", "", "data source name")

func SkipWithoutDataSource(t *testing.T) {
	if *DataSource == "" {
		t.Skipf("flag %q is not specified", "data-source")
	}
}

func Setup(t *testing.T, connectionString string, dbName string) (db *pgx.Conn, teardown func()) {
	t.Helper()

	SkipWithoutDataSource(t)
	ctx := context.Background()
	{
		conn, err := pgx.Connect(ctx, connectionString)
		if err != nil {
			t.Fatalf(`fail to create postgres admin client: %v`, err)
		}
		defer conn.Close(ctx)
		_, err = conn.Exec(ctx, `CREATE DATABASE `+dbName)
		if err != nil {
			t.Fatalf(`fail to create database: %v`, err)
		}
	}

	config, err := pgx.ParseConfig(connectionString)
	if err != nil {
		t.Fatalf(`fail to parse config: %v`, err)
	}
	config.Database = dbName
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		t.Fatalf(`fail to create postgres admin client: %v`, err)
	}

	teardown = func() {
		conn.Close(context.Background())
	}

	return conn, teardown
}

func InitDDLs(t *testing.T, db *pgx.Conn, ddls []string) {
	t.Helper()

	ctx := context.Background()

	for _, stmt := range ddls {
		if _, err := db.Exec(ctx, stmt); err != nil {
			t.Fatalf(`fail to execute ddl: %v`, err)
		}
	}
}
