# Review 0015: A5's docs, F2, and review 0014's fixes

| | |
|---|---|
| Date | 2026-11-24 |
| Reviewer | Gemini 4 Hydrogen (Medium), through agy (`adversarial_review`), about 52k tokens |
| Under review | What review 0014 didn't see: F2's diff (#15), A5's changes to interface.md and ADR 0008, and the commit that answered review 0014 |
| Scope | Contradictions inside the contract, claims a reader would act on wrongly, bugs in `normalize`, the Valuer fast path and the enums refusal |

There were 13 findings. Claude checked each one; ten changed the code or the
docs, and they land after v0.5.0, in v0.5.1.

| # | Finding | Verdict | Resolution |
|---|---|---|---|
| 1 | `[16]` of a named byte type panics in `reflect.Copy`, in `normalize` and at init | Valid | Copied byte by byte (`bytes16`); tested |
| 2 | A pointer-receiver Valuer behind a pointer is refused | Rejected: the pointer is stripped, and the pointee's pointer method set is checked on an addressable value; tested since M7.2 | — |
| 3 | `normalize` widens a predeclared `int32`, so a panic names `int64` | Valid | Predeclared types aren't normalized; tested |
| 4 | §4.3 doesn't give sqlc's `db_type` spellings | Valid | §4.3 gives them |
| 5 | §4 omits `enums`' own refusal | Valid | Fixed |
| 6 | §2 reads as if `pgx` needs sqlc's `sql_package` | Valid | §2 says the module needs pgx v5.7.0 or later, and no sqlc option is read |
| 7 | JSON of a null `Null<Type>` fails | Rejected: documented, and accepted by epistemer in A5.2 | — |
| 8 | "pgtype never panics" is too strong | Valid for a set built by hand with a wrong kind | §4.6 says so |
| 9 | `TestPgxRefusals` stopped checking that `enums` doesn't make pgx's refusals | Valid | It checks again for the names the enums file doesn't use |
| 10 | D10 omits `column` | Valid | Fixed |
| 11 | A5.5 names `readable(v)` for `readableValue(v)` | Valid | Fixed |
| 12 | F2's note names `enums` for a relation's `rename` | Valid | It names `views` only |
| 13 | Two diffs add one hunk | Rejected: the review saw F2's diff on its own and inside A5's range | — |
