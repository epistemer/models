# Review 0001: ADR 0001 and the roadmap (M0.3)

| | |
|---|---|
| Date | 2026-11-14 |
| Reviewer | Gemini 4 Hydrogen (High) through agy, adversarial review, about 59k tokens. Argon was cooling down |
| Under review | ADR 0001 and the roadmap as drafted, the brief's open questions, and an excerpt of the standing rules, in a scratch folder |
| Triage | Claude, as owner |

Twelve findings: ten accepted, two in part; two rejected.

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Critical | Stripping the prefix from column names can merge two columns (`app_name` and `name`) | **In part.** zp strips it from columns too (its `gen/meta`), so D4 keeps that for zp's keys. The collision is caught: finding 2's checks |
| 2 | Critical | `zpPascal` and `strip` aren't one-to-one; only table-against-enum variable collisions were guarded; an empty name after `strip` gives the variable `L10n` | **Accepted.** D4: the `labels` generator checks every key, path and variable for uniqueness and refuses an empty stripped name. M2.1 tests each refusal |
| 3 | Critical | The spike covered tables only, so a mismatch in enums, texts or zp's sqlc options would show after the freeze | **Accepted.** M1.1 covers every label and audits zp's `sqlc.yaml`, views and comments |
| 4 | Major | sqlc's constant naming was underspecified: replacement, duplicates, initialisms in values, `overrides`, non-ASCII | **Accepted.** D9 lists the rules and the limits; the fixture (now M1.2, before the freeze) covers each case |
| 5 | Major | A variable could clash with a name sqlc declares in the same package (a table called `l10n_user`) | **Accepted, narrowed.** D6: a table, view or enum whose Go name starts with `L10n` stops the run. Checking against every name sqlc emits would mean reimplementing sqlc's struct naming (singular forms) |
| 6 | Major | D8 claimed a second schema or request could come later without changing a key | **Accepted.** D8 now says the other schema needs its own namespace, designed then; the default schema's keys don't change |
| 7 | Major | The freeze comes before the generators; parameters must reach `out`; unknown parameters; `package` from a `dir` like `.` | **Accepted, mostly.** D2 and M1.3: parameters reach every query (`out` is a query) and unknown ones are errors; the generators validate `dir` and `package`. The freeze stays before M2, as in genx, but the fixture moves before it, and the risk is listed |
| 8 | Major | ADR says `i18n.TKer`, the brief `i18n.Ter` | **Rejected.** i18n's interface §3.2: `Key` and the generated `L10n` return `TKer`, which embeds `Ter`. `Label()` returns what the variables are. The brief's `Ter` described zp |
| 9 | Major | Done-when gaps: sqlc not pinned, refusals with no fixtures, no zero-enum test, thin runtime tests | **Accepted.** M1.2 pins sqlc v1.29.0 as a tool in a nested module and adds a schema without enums; M2.1 writes a small request per refusal; M2.3 covers tables, views, nil and duplicates |
| 10 | Minor | A comment of white space only; dots in names make paths ambiguous | **Accepted.** D7 falls back when the trimmed comment is empty. Dots: the uniqueness check catches a real ambiguity, and keys must stay zp's |
| 11 | Minor | The skill lives outside version control | **Rejected.** epistemer keeps every neozp skill in `~/.claude/skills/`; M4.2 says so |
| 12 | Minor | ADR "Accepted" while its review is open | **Accepted.** The triage lands in the same PR as the ADR, which is accepted when that merges |
