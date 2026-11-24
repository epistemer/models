# Changelog

## v0.5.0 (2026-11-24)

Confab iteration 1's findings, amendment A5 (ADR 0008), locked by
epistemer on 2026-11-24.

- **Column sets read any field type sqlc generates** (A5.1, F5): with
  `emit_pointers_for_null_types: true` and `overrides` (`decimal.Decimal`,
  `time.Time`, `uuid.UUID`, `time.Duration`, an app's own types), each
  with a pointer when nullable. New `column.Field`, which every generated
  getter but an enum's calls; it learns the field's type when the set is
  declared and panics there if its kind can't read it.
- **Text methods on enums** (A5.2, F28): `MarshalText` and
  `UnmarshalText`, which refuse a value the enum doesn't have.
- **A fourth generator, `pgx`** (A5.3, F15): `zz_<package>_pgx.go` with
  `RegisterTypes(ctx, conn)`, registering every enum and its array, for
  `pgxpool.Config.AfterConnect`. Needs pgx v5.6.0 or later.
- **`option.Static(value, text)`** (A5.4, F8), for options whose label is
  data; its label's `K()` is `""`.
- **Enum values' English text is the value made readable** (A5.5, F19):
  `in_review` is "In Review", as handlers makes route keys readable. Keys
  don't change; i18n shows the new text wherever nobody edited it.
- Docs: `table_prefix`'s trade-off, that sqlc keeps the prefix in its Go
  names (F2).

**Breaking:**

- Enum values' English text changes for every app (A5.5).
- An app's own `MarshalText`, `UnmarshalText` or `RegisterTypes` in the
  models package clashes with the generated ones.
- `encoding/json` now goes through the text methods: a sqlc row whose
  nullable enum is null (value `""`) no longer marshals, unless `""` is a
  value of the enum.
- A column set whose field type `column.Field` can't read now panics at
  package initialization; under A4 the same set didn't compile.

## v0.4.0 (2026-11-22)

Column sets for views, amendment A4 (ADR 0007), for ui's tables.

- **Package `column`**: `Kind`, `Value`, `Column[R]`, and `UUID`,
  `Numeric` and `Infinite` for generated getters. A column
  holds what the schema knows: its SQL name, its label, a kind from its SQL
  type, an enum's options, and a getter returning a plain value with nulls
  unwrapped, so a table package never imports the driver.
- **`views`** also writes a `<Type>Columns` set beside each `List<Type>`.
  A column of a type no kind takes (an array, json, bytea, …) is left out
  and named in the set's comment. A numeric keeps its exact digits; an
  infinite date or time is marked.
- Two refusals: a query whose constant hides `column`, and a column set
  whose name is already a query's, a relation's or an enum's.

**Breaking:**

- `views` now needs `labels` and i18n's `labels`, and `enums` when a view
  has an enum column; on its own it worked before. Without them the views
  file doesn't compile.
- A package-level `<Type>Columns` of the app's own, for a view's type,
  clashes with the generated one.

## v0.3.0 (2026-11-16)

Interfaces for UI packages, amendment A3 (ADR 0006): what a table or form
package needs from the models is an interface models owns.

- **Package `option`**: `SelectOption` (`Value() string`,
  `Label() i18n.TKer`). The models module now requires i18n `v0.2.0` or
  later.
- **`enums`** also writes `Value()` and `All<Type>Options() []option.SelectOption`
  for every enum type, so every enum value is a `SelectOption`.
- **`listing.Source`**: `List<Type>` takes a `Source` instead of a `Query`.
  `Query` implements `Source`, so calls that pass a `Query` don't change; a
  table package can pass its own type.
- **`listing.Labeler`**: the text of the id a link column is filtered by.
  Apps implement it; nothing generates it yet.
- A refusal: an enum whose `All<Type>Options` is another enum's Go type.

**Breaking:**

- `List<Type>`'s parameter type changes, so an interface, mock or fake
  written against `List<Type>(ctx, listing.Query)`, or a variable of that
  function type, must change to `listing.Source`. Calls that pass a
  `listing.Query` don't.
- An app's own method `Value` on an enum type, or its own package-level
  `All<Type>Options`, now clashes with the generated one. An enum type can't
  be a `driver.Valuer` either.

## v0.2.0 (2026-11-16)

Listings of views, interface amendment A2 (ADR 0005). Nothing changes for an
app that doesn't switch `views` on.

- **`views`**: a `List<Type>(ctx, listing.Query) ([]<Type>, int64, error)`
  method on sqlc's `Queries` for every relation whose name ends in `_view`,
  and a `Listings` interface of them all. pgx/v5 and PostgreSQL only.
- **Package `listing`**: `Query` (search, filters, sort, limit, offset),
  `Filter`, `Op`, `Sort`, `Build` and `ErrQuery`. Every column is checked
  against the view's own; identifiers are quoted; the search is a
  parenthesized, escaped `ILIKE`. The count runs only when the page can't
  tell the total.
- **Refusals** for what would make the views file fail to compile (§6),
  and the limits (§7).
- `examples/shop` lists a view; the generated methods are tested against
  PostgreSQL in CI.

Checked against a copy of zp: its 13 views get the same 13 types, with the
same columns in the same order as zp's `gen/view` (ADR 0004).

## v0.1.0 (2026-11-14)

The first release, against the interface locked on 2026-11-14.

- **`labels`**: a `labels/v1` file with a label for every enum value, table,
  view and column of sqlc's default schema: zp's keys, SQL names as paths,
  `L10n` variables, and the SQL comment or the readable name as English
  text.
- **`enums`**: `Label() i18n.TKer` on every enum type sqlc declares.
- **Parameters** `dir`, `package` and `table_prefix` (genx `v0.2.0`'s
  amendment A4).
- **Refusals** for everything that would give a key two meanings or a
  variable that can't compile (interface §6, with amendment A1: two enums
  with one Go type).
- `examples/shop`; the generators' tests in the nested module `gentest`;
  sqlc fixtures regenerated in CI on three OSes.

Checked against a copy of zp's schema: all 995 of zp's keys and variables
come out equal (ADR 0004).

Requires genx `v0.2.0` or later, i18n `v0.2.0` or later, and sqlc requests
of `v1.29.0`'s shape.
