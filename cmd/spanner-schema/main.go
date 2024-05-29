package main

import (
	"cloud.google.com/go/spanner"
	"context"
	"encoding/json"
	"fmt"
	cli2 "github.com/Jumpaku/gotaface/cli"
	"github.com/Jumpaku/gotaface/spanner/schema"
	"io"
	"os"
	"text/template"
)

//go:generate go run github.com/Jumpaku/cyamli@v1.0.0 generate golang -schema-path=cli.yaml -out-path=cli.gen.go
func main() {
	os.Exit(executor{}.Execute(os.Args, os.Stdin, os.Stdout, os.Stderr))
}

type executor struct{}

type implementation struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (e executor) Execute(args []string, stdin io.Reader, stdout, stderr io.Writer) (exitCode int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(stderr, `%+v`, err)
			exitCode = 1
		}
	}()

	var cli = NewCLI()

	i := implementation{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
	cli.FUNC = i.root
	cli.List.FUNC = i.list
	cli.Fetch.FUNC = i.fetch

	cli2.PanicIfErrorf(Run(cli, args), "Error")

	return 0
}

var _ cli2.Executor = executor{}

func (i implementation) root(subcommand []string, _ CLI_Input, _ error) (err error) {
	cli2.MustPrintf(i.stdout, GetDoc(subcommand))
	return nil
}

func (i implementation) list(subcommand []string, input CLI_List_Input, inputErr error) (err error) {
	if input.Opt_Help {
		cli2.MustPrintf(i.stdout, GetDoc(subcommand))
		return nil
	}
	cli2.PanicIfErrorf(inputErr, "fail to resolve command line arguments, to see help, run 'spanner-schema list -help'")

	c, err := spanner.NewClient(context.Background(), input.Arg_DataSource)
	cli2.PanicIfErrorf(err, "fail to create spanner client %q", input.Arg_DataSource)
	defer c.Close()

	tx := c.ReadOnlyTransaction()
	defer tx.Close()

	tables, err := schema.NewLister(tx).List(context.Background())
	cli2.PanicIfErrorf(err, "fail to list tables")

	for _, table := range tables {
		cli2.MustPrintf(i.stdout, "%s\n", table)
	}

	return nil
}

func (i implementation) fetch(subcommand []string, input CLI_Fetch_Input, inputErr error) (err error) {
	if input.Opt_Help {
		cli2.MustPrintf(i.stdout, GetDoc(subcommand))
		return nil
	}
	cli2.PanicIfErrorf(inputErr, "fail to resolve command line arguments, to see help, run 'spanner-schema fetch -help'")

	c, err := spanner.NewClient(context.Background(), input.Arg_DataSource)
	cli2.PanicIfErrorf(err, "fail to create spanner client %q", input.Arg_DataSource)
	defer c.Close()

	tx := c.ReadOnlyTransaction()
	defer tx.Close()

	schemas := schema.Schemas{}
	for _, table := range input.Arg_TargetTables {
		s, err := schema.NewFetcher(tx).Fetch(context.Background(), table)
		cli2.PanicIfErrorf(err, "fail to fetch table schema %q", table)

		schemas = append(schemas, s)
	}

	out := i.stdout
	if input.Opt_Output != "" {
		f, err := os.Create(input.Opt_Output)
		cli2.PanicIfErrorf(err, "fail to create file %q", input.Opt_Output)
		defer f.Close()
		out = f
	}

	switch input.Opt_Format {
	default:
		panic(fmt.Errorf("invalid option value %q for -format, which must be one of 'json', 'template'", input.Opt_Format))
	case "template":
		in := i.stdin
		if input.Opt_InputTemplate != "" {
			f, err := os.Open(input.Opt_InputTemplate)
			cli2.PanicIfErrorf(err, "fail to open file %q", input.Opt_InputTemplate)
			defer f.Close()
			in = f
		}

		b, err := io.ReadAll(in)
		cli2.PanicIfErrorf(err, "fail to read template from %q", input.Opt_InputTemplate)

		tmpl, err := template.New("schemas").Parse(string(b))
		cli2.PanicIfErrorf(err, "fail to read template from %q", input.Opt_InputTemplate)

		err = tmpl.Execute(out, schemas)
		cli2.PanicIfErrorf(err, "fail to execute template")
	case "json":
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		err := encoder.Encode(schemas)
		cli2.PanicIfErrorf(err, "fail to encode schemas into JSON")
	}

	return nil
}
