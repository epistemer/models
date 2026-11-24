# models — Developer Interface (v1, locked)

> Status: **locked by epistemer on 2026-11-14** (roadmap M1.5). Any change
> needs an amendment, recorded in §11 (A1–A5). This document covers what
> an app writes and what it gets back; why it is so is in
> [ADR 0001](adr/0001-models-on-sqlc.md), [ADR 0002](adr/0002-spike-s1-zp-labels.md)
> and, for `views`, [ADR 0005](adr/0005-view-listings.md),
> [ADR 0007](adr/0007-column-sets-for-views.md) and
> [ADR 0008](adr/0008-column-sets-read-any-field.md).

## 0. Decisions

| # | Decision | Choice |
|---|----------|--------|
| D1 | What it is | Four trivial generators (genx) over sqlc's codegen request, and three runtime packages: `listing`, for the code `views` generates (A2), `option` (A3) and `column` (A4). No program |
| D2 | App's choices | genx parameters (genx amendment A4): `dir`, `package`, `table_prefix` (§3) |
| D3 | Generators | `labels` writes a `labels/v1` file; `enums` writes `Label()`, `Value()` and `All<Type>Options()` for every enum type (A3); `views` writes a `List<Type>` method and a column set for every relation whose name ends in `_view` (§4); `pgx` writes `RegisterTypes` (A5) |
| D4 | Keys | zp's: `enum.<enum>.<value>`, `table.<table>`, `table.<table>.<T>_<C>` (§5.2) |
| D5 | Paths | The names in SQL: `enum.<enum>.<value>`, `table.<table>`, `table.<table>.<column>` (§5.2) |
| D6 | Variables | `L10n<Const>` for an enum value (`<Const>` is sqlc's constant), `L10n<T>` and `L10n<T>_<C>` for tables and columns (§5.2) |
| D7 | English text | The SQL comment, else the name made readable; an enum value made readable (A5, §5.3) |
| D8 | Scope | The request's default schema, one request per app, sqlc's default naming options (§7) |
| D9 | Versions | genx `v0.2.0` or later, i18n `v0.2.0` or later, sqlc requests of `v1.29.0`'s shape (`sqlc.request/v1`) |
| D10 | Module path | `github.com/epistemer/models`, package `models`, with only the metafile, its templates and a `doc.go`; packages `listing` (§4.4) and `option` (§4.5) |

## 1. Principles

1. **Keys are permanent.** Every key is a function of names in SQL and the
   `table_prefix` parameter, and nothing else. A key that would mean two
   things stops the run; it is never written.
2. **Data, not Go, carries the labels.** The `labels` generator writes a
   `labels/v1` file; i18n's `labels` generator turns it into Go. models
   never registers an i18n key.
3. **sqlc stays the source of the models.** models adds methods to sqlc's
   types and declares nothing sqlc declares.
4. **Loud over silent.** What models can't match (§7) stops the run with a
   message naming the source, rather than writing a wrong name.

## 2. What an app writes

An app that already generates its models with sqlc adds four things.

**sqlc's request, next to its Go.** A `gen: json:` block beside `gen: go:`
in `sqlc.yaml`, writing into the same directory:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "schema/"
    queries: "queries/"
    gen:
      go:
        package: "models"
        out: "models"
        sql_package: "pgx/v5"
        emit_all_enum_values: true   # optional: lists of an enum's values
      json:
        out: "models"
        filename: "codegen_request.json"
```

**The manifest, `genx.yaml`.** sqlc as a producer that emits the request as
a foreign file, models switched on with its parameters, and i18n's
generator after it:

```yaml
version: 1
producers:
  - name: sqlc
    run:   [go, tool, sqlc, generate, -f, sqlc.yaml]
    check: [go, tool, sqlc, diff, -f, sqlc.yaml]
    emits:
      - path: models/codegen_request.json
        domain: sqlc.request/v1
use:
  - github.com/epistemer/models:
      generators: [labels, enums, views, pgx]
      params:
        table_prefix: app_
  - github.com/epistemer/i18n: [labels]
```

The list form, `- github.com/epistemer/models: [labels, enums]`, uses every
default. `enums` needs `labels` and i18n's `labels`: its methods return the
variables those two declare. `views` needs them too, for its column sets'
labels (A4), and `enums` when a view has an enum column, for its options;
it needs sqlc's `sql_package: "pgx/v5"` (§7). `pgx` needs pgx/v5 too,
v5.6.0 or later, and nothing else of models'. `labels` works alone.

**A blank import**, so the module stays required and genx can find the
metafile (genx ADR 0003). The file goes in any package of the app and
carries that package's name; `main` here:

```go
//go:build genx

package main

import (
	_ "github.com/epistemer/i18n"
	_ "github.com/epistemer/models"
)
```

**Committed files.** `models/codegen_request.json`, the generated files
(§4) and `genx.lock` are committed, as genx expects. `go tool genx`
regenerates everything after a change to the schema, the queries or either
YAML file; `go tool genx check` in CI runs `sqlc diff` and compares every
generated file.

## 3. Parameters

Declared in models' metafile and set on its `use` entry (genx A4).

| Parameter | Default | Meaning |
|-----------|---------|---------|
| `dir` | `"models"` | The Go package's directory, relative to the module root: where sqlc writes its models and models writes its files. A string, slash-separated and clean: not empty, not absolute, no `.` or `..` element, no empty element, no trailing `/`. genx's rules for an output path (its §7) apply to it as well |
| `package` | `null` | The Go package's name: a string matching `^[A-Za-z_][A-Za-z0-9_]*$`, not `_` and not a Go keyword. `null` means the last element of `dir`, which must then meet the same rule. Set it when the directory's name doesn't (`my-models`), or the package's name differs from it |
| `table_prefix` | `""` | Removed from the front of table and column names, for variables, English text and column keys (§5.2). `""` removes nothing |

- `dir` and `package` must agree with sqlc's `out` and `package`; models
  can't see sqlc's Go options (§7).
- Both are checked on every run, with or without a request.
- **`table_prefix` is part of every column key.** Changing it after
  translations exist orphans them, as renaming a column would.
- **`table_prefix` reaches models' names only, not sqlc's.** sqlc keeps
  the prefix in every Go name: with tables `app_order` and
  `app_order_view` and `table_prefix: app_`, the labels' variables read
  `L10nOrder_Total`, but sqlc's types are `AppOrder` and `AppOrderView`, so
  the listing is `ListAppOrderView` and the column set
  `AppOrderViewColumns` (§4.3). One relation then has two names in Go.
  sqlc's own way out, a `rename` per relation, isn't followed by `enums`
  or `views` (§7). So the trade-off is: a prefix keeps keys and labels
  short while Go names keep it, or tables without a prefix give one name
  everywhere (a table named after a reserved word, such as `user`, then
  needs another name or quoting in every query).

## 4. Generators

| Generator | Reads | Writes |
|-----------|-------|--------|
| `labels` | the one `sqlc.request/v1` file | `<dir>/zz_<package>.labels.json`, a `labels/v1` file |
| `enums` | the one `sqlc.request/v1` file | `<dir>/zz_<package>_enums.go` |
| `views` | the one `sqlc.request/v1` file | `<dir>/zz_<package>_views.go` |
| `pgx` | the one `sqlc.request/v1` file | `<dir>/zz_<package>_pgx.go` (A5) |

With no request, neither writes anything. When the default schema gives no
label at all (no table or view, and no enum value), `labels` writes
nothing; with no enum type, `enums` and `pgx` write nothing; with no relation whose name
ends in `_view`, `views` writes nothing. Nothing empty is ever
written, because i18n would turn an empty labels file into Go with an
unused import. An enum with no values still gets its `Label()`, which
always returns nil.

All four make every check of §6 that isn't marked as `views`' or `pgx`'
own, so `enums`, `views` or `pgx` run without `labels` refuses what
`labels` would.

i18n's `labels` then writes `<dir>/zz_<package>_i18n.go`, which declares
every label's variable and `L10n(path)`.

### 4.1 The labels file

```json
{
  "domain": "labels/v1",
  "labels": [
    {"key": "enum.order_status.paid", "path": "enum.order_status.paid",
     "text": "paid", "var": "L10nOrderStatusPaid"},
    {"key": "table.app_order_line", "path": "table.app_order_line",
     "text": "Order Line", "var": "L10nOrderLine"},
    {"key": "table.app_order_line.OrderLine_Title", "path": "table.app_order_line.title",
     "text": "Title", "var": "L10nOrderLine_Title"}
  ],
  "package": {"dir": "<dir>", "name": "<package>"},
  "producer": "models"
}
```

Every label has all four fields. The order: each enum of the default schema
in the request's order with its values in theirs, then each table and view
in the request's order, each followed by its columns in the request's
order. genx writes object keys sorted, as above, and indents its own way.

### 4.2 The enums file

```go
// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
	"fmt"

	"github.com/epistemer/i18n"
	"github.com/epistemer/models/option"
)

