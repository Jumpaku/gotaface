package schema_test

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/Jumpaku/gotaface/spanner/schema"
	"github.com/Jumpaku/gotaface/spanner/schema/testdata"
	"github.com/Jumpaku/gotaface/spanner/test"
	"github.com/stretchr/testify/assert"
)

var ddls = map[string][]string{
	"ddl_00_all_types":      test.Split(testdata.DDL00AllTypesSQL),
	"ddl_01_interleave":     test.Split(testdata.DDL01InterleaveSQL),
	"ddl_02_foreign_keys":   test.Split(testdata.DDL02ForeignKeysSQL),
	"ddl_03_foreign_loop_1": test.Split(testdata.DDL03ForeignLoop1SQL),
	"ddl_04_foreign_loop_2": test.Split(testdata.DDL04ForeignLoop2SQL),
	"ddl_05_foreign_loop_3": test.Split(testdata.DDL05ForeignLoop3SQL),
}

func TestFetcher(t *testing.T) {
	testcases := []struct {
		ddl   string
		table string
		want  schema.SchemaTable
	}{
		{
			ddl:   "ddl_00_all_types",
			table: "T",
			want: schema.SchemaTable{
				Name: "T",
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
			ddl:   "ddl_01_interleave",
			table: "T_1",
			want: schema.SchemaTable{
				Name: "T_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11"},
			},
		},
		{
			ddl:   "ddl_01_interleave",
			table: "T_2",
			want: schema.SchemaTable{
				Name: "T_2",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_21", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_21"},
				Parent:     "T_1",
			},
		},
		{
			ddl:   "ddl_01_interleave",
			table: "T_3",
			want: schema.SchemaTable{
				Name: "T_3",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_31", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_21", "PK_31"},
				Parent:     "T_2",
			},
		},
		{
			ddl:   "ddl_01_interleave",
			table: "T_4",
			want: schema.SchemaTable{
				Name: "T_4",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_41", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_21", "PK_41"},
				Parent:     "T_2",
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "T_1",
			want: schema.SchemaTable{
				Name: "T_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "T_2",
			want: schema.SchemaTable{
				Name: "T_2",
				Columns: []schema.SchemaColumn{
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_22", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_21", "PK_22"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_2_1",
						ReferencedTable: "T_1",
						ReferencedKey:   []string{"PK_11", "PK_12"},
						ReferencingKey:  []string{"PK_21", "PK_22"},
					},
				},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "T_3",
			want: schema.SchemaTable{
				Name: "T_3",
				Columns: []schema.SchemaColumn{
					{Name: "PK_31", Type: "INT64"},
					{Name: "PK_32", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_31", "PK_32"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_3_2",
						ReferencedTable: "T_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_31", "PK_32"},
					},
				},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "T_4",
			want: schema.SchemaTable{
				Name: "T_4",
				Columns: []schema.SchemaColumn{
					{Name: "PK_41", Type: "INT64"},
					{Name: "PK_42", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_41", "PK_42"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_4_2",
						ReferencedTable: "T_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_41", "PK_42"},
					},
				},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "T_5",
			want: schema.SchemaTable{
				Name: "T_5",
				Columns: []schema.SchemaColumn{
					{Name: "PK_51", Type: "INT64"},
					{Name: "PK_52", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_51", "PK_52"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_5_3",
						ReferencedTable: "T_3",
						ReferencedKey:   []string{"PK_31", "PK_32"},
						ReferencingKey:  []string{"PK_51", "PK_52"},
					},
					{
						Name:            "FK_5_4",
						ReferencedTable: "T_4",
						ReferencedKey:   []string{"PK_41", "PK_42"},
						ReferencingKey:  []string{"PK_51", "PK_52"},
					},
				},
			},
		},
		{
			ddl:   "ddl_03_foreign_loop_1",
			table: "T_1",
			want: schema.SchemaTable{
				Name: "T_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_1_1",
						ReferencedTable: "T_1",
						ReferencedKey:   []string{"PK_12"},
						ReferencingKey:  []string{"PK_11"},
					},
				},
			},
		},
		{
			ddl:   "ddl_04_foreign_loop_2",
			table: "T_1",
			want: schema.SchemaTable{
				Name: "T_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_1_2",
						ReferencedTable: "T_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_11", "PK_12"},
					},
				},
			},
		},
		{
			ddl:   "ddl_04_foreign_loop_2",
			table: "T_2",
			want: schema.SchemaTable{
				Name: "T_2",
				Columns: []schema.SchemaColumn{
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_22", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_21", "PK_22"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_2_1",
						ReferencedTable: "T_1",
						ReferencedKey:   []string{"PK_11", "PK_12"},
						ReferencingKey:  []string{"PK_21", "PK_22"},
					},
				},
			},
		},
		{
			ddl:   "ddl_05_foreign_loop_3",
			table: "T_1",
			want: schema.SchemaTable{
				Name: "T_1",
				Columns: []schema.SchemaColumn{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_1_3",
						ReferencedTable: "T_3",
						ReferencedKey:   []string{"PK_31", "PK_32"},
						ReferencingKey:  []string{"PK_11", "PK_12"},
					},
				},
			},
		},
		{
			ddl:   "ddl_05_foreign_loop_3",
			table: "T_2",
			want: schema.SchemaTable{
				Name: "T_2",
				Columns: []schema.SchemaColumn{
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_22", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_21", "PK_22"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_2_1",
						ReferencedTable: "T_1",
						ReferencedKey:   []string{"PK_11", "PK_12"},
						ReferencingKey:  []string{"PK_21", "PK_22"},
					},
				},
			},
		},
		{
			ddl:   "ddl_05_foreign_loop_3",
			table: "T_3",
			want: schema.SchemaTable{
				Name: "T_3",
				Columns: []schema.SchemaColumn{
					{Name: "PK_31", Type: "INT64"},
					{Name: "PK_32", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_31", "PK_32"},
				ForeignKeys: []schema.SchemaForeignKey{
					{
						Name:            "FK_3_2",
						ReferencedTable: "T_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_31", "PK_32"},
					},
				},
			},
		},
	}

	for number, testcase := range testcases {
		t.Run(fmt.Sprintf("%03d:%s[%s]", number, testcase.ddl, testcase.table), func(t *testing.T) {
			database := fmt.Sprintf("fetcher_%0d", number)
			admin, client, teardown := test.Setup(t, database)
			defer teardown()
			test.InitDDLs(t, admin, client.DatabaseName(), ddls[testcase.ddl])

			ctx := context.Background()
			sut := schema.NewFetcher(client.ReadOnlyTransaction())
			got, err := sut.Fetch(ctx, testcase.table)
			assert.Nil(t, err)
			assert.Equal(t, testcase.want, got)
		})
	}
}
