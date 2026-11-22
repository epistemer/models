// Package column describes a view's columns for a table package: what the
// schema knows of each, and a getter from a row to a plain value with
// nulls unwrapped, so a table package never imports the driver
// (interface §4.6, amendment A4, ADR 0007). The views generator writes a
// <Type>Columns set per view; a column can also be built by hand, for a
// test or rows that don't come from sqlc.
package column

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/epistemer/i18n"

	"github.com/epistemer/models/option"
)

// Kind is a column's kind, from its SQL type.
type Kind uint8

const (
	Text    Kind = iota + 1 // text, varchar, char, citext, uuid
	Integer                 // int2, int4, int8
	Number                  // float4, float8, numeric
	Bool                    // bool
	Date                    // date
	Time                    // timestamp, timestamptz
	Enum                    // an enum of the schema
)

var kindNames = [...]string{"", "text", "integer", "number", "bool", "date", "time", "enum"}

// String is the kind's name: "text", "integer", …
func (k Kind) String() string {
	if int(k) < len(kindNames) && k != 0 {
		return kindNames[k]
	}
	return "unknown"
}

// Value is a cell's plain value: the field of its kind holds it, unless
// Null.
type Value struct {
	Null bool
	Text string              // Text; a numeric's exact digits; "infinity" or "-infinity" for an infinite Date or Time
	Int  int64               // Integer
	Num  float64             // Number
	Bool bool                // Bool
	Time time.Time           // Date (midnight UTC) and Time
	Enum option.SelectOption // Enum
}

// Column is one column of a relation whose rows are R.
type Column[R any] struct {
	Name    string                // its SQL name
	Label   i18n.TKer             // table.<relation>.<Field>, the label models generates
	Kind    Kind                  //
	Options []option.SelectOption // an enum's, in schema order
	Get     func(R) Value
}

// Numeric is a numeric's value from its exact digits, as PostgreSQL
// writes them: Num parsed from them (NaN, ±Inf, and ±Inf past float64's
// range), Text the digits. For generated getters.
func Numeric(digits string) Value {
	f, _ := strconv.ParseFloat(digits, 64)
	return Value{Num: f, Text: digits}
}

// Infinite is an infinite date or time: Time is zero and Text is
// "infinity", or "-infinity" when sign is negative. For generated getters.
func Infinite(sign int8) Value {
	if sign < 0 {
		return Value{Text: "-infinity"}
	}
	return Value{Text: "infinity"}
}

// UUID formats a UUID's bytes canonically, for generated getters.
func UUID(b [16]byte) string {
	var s [36]byte
	hex.Encode(s[0:8], b[0:4])
	s[8] = '-'
	hex.Encode(s[9:13], b[4:6])
	s[13] = '-'
	hex.Encode(s[14:18], b[6:8])
	s[18] = '-'
	hex.Encode(s[19:23], b[8:10])
	s[23] = '-'
	hex.Encode(s[24:36], b[10:16])
	return string(s[:])
}
