# Review 0004: the interface draft (M1.5)

| | |
|---|---|
| Date | 2026-11-14 |
| Reviewer | Gemini 4 Hydrogen (High) through agy, adversarial review, about 55k tokens. Argon was cooling down |
| Under review | The draft of `interface.md`, the fixture's README, genx's A4 and §7, i18n's labels template and spike S1's metafile, in a scratch folder |
| Triage | Claude, as owner |

Fifteen findings: twelve accepted (three of them in part), three rejected.

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Critical | A quoted name like `"first name"`, or one that `zpPascal` empties (`_`), gives a variable that isn't a Go identifier, or `L10n` itself | **Accepted.** §6 refuses a `<T>` or `<C>` that is empty or not a Go identifier |
| 2 | Critical | "First letter" in `structName` misreads `2fa`; an enum named `_` gives an empty type name | **Accepted.** §5.1 says "first character, if it is a letter"; §6 refuses an empty `structName` |
| 3 | Critical | The numbering rule was ambiguous, and `value_1` then `!!` gives two `value_1` | **Accepted.** §5.2 states sqlc's algorithm step by step. Checked on sqlc v1.29.0: it declares `EValue1` twice and its own Go breaks; models refuses the duplicate variable |
| 4 | Critical | `package` as "a Go identifier" can't be checked in jq, allows `_`, clashes with genx §7, and a `null` default accepts any JSON type | **Accepted.** §3: a string matching `^[A-Za-z_][A-Za-z0-9_]*$`, not `_`, not a keyword |
| 5 | Major | An enum with no values (valid in PostgreSQL) made "no tables or enums" the wrong test for an empty labels file | **Accepted.** §4: no file when the schema gives no label at all; an enum with no values still gets a `Label()` that returns nil |
| 6 | Major | "`\s` and Unicode `Z`" is read differently by gojq, jq and Go | **Accepted.** §5.3 lists the characters and the regular expression |
| 7 | Major | §7 was wrong about a second `sql:` block, a mismatched `dir` with no enums, and enum `overrides`; `omit_unused_structs` was missing | **Accepted.** §7 rewritten row by row |
| 8 | Major | The `L10n` check missed composite types; a missing default schema went unnoticed; `enums` alone skipped the checks | **In part.** A missing default schema is refused, and both generators make every check (§4). Composite types: rejected, because sqlc v1.29.0 declares no Go type for them (checked: a composite column became `sql.NullString`) |
| 9 | Major | `dir` must also meet genx's path rules; parameters must be checked even with no request | **Accepted.** §3 |
| 10 | Major | The blank-import file says `package main` | **In part.** It may sit in any package and carries that package's name; §2 says so, and how to regenerate. A file excluded by its build tag doesn't clash with the directory's package |
| 11 | Minor | Column order unstated; the JSON example's keys weren't sorted | **Accepted** |
| 12 | Minor | The `pay_pal` example could be read as index 3 | **Accepted** |
| 13 | Minor | The two-requests message names no files | **Rejected.** A query doesn't see paths; genx adds the input files to every query error (its §5) |
| 14 | Minor | The `null` rule for `package` skipped the keyword check | **Accepted** with 4 |
| 15 | Minor | Unguarded `.enums[]` and friends in the spike | **Rejected for the interface**: an implementation matter, already planned for M2.1 (review 0002, finding 1) |
