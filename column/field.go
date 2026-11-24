package column

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"
)

var (
	valuerType = reflect.TypeFor[driver.Valuer]()
	timeType   = reflect.TypeFor[time.Time]()

	stringType  = reflect.TypeFor[string]()
	pStringType = reflect.TypeFor[*string]()
	int64Type   = reflect.TypeFor[int64]()
	pInt64Type  = reflect.TypeFor[*int64]()
	int32Type   = reflect.TypeFor[int32]()
	pInt32Type  = reflect.TypeFor[*int32]()
	int16Type   = reflect.TypeFor[int16]()
	pInt16Type  = reflect.TypeFor[*int16]()
	float64Type = reflect.TypeFor[float64]()
	pFloat64    = reflect.TypeFor[*float64]()
	boolType    = reflect.TypeFor[bool]()
	pBoolType   = reflect.TypeFor[*bool]()
	pTimeType   = reflect.TypeFor[*time.Time]()
)

// Field returns a getter that reads a row's field as kind k, for
// generated column sets (A5). It calls field once, on R's zero value, to
// learn the field's type, and panics there if k can't read that type.
func Field[R any](k Kind, field func(R) any) func(R) Value {
	var zero R
	ft := reflect.TypeOf(field(zero))
	if ft == nil {
		panic(fmt.Sprintf("column: %s can't read a field of interface type", k))
	}
	if k == Enum {
		panic("column: Field doesn't read enums; generated enum getters read the field themselves")
	}

	if fast, ok := resolveFastPath(k, ft, field); ok {
		return fast
	}
	// A Valuer by value, such as every pgtype type: no reflection per row.
	if ft.Kind() != reflect.Pointer && ft.Implements(valuerType) {
		return func(r R) Value {
			res, err := field(r).(driver.Valuer).Value()
			return finishValuer(k, ft, res, err)
		}
	}

	reader := resolveReflectReader(k, ft, ft)
	return func(r R) Value {
		return reader(reflect.ValueOf(field(r)))
	}
}

func resolveFastPath[R any](k Kind, ft reflect.Type, field func(R) any) (func(R) Value, bool) {
	switch ft {
	case stringType, pStringType:
		var conv func(string) Value
		switch k {
		case Text:
			conv = func(s string) Value { return Value{Text: s} }
		case Number:
			conv = Numeric
		default:
			panicUnsupported(k, ft)
		}
		return bindFast(ft == pStringType, field, conv), true

	case int64Type, pInt64Type:
		return bindIntFast(k, ft, ft == pInt64Type, field, func(v int64) int64 { return v }), true

	case int32Type, pInt32Type:
		return bindIntFast(k, ft, ft == pInt32Type, field, func(v int32) int64 { return int64(v) }), true

	case int16Type, pInt16Type:
		return bindIntFast(k, ft, ft == pInt16Type, field, func(v int16) int64 { return int64(v) }), true

	case float64Type, pFloat64:
		if k != Number {
			panicUnsupported(k, ft)
		}
		return bindFast(ft == pFloat64, field, func(v float64) Value { return Value{Num: v} }), true

	case boolType, pBoolType:
		if k != Bool {
			panicUnsupported(k, ft)
		}
		return bindFast(ft == pBoolType, field, func(v bool) Value { return Value{Bool: v} }), true

	case timeType, pTimeType:
		if k != Date && k != Time {
			panicUnsupported(k, ft)
		}
		return bindFast(ft == pTimeType, field, func(t time.Time) Value { return Value{Time: t} }), true
	}
	return nil, false
}

func bindIntFast[R any, T any](k Kind, ft reflect.Type, isPtr bool, field func(R) any, toInt64 func(T) int64) func(R) Value {
	var conv func(T) Value
	switch k {
	case Integer:
		conv = func(v T) Value { return Value{Int: toInt64(v)} }
	case Number:
		conv = func(v T) Value { return Value{Num: float64(toInt64(v))} }
	default:
		panicUnsupported(k, ft)
	}
	return bindFast(isPtr, field, conv)
}

func bindFast[R any, T any](isPtr bool, field func(R) any, conv func(T) Value) func(R) Value {
	if isPtr {
		return func(r R) Value {
			p := field(r).(*T)
			if p == nil {
				return Value{Null: true}
			}
			return conv(*p)
		}
	}
	return func(r R) Value {
		return conv(field(r).(T))
	}
}

