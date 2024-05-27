package schema_test

import (
	_ "embed"
	"github.com/Jumpaku/gotaface/spanner/test"
)

//go:embed testdata/ddl_00_all_types.sql
var DDL00AllTypesSQL string

//go:embed testdata/ddl_01_interleave.sql
var DDL01InterleaveSQL string

//go:embed testdata/ddl_02_foreign_keys.sql
var DDL02ForeignKeysSQL string

//go:embed testdata/ddl_03_foreign_loop_1.sql
var DDL03ForeignLoop1SQL string

//go:embed testdata/ddl_04_foreign_loop_2.sql
var DDL04ForeignLoop2SQL string

//go:embed testdata/ddl_05_foreign_loop_3.sql
var DDL05ForeignLoop3SQL string

//go:embed testdata/ddl_06_unique_keys.sql
var DDL06UniqueKeysSQL string

//go:embed testdata/ddl_07_generated_col.sql
var DDL07GeneratedColSQL string

//go:embed testdata/ddl_08_index.sql
var DDL08IndexSQL string

//go:embed testdata/ddl_09_view.sql
var DDL09ViewSQL string

var TestDDLs = map[string][]string{
	"ddl_00_all_types":      test.Split(DDL00AllTypesSQL),
	"ddl_01_interleave":     test.Split(DDL01InterleaveSQL),
	"ddl_02_foreign_keys":   test.Split(DDL02ForeignKeysSQL),
	"ddl_03_foreign_loop_1": test.Split(DDL03ForeignLoop1SQL),
	"ddl_04_foreign_loop_2": test.Split(DDL04ForeignLoop2SQL),
	"ddl_05_foreign_loop_3": test.Split(DDL05ForeignLoop3SQL),
	"ddl_06_unique_keys":    test.Split(DDL06UniqueKeysSQL),
	"ddl_07_generated_col":  test.Split(DDL07GeneratedColSQL),
	"ddl_08_index":          test.Split(DDL08IndexSQL),
	"ddl_09_view":           test.Split(DDL09ViewSQL),
}
