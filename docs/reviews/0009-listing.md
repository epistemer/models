# Review 0009: package listing (M5.3)

| | |
|---|---|
| Date | 2026-11-15 |
| Reviewer | Gemini 4 Argon (Medium) through agy, adversarial review, about 20k tokens |
| Under review | The diff of `listing/` and interface §4.4, in a scratch folder |
| Triage | Claude, as owner |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Major | A typed nil (`(*string)(nil)`) passes the `nil` check and matches nothing | **Accepted.** Nil pointers, slices, maps, channels, functions and interfaces are refused too; §4.4 says so, and that a non-nil value the driver encodes as NULL still matches nothing |
| 2 | Minor | A NUL in `Search` fails in PostgreSQL instead of as `ErrQuery` | **Accepted** |
| 3 | Minor | Duplicates in `columns`, `SearchIn` and `Sort` are accepted | **In part.** A column listed twice in `columns` is refused (the scan would be off by one). Duplicates in `SearchIn` and `Sort` repeat a condition and stay |
| 4 | Minor | Tests miss an unknown `SearchIn` column with a search, the column named in filter errors, an error after arguments were added | **Accepted** |
