# ADR 0004: zp on models v0 (M3.2)

| | |
|---|---|
| Status | Accepted |
| Date | 2026-11-14 |
| Roadmap item | M3.2 |
| Deciders | Claude, as owner. Run on a copy of zp's schema and queries in a scratch folder; nothing of it is committed |

## Context

The generators on `main` (M2) had to be checked where they will run first:
zp, whose translations are stored under its own keys. ADR 0001's D12 also
left open what committing sqlc's request costs a real repository.

## What was run

A scratch module with zp's `schema/` and `sql/`, sqlc v1.29.0 and genx
`v0.2.0` as tools, i18n `v0.2.0`, and models from this repository through a
`replace`. The manifest is interface §2's, with `table_prefix` set to zp's
prefix and the other parameters at their defaults. sqlc used zp's Go
options without its type overrides, which point at zp's own packages.

## Results

| Check | Result |
|-------|--------|
| Keys and variables | **995 of 995 equal** to zp's (its `i18n.go` and `enums.go`), as in spike S1, now from the released metafile with the parameter |
| Refusals | None: zp's schema passes every check of interface §6 |
| Compile | sqlc's models, the 28 `Label()` methods and i18n's Go pass `go vet` |
| `go tool genx`, warm | 1.9 s, of which sqlc is 1.2 s |
| `go tool genx`, first run after `go mod tidy` | 6.2 s (tool builds cached, generated files absent) |
| `go tool genx check`, warm | 2.0 s, `sqlc diff` included |

**Committing the request (D12).** The request is 5.6 MB. Committed with the
rest of the scratch module, the repository's objects took 344 KB after
`git gc`. Adding one column to one table changed 161 lines of the request
(the column appears in the table and in every query that selects it), 6 of
the labels file and 2 of i18n's Go, and the objects grew by 4 KB. Git's
compression and deltas absorb the size; reviews see a readable diff.

## Decision

D12 stands: the request is committed. Nothing goes to genx's friction log:
parameters (A4), a foreign file and `sqlc diff` as a producer's check
worked as designed.

## Consequences

- zp's move needs no key change, as planned (ADR 0002's consequences).
- M3.1's example and the README quickstart can show interface §2 as is.