// Label returns the label of t, or nil when t isn't a value of OrderStatus.
func (t OrderStatus) Label() i18n.TKer {
	switch t {
	case OrderStatusPending:
		return L10nOrderStatusPending
	case OrderStatusPaid:
		return L10nOrderStatusPaid
	case OrderStatusShipped:
		return L10nOrderStatusShipped
	}
	return nil
}

// Value returns t as a string, for option.SelectOption.
func (t OrderStatus) Value() string { return string(t) }

// MarshalText returns t's value, or an error when t isn't a value of
// OrderStatus.
func (t OrderStatus) MarshalText() ([]byte, error) {
	if t.Label() == nil {
		return nil, fmt.Errorf("models: %q is not a value of OrderStatus", string(t))
	}
	return []byte(t), nil
}

// UnmarshalText sets t to text, or returns an error when text isn't a
// value of OrderStatus.
func (t *OrderStatus) UnmarshalText(text []byte) error {
	if OrderStatus(text).Label() == nil {
		return fmt.Errorf("models: %q is not a value of OrderStatus", string(text))
	}
	*t = OrderStatus(text)
	return nil
}

// AllOrderStatusOptions returns every value of OrderStatus, in the
// schema's order, as select options.
func AllOrderStatusOptions() []option.SelectOption {
	return []option.SelectOption{
		OrderStatusPending,
		OrderStatusPaid,
		OrderStatusShipped,
	}
}
```

- For each enum type, in the order of §4.1, on the type sqlc declares and
  with sqlc's constants: `Label()`, `Value()` (A3),
  `MarshalText()`, `UnmarshalText()` (A5) and `All<Type>Options()` (A3).
  Nothing else is declared. Every enum value is an `option.SelectOption`
  (§4.5).
- The text methods (A5) take exactly the enum's values in SQL, `""` only
  when the enum has it, and refuse any other text, case included, with
  `models: "bogus" is not a value of OrderStatus` and no sentinel.
  handlers' `ParamText` turns that into `ErrBadParam`, a 400, so a bad path
  parameter never reaches PostgreSQL. handlers' URL builders panic on a
  value that isn't the enum's (`handlers.Text`), a programming error.
  `encoding/json` uses the methods too, so marshalling a sqlc row whose
  nullable enum is null (`NullOrderStatus{}`, value `""`) fails unless
  `""` is a value. pgx doesn't use them.
- `All<Type>Options()` returns a new slice on every call; an enum with no
  values returns an empty one.
- `Null<Type>` gets no method: call `n.<Type>.Label()` when `n.Valid`.
- The value's string is `string(t)`, also `t.Value()`; the list of values
  is sqlc's `All<Type>Values()` (`emit_all_enum_values`).
- `i18n.TKer` is what i18n's `Key` returns; it renders with
  `.T(ctx)` and names its key with `.K()`.

### 4.3 The views file

Every table or view of the default schema whose SQL name ends in `_view`,
in the request's order, gets a method on sqlc's `Queries`. `<Type>` is
sqlc's type for the relation, `structName` of its name (sqlc makes names
singular, which leaves a name ending in `_view` alone except for its one
rule that ignores the end: a leading `oxen` becomes `ox`, so `oxen_view`
is `OxView`; a name starting with a digit gets sqlc's `_`, as in
`_1View`); each field is `structName` of a column's name. For a view `order_summary_view`:

```go
// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
	"context"

	"github.com/epistemer/models/listing"
)

