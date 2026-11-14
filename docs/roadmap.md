# models — Roadmap to the first release

> Status: **accepted plan**, written 2026-11-14. Module:
> `github.com/epistemer/models`.
>
> | Document | Answers |
> |----------|---------|
> | [ADR 0001](adr/0001-models-on-sqlc.md) | *What* gets generated, and why that way |
> | `interface.md` (written in M1.5) | *What* an app writes and gets back |
> | this file | *In what order*, and *when each piece counts as done* |

## What gets built

Two trivial generators over sqlc's request, and no program:

```
sqlc generate  →  models/codegen_request.json (sqlc.request/v1, foreign)
                    → models: labels  →  models/zz_models.labels.json  →  i18n: labels  →  models/zz_models_i18n.go
                    → models: enums   →  models/zz_models_enums.go       (Label() per enum type)
```

## How to use this document

- **Ownership.** Claude owns models until `v1.0.0`, on the terms of cfgx,
  i18n and genx ([ownership.md](ownership.md), 2026-11-14).
- **IDs.** Every item has an ID (`M2.1`) and fits one PR or one focused
  session. Branch names, PR titles and commit subjects start with the ID.
- **Interface first.** M1 ends by freezing `docs/interface.md`. After that a
  change to what an app writes or gets back is an amendment.
- **Versions.** The first usable release is `v0.1.0`. `v1.0.0` waits until
  zp's business code runs on it.
- **Sizes.** S is up to 1 day, M is 2–3 days, for one developer without AI
  help. They order the work; they don't promise dates.

## Overview

