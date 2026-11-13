# Spike S1: every zp label from sqlc's request

A record, kept as the starting point for M2.1 and M2.2. The decision is
[ADR 0002](../../adr/0002-spike-s1-zp-labels.md). Nothing loads the files
here: a metafile only counts in the directory of a Go package.

| File | What it is |
|------|------------|
| [genx.meta.yaml](genx.meta.yaml) | `labels` (enum values, tables, views, columns, as `labels/v1`) and `enums` (Go), with the naming functions in `defs` and the parameters hard-coded |
| [enums.go.tmpl](enums.go.tmpl) | `Label()` for every enum type |

## Running it again

zp's schema is private, so no request is committed here. In a scratch
folder, never in zp:

1. A module with zp's `schema/` and `sql/`, a package directory `gen/`
   holding these two files and a `gen.go` (`package gen`), and the
   real prefix in `defs`.
2. A `config/sqlc.yaml` with a `gen: go:` block (package `models`, out
   `../models`, `pgx/v5`, `emit_all_enum_values`) and a `gen: json:` block
   (out `../models`, filename `codegen_request.json`).
3. A `genx.yaml` with one producer, `sqlc`, run as
   `go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0 generate -f config/sqlc.yaml`,
   emitting `models/codegen_request.json` as `sqlc.request/v1`, and
   `use` switching on `<module>/gen: [labels, enums]` and
   `github.com/epistemer/i18n: [labels]` (`v0.2.0`, with a `//go:build genx`
   blank import).
4. `go run github.com/epistemer/genx/cmd/genx@v0.1.0`, then `go vet ./models`,
   and compare the `key` and `var` of every label in
   `models/zz_models.labels.json` with the `i.Key` declarations of zp's
   `models/i18n.go` and `models/enums.go`.
