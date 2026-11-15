# Review 0008: listings of views, the design (M5.1)

| | |
|---|---|
| Date | 2026-11-15 |
| Reviewer | Gemini 4 Argon (High) through agy, adversarial review, about 49k tokens |
| Under review | ADR 0005, the draft of amendment A2, and interface §6–§7, in a scratch folder |
| Triage | Claude, as owner. Claims about sqlc were checked against sqlc v1.29.0 in a scratch folder |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Critical | Two columns of a view with one `structName` scan into one field, silently | **Accepted as a refusal; the claim is wrong.** sqlc v1.29.0 declares the field twice (no `_2`), so the app never compiled. Refused anyway, as A1 does for enums |
| 2 | Critical | Without parentheses, `AND` binds the filters to the last search column only, leaking rows past a filter | **Accepted.** The search is a parenthesized group; a unit test will pin it |
| 3 | Critical | A `"` or NUL in a name breaks the SQL or the generated Go | **Accepted.** `Build` doubles `"` and refuses NUL, an empty relation or no columns; the template writes names as quoted Go strings |
| 4 | Critical | `emit_methods_with_db_argument` compiles and panics on a nil `q.db` | **Accepted as a limit; the claim is wrong.** sqlc then declares `Queries` with no `db`, so the file doesn't compile. Listed in §7 |
| 5 | Major | A `nil` value matches nothing for every `Op`, and the contract would block a later `IS NULL` | **Accepted.** A `nil` value is an error |
| 6 | Major | A query named `Listing` gives sqlc the constant `listing`, which clashes with the import | **Accepted.** Checked: sqlc writes `const listing`. Refused |
| 7 | Major | Missing refusals: a digit-led name, two `_view` relations with one type, no columns; `Listings` vs a singular table name; query names compared by `structName` | **In part.** Two views with one type and a view with no columns are refused. A table named `listings` gives sqlc `Listing`, so a table whose type is `Listings` becomes a limit. Rejected: a digit-led name gets sqlc's `_` (`_1View`, checked), which is a Go identifier; and sqlc names a query's method exactly as written (`list_things`, checked), so the raw name is right |
| 8 | Major | No `IN`; driver errors for bad values can't be told from database failures; the column list isn't exported; `SearchIn` can't default to every column | **In part.** Documented: `ErrQuery` covers the shape, and the caller converts values by column type. Rejected for v0.2.0: `IN` (zp's filters take one value; a new `Op` is a minor amendment later), an exported column list (the admin package defines its columns, as zp's does), and searching every column by default |
| 9 | Major | Counting before the page doubles the queries of every short listing; the total can come out lower than the rows seen | **Accepted.** The page runs first; the count runs only when the rows can't tell the total, and the total is never less than `Offset` plus the rows |
| 10 | Minor | `SearchIn` should be checked even when `Search` is empty | **Accepted** |
| 11 | Minor | Paging without a unique sort column is unstable | **Accepted** as documentation |
| 12 | Minor | Escaping order and an explicit `ESCAPE '\'`; the text form of non-text columns | **Accepted** |
