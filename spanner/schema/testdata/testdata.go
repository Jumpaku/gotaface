package testdata

import _ "embed"

//go:embed ddl_00_all_types.sql
var DDL00AllTypesSQL string

//go:embed ddl_01_interleave.sql
var DDL01InterleaveSQL string

//go:embed ddl_02_foreign_keys.sql
var DDL02ForeignKeysSQL string

//go:embed ddl_03_foreign_loop_1.sql
var DDL03ForeignLoop1SQL string

//go:embed ddl_04_foreign_loop_2.sql
var DDL04ForeignLoop2SQL string

//go:embed ddl_05_foreign_loop_3.sql
var DDL05ForeignLoop3SQL string

//go:embed ddl_06_unique_keys.sql
var DDL06UniqueKeysSQL string
