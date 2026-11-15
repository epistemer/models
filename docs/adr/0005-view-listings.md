# ADR 0005: Listings of views (v0.2.0)

| | |
|---|---|
| Status | Accepted: D1–D4 by epistemer on 2026-11-15, D5 with interface amendment A2 on the same day; D6–D8 decided by Claude as owner |
| Date | 2026-11-15 |
| Roadmap item | M5.1 |
| Deciders | epistemer decided D1–D4 (each as recommended, except D3 and D4, where epistemer chose zp's `_view` rule and pgx/v5 only). Claude proposed D5 and decided the rest |

## Context

zp's business code can't move onto models until models replaces zp's view
builder (`gen/view`), which ADR 0001's D13 left out of v0.1.0. 38 of zp's
files call it. What it does, from zp's code:

- For every Go type in sqlc's models whose name ends in `View`, it adds a
  method `View<Base>(ctx, Builder) (*ViewResult[T], error)` on `Queries`,
  and an interface `Q` that is sqlc's `Querier` plus those methods. zp's
  handlers depend on `Q` (555 references).
- The query goes through the app's `Builder`, which gets a squirrel
  `SelectBuilder` and adds search, filters, sorting, limit and offset. zp's
  `view.View`, an admin-table component, implements it from the request's
  query string, with column names taken from its own column definitions.
- The method runs a `COUNT(*)` with the filters, then the page, and returns
  the rows through a channel filled by a goroutine, plus an error channel and
  the total.

Facts checked for this record (spike S2, on sqlc v1.29.0, pgx v5 and
PostgreSQL 16, in a scratch folder):

- sqlc's request doesn't mark a view as a view: views are listed with the
  tables. A name is all a generator can go by.
- sqlc makes a relation's name singular before naming its type
  (`order_views` → `OrderView`). A name ending in `_view` stays the same
  (`order_summary_view` → `OrderSummaryView`, `people_view` →
  `PeopleView`).
- sqlc's DBTX differs by driver (`Query` for pgx, `QueryContext` for
  `database/sql`), and the request carries no Go option to tell which.
- pgx v5 encodes a filter's value of type `any` by the column's type: a
  sqlc enum type, a plain string, a number for `numeric`, a `time.Time` or a
  string for `date` all work. A value of the wrong kind fails with pgx's
  message, and an enum value that doesn't exist with PostgreSQL's.
- `CAST(col AS text) ILIKE $1` searches columns of any type, and one
  argument can serve several columns. `\%` and `\_` escape the search text.

## Decisions

**D1. The query is models' own value** (decided by epistemer). A runtime
package, `github.com/epistemer/models/listing`, holds a plain struct of
search, filters, sort, limit and offset. Generated code checks every column
it names against the view's own columns before any SQL is built, so a
column name never reaches SQL unchecked. No squirrel in an app's API. The
package that renders admin tables (not yet started) will fill the struct
from a request, as zp's `view.View` does today.

**D2. A method returns `([]T, int64, error)`** (decided by epistemer): the
page's rows, how many rows the search and filters match, and an error. No
channels, no goroutine, no generic type at the call site.

**D3. Relations whose SQL name ends in `_view`** (decided by epistemer):
zp's convention, in the default schema, tables included. No configuration.

**D4. pgx/v5 and PostgreSQL only** (decided by epistemer). An app on
another `sql_package` gets a compile error in the generated file (a limit,
interface §7); a request of another engine is refused when it has a
`_view` relation.

**D5. Names** (approved by epistemer with amendment A2): a third
generator, `views`, writing `<dir>/zz_<package>_views.go`; a method
`List<Type>` per relation (`ListOrderSummaryView`); an interface
`Listings` with every such method, for an app's own interface next to
sqlc's `Querier`; and in package `listing`, the types `Query`, `Filter`,
`Op` and `Sort`, the function `Build` and the error `ErrQuery`.

**D6. The SQL** (Claude). Two statements, both from `listing.Build`: the
page, then a `count(*)` with the search and filters, run only when the
page's rows can't tell the total (a short last page or a single page
tells it). Identifiers are always double-quoted with `"` doubled, so a
column named `order` works. The search is a parenthesized group of
`CAST(col AS text) ILIKE $n ESCAPE '\'` joined by `OR`, with `\`, then `%`
and `_`, escaped, and is `AND`ed with the filters. `LIMIT` and `OFFSET` are
written as integers, so both statements take the same arguments. The two
statements don't share a transaction: under concurrent writes the total and
the page may disagree, as in zp, but when rows come back the total is never
less than `Offset` plus their number. A `nil` filter value is an error,
which leaves `IS NULL` free for later.

**D7. Tests** (Claude). `listing.Build` has unit tests on the whole CI
matrix. The generated file is golden-tested and compiled in `gentest`,
whose sqlc output moves to pgx/v5 for it. A test runs the generated methods
against PostgreSQL on Linux CI (a service container), and is skipped where
no database is configured.

**D8. Not in v0.2.0** (Claude): building a `Query` from an HTTP request,
and column definitions with labels and types (the admin-table package's
job; a column's label is already `L10n("table.<view>.<column>")`);
filters on NULL; an `IN` operator (zp's filters take one value; a new
`Op` is a later minor amendment); exporting a view's column list; joins;
option lists, which still wait for the forms package.

## Consequences

- models gets its first runtime package, with no dependency. The root
  package still holds only the metafile (ADR 0001 D11).
- zp's move: `q.View<Base>(ctx, v)` becomes `q.List<Base>View(ctx, query)`
  with the rows as a slice; zp's `Q` becomes
  `interface { Querier; Listings }`; `ViewResult` and `gen/view` go. zp's
  13 views end in `_view`, so all of them get a method.
- `gentest`'s fixture moves to pgx/v5, so its goldens of sqlc's output
  change once.
- The roadmap gets M5, ending in `v0.2.0`.
