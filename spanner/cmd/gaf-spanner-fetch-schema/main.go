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
	"github.com/Jumpaku/cyamli/description"
	"github.com/Jumpaku/gotaface/spanner/schema"
)

//go:generate go run "github.com/Jumpaku/cyamli/cmd/cyamli@latest" golang -schema-path=cli.yaml -out-path=cli.gen.go
func main() {
	cli := NewCLI()
	cli.FUNC = fetch
	if err := Run(cli, os.Args); err != nil {
		log.Fatalf("%+v\n", err)
	}
}

//go:embed schema.sql.tpl
var schemaSQL string

func fetch(subcommand []string, input CLI_Input, inputErr error) (err error) {
	if inputErr != nil {
		showSimpleDescription(subcommand, os.Stderr)
		return fmt.Errorf("fail to resolve command line arguments: %w", inputErr)
	}
	if input.Opt_Help {
		showDetailDescription(subcommand, os.Stdout)
		return nil
	}
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

func showSimpleDescription(subcommand []string, writer io.Writer) {
	schema := LoadSchema()
	_ = description.DescribeCommand(
		description.SimpleExecutor(),
		description.CreateCommandData(schema.Program.Name, schema.Program.Version, subcommand, schema.Find(subcommand)),
		writer,
	)
}

func showDetailDescription(subcommand []string, writer io.Writer) {
	schema := LoadSchema()
	_ = description.DescribeCommand(
		description.DetailExecutor(),
		description.CreateCommandData(schema.Program.Name, schema.Program.Version, subcommand, schema.Find(subcommand)),
		writer,
	)
}
