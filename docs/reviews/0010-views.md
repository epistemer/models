# Review 0010: the views generator (M5.4)

| | |
|---|---|
| Date | 2026-11-15 |
| Reviewer | Gemini 4 Argon (Medium) through agy, adversarial review, about 80k tokens |
| Under review | The diff of the metafile, `views.go.tmpl`, the tests and CI, the generated file, and interface §4.3 and §6's views rows, in a scratch folder |
| Triage | Claude, as owner. Claims about sqlc were checked against sqlc v1.29.0 |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Critical | The shared checks run first, so §6's own examples (`user_1`/`user1`) get the shared message | **Accepted as documentation.** The run still stops; §6's examples now use `id`/`ID`, which only sqlc's initialism merges, and say the shared refusal comes first otherwise |
| 2 | Critical | Generated Go that doesn't compile: a table typed `Listings`; a `_view` type equal to an enum's; query constants `context`, `max`, `len` and the like; `oxen_view`, which sqlc makes `OxView`; a newline in a name reaching a comment | **In part.** Checked: sqlc's `^(ox)en` rule ignores the end of the name, so models now follows it; a query `Max` gives `const max`, so the template no longer calls `max`, and a query constant hiding any name the file uses (`listing`, `context`, `append`, `len`, `nil`, `int64`, `string`, `error`) is refused; control characters are refused (the shared rule already stops a newline; tested). Rejected: a table typed `Listings` stays a §7 limit (models can't compute sqlc's singular), and a `_view` type equal to an enum's is a duplicate sqlc declares itself |
| 3 | Major | `group_by` reports the alphabetically first colliding pair, not the first in the schema | **Rejected.** The pair it names really collides; the existing checks behave the same |
| 4 | Major | A query clash is reported before a later view's problems | **Rejected.** One true problem per run, as everywhere in §6 |
| 5 | Major | `rows` is only closed by `defer`, before the count runs | **Accepted.** Closed explicitly after the loop, as sqlc does. (pgx closes rows when `Next` returns false, so it worked) |
| 6 | Major | The PostgreSQL test doesn't check whether the count runs, the total's floor, three page boundaries, or a non-null nullable field; refusal tests miss lowercase `listing` and another engine without views; an exemption contradicted the test's comment | **Accepted.** A counting DBTX checks every case; a deletion between page and count tests the floor; the cases and the exemption are fixed |
| 7 | Minor | `index` on a one-element array; an unreachable check; a doc comment wrapped unlike §4.3 | **Accepted** (`IN`, the check removed, §4.3 matches the template) |
