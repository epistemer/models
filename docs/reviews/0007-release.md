# Review 0007: the release, v0.1.0 (M4.3)

| | |
|---|---|
| Date | 2026-11-14 |
| Reviewer | Gemini 4 Hydrogen (High) through agy, adversarial review, about 161k tokens. Argon was cooling down |
| Under review | interface.md, the metafile and template, README, CHANGELOG, the `models` skill and the refusal tests, in a scratch folder |
| Triage | Claude, as owner |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Major | `dir` wasn't checked against genx's path rules (§3 says they apply), so `"my models"` or `aux` passed on runs that write nothing | **Accepted.** models checks every element as genx §7 does; four tests added, one without a request |
| 2 | Minor | README: `go:generate` in the `genx`-tagged file is skipped; `genxtest`'s import path not named | **Accepted** |
| 3 | Minor | README didn't say `table_prefix` also strips column names, which are in keys | **Accepted** |
| 4 | Minor | Two enums with one `structName` (`user1`, `user_1`) and different values pass, and `enums` writes `Label()` twice | **Accepted** as interface amendment A1: §6 refuses it. sqlc declares the type twice too, so no app that compiled is affected |
| 5 | Minor | The skill says `enum.<type>`, leaves views out of two refusals | **Accepted.** The local source path stays: every neozp skill names its repository's checkout |
| 6 | Minor | The template leaves a trailing blank line | **Rejected.** gofmt removes it, as for every Go output (genx §6) |
| 7 | Minor | Test gaps (`dir` path rules, `package: ""`); the failure wrapper | **In part.** The cases are added; the wrapper stays (review 0005, finding 2) |

After the fixes, the README's quickstart was followed in an empty module,
from `go mod init` to `go tool genx check`, with the released genx and i18n.
