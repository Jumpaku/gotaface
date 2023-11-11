// Code generated by cyamli v0.0.8, DO NOT EDIT.
package main

import (
	"fmt"
	"bytes"
	"strings"

	cyamli_schema "github.com/Jumpaku/cyamli/schema"
	cyamli_golang "github.com/Jumpaku/cyamli/golang"
)

func LoadSchema() *cyamli_schema.Schema {
	var schema, _ = cyamli_schema.Load(bytes.NewBufferString("name: gaf-spanner-fetch-schema\nversion: v0.0.2\ndescription: Fetches schema data from a table in a Spanner database.\noptions:\n  -format:\n    short: \"\"\n    description: \"Specifies output format:\\n * json: outputs in JSON format.\\n * sql: outputs as DML statements, which does not necessarily reproduce the identical tables.\\n * txt.tpl: processes template text from stdin in form Go's text/template and outputs result. Available data in the template is described in https://github.com/Jumpaku/gotaface/blob/main/spanner/schema/fetch.go#L25\\n\"\n    type: \"\"\n    default: json\n  -help:\n    short: -h\n    description: Shows help.\n    type: boolean\n    default: \"\"\narguments:\n- name: data_source\n  description: Specifies data source in form \"projects/<project>/instances/<instance>/databases/<database>\".\n  type: \"\"\n  variadic: false\n- name: target_tables\n  description: Specify target tables to be fetched schemas.\n  type: \"\"\n  variadic: true\nsubcommands: {}\n"))
	return schema
}


type Func[Input any] func(subcommand []string, input Input, inputErr error) (err error)




type CLI struct {

	Func Func[CLI_Input]
}

type CLI_Input struct {
	Opt_Format string
	Opt_Help bool

	Arg_DataSource string
	Arg_TargetTables []string

}








func NewCLI() CLI {
	return CLI{}
}


func Run(cli CLI, args []string) error {
	s := LoadSchema()
	cmd, subcommand, restArgs := cyamli_golang.ResolveSubcommand(s, args)
	switch strings.Join(subcommand, " ") {

	case "":
		input := CLI_Input{
			Opt_Format: "json",
			Opt_Help: false,

		}
		funcMethod := cli.Func
		if funcMethod == nil {
			return fmt.Errorf("%q is unsupported: cli.Func not assigned", "")
		}
		err := cyamli_golang.ResolveInput(cmd, restArgs, &input)
		return funcMethod(subcommand, input, err)


	}
	return nil
}
