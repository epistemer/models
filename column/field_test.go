package column_test

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/epistemer/models/column"
)

type fakeText struct {
	String string
	Valid  bool
}

func (t fakeText) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.String, nil
}

type fakeNumeric struct {
	Digits string
	Valid  bool
}

func (n fakeNumeric) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Digits, nil
}

type fakeDate struct {
	Time     time.Time
	Infinity int8
	Valid    bool
}

func (d fakeDate) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	if d.Infinity > 0 {
		return "infinity", nil
	}
	if d.Infinity < 0 {
		return "-infinity", nil
	}
	return d.Time, nil
}

type fakeInt4 struct {
	Int64 int64
	Valid bool
}

func (i fakeInt4) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}

type fakeFloatValuer struct {
	F     float64
	Valid bool
}

func (f fakeFloatValuer) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	return f.F, nil
}

type fakeBytesValuer struct {
	B []byte
}

func (b fakeBytesValuer) Value() (driver.Value, error) {
	return b.B, nil
}

type fakeBoolValuer struct {
	B bool
}

func (b fakeBoolValuer) Value() (driver.Value, error) {
	return b.B, nil
}

type ptrOnlyValuer struct {
	S string
}

func (p *ptrOnlyValuer) Value() (driver.Value, error) {
	return p.S, nil
}

type errValuer struct{}

func (errValuer) Value() (driver.Value, error) {
	return nil, errors.New("boom")
}

type notAValuer struct {
	X int
}

type (
	eventID  int64
	code     string
	fakeUUID [16]byte
)

