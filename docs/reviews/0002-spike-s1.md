# Review 0002: spike S1 (M1.1)

| | |
|---|---|
| Date | 2026-11-14 |
| Reviewer | Gemini 3.8 Flash (High) through agy, adversarial review, about 247k tokens |
| Under review | ADR 0002, the spike's metafile and template, and ADR 0001's D4–D10, in a scratch folder |
| Triage | Claude, as owner |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Critical | `.enums[]`, `.columns[]` and `.vals` fail on null | **Not seen, guarded anyway.** sqlc's JSON writes empty arrays (checked on a toy request: `"enums": []`). M2.1 iterates with `[]?` so a future sqlc can't break it |
| 2 | Critical | Non-ASCII enum names give a wrong type name; the refusal isn't implemented | **Accepted.** ADR 0002 now says the spike has none of the checks; M2.1 writes them |
| 3 | Critical | No uniqueness, empty-strip or `L10n` checks | **Accepted.** Same as 2 |
| 4 | Major | ADR 0002 wrongly said a table or column variable could be unexported: every variable starts with `L10n` | **Accepted.** The non-ASCII refusal is now for enum names only (ADR 0001 D9, ADR 0002) |
| 5 | Major | `\s` in gojq is ASCII only, so a comment of non-breaking spaces isn't treated as empty | **Accepted.** M2.1 trims `[\s\p{Z}]` |
| 6 | Minor | Code point slicing splits decomposed letters | **Rejected.** zp's casing works on runes too, and keys must stay zp's |
| 7 | Minor | `strip` fails on null | **Rejected.** Names in the request are always strings |
