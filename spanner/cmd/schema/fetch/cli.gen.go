// Code generated by cyamli 0.0.0, DO NOT EDIT.
package main

import (
	"fmt"
	"bytes"
	"strings"
	"os"

	cyamli_schema "github.com/Jumpaku/cyamli/schema"
	cyamli_golang "github.com/Jumpaku/cyamli/golang"
	cyamli_description "github.com/Jumpaku/cyamli/description"
)

func newSchema() *cyamli_schema.Schema {
	var schema, _ = cyamli_schema.Load(bytes.NewBufferString("name: gf-spanner-schema-fetch\nversion: v0.0.1\ndescription: Fetches schema data from a table in a Spanner database.\noptions:\n  -format:\n    short: -o\n    description: \"Specifies output format:\\n * json: outputs in JSON format.\\n * sql: outputs as DML statements, which does not necessarily reproduce the identical tables.\\n * txt.tpl: processes template text from stdin in form Go's text/template and outputs result. Available data in the template is described in https://github.com/Jumpaku/gotaface/blob/main/spanner/schema/fetch.go#L25\\n\"\n    type: \"\"\n    default: json\narguments:\n- name: data_source\n  description: Specifies data source in form \"projects/<project>/instances/<instance>/databases/<database>\"\n  type: \"\"\n  variadic: false\n- name: target_tables\n  description: Specify target tables to be fetched schemas\n  type: \"\"\n  variadic: true\nsubcommands: {}\n"))
	return schema
}


type Func[Input any] func(cmdSchema *cyamli_schema.Command, subcommand []string, input Input) (err error)




type CLI struct {

	Func Func[CLI_Input]
}

type CLI_Input struct {
	Opt_Format string

	Arg_DataSource string
	Arg_TargetTables []string

}








func NewCLI() CLI {
	cli := CLI{}
	s := newSchema()

	cli.Func = cyamli_golang.NewDefaultFunc[CLI_Input](s.Program.Name)


	return cli
}


func Run(cli CLI, args []string) error {
	s := newSchema()
	cmd, subcommand, restArgs := cyamli_golang.ResolveSubcommand(s, args)
	switch strings.Join(subcommand, " ") {

	case "":
		input := CLI_Input{
			Opt_Format: "json",

		}
		if err := cyamli_golang.ResolveInput(cmd, restArgs, &input); err != nil {
			descData := cyamli_description.CreateCommandData(s.Program.Name, subcommand, cmd)
			if err := cyamli_description.DescribeCommand(cyamli_description.SimpleExecutor(), descData, os.Stderr); err != nil {
				panic(fmt.Errorf("fail to create command description: %w", err))
			}
			fmt.Fprintln(os.Stderr, "")
			return fmt.Errorf("fail to resolve input: %w", err)
		}
		funcMethod := cli.Func
		if funcMethod == nil {
			return fmt.Errorf("%q is unsupported: cli.Func not assigned", "")
		}
		if err := funcMethod(cmd, subcommand, input); err != nil {
			return fmt.Errorf("cli.Func(input) failed: %w", err)
		}


	}
	return nil
}
