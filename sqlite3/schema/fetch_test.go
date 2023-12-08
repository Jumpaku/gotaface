package schema_test

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/Jumpaku/gotaface/sqlite3/schema"
	"github.com/Jumpaku/gotaface/sqlite3/schema/testdata"
	"github.com/Jumpaku/gotaface/sqlite3/test"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var ddls = map[string]string{
	"ddl_00_all_types":              testdata.DDL00AllTypesSQL,
	"ddl_02_foreign_keys":           testdata.DDL02ForeignKeysSQL,
	"ddl_03_foreign_loop_1":         testdata.DDL03ForeignLoop1SQL,
	"ddl_04_foreign_loop_2":         testdata.DDL04ForeignLoop2SQL,
	"ddl_05_foreign_loop_3":         testdata.DDL05ForeignLoop3SQL,
	"ddl_06_unique_keys_index":      testdata.DDL06UniqueKeysIndexSQL,
	"ddl_07_unique_keys_constraint": testdata.DDL07UniqueKeysConstraintSQL,
	"ddl_08_unique_keys_column":     testdata.DDL08UniqueKeysColumnSQL,
}

func TestFetcher(t *testing.T) {
	testcases := []struct {
		ddl   string
		table string
		want  schema.SchemaTable
	}{
		{
			ddl:   "ddl_00_all_types",
			table: "A",
			want: schema.SchemaTable{
				Name: "A",
				Columns: []schema.SchemaColumn{
					{Name: "PK", Type: "INT64", Nullable: false},
					{Name: "Col_01", Type: "BOOL", Nullable: true},
					{Name: "Col_02", Type: "BOOL", Nullable: false},
					{Name: "Col_03", Type: "BYTES(50)", Nullable: true},
					{Name: "Col_04", Type: "BYTES(50)", Nullable: false},
					{Name: "Col_05", Type: "DATE", Nullable: true},
					{Name: "Col_06", Type: "DATE", Nullable: false},
					{Name: "Col_07", Type: "FLOAT64", Nullable: true},
					{Name: "Col_08", Type: "FLOAT64", Nullable: false},
					{Name: "Col_09", Type: "INT64", Nullable: true},
					{Name: "Col_10", Type: "INT64", Nullable: false},
					{Name: "Col_11", Type: "JSON", Nullable: true},
					{Name: "Col_12", Type: "JSON", Nullable: false},
					{Name: "Col_13", Type: "NUMERIC", Nullable: true},
					{Name: "Col_14", Type: "NUMERIC", Nullable: false},
					{Name: "Col_15", Type: "STRING(50)", Nullable: true},
					{Name: "Col_16", Type: "STRING(50)", Nullable: false},
					{Name: "Col_17", Type: "TIMESTAMP", Nullable: true},
					{Name: "Col_18", Type: "TIMESTAMP", Nullable: false},
				},
				PrimaryKey: []string{"PK"},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_1",
			want: schema.SchemaTable{
				Name: "C_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_2",
			want: schema.SchemaTable{
				Name: "C_2",
				Columns: []schema.SchemaColumn{
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_22", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_21", "PK_22"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "C_1",
						ReferencedKey:   []string{"PK_11", "PK_12"},
						ReferencingKey:  []string{"PK_21", "PK_22"},
					},
				},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_3",
			want: schema.SchemaTable{
				Name: "C_3",
				Columns: []schema.SchemaColumn{
					{Name: "PK_31", Type: "INT64"},
					{Name: "PK_32", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_31", "PK_32"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "C_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_31", "PK_32"},
					},
				},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_4",
			want: schema.SchemaTable{
				Name: "C_4",
				Columns: []schema.SchemaColumn{
					{Name: "PK_41", Type: "INT64"},
					{Name: "PK_42", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_41", "PK_42"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "C_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_41", "PK_42"},
					},
				},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_5",
			want: schema.SchemaTable{
				Name: "C_5",
				Columns: []schema.SchemaColumn{
					{Name: "PK_51", Type: "INT64"},
					{Name: "PK_52", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_51", "PK_52"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "C_4",
						ReferencedKey:   []string{"PK_41", "PK_42"},
						ReferencingKey:  []string{"PK_51", "PK_52"},
					},
					{
						ReferencedTable: "C_3",
						ReferencedKey:   []string{"PK_31", "PK_32"},
						ReferencingKey:  []string{"PK_51", "PK_52"},
					},
				},
			},
		},
		{
			ddl:   "ddl_03_foreign_loop_1",
			table: "D_1",
			want: schema.SchemaTable{
				Name: "D_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "D_1",
						ReferencedKey:   []string{"PK_12"},
						ReferencingKey:  []string{"PK_11"},
					},
				},
			},
		},
		{
			ddl:   "ddl_04_foreign_loop_2",
			table: "E_1",
			want: schema.SchemaTable{
				Name: "E_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "E_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_11", "PK_12"},
					},
				},
			},
		},
		{
			ddl:   "ddl_04_foreign_loop_2",
			table: "E_2",
			want: schema.SchemaTable{
				Name: "E_2",
				Columns: []schema.SchemaColumn{
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_22", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_21", "PK_22"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "E_1",
						ReferencedKey:   []string{"PK_11", "PK_12"},
						ReferencingKey:  []string{"PK_21", "PK_22"},
					},
				},
			},
		},
		{
			ddl:   "ddl_05_foreign_loop_3",
			table: "F_1",
			want: schema.SchemaTable{
				Name: "F_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "F_3",
						ReferencedKey:   []string{"PK_31", "PK_32"},
						ReferencingKey:  []string{"PK_11", "PK_12"},
					},
				},
			},
		},
		{
			ddl:   "ddl_05_foreign_loop_3",
			table: "F_2",
			want: schema.SchemaTable{
				Name: "F_2",
				Columns: []schema.SchemaColumn{
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_22", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_21", "PK_22"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "F_1",
						ReferencedKey:   []string{"PK_11", "PK_12"},
						ReferencingKey:  []string{"PK_21", "PK_22"},
					},
				},
			},
		},
		{
			ddl:   "ddl_05_foreign_loop_3",
			table: "F_3",
			want: schema.SchemaTable{
				Name: "F_3",
				Columns: []schema.SchemaColumn{
					{Name: "PK_31", Type: "INT64"},
					{Name: "PK_32", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_31", "PK_32"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "F_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_31", "PK_32"},
					},
				},
			},
		},
		{
			ddl:   "ddl_06_unique_keys_index",
			table: "G",
			want: schema.SchemaTable{
				Name: "G",
				Columns: []schema.SchemaColumn{
					{Name: "PK", Type: "INT64"},
					{Name: "C1", Type: "INT64"},
					{Name: "C2", Type: "INT64"},
					{Name: "C3", Type: "INT64"},
				},
				PrimaryKey: []string{"PK"},
				UniqueKeys: []schema.SchemaUniqueKey{
					{Name: "UQ_G_C1", Key: []string{"C1"}},
					{Name: "UQ_G_C1_C2", Key: []string{"C1", "C2"}},
					{Name: "UQ_G_C1_C2_C3", Key: []string{"C1", "C2", "C3"}},
					{Name: "UQ_G_C1_C3", Key: []string{"C1", "C3"}},
					{Name: "UQ_G_C1_C3_C2", Key: []string{"C1", "C3", "C2"}},
					{Name: "UQ_G_C2", Key: []string{"C2"}},
					{Name: "UQ_G_C2_C1", Key: []string{"C2", "C1"}},
					{Name: "UQ_G_C2_C1_C3", Key: []string{"C2", "C1", "C3"}},
					{Name: "UQ_G_C2_C3", Key: []string{"C2", "C3"}},
					{Name: "UQ_G_C2_C3_C1", Key: []string{"C2", "C3", "C1"}},
					{Name: "UQ_G_C3", Key: []string{"C3"}},
					{Name: "UQ_G_C3_C1", Key: []string{"C3", "C1"}},
					{Name: "UQ_G_C3_C1_C2", Key: []string{"C3", "C1", "C2"}},
					{Name: "UQ_G_C3_C2", Key: []string{"C3", "C2"}},
					{Name: "UQ_G_C3_C2_C1", Key: []string{"C3", "C2", "C1"}},
				},
			},
		},
		{
			ddl:   "ddl_07_unique_keys_constraint",
			table: "H",
			want: schema.SchemaTable{
				Name: "H",
				Columns: []schema.SchemaColumn{
					{Name: "PK", Type: "INT64"},
					{Name: "C1", Type: "INT64"},
					{Name: "C2", Type: "INT64"},
					{Name: "C3", Type: "INT64"},
				},
				PrimaryKey: []string{"PK"},
				UniqueKeys: []schema.SchemaUniqueKey{
					{Name: "", Key: []string{"C1"}},
					{Name: "", Key: []string{"C1", "C2"}},
					{Name: "", Key: []string{"C1", "C2", "C3"}},
					{Name: "", Key: []string{"C1", "C3"}},
					{Name: "", Key: []string{"C1", "C3", "C2"}},
					{Name: "", Key: []string{"C2"}},
					{Name: "", Key: []string{"C2", "C1"}},
					{Name: "", Key: []string{"C2", "C1", "C3"}},
					{Name: "", Key: []string{"C2", "C3"}},
					{Name: "", Key: []string{"C2", "C3", "C1"}},
					{Name: "", Key: []string{"C3"}},
					{Name: "", Key: []string{"C3", "C1"}},
					{Name: "", Key: []string{"C3", "C1", "C2"}},
					{Name: "", Key: []string{"C3", "C2"}},
					{Name: "", Key: []string{"C3", "C2", "C1"}},
				},
			},
		},
		{
			ddl:   "ddl_08_unique_keys_column",
			table: "I",
			want: schema.SchemaTable{
				Name: "I",
				Columns: []schema.SchemaColumn{
					{Name: "PK", Type: "INT64"},
					{Name: "C1", Type: "INT64"},
					{Name: "C2", Type: "INT64"},
					{Name: "C3", Type: "INT64"},
				},
				PrimaryKey: []string{"PK"},
				UniqueKeys: []schema.SchemaUniqueKey{
					{Name: "", Key: []string{"C1"}},
					{Name: "", Key: []string{"C2"}},
					{Name: "", Key: []string{"C3"}},
				},
			},
		},
	}

	for number, testcase := range testcases {
		t.Run(fmt.Sprintf("%03d:%s[%s]", number, testcase.ddl, testcase.table), func(t *testing.T) {
			db, teardown := test.Setup(t, fmt.Sprintf("fetcher_%0d.sqlite", number))
			defer teardown()

			test.InitDDLs(t, db, []string{ddls[testcase.ddl]})

			sut := schema.NewFetcher(db)
			got, err := sut.Fetch(context.Background(), testcase.table)
			assert.Nil(t, err)
			assertEqualSchemaTable(t, testcase.want, got)
		})
	}
}
func assertEqualSchemaTable(t *testing.T, want schema.SchemaTable, got schema.SchemaTable) {
	t.Helper()
	assert.Equal(t, want.Name, got.Name)
	assert.Equal(t, want.PrimaryKey, got.PrimaryKey)
	assert.Equal(t, want.Columns, got.Columns)
	assert.ElementsMatch(t, want.ForeignKeys, got.ForeignKeys)
	assert.ElementsMatch(t, want.UniqueKeys, got.UniqueKeys)
}
