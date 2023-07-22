package dml

import (
	"math/big"
	"reflect"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner"
	"cloud.google.com/go/spanner/apiv1/spannerpb"
	"github.com/Jumpaku/go-assert"
	"github.com/Jumpaku/gotaface/dml"
	"github.com/davecgh/go-spew/spew"

	proto3 "github.com/golang/protobuf/ptypes/struct"
)

func IsSupported(val any) bool {
	switch val.(type) {
	default:
		return false
	case string, *string, spanner.NullString:
		return true
	case []string, []*string, []spanner.NullString:
		return true
	case []byte:
		return true
	case [][]byte:
		return true
	case int, int8, int16, int32, int64,
		*int, *int8, *int16, *int32, *int64,
		spanner.NullInt64:
		return true
	case []int, []int8, []int16, []int32, []int64,
		[]*int, []*int8, []*int16, []*int32, []*int64,
		[]spanner.NullInt64:
		return true
	case bool, *bool, spanner.NullBool:
		return true
	case []bool, []*bool, []spanner.NullBool:
		return true
	case float32, float64, *float32, *float64, spanner.NullFloat64:
		return true
	case []float32, []float64, []*float32, []*float64, []spanner.NullFloat64:
		return true
	case big.Rat, *big.Rat, spanner.NullNumeric:
		return true
	case []big.Rat, []*big.Rat, []spanner.NullNumeric:
		return true
	case time.Time, *time.Time, spanner.NullTime:
		return true
	case []time.Time, []*time.Time, []spanner.NullTime:
		return true
	case civil.Date, *civil.Date, spanner.NullDate:
		return true
	case []civil.Date, []*civil.Date, []spanner.NullDate:
		return true
	case spanner.NullJSON:
		return true
	case []spanner.NullJSON:
		return true
	case spanner.GenericColumnValue:
		return true
	case spanner.Row, *spanner.Row, spanner.NullRow:
		return true
	case []spanner.Row, []*spanner.Row, []spanner.NullRow:
		return true
	}
}

type DBValue struct {
	Val any
}

var _ dml.DBValue = (*DBValue)(nil)

func NewDBValue(val any) *DBValue {
	assert.State(IsSupported(val), `val not supported: %v`, spew.Sdump(val))

	return &DBValue{Val: convert(val)}
}

func (v *DBValue) GoType() reflect.Type {
	assert.State(IsSupported(v.Val), `value not supported`)

	return reflect.TypeOf(v.Val)
}

func (v *DBValue) Get(ptr any) {
	dstPtr := reflect.ValueOf(ptr)
	assert.Params(dstPtr.IsValid(), `ptr must be valid`)
	assert.Params(dstPtr.Kind() == reflect.Pointer, `ptr must be a pointer`)
	assert.Params(!dstPtr.IsNil(), `ptr must be not nil`)
	assert.Params(IsSupported(dstPtr.Elem().Interface()), `ptr must references a value of a supported type`)
	assert.State(IsSupported(v.Val), `value not supported`)
	src := reflect.ValueOf(v.Val)
	assert.State(src.CanConvert(dstPtr.Type().Elem()), `ptr references an incompatible value`)

	dstPtr.Elem().Set(src.Convert(dstPtr.Type().Elem()))
}

func (v *DBValue) Set(val any) {
	assert.State(IsSupported(val), `val not supported`)
	v.Val = convert(v)
}

func mustInt64(v any) int64 {
	switch v := v.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return int64(v)
	default:
		panic(`cannot convert to int64`)
	}
}
func mustFloat64(v any) float64 {
	switch v := v.(type) {
	case float32:
		return float64(v)
	case float64:
		return float64(v)
	default:
		panic(`cannot convert to float64`)
	}
}

func decodeGeneralColumnValue[T any](src spanner.GenericColumnValue) T {
	var dst T
	if err := src.Decode(&dst); err != nil {
		assert.Unexpected(`failure to decode %v is unexpected: %w`, spew.Sdump(src), err)
	}
	return dst
}

