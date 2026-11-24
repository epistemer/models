# ADR 0008: column sets read any field type sqlc generates

| | |
|---|---|
| Status | Accepted by epistemer on 2026-11-24, with amendment A5's defaults |
| Date | 2026-11-24 |
| Roadmap item | M7.1 (A5.1); Confab iteration 1, finding F5 |
| Deciders | epistemer decided the goal on 2026-11-24 (board item "models: API questions from Confab iteration 1", M1) and decides this text; Claude proposes it |

## Context

A4's column sets read each field the way sqlc v1.29.0 types it with its
defaults and `sql_package: pgx/v5`: `r.Note.Valid`, `r.Note.String`,
`r.Total.Value()`. With `emit_pointers_for_null_types: true` or an
`overrides` entry, the field has another type and the views file doesn't
compile. A4 documented that as a limit (interface §7), and Confab followed
it, converting every nullable column by hand both ways (`internal/pg`).

epistemer decided on 2026-11-24 that column sets support zp's sqlc setup:
pointers for null types, and overrides (`numeric` to `decimal.Decimal`,
`timestamptz` to `time.Time`, `interval` to `time.Duration`, and an
app's own types), each with a pointer when nullable.

The constraint is what models sees. **sqlc's request carries none of its
Go options**: `settings.codegen.options` is the JSON plugin's, and it is
empty. One request serves both configurations, so models' output for
them must be the same file, and that file must compile against either.

A spike on 2026-11-24 (sqlc v1.29.0, one view of every kind, nullable and
not) gave these field types:

| SQL | Defaults | Pointers and zp-shaped overrides |
|---|---|---|
| `text` null | `pgtype.Text` | `*string` |
| `int4` null, `int2` null | `pgtype.Int4`, `pgtype.Int2` | `*int32`, `*int16` |
| `float8` null, `bool` null | `pgtype.Float8`, `pgtype.Bool` | `*float64`, `*bool` |
| `numeric`, null | `pgtype.Numeric` | `decimal.Decimal`, `*decimal.Decimal` |
| `uuid`, null | `pgtype.UUID` | `uuid.UUID`, `*uuid.UUID` (each needs its override; without one, `pgtype.UUID`) |
| `date`, `timestamp`, `timestamptz` | `pgtype.Date`, `pgtype.Timestamp`, `pgtype.Timestamptz` | `time.Time`, `*time.Time` with overrides; **pgtype without them, even with pointers** |
| `interval`, null | `pgtype.Interval` | `time.Duration`, `*time.Duration` |
| `int8` with a column override | `int64` | the app's `type ID int64` |
| enum null | `NullStatus` | `NullStatus`: pointers don't change enums |

sqlc matched `db_type` as `pg_catalog.numeric`, `pg_catalog.interval`,
`pg_catalog.timestamp`, `timestamptz`, `date` and `uuid`; it ignored
`pg_catalog.timestamptz` and `timestamp` without a word.

## Options

1. **Getters that read the field's type at run time** (this proposal).
   The generated getter hands the field to `column.Field`, which learns
   the field's type once, when the set is declared, and reads every value
   of it through the rules in A5.1: pointers, `driver.Valuer` (every
   pgtype type, `decimal.Decimal`, `uuid.UUID`), `time.Time`, and Go's
   basic kinds, named or not. Nothing for an app to write. One views
   file serves every sqlc configuration. *Cost:* a getter that can't read
   its field fails when the package is initialized, not when it compiles.
2. **genx parameters that mirror sqlc's Go options** (`null_types:
   pointers`, a list of overrides). Getters stay compiled code. *Cost:*
   the app writes its overrides twice, in `sqlc.yaml` and `genx.yaml`,
   and keeps them in step. Nor can a parameter say how to read an app's
   own type, so it needs a vocabulary of readers per Go type.
3. **sqlc's request with the Go options**, through a sqlc plugin. *Cost:*
   models would become a program, which ADR 0001 rules out.

## Decision

Option 1. `column.Field[R any](k Kind, field func(R) any)
func(R) Value`. It calls `field` on R's zero value when the set is
declared: the interface it returns holds the field's static type even
when the field is a nil pointer. Field resolves a reader for that type
there, and panics if kind `k` can't read it. Enum columns keep A4's
compiled getters, since sqlc types enums the same way in both
configurations.

## Consequences

- An app may switch on `emit_pointers_for_null_types` and the overrides
  above, and delete its conversion helpers. Confab's `internal/pg` goes,
  once its `sqlc.yaml` also overrides `date` and `timestamptz` (pgtype
  otherwise, as the spike shows).
- The failure moves from compile time to package initialization: an app's
  start, or any test that imports the models package. That is still before
  a row is read (Principle 4). A `driver.Valuer` whose result its kind
  can't take (an app's own type) panics on that row, naming both types.
- Default configurations go through `Field` too. Their goldens change and
  their values don't: the PostgreSQL test runs both configurations.
- `gentest` gains a second sqlc Go config over the shop schema, with
  `github.com/shopspring/decimal` and `github.com/google/uuid` in its own
  nested module only. The root module gains no dependency.
- `interval` stays out of every kind (A5.1's sub-question).
- Interface §4.3, §4.6 and §7 change (A5.1).
