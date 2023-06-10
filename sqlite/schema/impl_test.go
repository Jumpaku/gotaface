package schema_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Jumpaku/gotaface/sqlite/schema"
)

func refType[T any]() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}
func TestColumn_GoType(t *testing.T) {
	type testCase struct {
		sut  schema.Column
		want reflect.Type
	}
	testCases := []testCase{
		{
			sut:  schema.Column{TypeVal: "INT"},
			want: refType[sql.NullInt64](),
		}, {
			sut:  schema.Column{TypeVal: "INTEGER"},
			want: refType[sql.NullInt64](),
		}, {
			sut:  schema.Column{TypeVal: "TINYINT"},
			want: refType[sql.NullInt64](),
		}, {
			sut:  schema.Column{TypeVal: "SMALLINT"},
			want: refType[sql.NullInt64](),
		}, {
			sut:  schema.Column{TypeVal: "MEDIUMINT"},
			want: refType[sql.NullInt64](),
		}, {
			sut:  schema.Column{TypeVal: "BIGINT"},
			want: refType[sql.NullInt64](),
		}, {
			sut:  schema.Column{TypeVal: "UNSIGNED BIG INT"},
			want: refType[sql.NullInt64](),
		}, {
			sut:  schema.Column{TypeVal: "INT2"},
			want: refType[sql.NullInt64](),
		}, {
			sut:  schema.Column{TypeVal: "INT8"},
			want: refType[sql.NullInt64](),
		}, {
			sut:  schema.Column{TypeVal: "CHARACTER(20)"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "VARCHAR(255)"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "VARYING CHARACTER(255)"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "NCHAR(55)"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "NATIVE CHARACTER(70)"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "NVARCHAR(100)"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "TEXT"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "CLOB"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "BLOB"},
			want: refType[[]byte](),
		}, {
			sut:  schema.Column{TypeVal: ""},
			want: refType[[]byte](),
		}, {
			sut:  schema.Column{TypeVal: "REAL"},
			want: refType[sql.NullFloat64](),
		}, {
			sut:  schema.Column{TypeVal: "DOUBLE"},
			want: refType[sql.NullFloat64](),
		}, {
			sut:  schema.Column{TypeVal: "DOUBLE PRECISION"},
			want: refType[sql.NullFloat64](),
		}, {
			sut:  schema.Column{TypeVal: "FLOAT"},
			want: refType[sql.NullFloat64](),
		}, {
			sut:  schema.Column{TypeVal: "NUMERIC"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "DECIMAL(10,5)"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "BOOLEAN"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "DATE"},
			want: refType[sql.NullString](),
		}, {
			sut:  schema.Column{TypeVal: "DATETIME"},
			want: refType[sql.NullString](),
		},
	}

	for _, testCase := range testCases {
		got := testCase.sut.GoType()
		if got.Name() != testCase.want.Name() {
			t.Errorf("GoType() mismatch\n  got  = %v\n  want = %v", got, testCase.want)
		}
	}
}
