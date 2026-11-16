package listing_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/epistemer/models/listing"
)

var columns = []string{"id", "status", "order", "note"}

func TestBuild(t *testing.T) {
	cases := []struct {
		name        string
		q           listing.Query
		rows, count string
		args        []any
	}{
		{
			name:  "everything",
			rows:  `SELECT "id", "status", "order", "note" FROM "order_view"`,
			count: `SELECT count(*) FROM "order_view"`,
		},
		{
			name: "search, filters, sort and page",
			q: listing.Query{
				Search:   "50%_off\\",
				SearchIn: []string{"note", "status"},
				Filters: []listing.Filter{
					{Column: "status", Op: listing.Eq, Value: "paid"},
					{Column: "id", Op: listing.Ge, Value: 10},
				},
				Sort:   []listing.Sort{{Column: "order", Desc: true}, {Column: "id"}},
				Limit:  20,
				Offset: 40,
			},
			rows: `SELECT "id", "status", "order", "note" FROM "order_view"` +
				` WHERE (CAST("note" AS text) ILIKE $1 ESCAPE '\' OR CAST("status" AS text) ILIKE $1 ESCAPE '\')` +
				` AND "status" = $2 AND "id" >= $3 ORDER BY "order" DESC, "id" LIMIT 20 OFFSET 40`,
			count: `SELECT count(*) FROM "order_view"` +
				` WHERE (CAST("note" AS text) ILIKE $1 ESCAPE '\' OR CAST("status" AS text) ILIKE $1 ESCAPE '\')` +
				` AND "status" = $2 AND "id" >= $3`,
			args: []any{`%50\%\_off\\%`, "paid", 10},
		},
		{
			name: "every op",
			q: listing.Query{Filters: []listing.Filter{
				{"id", listing.Eq, 1}, {"id", listing.Ne, 2}, {"id", listing.Lt, 3},
				{"id", listing.Le, 4}, {"id", listing.Gt, 5}, {"id", listing.Ge, 6},
			}},
			rows: `SELECT "id", "status", "order", "note" FROM "order_view"` +
				` WHERE "id" = $1 AND "id" <> $2 AND "id" < $3 AND "id" <= $4 AND "id" > $5 AND "id" >= $6`,
			count: `SELECT count(*) FROM "order_view"` +
				` WHERE "id" = $1 AND "id" <> $2 AND "id" < $3 AND "id" <= $4 AND "id" > $5 AND "id" >= $6`,
			args: []any{1, 2, 3, 4, 5, 6},
		},
		{
			name:  "search columns without a search",
			q:     listing.Query{SearchIn: []string{"note"}, Offset: 5},
			rows:  `SELECT "id", "status", "order", "note" FROM "order_view" OFFSET 5`,
			count: `SELECT count(*) FROM "order_view"`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rows, count, args, err := listing.Build(c.q, "order_view", columns)
			if err != nil {
				t.Fatal(err)
			}
			if rows != c.rows {
				t.Errorf("rows:\n got %s\nwant %s", rows, c.rows)
			}
			if count != c.count {
				t.Errorf("count:\n got %s\nwant %s", count, c.count)
			}
			if !reflect.DeepEqual(args, c.args) {
				t.Errorf("args: got %#v, want %#v", args, c.args)
			}
		})
	}
}

func TestBuildQuotesNames(t *testing.T) {
	rows, _, _, err := listing.Build(listing.Query{Sort: []listing.Sort{{Column: `a"b`}}}, `we"ird_view`, []string{`a"b`})
	if err != nil {
		t.Fatal(err)
	}
	if want := `SELECT "a""b" FROM "we""ird_view" ORDER BY "a""b"`; rows != want {
		t.Errorf("got %s, want %s", rows, want)
	}
}

