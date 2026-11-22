package column_test

import (
	"math"
	"strings"
	"testing"

	"github.com/epistemer/models/column"
)

func TestUUID(t *testing.T) {
	b := [16]byte{0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x14, 0x17, 0x40, 0x00}
	if got := column.UUID(b); got != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("UUID = %q", got)
	}
}

func TestKindString(t *testing.T) {
	for k, want := range map[column.Kind]string{column.Text: "text", column.Enum: "enum", 0: "unknown", 99: "unknown"} {
		if got := k.String(); got != want {
			t.Errorf("%d.String() = %q, want %q", k, got, want)
		}
	}
}

func TestNumeric(t *testing.T) {
	if v := column.Numeric("123.45"); v.Num != 123.45 || v.Text != "123.45" {
		t.Errorf("%+v", v)
	}
	if v := column.Numeric("NaN"); v.Num == v.Num {
		t.Errorf("NaN: %+v", v)
	}
	if v := column.Numeric("-Infinity"); !math.IsInf(v.Num, -1) {
		t.Errorf("-Infinity: %+v", v)
	}
	if v := column.Numeric("1" + strings.Repeat("0", 400)); !math.IsInf(v.Num, 1) {
		t.Errorf("past float64: %+v", v)
	}
	if column.Infinite(-1).Text != "-infinity" || column.Infinite(1).Text != "infinity" {
		t.Error("Infinite")
	}
}