| Milestone | Goal | Needs |
|-----------|------|-------|
| [M0](#m0--foundation) | Ownership, design record, this plan, reviewed | – |
| [M1](#m1--unknowns-genx-parameters-and-the-freeze) | Every zp key matched; a fixture; genx gets parameters; interface frozen | M0 |
| [M2](#m2--the-generators) | `labels` and `enums`, with tests | M1 |
| [M3](#m3--an-example-and-zp) | An example app; a comparison with zp | M2 |
| [M4](#m4--v010) | Docs, the skill, `v0.1.0` | M3 |

## Definition of done (every work item)

1. `go test -race ./...` passes on the whole CI matrix.
2. Generated files are golden-tested with `genxtest.Render`, and a second
   run gives the same bytes.
3. Every labels file a test produces passes genx's conformance checks.
4. If the behaviour differs from interface.md, that document is amended in
   the same PR.
5. A new dependency needs an ADR.
6. Nothing of zp's schema, table names or enum names enters the repository.

---

## M0 — Foundation

| ✓ | ID | Item | Size | Done when |
|---|----|------|------|-----------|
| ☑ | M0.1 | Ownership delegated; `docs/ownership.md` | S | Merged |
| ☑ | M0.2 | ADR 0001 and this roadmap | S | Merged |
| ☑ | M0.3 | Cross-model review of ADR 0001 and the roadmap. Every finding becomes a change or a recorded rejection | S | Triage in [`docs/reviews/0001-design.md`](reviews/0001-design.md) merged |

## M1 — Unknowns, genx parameters and the freeze

| ✓ | ID | Item | Size | Done when |
|---|----|------|------|-----------|
| ☑ | M1.1 | **S1: every label as a trivial generator.** On a copy of zp in a scratch folder, write the labels of enum values, tables, views and columns with D4's keys and D6's variables, a hard-coded prefix standing in for the parameter. Audit zp's `sqlc.yaml` (`rename`, `initialisms`, enum `overrides`), its views (materialized ones included) and its SQL comments against D7–D9 | M | [ADR 0002](adr/0002-spike-s1-zp-labels.md) with the comparison: the key set equals zp's keys (its `i18n.go` and `enums.go`), every one, and the variables match. **Kill criterion:** if a key can't be matched in jq, a producer program writes the labels, and ADR 0001 is amended |
| ☐ | M1.2 | A fixture: a public toy schema (`testdata/schema/`) and the request sqlc v1.29.0 makes from it, committed under an explicit name. It covers what D4 and D9 name: duplicate, empty and punctuation-only enum values, a value that gives the same name as a numbered duplicate, an initialism in a type and in a value, digits, a view, comments (one only white space), a prefix, and no enums at all in a second fixture | S | A Linux-only CI job regenerates the request with sqlc pinned at v1.29.0 (a `tool` in a nested module) and fails if it changed |
| ☑ | M1.3 | **In genx: parameters for a generator** (on genx's list of what's "not in this version"). An interface amendment: the metafile declares parameters with defaults, a `use` entry sets them, every query (`units` and `out`) reads them, an unknown parameter is an error, `genx list` shows them, and `genxtest.Render` takes them. Written with this package as its first user | M | The amendment's text approved by epistemer; implemented with e2e scripts and failure messages; merged in genx |
| ☑ | M1.4 | In genx: tag the release that carries M1.3 (`v0.2.0`) | S | `go get` resolves it |
| ☐ | M1.5 | **Freeze the interface.** `docs/interface.md`: the generators, their parameters and defaults, the files they write, the keys, paths, variables and English text, the generated Go, the limits (D8, D9), and what an app writes to switch it on | M | Locked by epistemer |

## M2 — The generators

| ✓ | ID | Item | Size | Done when |
|---|----|------|------|-----------|
| ☐ | M2.1 | `labels` generator, with the naming functions in `defs` | M | Golden files with and without a prefix, and with `dir` and `package` set; the conformance check passes. Each refusal has a test with its message, on a small request written for it: another schema, two requests, a duplicate key, path or variable, an empty name after `strip`, a Go name starting with `L10n`, a letter outside ASCII, a bad `dir` or `package` |
| ☐ | M2.2 | `enums` generator | S | Golden file; the toy schema's sqlc models, the generated labels and i18n's Go compile and vet together in a nested test module; the schema without enums writes no enums file; the same refusals as `labels` |
| ☐ | M2.3 | `i18ntest` shows a translation set on generated keys: an enum value through `Label()`, a table, a view and a column through their variables and `L10n(path)`. `Label()` is nil for a value the enum doesn't have, and right for numbered duplicates | S | Tests in the nested module |

## M3 — An example and zp

| ✓ | ID | Item | Size | Done when |
|---|----|------|------|-----------|
| ☐ | M3.1 | `examples/shop`: an app with sqlc, genx, i18n and models, its manifest, and `genxtest.AssertGenerated` | M | CI green on all three OSes, without building sqlc (no `check` for it) |
| ☐ | M3.2 | zp, on a copy in a scratch folder: generate with models, compare keys and variables with zp's files, build zp's models package, and time a run. Measure what committing the request costs (size, diff of a one-column change) | S | Numbers and findings in the status log and ADR 0002 or 0003; friction goes to genx's friction log |

## M4 — v0.1.0

| ✓ | ID | Item | Size | Done when |
|---|----|------|------|-----------|
| ☐ | M4.1 | README quickstart: add models to an app on sqlc, genx and i18n | S | A newcomer can follow it |
| ☐ | M4.2 | The `models` skill in `~/.claude/skills/models/`, where epistemer keeps the other packages' skills | S | Written against the frozen interface |
| ☐ | M4.3 | Cross-model review of the release, then CHANGELOG and tag `v0.1.0` | S | Tag pushed |

## After v0.1.0 (backlog)

1. Parameters that mirror sqlc's `rename` and `initialisms`, for apps that
   set them (D9).
2. More than one request, for apps with several databases (D8).
3. Schemas other than the default one (D8).
4. Views' query builder: zp's `gen/view`, a runtime library (D13).
5. Option types, with the package that renders forms (D10).

## Risks

| Risk | Likelihood / impact | Mitigation |
|------|---------------------|------------|
| A zp key can't be matched by a query | Low / high | S1 compares every key before the freeze, with a kill criterion |
| The freeze comes before the generators exist | Medium / low | S1 and the fixture (M1.1, M1.2) run every naming rule before it; what M2 still finds is an amendment, as in genx |
| genx's parameters take longer or come out different | Medium / medium | M1.3 is designed with this package as the user; the freeze waits for it |
| sqlc changes its request between versions | Low / medium | The fixture is regenerated in CI with the pinned sqlc; the domain is `sqlc.request/v1` and a new major is a new domain |
| An app's sqlc naming options give a type name the query doesn't know | Medium / low | A compile error names the type (D9); backlog item 1 |
| The committed request is large | Medium / low | Measured in M3.2; friction for genx |