func TestBuildRefuses(t *testing.T) {
	cases := []struct {
		name     string
		q        listing.Query
		relation string
		columns  []string
		want     string
	}{
		{"unknown filter column", listing.Query{Filters: []listing.Filter{{"nope", listing.Eq, 1}}}, "v", columns, `"v" has no column "nope"`},
		{"unknown sort column", listing.Query{Sort: []listing.Sort{{Column: "nope"}}}, "v", columns, `"v" has no column "nope"`},
		{"unknown search column, no search", listing.Query{SearchIn: []string{"nope"}}, "v", columns, `"v" has no column "nope"`},
		{"injected column", listing.Query{Sort: []listing.Sort{{Column: `id"; DROP TABLE x; --`}}}, "v", columns, `has no column`},
		{"unknown search column", listing.Query{Search: "x", SearchIn: []string{"note", "nope"}}, "v", columns, `"v" has no column "nope"`},
		{"NUL in search", listing.Query{Search: "a\x00", SearchIn: []string{"note"}}, "v", columns, `NUL`},
		{"search without columns", listing.Query{Search: "x"}, "v", columns, `no column to search in`},
		{"unknown op", listing.Query{Filters: []listing.Filter{{"id", "LIKE", 1}}}, "v", columns, `filter on "id": unknown op "LIKE"`},
		{"zero op", listing.Query{Filters: []listing.Filter{{Column: "id", Value: 1}}}, "v", columns, `filter on "id": unknown op ""`},
		{"nil value", listing.Query{Filters: []listing.Filter{{"id", listing.Eq, nil}}}, "v", columns, `filter on "id": nil value`},
		{"typed nil value", listing.Query{Filters: []listing.Filter{{"note", listing.Ne, (*string)(nil)}}}, "v", columns, `filter on "note": nil value`},
		{"nil slice value", listing.Query{Filters: []listing.Filter{{"note", listing.Eq, []byte(nil)}}}, "v", columns, `nil value`},
		{"error after arguments", listing.Query{Search: "x", SearchIn: []string{"note"}, Filters: []listing.Filter{{"id", listing.Eq, 1}}, Sort: []listing.Sort{{Column: "nope"}}}, "v", columns, `"v" has no column "nope"`},
		{"negative limit", listing.Query{Limit: -1}, "v", columns, `negative limit`},
		{"negative offset", listing.Query{Offset: -1}, "v", columns, `negative offset`},
		{"empty relation", listing.Query{}, "", columns, `relation ""`},
		{"NUL in relation", listing.Query{}, "v\x00", columns, `relation`},
		{"no columns", listing.Query{}, "v", nil, `has no columns`},
		{"empty column", listing.Query{}, "v", []string{""}, `column ""`},
		{"NUL in column", listing.Query{}, "v", []string{"a\x00"}, `column`},
		{"column twice", listing.Query{}, "v", []string{"id", "id"}, `column "id"`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rows, count, args, err := listing.Build(c.q, c.relation, c.columns)
			if !errors.Is(err, listing.ErrQuery) {
				t.Fatalf("err = %v, want one wrapping ErrQuery", err)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Errorf("err = %q, want it to contain %q", err, c.want)
			}
			if rows != "" || count != "" || args != nil {
				t.Errorf("got %q, %q, %v with the error", rows, count, args)
			}
		})
	}
}

func TestQueryIsItsOwnSource(t *testing.T) {
	q := listing.Query{Search: "x", SearchIn: []string{"note"}, Limit: 3}
	var src listing.Source = q
	if got := src.ListingQuery(); !reflect.DeepEqual(got, q) {
		t.Errorf("got %+v, want %+v", got, q)
	}
}

// names is an app's Labeler, as a link column would call it.
type names map[string]string

func (n names) LabelOf(_ context.Context, id string) (string, bool, error) {
	s, ok := n[id]
	return s, ok, nil
}

func TestLabeler(t *testing.T) {
	var l listing.Labeler = names{"7": "Ann"}
	if s, ok, err := l.LabelOf(t.Context(), "7"); s != "Ann" || !ok || err != nil {
		t.Errorf("got %q, %v, %v", s, ok, err)
	}
	if _, ok, _ := l.LabelOf(t.Context(), "8"); ok {
		t.Error("found a row that isn't there")
	}
}
