# models — Developer Interface (v1, locked)

> Status: **locked by epistemer on 2026-11-14** (roadmap M1.5). Any change
> needs an amendment, recorded in §11 (A1, A2). This document covers what
> an app writes and what it gets back; why it is so is in
> [ADR 0001](adr/0001-models-on-sqlc.md), [ADR 0002](adr/0002-spike-s1-zp-labels.md)
> and, for `views`, [ADR 0005](adr/0005-view-listings.md).

## 0. Decisions

| # | Decision | Choice |
|---|----------|--------|
| D1 | What it is | Three trivial generators (genx) over sqlc's codegen request, and one runtime package, `listing`, for the code `views` generates (A2). No program |
| D2 | App's choices | genx parameters (genx amendment A4): `dir`, `package`, `table_prefix` (§3) |
| D3 | Generators | `labels` writes a `labels/v1` file; `enums` writes `Label()` for every enum type; `views` writes a `List<Type>` method for every relation whose name ends in `_view` (§4) |
| D4 | Keys | zp's: `enum.<enum>.<value>`, `table.<table>`, `table.<table>.<T>_<C>` (§5.2) |
| D5 | Paths | The names in SQL: `enum.<enum>.<value>`, `table.<table>`, `table.<table>.<column>` (§5.2) |
| D6 | Variables | `L10n<Const>` for an enum value (`<Const>` is sqlc's constant), `L10n<T>` and `L10n<T>_<C>` for tables and columns (§5.2) |
| D7 | English text | The SQL comment, else the name made readable; an enum value is itself (§5.3) |
| D8 | Scope | The request's default schema, one request per app, sqlc's default naming options (§7) |
| D9 | Versions | genx `v0.2.0` or later, i18n `v0.2.0` or later, sqlc requests of `v1.29.0`'s shape (`sqlc.request/v1`) |
| D10 | Module path | `github.com/epistemer/models`, package `models`, with only the metafile, its templates and a `doc.go`; package `listing` (§4.4) |

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
      generators: [labels, enums, views]
      params:
        table_prefix: app_
  - github.com/epistemer/i18n: [labels]
```

The list form, `- github.com/epistemer/models: [labels, enums]`, uses every
default. `enums` needs `labels` and i18n's `labels`: its methods return the
variables those two declare. `labels` works alone, and so does `views`,
which needs sqlc's `sql_package: "pgx/v5"` (§7).

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

## 4. Generators

| Generator | Reads | Writes |
|-----------|-------|--------|
| `labels` | the one `sqlc.request/v1` file | `<dir>/zz_<package>.labels.json`, a `labels/v1` file |
| `enums` | the one `sqlc.request/v1` file | `<dir>/zz_<package>_enums.go` |
| `views` | the one `sqlc.request/v1` file | `<dir>/zz_<package>_views.go` |

With no request, neither writes anything. When the default schema gives no
label at all (no table or view, and no enum value), `labels` writes
nothing; with no enum type, `enums` writes nothing; with no relation whose name
ends in `_view`, `views` writes nothing. Nothing empty is ever
written, because i18n would turn an empty labels file into Go with an
unused import. An enum with no values still gets its `Label()`, which
always returns nil.

All three make every check of §6 that isn't marked as `views`' own, so
`enums` or `views` run without `labels` refuses what `labels` would.

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

import "github.com/epistemer/i18n"

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
```

- One method per enum type, in the order of §4.1, on the type sqlc declares
  and with sqlc's constants. Nothing else is declared.
- `Null<Type>` gets no method: call `n.<Type>.Label()` when `n.Valid`.
- The value's string is `string(t)`; the list of values is sqlc's
  `All<Type>Values()` (`emit_all_enum_values`).
- `i18n.TKer` is what i18n's `Key` returns; it renders with
  `.T(ctx)` and names its key with `.K()`.

### 4.3 The views file

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

### 4.4 Package listing

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

## 5. Labels

### 5.1 Naming functions

| Function | Does | Example |
|----------|------|---------|
| `strip(s)` | Removes `table_prefix` from the front of `s`, if `s` starts with it | `app_order_line` → `order_line` |
| `zpPascal(s)` | Splits `s` on `_`, drops empty parts, upper-cases the first character of each part if it is an ASCII letter, and joins them | `order_line` → `OrderLine`, `user_id` → `UserId`, `2fa` → `2fa` |
| `readable(s)` | The same split, joined with spaces | `order_line` → `Order Line` |
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
| Enum value | The value itself, even when empty |
| Table, view, column | Its SQL comment with white space trimmed from both ends; if that leaves nothing, `readable(strip(name))`. White space is `\t`, `\n`, `\v`, `\f`, `\r`, U+0085 and the Unicode categories `Zs`, `Zl` and `Zp`, as the regular expression `[\t\n\v\f\r\x{85}\p{Z}]` |

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
| A `_view` relation or one of its columns has a name with a character outside ASCII, or whose `structName` is empty | `models: view "café_view": sqlc's Go name for it can't be derived` |
| A `_view` relation with no columns | `models: view "empty_view" has no columns` |
| Two columns of a `_view` relation with one `structName` (sqlc declares the field twice) | `models: view "dup_view": columns "user_1" and "user1" both give the field "User1"` |
| Two `_view` relations with one `<Type>` | `models: views "user_1_view" and "user1_view" both give sqlc the Go type "User1View"` |
| A query whose name, as written, is a method `List<Type>` | `models: query "ListOrderSummaryView" clashes with the listing method of view "order_summary_view"` |
| A query or enum that gives sqlc the name `Listings`: a query named so, or an enum whose `structName` is it | `models: enum "listings" gives sqlc the Go name "Listings", which views declares` |
| A query whose sqlc constant is `listing` (a query named `Listing` or `listing`), which hides the package `listing` | `models: query "Listing" gives sqlc the constant "listing", which views imports as a package` |

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
| `emit_methods_with_db_argument: true` | `Queries` has no `db`: `zz_<package>_views.go` doesn't compile |
| A table whose sqlc type is `Listings` (sqlc's singular of its name, which models can't compute) | `Listings redeclared` |

## 8. Testing an app

- `genxtest.AssertGenerated(t)` fails when any generated file, the request
  included (through `sqlc diff`), is stale.
- `i18ntest` sets a translation on a generated key and checks the rendered
  text, through `Label()` or a variable or `L10n(path)`, like any other key.

## 9. Not in v1

Parameters that mirror sqlc's `rename` and `initialisms`; several requests
or schemas; labels for an enum type itself; option lists for forms. The
roadmap's backlog orders them.

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
