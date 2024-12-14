package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/Jumpaku/gotaface/postgres/schema"
)

//go:generate go run "github.com/Jumpaku/cyamli/cmd/cyamli@latest" golang -schema-path=cli.yaml -out-path=cli.gen.go
var cli = NewCLI()

func main() {
	cli.FUNC = fetch
	if err := Run(cli, os.Args); err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func fetch(subcommand []string, input CLI_Input, inputErr error) (err error) {
	if inputErr != nil {
		fmt.Fprintln(os.Stderr, cli.DESC_Simple())
		return fmt.Errorf("fail to resolve command line arguments: %w", inputErr)
	}
	if input.Opt_Help {
		fmt.Fprintln(os.Stderr, cli.DESC_Detail())
		return nil
	}
	ctx := context.Background()
	dbx, err := pgx.Connect(ctx, input.Arg_DataSource)
	if err != nil {
		return fmt.Errorf("fail to open SQLite3 database: %w", err)
	}
	defer dbx.Close(ctx)

	fetcher := schema.NewFetcher(dbx)

	schemas := []schema.SchemaTable{}
	for _, targetTable := range input.Arg_TargetTables {
		result, err := fetcher.Fetch(ctx, targetTable)
		if err != nil {
			return fmt.Errorf("fail to fetch schema of %q in Spanner database: %w", targetTable, err)
		}
		schemas = append(schemas, result)
	}

	var out io.Writer = os.Stdout
	if input.Opt_Output != "" {
		f, err := os.Create(input.Opt_Output)
		if err != nil {
			return fmt.Errorf("fail to open output file %q: %w", input.Opt_Output, err)
		}
		defer f.Close()

		out = f
	}

	switch input.Opt_Format {
	default:
		return fmt.Errorf("invalid option value for format, which must be one of json, txt.tpl, sql")
	case "json":
		encoder := json.NewEncoder(out)
		for _, schema := range schemas {
			if err := encoder.Encode(schema); err != nil {
				return fmt.Errorf("fail to encode schema of %q into JSON: %w", schema.Name, err)
			}
		}
	case "txt.tpl":
		var in io.Reader = os.Stdin
		if input.Opt_InputTxtTpl != "" {
			f, err := os.Open(input.Opt_InputTxtTpl)
			if err != nil {
				return fmt.Errorf("fail to open output file %q: %w", input.Opt_InputTxtTpl, err)
			}
			defer f.Close()

			in = f
		}

		inBytes, err := io.ReadAll(in)
		if err != nil {
			return fmt.Errorf("fail to read from stdin: %w", err)
		}
		executor, err := template.New("txt.tpl").Parse(string(inBytes))
		if err != nil {
			return fmt.Errorf("fail to parse text template: %w", err)
		}
		for _, schema := range schemas {
			if err := executor.Execute(out, schema); err != nil {
				return fmt.Errorf("fail to process template with schema of %q: %w", schema.Name, err)
			}
		}
	}

	return nil
}
