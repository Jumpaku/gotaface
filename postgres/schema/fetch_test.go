package schema_test

import (
	"context"
	_ "embed"
	"fmt"
	"testing"
	"time"

	"github.com/Jumpaku/gotaface/postgres/schema"
	"github.com/Jumpaku/gotaface/postgres/schema/testdata"
	"github.com/Jumpaku/gotaface/postgres/test"
	"github.com/stretchr/testify/assert"
)

var ddls = map[string]string{
	"ddl_00_all_types":              testdata.DDL00AllTypesSQL,
	"ddl_02_foreign_keys":           testdata.DDL02ForeignKeysSQL,
	"ddl_03_foreign_loop_1":         testdata.DDL03ForeignLoop1SQL,
	"ddl_04_foreign_loop_2":         testdata.DDL04ForeignLoop2SQL,
	"ddl_05_foreign_loop_3":         testdata.DDL05ForeignLoop3SQL,
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
					{Name: "PK", Type: "integer", Nullable: false},
					{Name: "Col_01", Type: "bigint", Nullable: true},
					{Name: "Col_02", Type: "bigint", Nullable: false},
					{Name: "Col_04", Type: "bigint", Nullable: false},
					{Name: "Col_05", Type: "bit", Nullable: true},
					{Name: "Col_06", Type: "bit", Nullable: false},
					{Name: "Col_07", Type: "bit varying", Nullable: true},
					{Name: "Col_08", Type: "bit varying", Nullable: false},
					{Name: "Col_09", Type: "boolean", Nullable: true},
					{Name: "Col_10", Type: "boolean", Nullable: false},
					{Name: "Col_11", Type: "bytea", Nullable: true},
					{Name: "Col_12", Type: "bytea", Nullable: false},
					{Name: "Col_13", Type: "character", Nullable: true},
					{Name: "Col_14", Type: "character", Nullable: false},
					{Name: "Col_15", Type: "character varying", Nullable: true},
					{Name: "Col_16", Type: "character varying", Nullable: false},
					{Name: "Col_17", Type: "date", Nullable: true},
					{Name: "Col_18", Type: "date", Nullable: false},
					{Name: "Col_19", Type: "double precision", Nullable: true},
					{Name: "Col_20", Type: "double precision", Nullable: false},
					{Name: "Col_21", Type: "integer", Nullable: true},
					{Name: "Col_22", Type: "integer", Nullable: false},
					{Name: "Col_23", Type: "json", Nullable: true},
					{Name: "Col_24", Type: "json", Nullable: false},
					{Name: "Col_25", Type: "money", Nullable: true},
					{Name: "Col_26", Type: "money", Nullable: false},
					{Name: "Col_27", Type: "numeric", Nullable: true},
					{Name: "Col_28", Type: "numeric", Nullable: false},
					{Name: "Col_29", Type: "real", Nullable: true},
					{Name: "Col_30", Type: "real", Nullable: false},
					{Name: "Col_31", Type: "smallint", Nullable: true},
					{Name: "Col_32", Type: "smallint", Nullable: false},
					{Name: "Col_34", Type: "smallint", Nullable: false},
					{Name: "Col_36", Type: "integer", Nullable: false},
					{Name: "Col_37", Type: "text", Nullable: true},
					{Name: "Col_38", Type: "text", Nullable: false},
					{Name: "Col_39", Type: "time without time zone", Nullable: true},
					{Name: "Col_40", Type: "time without time zone", Nullable: false},
					{Name: "Col_41", Type: "time with time zone", Nullable: true},
					{Name: "Col_42", Type: "time with time zone", Nullable: false},
					{Name: "Col_43", Type: "timestamp without time zone", Nullable: true},
					{Name: "Col_44", Type: "timestamp without time zone", Nullable: false},
					{Name: "Col_45", Type: "timestamp with time zone", Nullable: true},
					{Name: "Col_46", Type: "timestamp with time zone", Nullable: false},
					{Name: "Col_47", Type: "uuid", Nullable: true},
					{Name: "Col_48", Type: "uuid", Nullable: false},
					{Name: "Col_49", Type: "xml", Nullable: true},
					{Name: "Col_50", Type: "xml", Nullable: false},
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
					{Name: "PK_11", Type: "integer"},
					{Name: "PK_12", Type: "integer"},
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
					{Name: "PK_21", Type: "integer"},
					{Name: "PK_22", Type: "integer"},
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
					{Name: "PK_31", Type: "integer"},
					{Name: "PK_32", Type: "integer"},
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
					{Name: "PK_41", Type: "integer"},
					{Name: "PK_42", Type: "integer"},
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
					{Name: "PK_51", Type: "integer"},
					{Name: "PK_52", Type: "integer"},
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
					{Name: "PK_11", Type: "integer"},
					{Name: "PK_12", Type: "integer"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						ReferencedTable: "D_1",
						ReferencedKey:   []string{"PK_12"},
						ReferencingKey:  []string{"PK_11"},
					},
				},
				UniqueKeys: []schema.SchemaUniqueKey{
					{Name: "", Key: []string{"PK_12"}},
				},
			},
		},
		{
			ddl:   "ddl_04_foreign_loop_2",
			table: "E_1",
			want: schema.SchemaTable{
				Name: "E_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "integer"},
					{Name: "PK_12", Type: "integer"},
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
					{Name: "PK_21", Type: "integer"},
					{Name: "PK_22", Type: "integer"},
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
					{Name: "PK_11", Type: "integer"},
					{Name: "PK_12", Type: "integer"},
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
					{Name: "PK_21", Type: "integer"},
					{Name: "PK_22", Type: "integer"},
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
					{Name: "PK_31", Type: "integer"},
					{Name: "PK_32", Type: "integer"},
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
			ddl:   "ddl_07_unique_keys_constraint",
			table: "H",
			want: schema.SchemaTable{
				Name: "H",
				Columns: []schema.SchemaColumn{
					{Name: "PK", Type: "integer"},
					{Name: "C1", Type: "integer"},
					{Name: "C2", Type: "integer"},
					{Name: "C3", Type: "integer"},
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
					{Name: "PK", Type: "integer"},
					{Name: "C1", Type: "integer"},
					{Name: "C2", Type: "integer"},
					{Name: "C3", Type: "integer"},
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
			now := time.Now().Unix()
			dbName := fmt.Sprintf("test_fetcher_%03d_%d", number, now)
			db, teardown := test.Setup(t, *test.DataSource, dbName)
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
