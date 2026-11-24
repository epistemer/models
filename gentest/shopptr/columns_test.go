package shopptr_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/epistemer/models/column"
	"github.com/epistemer/models/gentest/ids"
	"github.com/epistemer/models/gentest/shop"
	"github.com/epistemer/models/gentest/shopptr"
	"github.com/epistemer/models/listing"
)

func ptr[T any](v T) *T { return &v }

// Every kind's getter over sqlc's pointers and overrides (A5): set, and
// null as a nil pointer.
func TestColumns(t *testing.T) {
	c := shopptr.ShopKindsViewColumns
	day := time.Date(2026, 11, 22, 0, 0, 0, 0, time.UTC)
	at := time.Date(2026, 11, 22, 9, 30, 0, 0, time.UTC)
	full := shopptr.ShopKindsView{
		ID: ids.KindID(7), Name: "Bear", Note: ptr("soft"),
		Ref:   ptr(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
		Small: 3, Medium: ptr(int32(40)), Ratio: ptr(float32(0.5)), Score: 9.25,
		Price:  ptr(decimal.RequireFromString("123.45")),
		Active: ptr(true), Day: day, At: &at, Seen: &at, Span: ptr(time.Minute),
		State: shopptr.NullOrderStatus{OrderStatus: shopptr.OrderStatusPaid, Valid: true},
	}
	for _, tc := range []struct {
		col  column.Column[shopptr.ShopKindsView]
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
		{c.Seen, column.Time, func(v column.Value) bool { return v.Time.Equal(at) }},
		{c.State, column.Enum, func(v column.Value) bool { return v.Enum != nil && v.Enum.Value() == "paid" }},
	} {
		if tc.col.Kind != tc.kind {
			t.Errorf("%s: kind %v, want %v", tc.col.Name, tc.col.Kind, tc.kind)
		}
		if v := tc.col.Get(full); v.Null || !tc.ok(v) {
			t.Errorf("%s: %+v", tc.col.Name, v)
		}
	}
	var empty shopptr.ShopKindsView
	for _, col := range []column.Column[shopptr.ShopKindsView]{c.Note, c.Ref, c.Medium, c.Ratio, c.Price, c.Active, c.At, c.Seen, c.State} {
		if !col.Get(empty).Null {
			t.Errorf("%s of an empty row isn't null", col.Name)
		}
	}
	// A not-null time.Time override is a value, the zero time included.
	if v := c.Day.Get(empty); v.Null || !v.Time.IsZero() {
		t.Errorf("day of an empty row: %+v", v)
	}
}

// connect makes a schema of its own with the shop's tables, and drops it
// at the end. Skipped without MODELS_TEST_DATABASE_URL.
func connect(t *testing.T) *pgx.Conn {
	t.Helper()
	url := os.Getenv("MODELS_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("MODELS_TEST_DATABASE_URL isn't set")
	}
	ctx := t.Context()
	schema := fmt.Sprintf("models_test_%d", rand.Uint64())
	cfg, err := pgx.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	cfg.RuntimeParams["search_path"] = schema
	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		if _, err := conn.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Errorf("dropping %s: %v", schema, err)
		}
		conn.Close(ctx)
	})
	ddl, err := os.ReadFile("../../testdata/shop/schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{"CREATE SCHEMA " + schema, string(ddl)} {
		if _, err := conn.Exec(ctx, sql); err != nil {
			t.Fatalf("%.60s: %v", sql, err)
		}
	}
	return conn
}

// The same rows, listed through sqlc's defaults (package shop) and through
// pointers and overrides (package shopptr), give the same values in every
// column of the set (A5).
func TestBothConfigurations(t *testing.T) {
	conn := connect(t)
	ctx := t.Context()
	if _, err := conn.Exec(ctx, `INSERT INTO shop_kinds_view
	    (id, name, note, ref, small, medium, ratio, score, price, active, day, at, seen, span, state, tags, meta) VALUES
	    (1, 'Bear', 'soft', '123e4567-e89b-12d3-a456-426614174000', 3, 40, 0.5, 9.25, 123.45, true,
	     '2026-11-22', '2026-11-22 09:30', '2026-11-22 09:30+01', '90 minutes', 'paid', '{a,b}', '{"k": 1}'),
	    (2, 'Owl', NULL, NULL, -1, NULL, NULL, -0.5, NULL, NULL, '1999-12-31', NULL, NULL, NULL, NULL, NULL, NULL)`); err != nil {
		t.Fatal(err)
	}
	q := listing.Query{Sort: []listing.Sort{{Column: "id"}}}
	def, _, err := shop.New(conn).ListShopKindsView(ctx, q)
	if err != nil {
		t.Fatal(err)
	}
	ptrs, _, err := shopptr.New(conn).ListShopKindsView(ctx, q)
	if err != nil {
		t.Fatal(err)
	}
	if len(def) != 2 || len(ptrs) != 2 {
		t.Fatalf("rows: %d and %d", len(def), len(ptrs))
	}
	if *ptrs[0].Span != 90*time.Minute {
		t.Errorf("span: %v", *ptrs[0].Span)
	}
	dc, pc := shop.ShopKindsViewColumns, shopptr.ShopKindsViewColumns
	pairs := []struct {
		d column.Column[shop.ShopKindsView]
		p column.Column[shopptr.ShopKindsView]
	}{
		{dc.ID, pc.ID}, {dc.Name, pc.Name}, {dc.Note, pc.Note}, {dc.Ref, pc.Ref},
		{dc.Small, pc.Small}, {dc.Medium, pc.Medium}, {dc.Ratio, pc.Ratio}, {dc.Score, pc.Score},
		{dc.Price, pc.Price}, {dc.Active, pc.Active}, {dc.Day, pc.Day}, {dc.At, pc.At},
		{dc.Seen, pc.Seen}, {dc.State, pc.State},
	}
	for i := range def {
		for _, pr := range pairs {
			dv, pv := pr.d.Get(def[i]), pr.p.Get(ptrs[i])
			same := dv.Null == pv.Null && dv.Text == pv.Text && dv.Int == pv.Int &&
				dv.Num == pv.Num && dv.Bool == pv.Bool && dv.Time.Equal(pv.Time) &&
				(dv.Enum == nil) == (pv.Enum == nil) && (dv.Enum == nil || dv.Enum.Value() == pv.Enum.Value())
			if !same || pr.d.Kind != pr.p.Kind || pr.d.Name != pr.p.Name {
				t.Errorf("row %d, %s: defaults %+v, pointers %+v", i, pr.d.Name, dv, pv)
			}
		}
	}
}