func isNilSlice(s any) bool {
	return reflect.ValueOf(s).IsNil()
}
func nilSlice[NullVal any]() []NullVal {
	return []NullVal(nil)
}
func convertSlice[NullVal any](src any) []NullVal {
	if isNilSlice(src) {
		return nilSlice[NullVal]()
	}
	dst := []NullVal{}
	rv := reflect.ValueOf(src)
	for i := 0; i < rv.Len(); i++ {
		e := rv.Index(i).Interface()
		v, ok := convert(e).(NullVal)
		if !ok {
			var t NullVal
			assert.Unexpected(`failure to convert to %T is unexpected: %v`, t, spew.Sdump(e))
		}
		dst = append(dst, v)
	}
	return dst
}
func convert(src any) any {
	assert.Params(IsSupported(src), `src not supported`)
	switch src := src.(type) {
	default:
		return assert.Unexpected1[any](`conversion of not supported type is unexpected: %v`, spew.Sdump(src))
	case string, *string, spanner.NullString:
		dst := &spanner.NullString{}
		if err := dst.Scan(src); err != nil {
			assert.Unexpected(`failure to scan is unexpected: %v`, spew.Sdump(src))
		}
		return *dst
	case []string, []*string, []spanner.NullString:
		return convertSlice[spanner.NullString](src)
	case []byte:
		return src
	case [][]byte:
		return src
	case int, int8, int16, int32, int64,
		*int, *int8, *int16, *int32, *int64,
		spanner.NullInt64:
		switch src := src.(type) {
		case int, int8, int16, int32:
			return convert(mustInt64(src))
		case *int, *int8, *int16, *int32, *int64:
			rv := reflect.ValueOf(src)
			if rv.IsNil() {
				return spanner.NullInt64{}
			}
			return convert(rv.Elem())
		default:
			dst := &spanner.NullInt64{}
			if err := dst.Scan(src); err != nil {
				assert.Unexpected(`failure to scan is unexpected: %v`, spew.Sdump(src))
			}
			return *dst
		}
	case []int, []int8, []int16, []int32, []int64,
		[]*int, []*int8, []*int16, []*int32, []*int64,
		[]spanner.NullInt64:
		return convertSlice[spanner.NullInt64](src)
	case bool, *bool, spanner.NullBool:
		dst := &spanner.NullBool{}
		if err := dst.Scan(src); err != nil {
			assert.Unexpected(`failure to scan is unexpected: %v`, spew.Sdump(src))
		}
		return *dst
	case []bool, []*bool, []spanner.NullBool:
		return convertSlice[spanner.NullBool](src)
	case float32, float64, *float32, *float64, spanner.NullFloat64:
		switch src := src.(type) {
		case float32:
			return convert(mustFloat64(src))
		case *float32, *float64:
			rv := reflect.ValueOf(src)
			if rv.IsNil() {
				return spanner.NullFloat64{}
			}
			return convert(rv.Elem())
		default:
			dst := &spanner.NullFloat64{}
			if err := dst.Scan(src); err != nil {
				assert.Unexpected(`failure to scan is unexpected: %v`, spew.Sdump(src))
			}
			return *dst
		}
	case []float32, []float64, []*float32, []*float64, []spanner.NullFloat64:
		return convertSlice[spanner.NullFloat64](src)
	case big.Rat, *big.Rat, spanner.NullNumeric:
		dst := &spanner.NullNumeric{}
		if err := dst.Scan(src); err != nil {
			assert.Unexpected(`failure to scan is unexpected: %v`, spew.Sdump(src))
		}
		return *dst
	case []big.Rat, []*big.Rat, []spanner.NullNumeric:
		return convertSlice[spanner.NullNumeric](src)
	case time.Time, *time.Time, spanner.NullTime:
		dst := &spanner.NullTime{}
		if err := dst.Scan(src); err != nil {
			assert.Unexpected(`failure to scan is unexpected: %v`, spew.Sdump(src))
		}
		return *dst
	case []time.Time, []*time.Time, []spanner.NullTime:
		return convertSlice[spanner.NullTime](src)
	case civil.Date, *civil.Date, spanner.NullDate:
		dst := &spanner.NullDate{}
		if err := dst.Scan(src); err != nil {
			assert.Unexpected(`failure to scan is unexpected: %v`, spew.Sdump(src))
		}
		return *dst
	case []civil.Date, []*civil.Date, []spanner.NullDate:
		return convertSlice[spanner.NullDate](src)
	case spanner.NullJSON:
		return src
	case []spanner.NullJSON:
		return src
	case spanner.GenericColumnValue:
		switch src.Type.Code {
		default:
			return assert.Unexpected1[any](`unexpected type: %v`, src.Type.String())
		case spannerpb.TypeCode_STRING:
			return decodeGeneralColumnValue[spanner.NullString](src)
		case spannerpb.TypeCode_BYTES:
			return decodeGeneralColumnValue[[]byte](src)
		case spannerpb.TypeCode_INT64:
			return decodeGeneralColumnValue[spanner.NullInt64](src)
		case spannerpb.TypeCode_BOOL:
			return decodeGeneralColumnValue[spanner.NullBool](src)
		case spannerpb.TypeCode_FLOAT64:
			return decodeGeneralColumnValue[spanner.NullFloat64](src)
		case spannerpb.TypeCode_NUMERIC:
			return decodeGeneralColumnValue[spanner.NullNumeric](src)
		case spannerpb.TypeCode_TIMESTAMP:
			return decodeGeneralColumnValue[spanner.NullTime](src)
		case spannerpb.TypeCode_DATE:
			return decodeGeneralColumnValue[spanner.NullDate](src)
		case spannerpb.TypeCode_JSON:
			return decodeGeneralColumnValue[spanner.NullJSON](src)
		case spannerpb.TypeCode_ARRAY:
			switch src.Type.ArrayElementType.Code {
			default:
				return assert.Unexpected1[any](`Array of Arrays is unexpected`)
			case spannerpb.TypeCode_STRING:
				return decodeGeneralColumnValue[[]spanner.NullString](src)
			case spannerpb.TypeCode_BYTES:
				return decodeGeneralColumnValue[[][]byte](src)
			case spannerpb.TypeCode_INT64:
				return decodeGeneralColumnValue[[]spanner.NullInt64](src)
			case spannerpb.TypeCode_BOOL:
				return decodeGeneralColumnValue[[]spanner.NullBool](src)
			case spannerpb.TypeCode_FLOAT64:
				return decodeGeneralColumnValue[[]spanner.NullFloat64](src)
			case spannerpb.TypeCode_NUMERIC:
				return decodeGeneralColumnValue[[]spanner.NullNumeric](src)
			case spannerpb.TypeCode_TIMESTAMP:
				return decodeGeneralColumnValue[[]spanner.NullTime](src)
			case spannerpb.TypeCode_DATE:
				return decodeGeneralColumnValue[[]spanner.NullDate](src)
			case spannerpb.TypeCode_JSON:
				return decodeGeneralColumnValue[[]spanner.NullJSON](src)
			case spannerpb.TypeCode_STRUCT:
				return decodeGeneralColumnValue[[]spanner.NullRow](src)
			}
		case spannerpb.TypeCode_STRUCT:
			_, isNull := src.Value.Kind.(*proto3.Value_NullValue)
			if isNull {
				return spanner.NullRow{}
			}

			fields := src.Type.StructType.Fields
			fieldCount := len(fields)
			listValues := src.Value.GetListValue()
			columnNames := make([]string, fieldCount)
			columnValues := make([]any, fieldCount)
			for index := 0; index < fieldCount; index++ {
				columnNames[index] = fields[index].Name
				columnValues[index] = convert(spanner.GenericColumnValue{
					Type:  fields[index].Type,
					Value: listValues.Values[index],
				})
			}

			row, err := spanner.NewRow(columnNames, columnValues)
			if err != nil {
				assert.Unexpected(`failure to create row from %v is unexpected: %w`, spew.Sdump(src), err)
			}

			return convert(row)
		}
	case spanner.Row, *spanner.Row, spanner.NullRow:
		switch src := src.(type) {
		case spanner.Row:
			return spanner.NullRow{Valid: true, Row: src}
		case *spanner.Row:
			if src == nil {
				return spanner.NullRow{}
			}
			return spanner.NullRow{Valid: true, Row: *src}
		default:
			return src
		}
	case []spanner.Row, []*spanner.Row, []spanner.NullRow:
		return convertSlice[spanner.NullRow](src)
	}
}