// Listings is the listing methods of Queries, one per view.
type Listings interface {
	ListOrderSummaryView(ctx context.Context, src listing.Source) ([]OrderSummaryView, int64, error)
}

var _ Listings = (*Queries)(nil)

// ListOrderSummaryView returns the rows of order_summary_view
// that src selects, and how many rows its search and filters match.
func (q *Queries) ListOrderSummaryView(ctx context.Context, src listing.Source) ([]OrderSummaryView, int64, error) {
	lq := src.ListingQuery()
	rowsSQL, countSQL, args, err := listing.Build(lq, "order_summary_view",
		[]string{"id", "status", "total", "placed", "note"})
	// … runs rowsSQL through q.db (pgx/v5's DBTX), scanning each row into
	// &i.ID, &i.Status, &i.Total, &i.Placed, &i.Note; then countSQL, when
	// the rows can't tell the total.
}
```

- `src` is a `listing.Source` (A3): a `listing.Query` value, or a table
  package's own type. It is read once, and must not be nil (a nil `src`
  panics). A filter on an enum column takes the enum value, or an
  `option.SelectOption`'s `Value()`, not the option itself.
- An error from `listing.Build` returns before any query runs.
- The rows tell the total, and the count isn't run, when they end the
  listing: `Limit` is 0 or more than the rows returned, and either
  `Offset` is 0 or at least one row came back. The total is then
  `Offset` plus the number of rows. Otherwise the count runs; when rows
  came back, the total is never less than `Offset` plus their number. The
  two statements share no transaction.
- An app embeds `Listings` in its own interface:
  `type Q interface { Querier; Listings }`.

**Column sets** (A4). Beside each `List<Type>`, the views file declares
`<Type>Columns`, a struct with one field per column a kind takes (§4.6),
named as the row's field and in the view's order:

```go
// OrderSummaryViewColumns are order_summary_view's columns, for a table package (A4).
// Left out, of a type no kind takes: meta (jsonb).
var OrderSummaryViewColumns = struct {
	ID     column.Column[OrderSummaryView]
	Status column.Column[OrderSummaryView]
	Note   column.Column[OrderSummaryView]
}{
	ID: column.Column[OrderSummaryView]{
		Name: "id", Label: L10nOrderSummaryView_Id, Kind: column.Integer,
		Get:  column.Field(column.Integer, func(r OrderSummaryView) any { return r.ID }),
	},
	Status: column.Column[OrderSummaryView]{
		Name: "status", Label: L10nOrderSummaryView_Status, Kind: column.Enum,
		Options: AllOrderStatusOptions(),
		Get:     func(r OrderSummaryView) column.Value { return column.Value{Enum: r.Status} },
	},
	Note: column.Column[OrderSummaryView]{
		Name: "note", Label: L10nOrderSummaryView_Note, Kind: column.Text,
		Get:  column.Field(column.Text, func(r OrderSummaryView) any { return r.Note }),
	},
}
```

- The kind comes from the column's SQL type:

  | SQL type | Kind | The value |
  |---|---|---|
  | `text`, `varchar`, `char`, `citext`, `name` | Text | `Text` |
  | `uuid` | Text | `Text`, canonical |
  | `int2`, `int4`, `int8` (and their aliases and serials) | Integer | `Int` |
  | `float4`, `float8` | Number | `Num` |
  | `numeric` | Number | `Num` parsed from its exact digits (NaN, ±Inf), and the digits in `Text` |
  | `bool` | Bool | `Bool` |
  | `date` | Date | `Time`, midnight UTC; an infinite one has zero `Time` and `Text` "infinity" or "-infinity" |
  | `timestamp`, `timestamptz` | Time | `Time`; infinite, as for `date` |
  | an enum of the default schema | Enum | `Enum`, the value; `Options` is `All<Type>Options()` |

- Every getter but an enum's is `column.Field` (§4.6, A5), which reads
  whatever Go type sqlc gives the field: pgtype's with sqlc's defaults,
  `*T` with `emit_pointers_for_null_types: true`, or an override's type
  (`decimal.Decimal`, `time.Time`, `uuid.UUID`, an app's
  `type ID int64`), each with a pointer when nullable. One views file
  serves every configuration, since sqlc's request carries none of its Go
  options. An enum column's getter reads `r.<Field>` or, when nullable,
  sqlc's `Null<Type>` (sqlc keeps it with pointers too).
- A null is `Value{Null: true}`: a nil pointer, pgtype's `Valid` false,
  or sqlc's `Null<Type>`.
- An array, or any other type (json, jsonb, bytea, interval, …), is left
  out, and the set's comment names it, with or without an override. It
  stays in the row and the listing.
- `Label` is the column's variable from the labels file (§5), so the set
  needs `labels` and i18n's `labels`; an enum column's `Options` needs
  `enums`.

### 4.4 Package listing

`github.com/epistemer/models/listing`, no dependencies outside the
standard library:

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

// Source gives the query of a listing (A3). Query implements it; so can a
// table package's own type.
type Source interface {
	ListingQuery() Query
}

// ListingQuery returns q, so a Query is its own Source.
func (q Query) ListingQuery() Query

// Labeler gives the text for the id a link column is filtered by, such as
// a customer's name for its id, and false when there is no such row. A
// table calls it once per link filter of a request, to show the chosen
// value; rows carry their own text through the view's joins (A3).
// Apps implement it; table packages call it. models generates none.
type Labeler interface {
	LabelOf(ctx context.Context, id string) (string, bool, error)
}
```

- Column names are SQL names and must be among `columns`, in `SearchIn`
  too when `Search` is `""`; a name that isn't is an error naming it. So is
  a `Search` with no `SearchIn` or holding a NUL, an `Op` not listed
  above, a `nil` `Value` (untyped, or a nil pointer, slice, map, channel,
  function or interface; kept free for a later `IS NULL`), a negative
  `Limit` or `Offset`, an empty `relation` or `columns`, a column listed
  twice in `columns`, and a name holding a NUL. A non-nil value that the
  driver encodes as NULL (`pgtype.Text{}`) matches nothing.
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

### 4.5 Package option

`github.com/epistemer/models/option` (A3), which imports only i18n:

