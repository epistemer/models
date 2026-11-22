package shop

import (
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/epistemer/models/column"
)

// Every kind's getter, set and null (interface §4.6, A4).
func TestColumns(t *testing.T) {
	c := ShopKindsViewColumns
	day := time.Date(2026, 11, 22, 0, 0, 0, 0, time.UTC)
	at := time.Date(2026, 11, 22, 9, 30, 0, 0, time.UTC)
	full := ShopKindsView{
		ID: 7, Name: "Bear", Note: pgtype.Text{String: "soft", Valid: true},
		Ref:   pgtype.UUID{Bytes: [16]byte{0x12, 0x3e, 0x45, 0x67, 0xe8, 0x9b, 0x12, 0xd3, 0xa4, 0x56, 0x42, 0x66, 0x14, 0x17, 0x40, 0x00}, Valid: true},
		Small: 3, Medium: pgtype.Int4{Int32: 40, Valid: true}, Ratio: pgtype.Float4{Float32: 0.5, Valid: true}, Score: 9.25,
		Price:  pgtype.Numeric{Int: big.NewInt(12345), Exp: -2, Valid: true},
		Active: pgtype.Bool{Bool: true, Valid: true}, Day: pgtype.Date{Time: day, Valid: true},
		At: pgtype.Timestamp{Time: at, Valid: true}, State: NullOrderStatus{OrderStatus: OrderStatusPaid, Valid: true},
	}
	for _, tc := range []struct {
		col  column.Column[ShopKindsView]
		kind column.Kind
		ok   func(column.Value) bool
	}{
		{c.ID, column.Integer, func(v column.Value) bool { return v.Int == 7 }},
		{c.Name, column.Text, func(v column.Value) bool { return v.Text == "Bear" }},
		{c.Note, column.Text, func(v column.Value) bool { return v.Text == "soft" }},
		{c.Ref, column.Text, func(v column.Value) bool { return v.Text == "123e4567-e89b-12d3-a456-426614174000" }},
		{c.Small, column.Integer, func(v column.Value) bool { return v.Int == 3 }},
		{c.Medium, column.Integer, func(v column.Value) bool { return v.Int == 40 }},
		{c.Ratio, column.Number, func(v column.Value) bool { return v.Num == 0.5 }},
		{c.Score, column.Number, func(v column.Value) bool { return v.Num == 9.25 }},
		{c.Price, column.Number, func(v column.Value) bool { return v.Num == 123.45 && v.Text == "123.45" }},
		{c.Active, column.Bool, func(v column.Value) bool { return v.Bool }},
		{c.Day, column.Date, func(v column.Value) bool { return v.Time.Equal(day) }},
		{c.At, column.Time, func(v column.Value) bool { return v.Time.Equal(at) }},
		{c.State, column.Enum, func(v column.Value) bool { return v.Enum != nil && v.Enum.Value() == "paid" }},
	} {
		if tc.col.Kind != tc.kind {
			t.Errorf("%s: kind %v, want %v", tc.col.Name, tc.col.Kind, tc.kind)
		}
		if v := tc.col.Get(full); v.Null || !tc.ok(v) {
			t.Errorf("%s: %+v", tc.col.Name, v)
		}
		if tc.col.Label == nil || tc.col.Label.K() != "table.shop_kinds_view.KindsView_"+zp(tc.col.Name) {
			t.Errorf("%s: label %v", tc.col.Name, tc.col.Label)
		}
	}
	// Nulls: every nullable column of an empty row.
	var empty ShopKindsView
	for _, col := range []column.Column[ShopKindsView]{c.Note, c.Ref, c.Medium, c.Ratio, c.Price, c.Active, c.Day, c.At, c.State} {
		if !col.Get(empty).Null {
			t.Errorf("%s of an empty row isn't null", col.Name)
		}
	}
	// Infinite dates and times aren't year 1.
	inf := ShopKindsView{Day: pgtype.Date{InfinityModifier: pgtype.Infinity, Valid: true}, At: pgtype.Timestamp{InfinityModifier: pgtype.NegativeInfinity, Valid: true}}
	if v := c.Day.Get(inf); v.Null || v.Text != "infinity" || !v.Time.IsZero() {
		t.Errorf("infinite day: %+v", v)
	}
	if v := c.At.Get(inf); v.Text != "-infinity" {
		t.Errorf("-infinite time: %+v", v)
	}
	if v := c.Price.Get(ShopKindsView{Price: pgtype.Numeric{NaN: true, Valid: true}}); v.Num == v.Num || v.Text != "NaN" {
		t.Errorf("NaN numeric: %+v", v)
	}
	if len(c.State.Options) != 3 || c.State.Options[1].Value() != "paid" {
		t.Errorf("state options %v", c.State.Options)
	}
	if ShopOrderListViewColumns.Payment.Kind != column.Enum || ShopOrderListViewColumns.Order.Name != "order" {
		t.Error("the order list's columns")
	}
}

func zp(name string) string {
	b := []byte(name)
	b[0] -= 'a' - 'A'
	return string(b)
}
