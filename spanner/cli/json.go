package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner"
	"github.com/Jumpaku/go-assert"
	gotaface_spanner "github.com/Jumpaku/gotaface/spanner"
	"github.com/Jumpaku/gotaface/spanner/dml"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func ToJSON(v *dml.DBValue) ([]byte, error) {
	assert.State(dml.IsSupported(v.Val), `val not supported`)

	return marshalJSON(v.Val)
}

func FromJSON(columnType string, b []byte) (*dml.DBValue, error) {
	v, err := unmarshalJSONColumnValue(columnType, b)
	if err != nil {
		return nil, fmt.Errorf(`fail to obtain DBValue from JSON: %w`, err)
	}

	return dml.NewDBValue(v), nil
}

func nilSliceOf(elemType string) any {
	switch {
	case gotaface_spanner.IsSpannerInt64(elemType):
		return ([]spanner.NullInt64)(nil)
	case gotaface_spanner.IsSpannerString(elemType):
		return ([]spanner.NullString)(nil)
	case gotaface_spanner.IsSpannerBool(elemType):
		return ([]spanner.NullBool)(nil)
	case gotaface_spanner.IsSpannerFloat64(elemType):
		return ([]spanner.NullFloat64)(nil)
	case gotaface_spanner.IsSpannerTimestamp(elemType):
		return ([]spanner.NullTime)(nil)
	case gotaface_spanner.IsSpannerDate(elemType):
		return ([]spanner.NullDate)(nil)
	case gotaface_spanner.IsSpannerNumeric(elemType):
		return ([]spanner.NullNumeric)(nil)
	case gotaface_spanner.IsSpannerBytes(elemType):
		return ([][]byte)(nil)
	case gotaface_spanner.IsSpannerJSON(elemType):
		return ([]spanner.NullJSON)(nil)
	case gotaface_spanner.IsSpannerStruct(elemType):
		return ([]spanner.NullRow)(nil)
	case gotaface_spanner.IsSpannerArray(elemType):
		t := reflect.SliceOf(gotaface_spanner.GoType(gotaface_spanner.SpannerArrayElemType(elemType)))
		return reflect.Zero(t)
	default:
		panic(`unsupported type`)
	}
}

var nullJSON = []byte("null")
var emptyArrayJSON = []byte("[]")

func marshalJSONSlice(v any) ([]byte, error) {
	assert.Params(dml.IsSupported(v), `v must be a supported type: %v`, spew.Sdump(v))

	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return nullJSON, nil
	}
	if rv.Len() == 0 {
		return emptyArrayJSON, nil
	}
	arr := []json.RawMessage{}
	for i := 0; i < rv.Len(); i++ {
		e := rv.Index(i).Interface()
		m, err := marshalJSON(dml.NewDBValue(e))
		if err != nil {
			return nil, fmt.Errorf(`fail to marshal %v to JSON: %w`, spew.Sdump(e), err)
		}

		arr = append(arr, m)
	}

	return json.Marshal(arr)
}

func marshalJSON(v any) ([]byte, error) {
	assert.Params(dml.IsSupported(v), `v must be a supported type: %v`, spew.Sdump(v))

	protoJSONMarshal := protojson.MarshalOptions{EmitUnpopulated: true, AllowPartial: true}.Marshal
	switch v := v.(type) {
	default:
		return nil, fmt.Errorf(`fail to marshal %v`, spew.Sdump(v))
	case spanner.NullString:
		return json.Marshal(v)
	case []spanner.NullString:
		return marshalJSONSlice(v)
	case []byte:
		if v == nil {
			return nullJSON, nil
		}
		if len(v) == 0 {
			return json.Marshal(v)
		}

		return protoJSONMarshal(wrapperspb.Bytes(v))
	case [][]byte:
		return marshalJSONSlice(v)
	case spanner.NullInt64:
		return json.Marshal(v)
	case []spanner.NullInt64:
		return marshalJSONSlice(v)
	case spanner.NullBool:
		return json.Marshal(v)
	case []spanner.NullBool:
		return marshalJSONSlice(v)
	case spanner.NullFloat64:
		return json.Marshal(v)
	case []spanner.NullFloat64:
		return marshalJSONSlice(v)
	case spanner.NullNumeric:
		return marshalJSONSlice(v)
	case []spanner.NullNumeric:
		return marshalJSONSlice(v)
	case spanner.NullTime:
		return json.Marshal(v)
	case []spanner.NullTime:
		return marshalJSONSlice(v)
	case spanner.NullDate:
		return json.Marshal(v)
	case []spanner.NullDate:
		return marshalJSONSlice(v)
	case spanner.NullJSON:
		return json.Marshal(v)
	case []spanner.NullJSON:
		return marshalJSONSlice(v)
	case spanner.GenericColumnValue:
		return protoJSONMarshal(v.Value)
	case spanner.NullRow:
		if !v.Valid {
			return nullJSON, nil
		}
		obj := map[string]json.RawMessage{}
		for i := 0; i < v.Row.Size(); i++ {
			name := v.Row.ColumnName(i)
			var u spanner.GenericColumnValue
			if err := v.Row.Column(i, &u); err != nil {
				return nil, fmt.Errorf(`fail to read column value at %v: %w`, name, err)
			}

			m, err := marshalJSON(u)
			if err != nil {
				return nil, fmt.Errorf(`fail to marshal spanner.GenericColumnValue %v at %v to JSON: : %w`, spew.Sdump(u), name, err)
			}

			obj[name] = m
		}
		return json.Marshal(obj)
	case []spanner.NullRow:
		return marshalJSONSlice(v)
	}
}

