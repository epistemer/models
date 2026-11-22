# Review 0013: A4, column sets for views

| | |
|---|---|
| Date | 2026-11-22 |
| Reviewer | Gemini 4 Argon (Medium), through agy (`adversarial_review`), about 27k tokens |
| Under review | The A4 diff (package `column`, the views unit's jq, the template, the fixture and its test), a generated sample, and ADR 0007 with interface §4.3 |
| Scope | sqlc's pgx/v5 types against the getters, nulls and edge values, name clashes in the generated Go, jq escaping |

There were 10 findings. Claude checked each, sqlc's type mapping against
its v1.29.0 source; six changed the code or the docs.

| # | Finding | Verdict | Resolution |
|---|---|---|---|
| 1 | `citext` and `name` are `interface{}` in sqlc, so their getters don't compile | Rejected: sqlc v1.29.0 maps `text`, `varchar`, `bpchar`, `citext` and `name` to `string` or `pgtype.Text` (`postgresql_type.go`) | — |
| 2 | Two columns giving one field, or one label | Rejected: `views` already refuses two columns with one field (§6), and `labels` two labels with one variable | — |
| 3 | A query named `Column` hides the import; `<Type>Columns` can clash with a query, a relation or an enum | Valid | `column` joins the names a query's constant can't take, and a clashing set is refused; both tested |
| 4 | `views` needs `labels` and `enums` without saying so | Valid, already documented as a need (§2, ADR) | A row in §7's limits too |
| 5 | A numeric NaN or ±Infinity, or one past float64's range, reads as 0 | Valid | `column.Numeric` parses the exact digits: NaN, ±Inf; tested |
| 6 | An infinite date or timestamp reads as year 1 | Valid | `column.Infinite`: zero `Time` and `Text` "infinity" or "-infinity"; tested |
| 7 | A type qualified with the default schema (an enum written `public.order_status`) is left out; other spellings | Valid for the qualification; the rest rejected: sqlc writes the catalog's names (`varchar`, `timestamptz`), lower-cased | The default schema's qualification is stripped, as pg_catalog's is |
| 8 | Building getters by `gsub` on a sentinel is fragile | Valid | `column_how` takes the field and writes it directly |
| 9 | `Value()` boxes the digits on every read | Rejected: one small allocation per numeric cell, for the exact digits pgx's own text encoding gives | — |
| 10 | The ADR's example says `L10n…_ID` where the code writes `_Id`; `Text`'s comment | Valid | The ADR follows the code |

One more limit found while checking: sqlc's `emit_pointers_for_null_types`
types nullable columns as pointers, which the getters don't read. It is
in §7.
