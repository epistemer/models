package shop_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/epistemer/models/gentest/shop"
	"github.com/epistemer/models/listing"
)

// queries connects to MODELS_TEST_DATABASE_URL, creates the shop fixture's
// schema in a schema of its own, fills it, and drops it at the end. The
// test is skipped without the variable; CI sets it on Linux.
func queries(t *testing.T) *shop.Queries {
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
	for _, sql := range []string{
		"CREATE SCHEMA " + schema,
		string(ddl),
		`INSERT INTO shop_customer (id, email, id_kind) VALUES
		   (1, 'ann@example.com', 'id'), (2, 'bob_50%@example.com', 'uuid')`,
		`INSERT INTO shop_order (id, customer_id, status, payment, shop_ref) VALUES
		   (10, 1, 'paid', 'pay-pal', 'r-10'),
		   (11, 1, 'pending', NULL, 'r-11'),
		   (12, 2, 'paid', 'cash:on:delivery', 'r-12'),
		   (13, 2, 'shipped', 'pay_pal', 'r-13'),
		   (14, 2, 'paid', '', 'r-14')`,
		`INSERT INTO shop_page_view (id, path, seen_at) VALUES
		   (1, '/', '2026-11-01T10:00:00Z'), (2, '/orders', '2026-11-02T10:00:00Z')`,
	} {
		if _, err := conn.Exec(ctx, sql); err != nil {
			t.Fatalf("%.60s: %v", sql, err)
		}
	}
	return shop.New(conn)
}

func filter(column string, op listing.Op, value any) listing.Filter {
	return listing.Filter{Column: column, Op: op, Value: value}
}

func ids(rows []shop.ShopOrderListView) []int64 {
	var out []int64
	for _, r := range rows {
		out = append(out, r.ID)
	}
	return out
}

func TestListings(t *testing.T) {
	q := queries(t)
	byID := []listing.Sort{{Column: "id"}}
	for _, c := range []struct {
		name  string
		lq    listing.Query
		ids   []int64
		total int64
	}{
		{"everything", listing.Query{Sort: byID}, []int64{10, 11, 12, 13, 14}, 5},
		{"enum filter, sqlc's type", listing.Query{Filters: []listing.Filter{filter("status", listing.Eq, shop.OrderStatusPaid)}, Sort: byID}, []int64{10, 12, 14}, 3},
		{"enum filter, a string", listing.Query{Filters: []listing.Filter{filter("status", listing.Ne, "paid")}, Sort: byID}, []int64{11, 13}, 2},
		{"every comparison", listing.Query{Filters: []listing.Filter{
			filter("id", listing.Gt, 10), filter("id", listing.Ge, 11), filter("id", listing.Lt, 14), filter("id", listing.Le, 13),
		}, Sort: byID}, []int64{11, 12, 13}, 3},
		{"search, case and type ignored", listing.Query{Search: "PAY", SearchIn: []string{"payment", "email"}, Sort: byID}, []int64{10, 13}, 2},
		{"search escapes % and _", listing.Query{Search: "b_50%", SearchIn: []string{"email"}, Sort: byID}, []int64{12, 13, 14}, 3},
		{"a wildcard is literal", listing.Query{Search: "_", SearchIn: []string{"order"}, Sort: byID}, nil, 0},
		{"search and a filter", listing.Query{
			Search: "r-1", SearchIn: []string{"email", "order"},
			Filters: []listing.Filter{filter("status", listing.Eq, "paid")}, Sort: byID,
		}, []int64{10, 12, 14}, 3},
		{"sort by a keyword column, descending", listing.Query{Sort: []listing.Sort{{Column: "order", Desc: true}}}, []int64{14, 13, 12, 11, 10}, 5},
		{"first page, counted", listing.Query{Sort: byID, Limit: 2}, []int64{10, 11}, 5},
		{"middle page, counted", listing.Query{Sort: byID, Limit: 2, Offset: 2}, []int64{12, 13}, 5},
		{"short last page, not counted", listing.Query{Sort: byID, Limit: 2, Offset: 4}, []int64{14}, 5},
		{"past the end, counted", listing.Query{Sort: byID, Limit: 2, Offset: 10}, nil, 5},
		{"offset without a limit", listing.Query{Sort: byID, Offset: 3}, []int64{13, 14}, 5},
	} {
		t.Run(c.name, func(t *testing.T) {
			rows, total, err := q.ListShopOrderListView(t.Context(), c.lq)
			if err != nil {
				t.Fatal(err)
			}
			if got := ids(rows); !slices.Equal(got, c.ids) {
				t.Errorf("ids %v, want %v", got, c.ids)
			}
			if total != c.total {
				t.Errorf("total %d, want %d", total, c.total)
			}
		})
	}
}

func TestListingFields(t *testing.T) {
	q := queries(t)
	rows, _, err := q.ListShopOrderListView(t.Context(), listing.Query{Filters: []listing.Filter{filter("id", listing.Eq, 11)}})
	if err != nil {
		t.Fatal(err)
	}
	want := shop.ShopOrderListView{ID: 11, Status: shop.OrderStatusPending, Email: "ann@example.com", Order: "r-11"}
	if len(rows) != 1 || rows[0] != want {
		t.Errorf("got %+v, want [%+v]", rows, want)
	}
	pages, total, err := q.ListShopPageView(t.Context(), listing.Query{Sort: []listing.Sort{{Column: "seen_at", Desc: true}}})
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 || len(pages) != 2 || pages[0].Path != "/orders" ||
		!pages[0].SeenAt.Time.Equal(time.Date(2026, 11, 2, 10, 0, 0, 0, time.UTC)) {
		t.Errorf("got %+v, %d", pages, total)
	}
}

func TestListingErrors(t *testing.T) {
	q := queries(t)
	_, _, err := q.ListShopOrderListView(t.Context(), listing.Query{Sort: []listing.Sort{{Column: "id; DROP TABLE shop_order"}}})
	if !errors.Is(err, listing.ErrQuery) {
		t.Errorf("unknown column: err = %v, want ErrQuery", err)
	}
	_, _, err = q.ListShopOrderListView(t.Context(), listing.Query{Filters: []listing.Filter{filter("status", listing.Eq, "bogus")}})
	if err == nil || errors.Is(err, listing.ErrQuery) {
		t.Errorf("bad enum value: err = %v, want PostgreSQL's", err)
	}
}