```go
// SelectOption is one choice of a select input or a table's filter: the
// value sent and stored, and its label for people. Every enum value
// implements it (§4.2).
type SelectOption interface {
	Value() string
	Label() i18n.TKer
}


// Static returns an option whose label is text from data, such as an
// event's name: its Label renders text in every locale, and its K is ""
// (A5).
func Static(value, text string) SelectOption
```

models owns this interface, and the form and table packages consume it,
so models depends on neither (ADR 0006). A data label has no key, so a
`Static` option's `Label().K()` is `""`.

### 4.6 Package column

`github.com/epistemer/models/column` (A4), which imports i18n and
`option` only, as `option` does:

```go
type Kind uint8 // Text, Integer, Number, Bool, Date, Time, Enum
func (k Kind) String() string

// Value is a cell's plain value: the field of its kind holds it, unless Null.
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
	Label   i18n.TKer
	Kind    Kind
	Options []option.SelectOption // an enum's, in schema order
	Get     func(R) Value
}

func UUID(b [16]byte) string     // canonical text, for generated getters
func Numeric(digits string) Value // a numeric from its exact digits, for generated getters
func Infinite(sign int8) Value    // an infinite date or time, for generated getters

// Field returns a getter that reads a row's field as kind k, for
// generated column sets (A5). It calls field once, on R's zero value, to
// learn the field's type, and panics there if k can't read that type.
func Field[R any](k Kind, field func(R) any) func(R) Value
```

What `Field` reads, for the field's type (A5, [ADR 0008](adr/0008-column-sets-read-any-field.md)):

| The field's type | Read as |
|---|---|
| a pointer | nil is `Null`; otherwise what it points to, by these rules |
| a `driver.Valuer`: every pgtype type, `decimal.Decimal`, `uuid.UUID`, an app's own | its `Value()`: nil is `Null`; otherwise that value, by these rules |
| `time.Time` | itself |
| a string, `[]byte`, bool, integer or float type, named or not (`type EventID int64`) | its value |
| `[16]byte`, named or not | a UUID, in canonical text |

What each kind takes:

| Kind | Takes |
|---|---|
| Text | a string, `[]byte`, a UUID |
| Integer | an integer |
| Number | a float or an integer; a string as a numeric's exact digits, as `Numeric` (`decimal.Decimal` and `pgtype.Numeric` both give one) |
| Bool | a bool |
| Date, Time | a `time.Time`; `"infinity"` or `"-infinity"` as `Infinite` |
| Enum | nothing: generated enum getters read the field themselves |

- A type `Field` can't read, or a kind that can't take it, panics when the
  set is declared, which is when the models package is initialized: at
  the app's start, or in any test that imports it. The message names the
  kind and the type, `column: number can't read ids.Money`, and the stack
  names the set's line, which names the column.
- A `driver.Valuer` whose value its kind can't take, or whose `Value`
  fails, panics on that row, naming both types. pgtype's types,
  `decimal.Decimal` and `uuid.UUID` never do.

A column holds what the schema knows, and nothing a page decides: whether
it is searched, sorted or filtered is the table package's (ADR 0007, ui's
S22). A column can be built by hand, for a test or for rows that don't come
from sqlc.

### 4.7 The pgx file

`pgx` (A5) writes `<dir>/zz_<package>_pgx.go` when the default schema has
an enum:

```go
// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// RegisterTypes registers every enum of the schema, and its array, with
// conn's type map, so pgx can encode and scan []OrderStatus. For
// pgxpool.Config.AfterConnect.
func RegisterTypes(ctx context.Context, conn *pgx.Conn) error {
	types, err := conn.LoadTypes(ctx, []string{
		"public.order_status", "public._order_status",
		"public.payment_method", "public._payment_method",
	})
	if err != nil {
		return err
	}
	conn.TypeMap().RegisterTypes(types)
	return nil
}
```

- Names are qualified with the request's default schema, in §4.1's order,
  each enum followed by its array, which PostgreSQL names `_` plus the
  enum's name. pgx matches names as written, so none is quoted.
- An app wires it as `cfg.AfterConnect = models.RegisterTypes` on its
  `pgxpool.Config`. It needs pgx v5.6.0 or later, which added `LoadTypes`.
- `pgx` keeps `enums` free of the driver, so `labels` and `enums` still
  work on any `sql_package`.

## 5. Labels

### 5.1 Naming functions

| Function | Does | Example |
|----------|------|---------|
| `strip(s)` | Removes `table_prefix` from the front of `s`, if `s` starts with it | `app_order_line` → `order_line` |
| `zpPascal(s)` | Splits `s` on `_`, drops empty parts, upper-cases the first character of each part if it is an ASCII letter, and joins them | `order_line` → `OrderLine`, `user_id` → `UserId`, `2fa` → `2fa` |
| `readable(s)` | The same split, joined with spaces | `order_line` → `Order Line` |
| `readableValue(v)` | For an enum value (A5): split at `.`, `_`, `-` and where lower case turns to upper, as handlers' §13.3; each word's first ASCII letter upper-cased; joined with spaces; `v` itself when no word is left | `in_review` → `In Review`, `userID` → `User ID`, `!!` → `!!` |
| `replace(v)` | sqlc's `EnumReplace`: `-`, `/`, `:` and `_` become `_`; other ASCII letters and digits stay; everything else goes | `cash:on:delivery` → `cash_on_delivery`, `café` → `caf` |
| `structName(s)` | sqlc's `StructName` with `initialisms: [id]`, for an ASCII `s`: every character that isn't an ASCII letter or digit becomes `_`; split on `_`; a part equal to `id` becomes `ID`, any other gets its first character upper-cased if it is a letter; join; prefix `_` if the result starts with a digit | `id_kind` → `IDKind`, `2fa_method` → `_2faMethod` |

### 5.2 Keys, paths and variables

For an enum `e` of the default schema, its Go type is `structName(e)`. Its
values are numbered as sqlc v1.29.0 does: with `seen` an empty set, for the
value `v` at index `i` (from 0), `r` is `replace(v)`; if `r` is empty or in
`seen`, `r` becomes `value_<i>`, without checking `seen` again; `r` is added
to `seen`. The constant is `structName(e + "_" + r)`.

