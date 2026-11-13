# models

neozp's package for sqlc-generated database models.

- Read [docs/roadmap.md](docs/roadmap.md), then the newest entry of
  [docs/status.md](docs/status.md). The design is
  [ADR 0001](docs/adr/0001-models-on-sqlc.md);
  [docs/brief.md](docs/brief.md) is the history before it.
- Claude owns this repository until v1.0.0, on the terms of
  [docs/ownership.md](docs/ownership.md) (delegated 2026-11-14). API design
  and the interface freeze are epistemer's.
- Load the `genx` skill before touching a manifest, a metafile or a labels
  file, and the `i18n` and `cfgx` skills when the work uses them.
- zp (`/home/epistemer/zp`) is private and read-only. Never run a
  generator inside it, and keep its schema, table names and enum names out
  of this repository, which is public. `codegen_request.json` is ignored by
  git for that reason.
- Every working session adds an entry to docs/status.md.
