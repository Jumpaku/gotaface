package testdata

import _ "embed"

//go:embed ddl_00_all_types.sql
var DDL00AllTypesSQL string

//go:embed ddl_02_foreign_keys.sql
var DDL02ForeignKeysSQL string

//go:embed ddl_03_foreign_loop_1.sql
var DDL03ForeignLoop1SQL string

//go:embed ddl_04_foreign_loop_2.sql
var DDL04ForeignLoop2SQL string

//go:embed ddl_05_foreign_loop_3.sql
var DDL05ForeignLoop3SQL string

//go:embed ddl_06_unique_keys_index.sql
var DDL06UniqueKeysIndexSQL string

//go:embed ddl_07_unique_keys_constraint.sql
var DDL07UniqueKeysConstraintSQL string

//go:embed ddl_08_unique_keys_column.sql
var DDL08UniqueKeysColumnSQL string