sqlc doesn't check the numbered name: for the values `value_1` and `!!` it
declares `EValue1` twice, and its own Go doesn't compile. models refuses
the duplicate variable (§6).

| Label | Key | Path | Variable |
|-------|-----|------|----------|
| Enum `e`, value `v` | `enum.<e>.<v>` | `enum.<e>.<v>` | `L10n` + the constant |
| Table or view `t` | `table.<t>` | `table.<t>` | `L10n<T>` |
| Its column `c` | `table.<t>.<T>_<C>` | `table.<t>.<c>` | `L10n<T>_<C>` |

where `<T>` is `zpPascal(strip(t))` and `<C>` is `zpPascal(strip(c))`.
`<e>`, `<v>`, `<t>` and `<c>` are the names in SQL, unchanged. Examples,
with `table_prefix: app_`:

| SQL | Key | Path | Variable |
|-----|-----|------|----------|
| `payment_method` value `pay_pal` at index 4, after `pay-pal` at index 3 | `enum.payment_method.pay_pal` | same | `L10nPaymentMethodValue4` |
| `app_order.customer_id` | `table.app_order.Order_CustomerId` | `table.app_order.customer_id` | `L10nOrder_CustomerId` |
| `app_order.app_ref` | `table.app_order.Order_Ref` | `table.app_order.app_ref` | `L10nOrder_Ref` |

`zpPascal` isn't sqlc's naming: `customer_id` gives `CustomerId` where
sqlc's field is `CustomerID`. It is what zp's keys were made with. The
fixture ([testdata](../testdata/README.md)) covers every rule.

### 5.3 English text

| Label | Text |
|-------|------|
| Enum value | `readableValue(v)` (A5): `v` split at `.`, `_` and `-`, and where lower case turns to upper (a run of capitals stays one word: a capital followed by a lower-case letter starts a new word), each word's first letter upper-cased if it is an ASCII letter, joined with spaces; the value itself when that leaves no word, so an empty value keeps an empty text |
| Table, view, column | Its SQL comment with white space trimmed from both ends; if that leaves nothing, `readable(strip(name))`. White space is `\t`, `\n`, `\v`, `\f`, `\r`, U+0085 and the Unicode categories `Zs`, `Zl` and `Zp`, as the regular expression `[\t\n\v\f\r\x{85}\p{Z}]` |

Examples of an enum value's text: `in_review` → `In Review`,
`credit-card` → `Credit Card`, `bank/transfer` → `Bank/transfer`,
`cash:on:delivery` → `Cash:on:delivery`, `userID` → `User ID`,
`u2f` → `U2f`, `level_3` → `Level 3`, `café` → `Café`, `!!` → `!!`,
`''` → `''`. It is handlers' §13.3 rule, but for the first letter of a
word outside ASCII, which stays as it is.

The text is the default in the source language. People change it, and
translate it, in the running app (i18n); models never writes translations.

## 6. Refusals

Each stops the run, through genx, with a message that starts `models:` and
names the source; genx adds the generator and the input file. The cases are
the contract; their wording may get clearer without an amendment.

| Case | Example message |
|------|-----------------|
| `dir` isn't a string meeting §3's rule | `models: parameter dir "models/": not a clean relative slash-separated path` |
| `package` isn't a string meeting §3's rule, or is `null` with a `dir` whose last element doesn't | `models: parameter package: "my-models" is not a Go package name; set package` |
| More than one `sqlc.request/v1` file | `models: 2 sqlc requests; one is supported` |
| The request names a default schema it doesn't hold | `models: the default schema "public" isn't in the request` |
| A table, view or enum in a schema other than the default one, `pg_catalog`, `information_schema` and `pg_temp` | `models: schema "audit" has tables or enums; only the default schema "public" is supported` |
| An enum name with a character outside ASCII, or whose `structName` is empty | `models: enum "état": sqlc's Go name for it can't be derived` |
| A table, view or column whose `<T>` or `<C>` (§5.2) is empty or has a character that isn't a letter or a decimal digit (`[\p{L}\p{Nd}]`), so `L10n<T>` or `L10n<T>_<C>` isn't a Go identifier | `models: column "app_order"."first name" gives "First name", which isn't a Go identifier` |
| Two labels with one key, one path or one variable | `models: tables "app_user" and "user" both give the variable "L10nUser"` |
| Two enums with one Go type, `structName` of their names (A1) | `models: enum "user1" and enum "user_1" both give sqlc the Go type "User1"` |
| A table, view, enum or query of the default schema whose Go name in sqlc starts with `L10n`, which could clash with a variable: `structName` of its SQL name (sqlc may make a table's name singular, which doesn't change its start) | `models: table "l10n_note" gives sqlc the Go name "L10nNote"; names starting with "L10n" are reserved for labels` |

Made by `views` only (A2), and only when the request has a relation whose name
ends in `_view`:

