package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"cloud.google.com/go/spanner"
	cyamli_schema "github.com/Jumpaku/cyamli/schema"
	"github.com/Jumpaku/gotaface/spanner/schema"
)

//go:generate go run "github.com/Jumpaku/cyamli/cmd/cyamli@latest" golang -schema-path=cli.yaml -out-path=cli.gen.go
func main() {
	cli := NewCLI()
	cli.Func = fetch
	if err := Run(cli, os.Args); err != nil {
		log.Fatalf("%+v\n", err)
	}
}

//go:embed schema.sql.tpl
var schemaSQL string

func fetch(cmdSchema *cyamli_schema.Command, subcommand []string, input CLI_Input) (err error) {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, input.Arg_DataSource)
	if err != nil {
		return fmt.Errorf("fail to create Spanner client: %w", err)
	}
	defer client.Close()

	fetcher := schema.NewFetcher(client.ReadOnlyTransaction())

	schemas := []schema.SchemaTable{}
	for _, targetTable := range input.Arg_TargetTables {
		result, err := fetcher.Fetch(ctx, targetTable)
		if err != nil {
			return fmt.Errorf("fail to fetch schema of %q in Spanner database: %w", targetTable, err)
		}
		schemas = append(schemas, result)
	}

	switch input.Opt_Format {
	default:
		return fmt.Errorf("invalid option value for format, which must be one of json, txt.tpl, sql")
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		for _, schema := range schemas {
			if err := encoder.Encode(schema); err != nil {
				return fmt.Errorf("fail to encode schema of %q into JSON: %w", schema.Name, err)
			}
		}
	case "txt.tpl":
		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("fail to read from stdin: %w", err)
		}
		executor, err := template.New("txt.tpl").Parse(string(stdinBytes))
		if err != nil {
			return fmt.Errorf("fail to parse text template: %w", err)
		}
		for _, schema := range schemas {
			if err := executor.Execute(os.Stdout, schema); err != nil {
				return fmt.Errorf("fail to process template with schema of %q: %w", schema.Name, err)
			}
		}
	case "sql":
		executor := template.Must(template.New("sql.tpl").Parse(schemaSQL))
		for _, schema := range schemas {
			if err := executor.Execute(os.Stdout, schema); err != nil {
				return fmt.Errorf("fail to process template with schema of %q: %w", schema.Name, err)
			}
		}
	}

	return nil
}
