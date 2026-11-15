# Changelog

## v0.2.0 (2026-11-15)

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