| Case | Example message |
|------|-----------------|
| The request's engine isn't `postgresql` | `models: views need PostgreSQL; the request's engine is "mysql"` |
| A `_view` relation or one of its columns has a name with a character outside printable ASCII (a control character too), or a column whose `structName` is empty | `models: view "café_view": sqlc's Go name for it can't be derived` |
| A `_view` relation with no columns | `models: view "empty_view" has no columns` |
| Two columns of a `_view` relation with one `structName` (sqlc declares the field twice) | `models: view "dup_view": columns "id" and "ID" both give the field "ID"` |
| Two `_view` relations with one `<Type>` | `models: views "id_view" and "ID_view" both give sqlc the Go type "IDView"` |
| A query whose name, as written, is a method `List<Type>` | `models: query "ListOrderSummaryView" clashes with the listing method of view "order_summary_view"` |
| An enum whose `structName` is `Listings`, the type sqlc would then declare (a query named `Listings` gives a method and the constant `listings`, which don't clash) | `models: enum "listings" gives sqlc the Go name "Listings", which views declares` |
| A query whose sqlc constant (its name with the first letter lower-cased) hides a name the views file uses: `listing`, `column` (A4), `context`, `append`, `len`, `nil`, `int64`, `string` or `error` | `models: query "Listing" gives sqlc the constant "listing", which hides a name the views file uses` |
| A view's column set `<Type>Columns` (A4) is already a query's, a relation's or an enum's Go name | `models: view "order_view": its column set "OrderViewColumns" is already a query's, a relation's or an enum's Go name` |

Made by `enums` only (A5), and only when the default schema has an enum:

| Case | Example message |
|------|-----------------|
| A query whose sqlc constant hides a name the enums file uses: `fmt`, `i18n`, `option`, `string`, `error`, `nil` or `byte` | `models: query "Fmt" gives sqlc the constant "fmt", which hides a name the enums file uses` |

Made by `pgx` only (A5), and only when the default schema has an enum:

| Case | Example message |
|------|-----------------|
| An enum whose Go type is `RegisterTypes` | `models: enum "register_types" gives sqlc the Go name "RegisterTypes", which pgx declares` |
| A query whose sqlc constant hides a name the pgx file uses: `context`, `pgx`, `string`, `error` or `nil` | `models: query "Pgx" gives sqlc the constant "pgx", which hides a name the pgx file uses` |

Made by every generator, like the shared cases above it (A3):

| Case | Example message |
|------|-----------------|
| An enum whose `All<Type>Options` is another enum's Go type or one of sqlc's enum constants | `models: enum "x" gives the function "AllXOptions", which is enum "all_x_options"'s Go type` |

Where zp's casing also collides (`user_1` and `user1`), the shared refusal
of two labels with one key or variable stops the run first.

genx itself refuses an unknown parameter, a value of the wrong JSON type
(genx A4), and an output path its §7 rejects.

## 7. Limits

What models doesn't see, so doesn't follow. Each ends in a compile error,
a refusal, or a method nothing calls; never in a wrong key.

| sqlc option or case | Effect |
|---------------------|--------|
| `rename` or a non-default `initialisms` that changes an enum type's name | `zz_<package>_enums.go` doesn't compile: `undefined: <Type>` |
| The same, changing a constant's name | `zz_<package>_enums.go` doesn't compile: `undefined: <Const>`. Keys are unaffected |
| `omit_unused_structs: true`, when an enum isn't used by any query | sqlc leaves the type out: `undefined: <Type>` |
| An `overrides` entry for an enum's database type | sqlc still declares the enum's type, so `Label()` compiles, but the fields and parameters it overrides have the other type, which has no `Label()` |
| `out` or `package` in `gen: go:` that differ from `dir` and `package` | models' files land in another directory or package: `zz_<package>_enums.go` doesn't compile there, and labels sit apart from the models |
| A second `sql:` block whose request is also declared in `genx.yaml` | Refused (§6). One whose request isn't declared is never seen |
| Tables or enums in another schema | Refused (§6) |
| `sql_package` other than `pgx/v5` | `zz_<package>_views.go` doesn't compile: sqlc's `DBTX` has no `Query` |
| `rename`, a non-default `initialisms` or `emit_exact_table_names` that change a `_view` relation's type or a field's name | `zz_<package>_views.go` doesn't compile: `undefined: <Type>` or `has no field` |
| `omit_unused_structs: true` and no query uses the view | `undefined: <Type>` |
| A column `overrides` entry | The field has the override's type; it scans if that type can scan the column |
| A method `Value` on an enum type, or a package-level `All<Type>Options`, of the app's own (A3); a table whose sqlc type, made singular, is `All<Type>Options` | `Value redeclared` or `All<Type>Options redeclared`. An enum type can't also be a `driver.Valuer`, whose `Value` has another signature |
| `emit_methods_with_db_argument: true` | `Queries` has no `db`: `zz_<package>_views.go` doesn't compile |
| A table whose sqlc type is `Listings` (sqlc's singular of its name, which models can't compute) | `Listings redeclared` |
| An `overrides` entry on a `_view` relation's enum column, or on an enum's type (A5) | the field isn't `<Type>` or `Null<Type>`: the enum column's getter in `zz_<package>_views.go` doesn't compile |
| An override over a kind-taking SQL type whose Go type `Field` can't read (a struct that isn't a `driver.Valuer`) (A5) | `column.Field` panics when the models package is initialized, naming the kind and the type |
| A method `MarshalText` or `UnmarshalText` on an enum type, of the app's own (A5) | `MarshalText redeclared` |
| An enum whose array PostgreSQL didn't name `_<enum>` (a name of 63 bytes, or a clash) (A5) | `RegisterTypes` fails on connect, from pgx's `LoadTypes` |
| A function `RegisterTypes` of the app's own in the package, or a table whose sqlc type, made singular, is `RegisterTypes` (A5) | `RegisterTypes redeclared` |
| `pgx` with pgx older than v5.6.0 (A5) | `conn.LoadTypes undefined` |
| `views` without `labels` and i18n's `labels`, or without `enums` for a view with an enum column (A4) | `undefined: L10n…` or `undefined: All<Type>Options` |

## 8. Testing an app

- `genxtest.AssertGenerated(t)` fails when any generated file, the request
  included (through `sqlc diff`), is stale.
- `i18ntest` sets a translation on a generated key and checks the rendered
  text, through `Label()` or a variable or `L10n(path)`, like any other key.

## 9. Not in v1

Parameters that mirror sqlc's `rename` and `initialisms`; several requests
or schemas; labels for an enum type itself. The roadmap's backlog orders
them.

## 10. What a change to this document means

| Change | Means |
|--------|-------|
| A key's rule (§5.2), the prefix's role, the English text's rule | A new major version: keys would move under existing translations |
| A variable's, a path's or a file's name | A new major version: apps' code refers to them |
| A new refusal | An amendment; it can break an app that ran before |
| A new parameter with a default that keeps today's output | An amendment, minor version |
| A message's wording | Nothing: §6 lists cases, not wording |

## 11. Amendments

### A1. Two enums with one Go type (2026-11-14, before v0.1.0)

The release review found that two enums whose names give one
`structName`, such as `user1` and `user_1`, pass every check of §6 when
their values differ, and `enums` then writes two `Label()` methods on one
type. sqlc declares that type twice too, so the app never compiled. §6 now
refuses it, naming both enums. A clarification under Principle 4, made by
Claude as owner before any release; it can be overturned.

`dir`'s element rules (§3, genx §7) are now also checked by models itself,
so a bad `dir` is refused on runs that write nothing, as §3 says.