func TestField(t *testing.T) {
	now := time.Date(2026, 11, 24, 12, 30, 0, 0, time.UTC)
	day := time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)
	uuidBytes := [16]byte{0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x14, 0x17, 0x40, 0x00}
	wantUUID := "123e4567-e89b-12d3-a456-426614174000"

	strVal := "hello"
	var nilStr *string
	i64 := int64(42)
	pi64 := &i64
	var nilPI64 *int64
	var nilPPI64 **int64
	i32 := int32(32)
	i16 := int16(16)
	f64 := 3.5
	bTrue := true

	tests := []struct {
		name string
		get  func() column.Value
		want column.Value
		nan  bool
	}{
		// Valuer: pgtype.Text-like
		{
			name: "valuer text valid",
			get: func() column.Value {
				return column.Field(column.Text, func(v fakeText) any { return v })(fakeText{String: "abc", Valid: true})
			},
			want: column.Value{Text: "abc"},
		},
		{
			name: "valuer text null",
			get: func() column.Value {
				return column.Field(column.Text, func(v fakeText) any { return v })(fakeText{Valid: false})
			},
			want: column.Value{Null: true},
		},
		{
			name: "valuer bytes to text",
			get: func() column.Value {
				return column.Field(column.Text, func(v fakeBytesValuer) any { return v })(fakeBytesValuer{B: []byte("raw")})
			},
			want: column.Value{Text: "raw"},
		},
		// Valuer: pgtype.Numeric-like
		{
			name: "valuer numeric digits",
			get: func() column.Value {
				return column.Field(column.Number, func(v fakeNumeric) any { return v })(fakeNumeric{Digits: "123.45", Valid: true})
			},
			want: column.Value{Num: 123.45, Text: "123.45"},
		},
		{
			name: "valuer numeric NaN",
			get: func() column.Value {
				return column.Field(column.Number, func(v fakeNumeric) any { return v })(fakeNumeric{Digits: "NaN", Valid: true})
			},
			want: column.Value{Text: "NaN"},
			nan:  true,
		},
		{
			name: "valuer numeric null",
			get: func() column.Value {
				return column.Field(column.Number, func(v fakeNumeric) any { return v })(fakeNumeric{Valid: false})
			},
			want: column.Value{Null: true},
		},
		// Valuer: pgtype.Date-like
		{
			name: "valuer date time.Time",
			get: func() column.Value {
				return column.Field(column.Date, func(v fakeDate) any { return v })(fakeDate{Time: day, Valid: true})
			},
			want: column.Value{Time: day},
		},
		{
			name: "valuer date +infinity",
			get: func() column.Value {
				return column.Field(column.Date, func(v fakeDate) any { return v })(fakeDate{Infinity: 1, Valid: true})
			},
			want: column.Value{Text: "infinity"},
		},
		{
			name: "valuer time -infinity",
			get: func() column.Value {
				return column.Field(column.Time, func(v fakeDate) any { return v })(fakeDate{Infinity: -1, Valid: true})
			},
			want: column.Value{Text: "-infinity"},
		},
		// Valuer: pgtype.Int4-like
		{
			name: "valuer int4 to integer",
			get: func() column.Value {
				return column.Field(column.Integer, func(v fakeInt4) any { return v })(fakeInt4{Int64: 7, Valid: true})
			},
			want: column.Value{Int: 7},
		},
		{
			name: "valuer int4 to number",
			get: func() column.Value {
				return column.Field(column.Number, func(v fakeInt4) any { return v })(fakeInt4{Int64: 7, Valid: true})
			},
			want: column.Value{Num: 7},
		},
		// Valuer: float64 & bool
		{
			name: "valuer float64",
			get: func() column.Value {
				return column.Field(column.Number, func(v fakeFloatValuer) any { return v })(fakeFloatValuer{F: 2.25, Valid: true})
			},
			want: column.Value{Num: 2.25},
		},
		{
			name: "valuer bool",
			get: func() column.Value {
				return column.Field(column.Bool, func(v fakeBoolValuer) any { return v })(fakeBoolValuer{B: true})
			},
			want: column.Value{Bool: true},
		},
		// Pointer-receiver-only Valuer
		{
			name: "pointer-receiver valuer by value",
			get: func() column.Value {
				return column.Field(column.Text, func(v ptrOnlyValuer) any { return v })(ptrOnlyValuer{S: "ptr"})
			},
			want: column.Value{Text: "ptr"},
		},
		{
			name: "pointer-receiver valuer by pointer non-nil",
			get: func() column.Value {
				return column.Field(column.Text, func(v *ptrOnlyValuer) any { return v })(&ptrOnlyValuer{S: "ptr2"})
			},
			want: column.Value{Text: "ptr2"},
		},
		{
			name: "pointer-receiver valuer by pointer nil",
			get:  func() column.Value { return column.Field(column.Text, func(v *ptrOnlyValuer) any { return v })(nil) },
			want: column.Value{Null: true},
		},
		// Fast paths & pointers
		{
			name: "string",
			get:  func() column.Value { return column.Field(column.Text, func(v string) any { return v })("hi") },
			want: column.Value{Text: "hi"},
		},
		{
			name: "*string non-nil",
			get:  func() column.Value { return column.Field(column.Text, func(v *string) any { return v })(&strVal) },
			want: column.Value{Text: "hello"},
		},
		{
			name: "*string nil",
			get:  func() column.Value { return column.Field(column.Text, func(v *string) any { return v })(nilStr) },
			want: column.Value{Null: true},
		},
		{
			name: "**int64 non-nil",
			get:  func() column.Value { return column.Field(column.Integer, func(v **int64) any { return v })(&pi64) },
			want: column.Value{Int: 42},
		},
		{
			name: "**int64 inner nil",
			get:  func() column.Value { return column.Field(column.Integer, func(v **int64) any { return v })(&nilPI64) },
			want: column.Value{Null: true},
		},
		{
			name: "**int64 outer nil",
			get:  func() column.Value { return column.Field(column.Integer, func(v **int64) any { return v })(nilPPI64) },
			want: column.Value{Null: true},
		},
		{
			name: "int32 fast path",
			get:  func() column.Value { return column.Field(column.Integer, func(v int32) any { return v })(i32) },
			want: column.Value{Int: 32},
		},
		{
			name: "*int16 fast path",
			get:  func() column.Value { return column.Field(column.Integer, func(v *int16) any { return v })(&i16) },
			want: column.Value{Int: 16},
		},
		{
			name: "*float64 fast path",
			get:  func() column.Value { return column.Field(column.Number, func(v *float64) any { return v })(&f64) },
			want: column.Value{Num: 3.5},
		},
		{
			name: "*bool fast path",
			get:  func() column.Value { return column.Field(column.Bool, func(v *bool) any { return v })(&bTrue) },
			want: column.Value{Bool: true},
		},
		// Named basic types
		{
			name: "named int64 as integer",
			get: func() column.Value {
				return column.Field(column.Integer, func(v eventID) any { return v })(eventID(99))
			},
			want: column.Value{Int: 99},
		},
		{
			name: "named int64 as number",
			get:  func() column.Value { return column.Field(column.Number, func(v eventID) any { return v })(eventID(99)) },
			want: column.Value{Num: 99},
		},
		{
			name: "named string as text",
			get:  func() column.Value { return column.Field(column.Text, func(v code) any { return v })(code("X1")) },
			want: column.Value{Text: "X1"},
		},
		{
			name: "named string as number",
			get:  func() column.Value { return column.Field(column.Number, func(v code) any { return v })(code("12.5")) },
			want: column.Value{Num: 12.5, Text: "12.5"},
		},
		{
			name: "uint16 as integer",
			get: func() column.Value {
				return column.Field(column.Integer, func(v uint16) any { return v })(uint16(60000))
			},
			want: column.Value{Int: 60000},
		},
		{
			name: "uint16 as number",
			get:  func() column.Value { return column.Field(column.Number, func(v uint16) any { return v })(uint16(12)) },
			want: column.Value{Num: 12},
		},
		{
			name: "float32 as number",
			get: func() column.Value {
				return column.Field(column.Number, func(v float32) any { return v })(float32(1.5))
			},
			want: column.Value{Num: 1.5},
		},
		{
			name: "bool",
			get:  func() column.Value { return column.Field(column.Bool, func(v bool) any { return v })(true) },
			want: column.Value{Bool: true},
		},
		{
			name: "[]byte as text",
			get: func() column.Value {
				return column.Field(column.Text, func(v []byte) any { return v })([]byte("bytes"))
			},
			want: column.Value{Text: "bytes"},
		},
		{
			name: "[16]byte as uuid",
			get:  func() column.Value { return column.Field(column.Text, func(v [16]byte) any { return v })(uuidBytes) },
			want: column.Value{Text: wantUUID},
		},
		{
			name: "named [16]byte as uuid",
			get: func() column.Value {
				return column.Field(column.Text, func(v fakeUUID) any { return v })(fakeUUID(uuidBytes))
			},
			want: column.Value{Text: wantUUID},
		},
		{
			name: "time.Time as time",
			get:  func() column.Value { return column.Field(column.Time, func(v time.Time) any { return v })(now) },
			want: column.Value{Time: now},
		},
		{
			name: "*time.Time non-nil as date",
			get:  func() column.Value { return column.Field(column.Date, func(v *time.Time) any { return v })(&day) },
			want: column.Value{Time: day},
		},
		{
			name: "*time.Time nil as date",
			get:  func() column.Value { return column.Field(column.Date, func(v *time.Time) any { return v })(nil) },
			want: column.Value{Null: true},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.get()
			if tc.nan {
				if !math.IsNaN(got.Num) || got.Text != tc.want.Text || got.Null != tc.want.Null {
					t.Errorf("got %+v, want NaN with Text=%q", got, tc.want.Text)
				}
				return
			}
			if got != tc.want {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestFieldPanics(t *testing.T) {
	t.Run("init non-valuer struct", func(t *testing.T) {
		msg := catchPanic(func() {
			column.Field(column.Number, func(v notAValuer) any { return v })
		})
		if msg != "column: number can't read column_test.notAValuer" {
			t.Errorf("panic = %q", msg)
		}
	})

	t.Run("init text over int", func(t *testing.T) {
		msg := catchPanic(func() {
			column.Field(column.Text, func(v int) any { return v })
		})
		if msg != "column: text can't read int" {
			t.Errorf("panic = %q", msg)
		}
	})

	t.Run("init enum kind", func(t *testing.T) {
		msg := catchPanic(func() {
			column.Field(column.Enum, func(v string) any { return v })
		})
		if msg != "column: Field doesn't read enums; generated enum getters read the field themselves" {
			t.Errorf("panic = %q", msg)
		}
	})

	t.Run("init interface field", func(t *testing.T) {
		msg := catchPanic(func() {
			column.Field(column.Text, func(v any) any { return v })
		})
		if msg != "column: text can't read a field of interface type" {
			t.Errorf("panic = %q", msg)
		}
	})

	t.Run("row valuer wrong type", func(t *testing.T) {
		get := column.Field(column.Text, func(v fakeBoolValuer) any { return v })
		msg := catchPanic(func() {
			get(fakeBoolValuer{B: true})
		})
		if msg != "column: text can't take a bool from column_test.fakeBoolValuer" {
			t.Errorf("panic = %q", msg)
		}
	})

	t.Run("row valuer error", func(t *testing.T) {
		get := column.Field(column.Text, func(v errValuer) any { return v })
		msg := catchPanic(func() {
			get(errValuer{})
		})
		if !strings.Contains(msg, "column_test.errValuer") || !strings.Contains(msg, "boom") {
			t.Errorf("panic = %q, want type and error", msg)
		}
	})
}

func TestFieldCallsProbeOnce(t *testing.T) {
	calls := 0
	get := column.Field(column.Text, func(s string) any {
		calls++
		return s
	})
	if calls != 1 {
		t.Fatalf("init calls = %d, want 1", calls)
	}
	_ = get("a")
	_ = get("b")
	if calls != 3 {
		t.Fatalf("total calls = %d, want 3 (1 init + 2 rows)", calls)
	}
}

func catchPanic(fn func()) (msg string) {
	defer func() {
		if r := recover(); r != nil {
			msg = fmt.Sprint(r)
		}
	}()
	fn()
	return ""
}

type codeText string

// namedValuer returns a named string, a named int64 or a [16]byte from
// Value, as an app's own Valuer may.
type namedValuer struct{ v driver.Value }

func (n namedValuer) Value() (driver.Value, error) { return n.v, nil }

type namedInt int64

// TestFieldReview covers cases a review of A5 found untested.
func TestFieldReview(t *testing.T) {
	text := column.Field(column.Text, func(r namedValuer) any { return r })
	if v := text(namedValuer{codeText("x")}); v.Text != "x" {
		t.Errorf("named string from a Valuer: %+v", v)
	}
	id := [16]byte{0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x14, 0x17, 0x40, 0x00}
	if v := text(namedValuer{id}); v.Text != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("[16]byte from a Valuer: %+v", v)
	}
	integer := column.Field(column.Integer, func(r namedValuer) any { return r })
	if v := integer(namedValuer{namedInt(5)}); v.Int != 5 {
		t.Errorf("named int64 from a Valuer: %+v", v)
	}
	number := column.Field(column.Number, func(r namedValuer) any { return r })
	if msg := catchPanic(func() { number(namedValuer{id}) }); msg != "column: number can't take a [16]uint8 from column_test.namedValuer" {
		t.Errorf("[16]byte to a number: %q", msg)
	}
	if msg := catchPanic(func() { column.Field(column.Text, func(r *driver.Valuer) any { return r }) }); msg != "column: text can't read *driver.Valuer" {
		t.Errorf("a pointer to an interface: %q", msg)
	}
	for _, c := range []struct {
		name string
		fn   func()
		want string
	}{
		{"string as bool", func() { column.Field(column.Bool, func(r string) any { return r }) }, "column: bool can't read string"},
		{"int64 as text", func() { column.Field(column.Text, func(r int64) any { return r }) }, "column: text can't read int64"},
		{"float64 as integer", func() { column.Field(column.Integer, func(r float64) any { return r }) }, "column: integer can't read float64"},
		{"bool as text", func() { column.Field(column.Text, func(r *bool) any { return r }) }, "column: text can't read *bool"},
		{"time as text", func() { column.Field(column.Text, func(r time.Time) any { return r }) }, "column: text can't read time.Time"},
	} {
		if msg := catchPanic(c.fn); msg != c.want {
			t.Errorf("%s: %q, want %q", c.name, msg, c.want)
		}
	}
	if v := column.Field(column.Number, func(r string) any { return r })("12.50"); v.Num != 12.5 || v.Text != "12.50" {
		t.Errorf("a string as a numeric's digits: %+v", v)
	}
	ptrText := column.Field(column.Text, func(r *fakeText) any { return r })
	if !ptrText(nil).Null || !ptrText(&fakeText{}).Null || ptrText(&fakeText{String: "a", Valid: true}).Text != "a" {
		t.Error("a pointer to a Valuer by value")
	}
}

type myByte uint8

// TestFieldReview2 covers review 0015's cases.
func TestFieldReview2(t *testing.T) {
	var b [16]myByte
	b[0] = 0x12
	text := column.Field(column.Text, func(r namedValuer) any { return r })
	if v := text(namedValuer{b}); v.Text != "12000000-0000-0000-0000-000000000000" {
		t.Errorf("[16]myByte from a Valuer: %+v", v)
	}
	if v := column.Field(column.Text, func(r [16]myByte) any { return r })(b); v.Text != "12000000-0000-0000-0000-000000000000" {
		t.Errorf("a [16]myByte field: %+v", v)
	}
	if msg := catchPanic(func() { text(namedValuer{int32(5)}) }); msg != "column: text can't take a int32 from column_test.namedValuer" {
		t.Errorf("an int32 to text names int32: %q", msg)
	}
}
