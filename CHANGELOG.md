# Changelog

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