### A2. Listings of views (2026-11-15, for v0.2.0)

Approved by epistemer on 2026-11-15; the design is
[ADR 0005](adr/0005-view-listings.md). A third generator, `views` (§4.3),
and package `listing` (§4.4), with their refusals (§6) and limits (§7).
It changes nothing an app on v0.1.0 writes or gets: `views` runs only when
switched on. A new generator and package: a minor version, `v0.2.0`.

Clarified while building it (M5.4, review 0010), by Claude as owner: a
query named `Listings` isn't refused, since sqlc declares no type for it;
the floor on the total (§4.3) holds only when rows came back, so a page
past the end reports the count; `oxen` names follow sqlc's singular rule;
a control character in a name is refused, like one outside ASCII; and the
refusal of a query constant covers every name the file uses, not only
`listing`. Each refuses only apps whose generated code couldn't compile.

### A3. Interfaces for UI packages (2026-11-16, for v0.3.0)

Decided by epistemer on 2026-11-16; the design is
[ADR 0006](adr/0006-interfaces-for-ui-packages.md). What zp's model layer
got from its table and form packages is now an interface models owns:
package `option` with `SelectOption` (§4.5); `Value()` and
`All<Type>Options()` on every enum type (§4.2); `listing.Source`, which
`List<Type>` now takes, and `listing.Labeler` (§4.3, §4.4); with one
refusal (§6) and one limit (§7). `List<Type>`'s parameter type changes, and
the enums file gains declarations: a minor version, `v0.3.0`, under v0's
rules.

### A4. Column sets for views (2026-11-22, for v0.4.0)

Accepted by epistemer on 2026-11-22 with the defaults of M1–M5, from ui's
design of tables (its D7.4):
package `column` (`Kind`, `Value`, `Column[R]`), and a `<Type>Columns`
set the `views` generator writes beside each `List<Type>`. The design is
[ADR 0007](adr/0007-column-sets-for-views.md). A new package and new
declarations: a minor version, `v0.4.0`. §4.3 and §4.6 state it.


### A5. Confab iteration 1 (2026-11-24, for v0.5.0)

> **Locked by epistemer on 2026-11-24, with every default** (A5.1a–c,
> A5.3a–b). From Confab's iteration 1 (findings F5, F28, F15, F8, F19)
> and epistemer's answers of 2026-11-24 on the board item "models: API
> questions from Confab iteration 1". The board's M1–M5 are A5.1–A5.5
> here, since M1–M6 name this repository's milestones. §0, §2, §4, §5.3,
> §6 and §7 state it; the text below is the record of what was locked.

New declarations and a changed English text: a minor version, `v0.5.0`,
under v0's rules. §10 calls a change to the English text's rule a new
major version; that line is for v1, and A5.5 ships as a minor by
epistemer's decision.

#### A5.1 Column sets read any field type sqlc generates (M1, F5)

The design is [ADR 0008](adr/0008-column-sets-read-any-field.md). sqlc's
request carries none of its Go options, so one views file must compile
whether or not the app sets `emit_pointers_for_null_types` or
`overrides`. Every getter but an enum's hands its field to a new function
of package `column`:

```go
// Field returns a getter that reads a row's field as kind k, for
// generated column sets (A5). It calls field once, on R's zero value, to
// learn the field's type, and panics there if k can't read that type.
func Field[R any](k Kind, field func(R) any) func(R) Value
```

```go
Note: column.Column[OrderSummaryView]{
	Name: "note", Label: L10nOrderSummaryView_Note, Kind: column.Text,
	Get:  column.Field(column.Text, func(r OrderSummaryView) any { return r.Note }),
},
```

What `Field` reads, for the field's type:

| The field's type | Read as |
|---|---|
| a pointer | nil is `Null`; otherwise what it points to, by these rules |
| a `driver.Valuer`: every pgtype type, `decimal.Decimal`, `uuid.UUID`, an app's own | its `Value()`: nil is `Null`; otherwise that value, by these rules |
| `time.Time` | itself |
| a string, `[]byte`, bool, integer or float type, named or not (`type EventID int64`) | its value |
| `[16]byte`, named or not | a UUID, in canonical text |

What each kind takes:

| Kind | Takes |
|---|---|
| Text | a string, `[]byte`, a UUID |
| Integer | an integer |
| Number | a float or an integer; a string as a numeric's exact digits, as `column.Numeric` (`decimal.Decimal` and `pgtype.Numeric` both give one) |
| Bool | a bool |
| Date, Time | a `time.Time`; `"infinity"` or `"-infinity"` as `column.Infinite` |
| Enum | nothing: enum getters stay A4's, since sqlc types a nullable enum `Null<Type>` with or without pointers |

- A type `Field` can't read, or a kind that can't take it, panics when the
  set is declared, which is when the models package is initialized: at
  the app's start, or in any test that imports it. The message names the
  kind and the type, `column: number can't read ids.Money`, and the stack
  names the set's line, which names the column.
- A `driver.Valuer` whose value its kind can't take panics on that row,
  naming both types. pgtype's types, `decimal.Decimal` and `uuid.UUID`
  never do.
- Tested in `gentest`, against PostgreSQL, with sqlc's defaults and with
  `emit_pointers_for_null_types: true` plus these overrides, each with
  `pointer: true` for its nullable entry: `numeric` to `decimal.Decimal`;
  `date`, `timestamp` and `timestamptz` to `time.Time`; `uuid` to
  `uuid.UUID`; `interval` to `time.Duration`; and a column override to an
  app's `type … int64`. With pointers but without the `date`, `timestamp`
  and `timestamptz` overrides, sqlc keeps pgtype for them, which `Field`
  reads too. sqlc v1.29.0 matches `db_type` spelled `pg_catalog.numeric`,
  `pg_catalog.interval`, `pg_catalog.timestamp`, `timestamptz`, `date` and
  `uuid`, and silently ignores `pg_catalog.timestamptz` and `timestamp`;
  the docs give the spellings the test uses.
- `interval` takes no kind, with or without an override: the set leaves it
  out and names it, as for `jsonb`.
