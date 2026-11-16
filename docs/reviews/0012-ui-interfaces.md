# Review 0012: interfaces for UI packages (M6)

| | |
|---|---|
| Date | 2026-11-16 |
| Reviewer | Gemini 4 Argon (Medium) through agy, adversarial review, about 46k tokens |
| Under review | The diff of M6 (without generated files and `go.sum`), ADR 0006 and interface §4.2 and §4.5, in a scratch folder |
| Triage | Claude, as owner |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Critical | `All<Type>Options` can clash with an enum constant (`all_x` with the value `options`) or a table's struct | **In part.** Constants are refused now, with a test. A table's struct is sqlc's name made singular, which models can't compute (as for `Listings`), so it is a §7 limit |
| 2 | Major | `go.mod` gains dependencies without `go.sum`; `listing` users inherit i18n's | **Rejected.** `go.sum` is committed (the diff left it out). Apps using models already require i18n (ADR 0006, D5) |
| 3 | Major | `Labeler.LabelOf` labels one id per call, so a page of link columns makes a query per row; ids are strings | **Rejected**, as epistemer pointed out: a table calls it once per link filter, for the chosen id; rows get their text from the view's joins, as every zp view does. Strings stay (zp's ids were strings, and a typed id would need generics) |
| 4 | Major | A table that keeps a filter as an `option.SelectOption` can't pass it as a `Filter.Value` | **Accepted** as documentation: pass the enum value or `Value()` (§4.3) |
| 5 | Major | The CHANGELOG understates what breaks: interfaces, mocks and fakes of `Listings`; `All<Type>Options` is a package-level name; an enum type can't be a `driver.Valuer` | **Accepted**, in the CHANGELOG and §7 |
| 6 | Minor | `All<Type>Options()` of an enum with no values returns a zero-capacity slice, not a new array | **Rejected.** Nothing can be written to it |
| 7 | Minor | Untested: a nil `Source`, `Labeler`, `Value()` outside the enum; the enums render | **In part.** The first three are tested now; the goldens already cover the enums file, a digit-led enum included |