func resolveReflectReader(k Kind, ft, fieldType reflect.Type) func(reflect.Value) Value {
	if ft.Kind() == reflect.Interface {
		panicUnsupported(k, fieldType)
	}
	if ft.Kind() == reflect.Ptr {
		elemReader := resolveReflectReader(k, ft.Elem(), fieldType)
		return func(rv reflect.Value) Value {
			if rv.IsNil() {
				return Value{Null: true}
			}
			return elemReader(rv.Elem())
		}
	}

	if ft.Implements(valuerType) {
		return func(rv reflect.Value) Value {
			res, err := rv.Interface().(driver.Valuer).Value()
			return finishValuer(k, fieldType, res, err)
		}
	}
	if reflect.PointerTo(ft).Implements(valuerType) {
		return func(rv reflect.Value) Value {
			var res driver.Value
			var err error
			if rv.CanAddr() {
				res, err = rv.Addr().Interface().(driver.Valuer).Value()
			} else {
				ptr := reflect.New(ft)
				ptr.Elem().Set(rv)
				res, err = ptr.Interface().(driver.Valuer).Value()
			}
			return finishValuer(k, fieldType, res, err)
		}
	}

	if ft == timeType {
		if k != Date && k != Time {
			panicUnsupported(k, fieldType)
		}
		return func(rv reflect.Value) Value {
			return Value{Time: rv.Interface().(time.Time)}
		}
	}

	switch ft.Kind() {
	case reflect.String:
		switch k {
		case Text:
			return func(rv reflect.Value) Value { return Value{Text: rv.String()} }
		case Number:
			return func(rv reflect.Value) Value { return Numeric(rv.String()) }
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch k {
		case Integer:
			return func(rv reflect.Value) Value { return Value{Int: rv.Int()} }
		case Number:
			return func(rv reflect.Value) Value { return Value{Num: float64(rv.Int())} }
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch k {
		case Integer:
			return func(rv reflect.Value) Value { return Value{Int: int64(rv.Uint())} }
		case Number:
			return func(rv reflect.Value) Value { return Value{Num: float64(rv.Uint())} }
		}

	case reflect.Float32, reflect.Float64:
		if k == Number {
			return func(rv reflect.Value) Value { return Value{Num: rv.Float()} }
		}

	case reflect.Bool:
		if k == Bool {
			return func(rv reflect.Value) Value { return Value{Bool: rv.Bool()} }
		}

	case reflect.Slice:
		if ft.Elem().Kind() == reflect.Uint8 && k == Text {
			return func(rv reflect.Value) Value { return Value{Text: string(rv.Bytes())} }
		}

	case reflect.Array:
		if ft.Len() == 16 && ft.Elem().Kind() == reflect.Uint8 && k == Text {
			return func(rv reflect.Value) Value {
				var b [16]byte
				reflect.Copy(reflect.ValueOf(&b).Elem(), rv)
				return Value{Text: UUID(b)}
			}
		}
	}

	panicUnsupported(k, fieldType)
	return nil
}

func finishValuer(k Kind, fieldType reflect.Type, res driver.Value, err error) Value {
	if err != nil {
		panic(fmt.Sprintf("column: %s.Value(): %v", fieldType, err))
	}
	if res == nil {
		return Value{Null: true}
	}
	return readValuerResult(k, fieldType, res)
}

func readValuerResult(k Kind, fieldType reflect.Type, res driver.Value) Value {
	switch k {
	case Text:
		switch v := res.(type) {
		case string:
			return Value{Text: v}
		case []byte:
			return Value{Text: string(v)}
		}

	case Integer:
		if n, ok := asInt64(res); ok {
			return Value{Int: n}
		}

	case Number:
		switch v := res.(type) {
		case float64:
			return Value{Num: v}
		case float32:
			return Value{Num: float64(v)}
		case string:
			return Numeric(v)
		default:
			if n, ok := asInt64(res); ok {
				return Value{Num: float64(n)}
			}
		}

	case Bool:
		if b, ok := res.(bool); ok {
			return Value{Bool: b}
		}

	case Date, Time:
		switch v := res.(type) {
		case time.Time:
			return Value{Time: v}
		case string:
			if v == "infinity" {
				return Infinite(1)
			}
			if v == "-infinity" {
				return Infinite(-1)
			}
		}
	}

	if n, ok := normalize(k, res); ok {
		return readValuerResult(k, fieldType, n)
	}
	panic(fmt.Sprintf("column: %s can't take a %T from %s", k, res, fieldType))
}

// normalize turns a Valuer's result of a named basic type, or a [16]byte
// for Text, into the unnamed type readValuerResult takes; false when it
// is neither or already unnamed.
func normalize(k Kind, res driver.Value) (driver.Value, bool) {
	rv := reflect.ValueOf(res)
	var n any
	switch rv.Kind() {
	case reflect.String:
		n = rv.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n = rv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n = rv.Uint()
	case reflect.Float32, reflect.Float64:
		n = rv.Float()
	case reflect.Bool:
		n = rv.Bool()
	case reflect.Slice:
		if rv.Type().Elem().Kind() != reflect.Uint8 {
			return nil, false
		}
		n = rv.Bytes()
	case reflect.Array:
		if k != Text || rv.Len() != 16 || rv.Type().Elem().Kind() != reflect.Uint8 {
			return nil, false
		}
		var b [16]byte
		reflect.Copy(reflect.ValueOf(&b).Elem(), rv)
		n = UUID(b)
	default:
		return nil, false
	}
	if reflect.TypeOf(n) == rv.Type() {
		return nil, false // already unnamed: the kind can't take it
	}
	return n, true
}

func asInt64(v any) (int64, bool) {
	switch n := v.(type) {
	case int:
		return int64(n), true
	case int8:
		return int64(n), true
	case int16:
		return int64(n), true
	case int32:
		return int64(n), true
	case int64:
		return n, true
	case uint:
		return int64(n), true
	case uint8:
		return int64(n), true
	case uint16:
		return int64(n), true
	case uint32:
		return int64(n), true
	case uint64:
		return int64(n), true
	}
	return 0, false
}

func panicUnsupported(k Kind, fieldType reflect.Type) {
	panic(fmt.Sprintf("column: %s can't read %s", k, fieldType))
}