func unmarshalJSON[NullVal any](b []byte) (NullVal, error) {
	var v NullVal
	if err := json.Unmarshal(b, &v); err != nil {
		return v, fmt.Errorf(`fail to unmarshal %T from JSON: %v`, v, err)
	}
	return v, nil
}

func unmarshalJSONColumnValue(columnType string, b []byte) (any, error) {
	lower := strings.ToLower(columnType)
	dec := json.NewDecoder(bytes.NewBuffer(bytes.Clone(b)))
	dec.UseNumber()
	var a any
	if err := dec.Decode(&a); err != nil {
		return nil, fmt.Errorf(`fail to unmarshal from JSON as %v value: %w`, columnType, err)
	}

	protoJSONUnmarshal := protojson.UnmarshalOptions{AllowPartial: true}.Unmarshal

	switch {
	case gotaface_spanner.IsSpannerInt64(columnType):
		switch a := a.(type) {
		case nil:
			return spanner.NullInt64{}, nil
		case string:
			val, err := strconv.ParseInt(a, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return val, nil
		case bool:
			if a {
				return 1, nil
			}
			return 0, nil
		case json.Number:
			val, err := a.Int64()
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return val, nil
		default:
			v, err := unmarshalJSON[spanner.NullInt64](b)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		}

	case gotaface_spanner.IsSpannerString(columnType):
		switch a := a.(type) {
		case nil:
			return spanner.NullString{}, nil
		case string:
			return a, nil
		case bool:
			return strconv.FormatBool(a), nil
		case json.Number:
			return a.String(), nil
		default:
			v, err := unmarshalJSON[spanner.NullString](b)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		}
	case gotaface_spanner.IsSpannerBool(columnType):
		switch a := a.(type) {
		case nil:
			return spanner.NullBool{}, nil
		case string:
			v, err := strconv.ParseBool(a)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		case bool:
			return a, nil
		case json.Number:
			if n, err := a.Int64(); err == nil {
				return n == 0, nil
			}
			if f, err := a.Float64(); err == nil {
				return f == 0, nil
			}
			return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value`, a, columnType)
		default:
			v, err := unmarshalJSON[spanner.NullBool](b)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		}
	case gotaface_spanner.IsSpannerFloat64(columnType):
		switch a := a.(type) {
		case nil:
			return spanner.NullFloat64{}, nil
		case string:
			v, err := strconv.ParseFloat(a, 64)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		case bool:
			if a {
				return 1.0, nil
			}
			return 0.0, nil
		case json.Number:
			v, err := a.Float64()
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		default:
			v, err := unmarshalJSON[spanner.NullFloat64](b)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		}
	case gotaface_spanner.IsSpannerTimestamp(columnType):
		switch a := a.(type) {
		case nil:
			return spanner.NullTime{}, nil
		case string:
			t := &timestamppb.Timestamp{}
			if err := protoJSONUnmarshal(b, t); err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return t.AsTime(), nil
		case json.Number:
			if val, err := a.Int64(); err == nil {
				return time.Unix(val, 0), nil
			}

			if val, err := a.Float64(); err == nil {
				sec := math.Floor(val)
				nsec := math.Floor((val - sec) * 1e9)
				return time.Unix(int64(sec), int64(nsec)), nil
			}

			return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value`, a, columnType)
		default:
			v, err := unmarshalJSON[spanner.NullFloat64](b)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		}
	case gotaface_spanner.IsSpannerDate(columnType):
		switch a := a.(type) {
		case nil:
			return spanner.NullTime{}, nil
		case string:
			v, err := civil.ParseDate(a)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		default:
			v, err := unmarshalJSON[spanner.NullFloat64](b)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		}
	case gotaface_spanner.IsSpannerNumeric(columnType):
		switch a := a.(type) {
		case nil:
			return dml.NewDBValue(spanner.NullNumeric{}), nil
		case string:
			if r, ok := (&big.Rat{}).SetString(a); ok {
				return r, nil
			}

			if i, err := strconv.ParseInt(a, 10, 64); err == nil {
				return (&big.Rat{}).SetInt64(i), nil
			}

			if f, err := strconv.ParseFloat(a, 64); err == nil {
				return (&big.Rat{}).SetFloat64(f), nil
			}

			return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value`, a, columnType)
		case bool:
			if a {
				return big.NewRat(1, 1), nil
			}
			return big.NewRat(0, 1), nil
		case json.Number:
			if i, err := a.Int64(); err == nil {
				return (&big.Rat{}).SetInt64(i), nil
			}

			if f, err := a.Float64(); err == nil {
				return (&big.Rat{}).SetFloat64(f), nil
			}

			return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value`, a, columnType)
		default:
			v, err := unmarshalJSON[spanner.NullNumeric](b)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		}
	case gotaface_spanner.IsSpannerBytes(columnType):
		switch a := a.(type) {
		default:
			return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value`, a, columnType)
		case nil:
			return ([]byte)(nil), nil
		case string:
			v := &wrapperspb.BytesValue{}
			if err := protoJSONUnmarshal(b, v); err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}

			return v.Value, nil
		case []any:
			v := []byte{}
			for _, e := range a {
				switch e := e.(type) {
				default:
					return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value`, a, columnType)
				case string:
					n, err := strconv.ParseUint(e, 10, 8)
					if err != nil {
						return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
					}
					v = append(v, byte(n))
				case json.Number:
					n, err := e.Int64()
					if err != nil {
						return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
					}
					v = append(v, byte(n))
				}
			}
			return v, nil
		}
	case gotaface_spanner.IsSpannerJSON(columnType):
		v, err := unmarshalJSON[spanner.NullFloat64](b)
		if err != nil {
			return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
		}
		return v, nil
	case gotaface_spanner.IsSpannerArray(columnType):
		elemType := gotaface_spanner.SpannerArrayElemType(columnType)
		switch a := a.(type) {
		default:
			return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value`, a, columnType)
		case nil:
			return nilSliceOf(elemType), nil
		case []any:
			m := []json.RawMessage{}
			if err := json.Unmarshal(b, &m); err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}

			v := make([]any, len(m))
			for i, e := range m {
				var err error
				if v[i], err = unmarshalJSONColumnValue(elemType, e); err != nil {
					return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
				}
			}

			return v, nil
		}
	case gotaface_spanner.IsSpannerStruct(columnType):
		switch a := a.(type) {
		case nil:
			return spanner.NullRow{}, nil
		case map[string]any:
			columnNames := []string{}
			columnValues := []any{}

			m := map[string]json.RawMessage{}
			if err := json.Unmarshal(b, &m); err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}

			columnNames, columnTypes := gotaface_spanner.SpannerStructElemTypes(columnType)
			for key, val := range m {
				columnNames = append(columnNames, key)
				if err := unmarshalJSONColumnValue(); err != nil {
					return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
				}

				columnValues = append(columnValues)
			}

			if val, err := a.Int64(); err == nil {
				return time.Unix(val, 0), nil
			}

			if val, err := a.Float64(); err == nil {
				sec := math.Floor(val)
				nsec := math.Floor((val - sec) * 1e9)
				return time.Unix(int64(sec), int64(nsec)), nil
			}

			return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value`, a, columnType)
		default:
			v, err := unmarshalJSON[spanner.NullFloat64](b)
			if err != nil {
				return nil, fmt.Errorf(`fail to unmarshal %v from JSON as %v value: %w`, a, columnType, err)
			}
			return v, nil
		}
	default:
		return RefType[spanner.GenericColumnValue]()
	}
	assert.State(IsSupported(v.Val), `val not supported`)

}
