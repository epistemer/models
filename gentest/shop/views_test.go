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
func queries(t *testing.T) *pgx.Conn {
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
	return conn
}

// counting is sqlc's DBTX, counting the count statements (the only
// QueryRow a listing runs) and running before ahead of each.
type counting struct {
	shop.DBTX
	counts int
	before func()
}

func (c *counting) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	c.counts++
	if c.before != nil {
		c.before()
	}
	return c.DBTX.QueryRow(ctx, sql, args...)
}

// table stands for a table package's own state: it implements
// listing.Source, as a table would (models' interface §4.4).
type table struct {
	status shop.OrderStatus
	page   int
}

func (t table) ListingQuery() listing.Query {
	return listing.Query{
		Filters: []listing.Filter{filter("status", listing.Eq, t.status)},
		Sort:    []listing.Sort{{Column: "id"}},
		Limit:   2,
		Offset:  2 * t.page,
	}
}

func TestListingSource(t *testing.T) {
	q := shop.New(queries(t))
	rows, total, err := q.ListShopOrderListView(t.Context(), table{status: shop.OrderStatusPaid, page: 1})
	if err != nil {
		t.Fatal(err)
	}
	if got := ids(rows); !slices.Equal(got, []int64{14}) || total != 3 {
		t.Errorf("got %v and %d, want [14] and 3", got, total)
	}
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
	db := &counting{DBTX: queries(t)}
	q := shop.New(db)
	byID := []listing.Sort{{Column: "id"}}
	for _, c := range []struct {
		name    string
		lq      listing.Query
		ids     []int64
		total   int64
		counted bool
	}{
		{"everything", listing.Query{Sort: byID}, []int64{10, 11, 12, 13, 14}, 5, false},
		{"enum filter, sqlc's type", listing.Query{Filters: []listing.Filter{filter("status", listing.Eq, shop.OrderStatusPaid)}, Sort: byID}, []int64{10, 12, 14}, 3, false},
		{"enum filter, a string", listing.Query{Filters: []listing.Filter{filter("status", listing.Ne, "paid")}, Sort: byID}, []int64{11, 13}, 2, false},
		{"every comparison", listing.Query{Filters: []listing.Filter{
			filter("id", listing.Gt, 10), filter("id", listing.Ge, 11), filter("id", listing.Lt, 14), filter("id", listing.Le, 13),
		}, Sort: byID}, []int64{11, 12, 13}, 3, false},
		{"search, case and type ignored", listing.Query{Search: "PAY", SearchIn: []string{"payment", "email"}, Sort: byID}, []int64{10, 13}, 2, false},
		{"search escapes % and _", listing.Query{Search: "b_50%", SearchIn: []string{"email"}, Sort: byID}, []int64{12, 13, 14}, 3, false},
		{"a wildcard is literal", listing.Query{Search: "_", SearchIn: []string{"order"}, Sort: byID}, nil, 0, false},
		{"search and a filter", listing.Query{
			Search: "r-1", SearchIn: []string{"email", "order"},
			Filters: []listing.Filter{filter("status", listing.Eq, "paid")}, Sort: byID,
		}, []int64{10, 12, 14}, 3, false},
		{"sort by a keyword column, descending", listing.Query{Sort: []listing.Sort{{Column: "order", Desc: true}}}, []int64{14, 13, 12, 11, 10}, 5, false},
		{"first page", listing.Query{Sort: byID, Limit: 2}, []int64{10, 11}, 5, true},
		{"middle page", listing.Query{Sort: byID, Limit: 2, Offset: 2}, []int64{12, 13}, 5, true},
		{"short last page", listing.Query{Sort: byID, Limit: 2, Offset: 4}, []int64{14}, 5, false},
		{"full last page", listing.Query{Sort: byID, Limit: 2, Offset: 3}, []int64{13, 14}, 5, true},
		{"one page that fits exactly", listing.Query{Sort: byID, Limit: 5}, []int64{10, 11, 12, 13, 14}, 5, true},
		{"past the end", listing.Query{Sort: byID, Limit: 2, Offset: 10}, nil, 5, true},
		{"offset without a limit", listing.Query{Sort: byID, Offset: 3}, []int64{13, 14}, 5, false},
		{"offset past the end without a limit", listing.Query{Sort: byID, Offset: 10}, nil, 5, true},
		{"no match with a limit", listing.Query{Filters: []listing.Filter{filter("id", listing.Gt, 99)}, Limit: 2}, nil, 0, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			db.counts = 0
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
			if counted := db.counts > 0; counted != c.counted {
				t.Errorf("counted %v, want %v", counted, c.counted)
			}
		})
	}
}

// TestListingTotalFloor: rows deleted between the page and the count don't
// make the total smaller than the rows already returned.
func TestListingTotalFloor(t *testing.T) {
	conn := queries(t)
	db := &counting{DBTX: conn, before: func() {
		if _, err := conn.Exec(context.Background(), "DELETE FROM shop_order WHERE id >= 11"); err != nil {
			t.Error(err)
		}
	}}
	rows, total, err := shop.New(db).ListShopOrderListView(t.Context(),
		listing.Query{Sort: []listing.Sort{{Column: "id"}}, Limit: 2, Offset: 2})
	if err != nil {
		t.Fatal(err)
	}
	if got := ids(rows); !slices.Equal(got, []int64{12, 13}) || total != 4 {
		t.Errorf("got %v and %d, want [12 13] and 4", got, total)
	}
}

func TestListingFields(t *testing.T) {
	q := shop.New(queries(t))
	rows, _, err := q.ListShopOrderListView(t.Context(), listing.Query{
		Filters: []listing.Filter{filter("id", listing.Le, 11)}, Sort: []listing.Sort{{Column: "id"}}})
	if err != nil {
		t.Fatal(err)
	}
	want := []shop.ShopOrderListView{
		{ID: 10, Status: shop.OrderStatusPaid, Payment: shop.NullPaymentMethod{PaymentMethod: shop.PaymentMethodPayPal, Valid: true}, Email: "ann@example.com", Order: "r-10"},
		{ID: 11, Status: shop.OrderStatusPending, Email: "ann@example.com", Order: "r-11"},
	}
	if !slices.Equal(rows, want) {
		t.Errorf("got %+v, want %+v", rows, want)
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
	q := shop.New(queries(t))
	_, _, err := q.ListShopOrderListView(t.Context(), listing.Query{Sort: []listing.Sort{{Column: "id; DROP TABLE shop_order"}}})
	if !errors.Is(err, listing.ErrQuery) {
		t.Errorf("unknown column: err = %v, want ErrQuery", err)
	}
	_, _, err = q.ListShopOrderListView(t.Context(), listing.Query{Filters: []listing.Filter{filter("status", listing.Eq, "bogus")}})
	if err == nil || errors.Is(err, listing.ErrQuery) {
		t.Errorf("bad enum value: err = %v, want PostgreSQL's", err)
	}
}
