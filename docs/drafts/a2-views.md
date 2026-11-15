# Draft: interface amendment A2, listings of views

> Status: **draft, for epistemer's approval** (roadmap M5.2). Once
> approved, its text goes into [interface.md](../interface.md) as written
> below, and this file is deleted. The why is
> [ADR 0005](../adr/0005-view-listings.md).

A2 adds a generator and a runtime package. It changes nothing an app on
v0.1.0 writes or gets: `views` runs only when switched on.

## Changes to existing sections

- **§0 D1**: "Three trivial generators (genx) over sqlc's codegen request,
  and one runtime package, `listing`, for the code `views` generates."
- **§0 D3** adds: "`views` writes a `List<Type>` method for every relation
  whose name ends in `_view` (§4.3)."
- **§0 D10**: "… package `models`, with only the metafile, its templates
  and a `doc.go`; package `listing` (§4.4)."
- **§2**: the manifest's example becomes `generators: [labels, enums,
  views]`, and the text adds: "`views` needs sqlc's `sql_package:
  "pgx/v5"` (§7)."
- **§4**'s table gets a row: `views` | the one `sqlc.request/v1` file |
  `<dir>/zz_<package>_views.go`. "With no relation whose name ends in
  `_view`, `views` writes nothing."
- **§9**: "views' query builder" is removed.

## New §4.3 The views file

Every table or view of the default schema whose SQL name ends in `_view`,
in the request's order, gets a method on sqlc's `Queries`. `<Type>` is
sqlc's type for the relation, `structName` of its name (a name ending in
`_view` isn't made singular; one starting with a digit gets sqlc's `_`,
as in `_1View`); each field is `structName` of a column's name. For a view `order_summary_view`:

```go
// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
	"context"

	"github.com/epistemer/models/listing"
)

// Listings is the listing methods of Queries, one per view.
type Listings interface {
	ListOrderSummaryView(ctx context.Context, lq listing.Query) ([]OrderSummaryView, int64, error)
}

var _ Listings = (*Queries)(nil)

// ListOrderSummaryView returns the rows of order_summary_view that lq
// selects, and how many rows its search and filters match.
func (q *Queries) ListOrderSummaryView(ctx context.Context, lq listing.Query) ([]OrderSummaryView, int64, error) {
	rowsSQL, countSQL, args, err := listing.Build(lq, "order_summary_view",
		[]string{"id", "status", "total", "placed", "note"})
	// … runs rowsSQL through q.db (pgx/v5's DBTX), scanning each row into
	// &i.ID, &i.Status, &i.Total, &i.Placed, &i.Note; then countSQL, when
	// the rows can't tell the total.
}
```

- An error from `listing.Build` returns before any query runs.
- The rows tell the total, and the count isn't run, when they end the
  listing: `Limit` is 0 or more than the rows returned, and either
  `Offset` is 0 or at least one row came back. The total is then
  `Offset` plus the number of rows. Otherwise the count runs, and the total
  is never less than `Offset` plus the rows returned. The two statements
  share no transaction.
- An app embeds `Listings` in its own interface:
  `type Q interface { Querier; Listings }`.

## New §4.4 Package listing

`github.com/epistemer/models/listing`, no dependencies:

```go
// Query selects a page of a view's rows.
type Query struct {
	Search   string   // rows where a SearchIn column, as text, contains it, case ignored
	SearchIn []string // column names; needed when Search isn't ""
	Filters  []Filter // every one must hold
	Sort     []Sort   // in order
	Limit    int      // 0: no limit
	Offset   int
}

type Filter struct {
	Column string
	Op     Op
	Value  any // passed to the driver, which encodes it by the column's type
}

type Op string

const (
	Eq Op = "="
	Ne Op = "<>"
	Lt Op = "<"
	Le Op = "<="
	Gt Op = ">"
	Ge Op = ">="
)

type Sort struct {
	Column string
	Desc   bool
}

// Build returns the statement for a page of relation's rows and the
// statement counting the rows that match, with their arguments. For
// generated code.
func Build(q Query, relation string, columns []string) (rows, count string, args []any, err error)

// ErrQuery is wrapped by every error of Build.
var ErrQuery = errors.New("listing: bad query")
```

- Column names are SQL names and must be among `columns`, in `SearchIn`
  too when `Search` is `""`; a name that isn't is an error naming it. So is
  a `Search` with no `SearchIn`, an `Op` not listed above, a `nil` `Value`
  (kept free for a later `IS NULL`), a negative `Limit` or `Offset`, an
  empty `relation` or `columns`, and a name holding a NUL.
- `ErrQuery` covers the query's shape only. A value its column can't take
  (`"bogus"` for an enum, `"x"` for a date) fails in pgx or PostgreSQL;
  the caller converts values by the column's type first.
- The search matches PostgreSQL's text form of a value: an enum's SQL
  value, not its label; a timestamp in the session's time zone.
- `Limit` and `Offset` page by `Sort`, so stable pages need a `Sort` ending
  on a column that is unique in the view.
- The SQL is ADR 0005's D6: every identifier double-quoted with `"`
  doubled; the search a parenthesized group of `CAST(c AS text) ILIKE $n
  ESCAPE '\'` joined by `OR`, its text with `\`, then `%` and `_`,
  escaped; filters `AND`ed with it; `LIMIT` and `OFFSET` as integers. Its
  exact text isn't part of the contract.

## Additions to §6 Refusals

Made by `views` only, and only when the request has a relation whose name
ends in `_view`:

| Case | Example message |
|------|-----------------|
| The request's engine isn't `postgresql` | `models: views need PostgreSQL; the request's engine is "mysql"` |
| A `_view` relation or one of its columns has a name with a character outside ASCII, or whose `structName` is empty | `models: view "café_view": sqlc's Go name for it can't be derived` |
| A `_view` relation with no columns | `models: view "empty_view" has no columns` |
| Two columns of a `_view` relation with one `structName` (sqlc declares the field twice) | `models: view "dup_view": columns "user_1" and "user1" both give the field "User1"` |
| Two `_view` relations with one `<Type>` | `models: views "user_1_view" and "user1_view" both give sqlc the Go type "User1View"` |
| A query whose name, as written, is a method `List<Type>` | `models: query "ListOrderSummaryView" clashes with the listing method of view "order_summary_view"` |
| A query or enum that gives sqlc the name `Listings`: a query named so, or an enum whose `structName` is it | `models: enum "listings" gives sqlc the Go name "Listings", which views declares` |
| A query whose sqlc constant is `listing` (a query named `Listing` or `listing`), which hides the package `listing` | `models: query "Listing" gives sqlc the constant "listing", which views imports as a package` |

## Additions to §7 Limits

| sqlc option or case | Effect |
|---------------------|--------|
| `sql_package` other than `pgx/v5` | `zz_<package>_views.go` doesn't compile: sqlc's `DBTX` has no `Query` |
| `rename`, a non-default `initialisms` or `emit_exact_table_names` that change a `_view` relation's type or a field's name | `zz_<package>_views.go` doesn't compile: `undefined: <Type>` or `has no field` |
| `omit_unused_structs: true` and no query uses the view | `undefined: <Type>` |
| A column `overrides` entry | The field has the override's type; it scans if that type can scan the column |
| `emit_methods_with_db_argument: true` | `Queries` has no `db`: `zz_<package>_views.go` doesn't compile |
| A table whose sqlc type is `Listings` (sqlc's singular of its name, which models can't compute) | `Listings redeclared` |

## §10 for this amendment

A new generator and package: an amendment, minor version (`v0.2.0`).
