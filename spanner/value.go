package spanner

import (
	"reflect"
	"strings"
	"unicode"

	"cloud.google.com/go/spanner"
)

func RefType[T any]() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

func IsSpannerInt64(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "int64")
}

func IsSpannerString(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "string")
}

func IsSpannerBool(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "bool")
}

func IsSpannerFloat64(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "float64")
}

func IsSpannerTimestamp(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "timestamp")
}

func IsSpannerDate(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "date")
}

func IsSpannerNumeric(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "numeric")
}

func IsSpannerBytes(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "bytes")
}

func IsSpannerJSON(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "json")
}

func IsSpannerArray(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "array<")
}

func SpannerArrayElemType(columnType string) string {
	lower := strings.ToLower(columnType)
	return lower[6 : len(lower)-1]
}
func SpannerStructElemTypes(columnType string) ([]string, []string) {
	depth := 0
	currentField := [][]rune{{}}
	fieldNames := []string{}
	fieldTypes := []string{}
	for _, c := range columnType[7 : len(columnType)-1] {
		switch {
		case c == '<':
			depth++
			currentField[len(currentField)-1] = append(currentField[len(currentField)-1], c)
		case c == '>':
			depth--
			currentField[len(currentField)-1] = append(currentField[len(currentField)-1], c)
		case unicode.IsSpace(c) && depth == 0:
			if len(currentField[len(currentField)-1]) != 0 {
				currentField = append(currentField, []rune{})
			}
		case c == ',' && depth == 0:
			if len(currentField) == 1 {
				fieldNames = append(fieldNames, "")
				fieldTypes = append(fieldTypes, string(currentField[0]))
			} else {
				fieldNames = append(fieldNames, string(currentField[0]))
				fieldTypes = append(fieldTypes, string(currentField[1]))
			}
			currentField = [][]rune{{}}
		default:
			currentField[len(currentField)-1] = append(currentField[len(currentField)-1], c)
		}
	}

	return fieldNames, fieldTypes
}

func IsSpannerStruct(columnType string) bool {
	return strings.HasPrefix(strings.ToLower(columnType), "struct")
}

func GoType(columnType string) reflect.Type {
	switch {
	case IsSpannerInt64(columnType):
		return RefType[spanner.NullInt64]()
	case IsSpannerString(columnType):
		return RefType[spanner.NullString]()
	case IsSpannerBool(columnType):
		return RefType[spanner.NullBool]()
	case IsSpannerFloat64(columnType):
		return RefType[spanner.NullFloat64]()
	case IsSpannerTimestamp(columnType):
		return RefType[spanner.NullTime]()
	case IsSpannerDate(columnType):
		return RefType[spanner.NullDate]()
	case IsSpannerNumeric(columnType):
		return RefType[spanner.NullNumeric]()
	case IsSpannerBytes(columnType):
		return RefType[[]byte]()
	case IsSpannerJSON(columnType):
		return RefType[spanner.NullJSON]()
	case IsSpannerArray(columnType):
		return reflect.SliceOf(GoType(SpannerArrayElemType(columnType)))
	case IsSpannerStruct(columnType):
		return RefType[spanner.NullRow]()
	default:
		return RefType[spanner.GenericColumnValue]()
	}
}
