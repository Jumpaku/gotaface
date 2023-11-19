# gaf-spanner-fetch-schema

Fetches schema data from a table in a Spanner database.

## Usage:

```sh
gaf-spanner-fetch-schema [<option>|<argument>]... [-- [<argument>]...]
```

## Options:

* `-format=<string>` (`default="json"`):
    Specifies output format:
    * `json`: outputs in JSON format.
    * `sql`: outputs as DML statements. Note that the statements do not guarantee to reproduce the identical tables.
    * `txt.tpl`: processes template text from stdin in form Go's text/template and outputs result. Available data in the template can be found in https://github.com/Jumpaku/gotaface/blob/main/spanner/schema/fetch.go#L25 .

## Arguments:
*  `[0]` (`<data_source:string>`):
    Specifies data source in form `projects/<project>/instances/<instance>/databases/<database>`.

* `[1:]` (`[<target_tables:string>]...`):
    Specify target tables to be fetched schemas.

## Example
```bash
make example-spanner
go run ./spanner/cmd/gaf-spanner-fetch-schema projects/gotaface/instances/example/databases/db -format=txt.tpl A B_1 B_2 B_3 B_4 C_1 C_2 C_3 C_4 C_5 D_1 E_1 E_2 F_1 F_2 G < spanner/cmd/gaf-spanner-fetch-schema/schema.sql.tpl

cat spanner/schema/testdata/*.sql | spanner-cli -d db -p gotaface -i example
```