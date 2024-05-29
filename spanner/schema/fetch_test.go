package schema_test

import (
	"context"
	_ "embed"
	"fmt"
	"testing"

	"github.com/Jumpaku/gotaface/spanner/schema"
	"github.com/Jumpaku/gotaface/spanner/test"
	"github.com/stretchr/testify/assert"
)

func TestFetcher(t *testing.T) {
	test.SkipWithoutDataSource(t)

	testcases := []struct {
		ddl   string
		table string
		want  schema.Table
	}{
		{
			ddl:   "ddl_00_all_types",
			table: "A",
			want: schema.Table{
				Name: "A",
				Columns: []schema.Column{
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
				PrimaryKey:  []string{"PK"},
				ForeignKeys: []schema.ForeignKey{},
				Indexes:     []schema.Index{},
			},
		},
		{
			ddl:   "ddl_01_interleave",
			table: "B_1",
			want: schema.Table{
				Name: "B_1",
				Columns: []schema.Column{
					{Name: "PK_11", Type: "INT64"},
				},
				PrimaryKey:  []string{"PK_11"},
				ForeignKeys: []schema.ForeignKey{},
				Indexes:     []schema.Index{},
			},
		},
		{
			ddl:   "ddl_01_interleave",
			table: "B_2",
			want: schema.Table{
				Name: "B_2",
				Columns: []schema.Column{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_21", Type: "INT64"},
				},
				PrimaryKey:  []string{"PK_11", "PK_21"},
				Parent:      "B_1",
				ForeignKeys: []schema.ForeignKey{},
				Indexes:     []schema.Index{},
			},
		},
		{
			ddl:   "ddl_01_interleave",
			table: "B_3",
			want: schema.Table{
				Name: "B_3",
				Columns: []schema.Column{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_31", Type: "INT64"},
				},
				PrimaryKey:  []string{"PK_11", "PK_21", "PK_31"},
				Parent:      "B_2",
				ForeignKeys: []schema.ForeignKey{},
				Indexes:     []schema.Index{},
			},
		},
		{
			ddl:   "ddl_01_interleave",
			table: "B_4",
			want: schema.Table{
				Name: "B_4",
				Columns: []schema.Column{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_41", Type: "INT64"},
				},
				PrimaryKey:  []string{"PK_11", "PK_21", "PK_41"},
				Parent:      "B_2",
				ForeignKeys: []schema.ForeignKey{},
				Indexes:     []schema.Index{},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_1",
			want: schema.Table{
				Name: "C_1",
				Columns: []schema.Column{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey:  []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.ForeignKey{},
				Indexes:     []schema.Index{},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_2",
			want: schema.Table{
				Name: "C_2",
				Columns: []schema.Column{
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_22", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_21", "PK_22"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_C_2_1",
						ReferencedTable: "C_1",
						ReferencedKey:   []string{"PK_11", "PK_12"},
						ReferencingKey:  []string{"PK_21", "PK_22"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_3",
			want: schema.Table{
				Name: "C_3",
				Columns: []schema.Column{
					{Name: "PK_31", Type: "INT64"},
					{Name: "PK_32", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_31", "PK_32"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_C_3_2",
						ReferencedTable: "C_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_31", "PK_32"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_4",
			want: schema.Table{
				Name: "C_4",
				Columns: []schema.Column{
					{Name: "PK_41", Type: "INT64"},
					{Name: "PK_42", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_41", "PK_42"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_C_4_2",
						ReferencedTable: "C_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_41", "PK_42"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_02_foreign_keys",
			table: "C_5",
			want: schema.Table{
				Name: "C_5",
				Columns: []schema.Column{
					{Name: "PK_51", Type: "INT64"},
					{Name: "PK_52", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_51", "PK_52"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_C_5_3",
						ReferencedTable: "C_3",
						ReferencedKey:   []string{"PK_31", "PK_32"},
						ReferencingKey:  []string{"PK_51", "PK_52"},
					},
					{
						Name:            "FK_C_5_4",
						ReferencedTable: "C_4",
						ReferencedKey:   []string{"PK_41", "PK_42"},
						ReferencingKey:  []string{"PK_51", "PK_52"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_03_foreign_loop_1",
			table: "D_1",
			want: schema.Table{
				Name: "D_1",
				Columns: []schema.Column{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_D_1_1",
						ReferencedTable: "D_1",
						ReferencedKey:   []string{"PK_12"},
						ReferencingKey:  []string{"PK_11"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_04_foreign_loop_2",
			table: "E_1",
			want: schema.Table{
				Name: "E_1",
				Columns: []schema.Column{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_E_1_2",
						ReferencedTable: "E_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_11", "PK_12"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_04_foreign_loop_2",
			table: "E_2",
			want: schema.Table{
				Name: "E_2",
				Columns: []schema.Column{
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_22", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_21", "PK_22"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_E_2_1",
						ReferencedTable: "E_1",
						ReferencedKey:   []string{"PK_11", "PK_12"},
						ReferencingKey:  []string{"PK_21", "PK_22"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_05_foreign_loop_3",
			table: "F_1",
			want: schema.Table{
				Name: "F_1",
				Columns: []schema.Column{
					{Name: "PK_11", Type: "INT64"},
					{Name: "PK_12", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_11", "PK_12"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_F_1_3",
						ReferencedTable: "F_3",
						ReferencedKey:   []string{"PK_31", "PK_32"},
						ReferencingKey:  []string{"PK_11", "PK_12"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_05_foreign_loop_3",
			table: "F_2",
			want: schema.Table{
				Name: "F_2",
				Columns: []schema.Column{
					{Name: "PK_21", Type: "INT64"},
					{Name: "PK_22", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_21", "PK_22"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_F_2_1",
						ReferencedTable: "F_1",
						ReferencedKey:   []string{"PK_11", "PK_12"},
						ReferencingKey:  []string{"PK_21", "PK_22"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_05_foreign_loop_3",
			table: "F_3",
			want: schema.Table{
				Name: "F_3",
				Columns: []schema.Column{
					{Name: "PK_31", Type: "INT64"},
					{Name: "PK_32", Type: "INT64"},
				},
				PrimaryKey: []string{"PK_31", "PK_32"},
				ForeignKeys: []schema.ForeignKey{
					{
						Name:            "FK_F_3_2",
						ReferencedTable: "F_2",
						ReferencedKey:   []string{"PK_21", "PK_22"},
						ReferencingKey:  []string{"PK_31", "PK_32"},
					},
				},
				Indexes: []schema.Index{},
			},
		},
		{
			ddl:   "ddl_06_unique_keys",
			table: "G",
			want: schema.Table{
				Name: "G",
				Columns: []schema.Column{
					{Name: "PK", Type: "INT64"},
					{Name: "C1", Type: "INT64"},
					{Name: "C2", Type: "INT64"},
					{Name: "C3", Type: "INT64"},
				},
				PrimaryKey: []string{"PK"},
				Indexes: []schema.Index{
					{Name: "UQ_G_C1", Unique: true, Key: []schema.IndexKey{{Name: "C1"}}},
					{Name: "UQ_G_C1_C2", Unique: true, Key: []schema.IndexKey{{Name: "C1"}, {Name: "C2"}}},
					{Name: "UQ_G_C1_C2_C3", Unique: true, Key: []schema.IndexKey{{Name: "C1"}, {Name: "C2"}, {Name: "C3"}}},
					{Name: "UQ_G_C1_C3", Unique: true, Key: []schema.IndexKey{{Name: "C1"}, {Name: "C3"}}},
					{Name: "UQ_G_C1_C3_C2", Unique: true, Key: []schema.IndexKey{{Name: "C1"}, {Name: "C3"}, {Name: "C2"}}},
					{Name: "UQ_G_C2", Unique: true, Key: []schema.IndexKey{{Name: "C2"}}},
					{Name: "UQ_G_C2_C1", Unique: true, Key: []schema.IndexKey{{Name: "C2"}, {Name: "C1"}}},
					{Name: "UQ_G_C2_C1_C3", Unique: true, Key: []schema.IndexKey{{Name: "C2"}, {Name: "C1"}, {Name: "C3"}}},
					{Name: "UQ_G_C2_C3", Unique: true, Key: []schema.IndexKey{{Name: "C2"}, {Name: "C3"}}},
					{Name: "UQ_G_C2_C3_C1", Unique: true, Key: []schema.IndexKey{{Name: "C2"}, {Name: "C3"}, {Name: "C1"}}},
					{Name: "UQ_G_C3", Unique: true, Key: []schema.IndexKey{{Name: "C3"}}},
					{Name: "UQ_G_C3_C1", Unique: true, Key: []schema.IndexKey{{Name: "C3"}, {Name: "C1"}}},
					{Name: "UQ_G_C3_C1_C2", Unique: true, Key: []schema.IndexKey{{Name: "C3"}, {Name: "C1"}, {Name: "C2"}}},
					{Name: "UQ_G_C3_C2", Unique: true, Key: []schema.IndexKey{{Name: "C3"}, {Name: "C2"}}},
					{Name: "UQ_G_C3_C2_C1", Unique: true, Key: []schema.IndexKey{{Name: "C3"}, {Name: "C2"}, {Name: "C1"}}},
				},
				ForeignKeys: []schema.ForeignKey{},
			},
		},
		{
			ddl:   "ddl_07_generated_col",
			table: "H",
			want: schema.Table{
				Name: "H",
				Columns: []schema.Column{
					{Name: "PK", Type: "INT64"},
					{Name: "Gen", Type: "INT64", Nullable: true, Generated: true},
				},
				PrimaryKey:  []string{"PK"},
				ForeignKeys: []schema.ForeignKey{},
				Indexes:     []schema.Index{},
			},
		},
		{
			ddl:   "ddl_08_index",
			table: "I",
			want: schema.Table{
				Name: "I",
				Columns: []schema.Column{
					{Name: "PK", Type: "INT64"},
					{Name: "C1", Type: "INT64"},
					{Name: "C2", Type: "INT64"},
				},
				PrimaryKey: []string{"PK"},
				Indexes: []schema.Index{
					{Name: "IDX_I_C1Asc_C2Desc", Key: []schema.IndexKey{{Name: "C1"}, {Name: "C2", Desc: true}}},
					{Name: "IDX_I_C1Desc_C2Asc", Key: []schema.IndexKey{{Name: "C1", Desc: true}, {Name: "C2"}}},
					{Name: "IDX_I_Storing", Key: []schema.IndexKey{{Name: "C1"}}},
				},
				ForeignKeys: []schema.ForeignKey{},
			},
		},
		{
			ddl:   "ddl_09_view",
			table: "K",
			want: schema.Table{
				Name: "K",
				View: true,
				Columns: []schema.Column{
					{Name: "PK_2", Type: "INT64", Nullable: true},
					{Name: "Col1_2", Type: "INT64", Nullable: true},
					{Name: "Col2_2", Type: "INT64", Nullable: true},
				},
			},
		},
	}

	for number, testcase := range testcases {
		t.Run(fmt.Sprintf("%03d:%s[%s]", number, testcase.ddl, testcase.table), func(t *testing.T) {
			admin, client := test.Setup(t)

			test.InitDDLs(t, admin, client.DatabaseName(), TestDDLs[testcase.ddl])

			ctx := context.Background()
			sut := schema.NewFetcher(client.ReadOnlyTransaction())
			got, err := sut.Fetch(ctx, testcase.table)
			assert.Nil(t, err)
			assert.Equal(t, testcase.want, got)
		})
	}
}