- §7's row on `emit_pointers_for_null_types` and column overrides goes. A
  new row: an `overrides` entry on an enum column, or on an enum's type,
  gives a field that isn't `<Type>` or `Null<Type>`, so the enum getter
  doesn't compile.
- `column.UUID`, `column.Numeric` and `column.Infinite` stay.

**Sub-questions, each with a default:**

- **A5.1a.** Getters resolved at run time, as above, *default*; or genx
  parameters mirroring sqlc's options (`emit_pointers_for_null_types`,
  the overrides) for compiled getters. The default needs nothing written
  twice, and a parameter couldn't say how to read an app's own type. It
  moves the failure from compile time to the package's initialization.
- **A5.1b.** The name `column.Field`, *default*.
- **A5.1c.** `interval` stays out of every kind, *default*. The
  alternative is a `column.Duration` kind with a `Value.Dur
  time.Duration` field. ui's tables switch on `Kind`, so that would also
  need a ui amendment to render one.

#### A5.2 Text methods on enums (M2, F28)

The enums file declares two more methods per enum type:

```go
// MarshalText returns t's value, or an error when t isn't a value of
// OrderStatus.
func (t OrderStatus) MarshalText() ([]byte, error)

// UnmarshalText sets t to text, or returns an error when text isn't a
// value of OrderStatus.
func (t *OrderStatus) UnmarshalText(text []byte) error
```

- The error reads `models: "bogus" is not a value of OrderStatus`, with no
  sentinel. handlers' `ParamText` turns it into `ErrBadParam` (a 400), so
  `/events/3/status/bogus` never reaches PostgreSQL.
- A value is one of the enum's values in SQL, `""` included only when the
  enum has it. No other text passes, case included.
- handlers' URL builders call `MarshalText` through `handlers.Text`, which
  panics on an error. A builder given a value that isn't the enum's then
  panics, which is a programming error.
- `encoding/json` uses them too, so marshalling a sqlc row whose nullable
  enum is null (`NullOrderStatus{}`, whose value is `""`) fails, unless
  `""` is a value. Apps send their own API types (oapi-codegen), not sqlc
  rows. pgx doesn't use text methods.
- A new §7 row: an app's own `MarshalText` or `UnmarshalText` on an enum
  type gives `redeclared`.

#### A5.3 Registering enums with pgx (M3, F15)

A fourth generator, **`pgx`**, writes `<dir>/zz_<package>_pgx.go` when the
default schema has an enum:

```go
// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// RegisterTypes registers every enum of the schema, and its array, with
// conn's type map, so pgx can encode and scan []OrderStatus. For
// pgxpool.Config.AfterConnect.
func RegisterTypes(ctx context.Context, conn *pgx.Conn) error {
	types, err := conn.LoadTypes(ctx, []string{
		"public.order_status", "public._order_status",
		"public.payment_method", "public._payment_method",
	})
	if err != nil {
		return err
	}
	conn.TypeMap().RegisterTypes(types)
	return nil
}
```

- Names are schema-qualified with the request's default schema, in §4.1's
  order, each enum followed by its array, which PostgreSQL names `_` plus
  the enum's name. pgx matches names as written, so no quoting is needed.
- It needs pgx v5.6.0 or later, which added `LoadTypes`.
- An app wires it as `cfg.AfterConnect = models.RegisterTypes`.
- New refusals, made by `pgx` only: an enum whose Go type is
  `RegisterTypes`; a query whose constant hides a name the file uses
  (`context`, `pgx`, `string`, `error`, `nil`).
- A new §7 row: the array's name differs from `_<enum>` when PostgreSQL
  shortened or renamed it (an enum name of 63 bytes, or a clash). Then
  `RegisterTypes` fails on connect, naming the type.
- **A5.3a.** A new generator, *default*, so that `enums` keeps working
  on any `sql_package`; or `RegisterTypes` in the enums file, which then
  needs pgx/v5.
- **A5.3b.** The generator's name `pgx`, the file `zz_<package>_pgx.go`
  and the function `RegisterTypes(ctx context.Context, conn *pgx.Conn)
  error`, *default*.

#### A5.4 Options whose label is data (M4, F8)

Package `option` gains one function. `SelectOption` stays as it is.

```go
// Static returns an option whose label is text from data, such as an
// event's name: its Label renders text in every locale, and its K is "".
func Static(value, text string) SelectOption
```

- A data label has no key, so `Label().K()` is `""`. No ui package reads
  `K()` of an option's label.

#### A5.5 An enum value's English text is the value made readable (M5, F19)

D7 and §5.3 change: an enum value's text is `readable(v)`, the value
made readable as handlers' interface §13.3 makes route keys readable. The
value is split at `.`, `_` and `-`, and where lower case turns to upper
(a run of capitals stays one word, and a capital followed by a lower-case
letter starts a new one: `HTTPServer` is `HTTP Server`). Each word's first letter is
upper-cased if it is an ASCII letter, and the words are joined with
spaces. When that leaves no word, the text is the value itself, so an
empty value keeps an empty text. Keys, paths and variables don't change.

| Value | Text |
|---|---|
| `in_review` | `In Review` |
| `credit-card` | `Credit Card` |
| `bank/transfer` | `Bank/transfer` |
| `cash:on:delivery` | `Cash:on:delivery` |
| `inReview`, `userID` | `In Review`, `User ID` |
| `u2f`, `level_3` | `U2f`, `Level 3` |
| `café` | `Café` |
| `!!` | `!!` |
| `''` | `''` (empty) |

- The rule for tables, views and columns (`readable` of §5.1, split at
  `_` only) stays.
- i18n shows an unedited value's new text at once. An edited one keeps
  its edit, which records the old text as the one it was based on (i18n
  §8.1).
- It changes neozp's labels convention, which said an enum value is
  itself (decisions.md D7, 2026-11-24).

**Clarified while building it** (M7, [review 0014](reviews/0014-a5.md)), by
Claude as owner: the enums file now imports `fmt`, so `enums` refuses a
query whose constant hides a name that file uses (§6), as `views` and
`pgx` do for theirs. A `driver.Valuer` whose value is of a named basic
type, or a `[16]byte` for Text, is read by the same rules as a field of
that type. A field that is a pointer to an interface is refused when the
set is declared. Each refuses only apps whose generated code couldn't
compile or would panic later.
