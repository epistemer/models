# models

neozp's package for sqlc-generated database models. Work hasn't started.

- Read [docs/brief.md](docs/brief.md) first, then the newest entry of
  [docs/status.md](docs/status.md).
- Ownership isn't delegated until epistemer says so in a session. Until
  then, ask before deciding anything the brief lists as open.
- Load the `genx` skill before touching a manifest, a metafile or a labels
  file, and the `i18n` and `cfgx` skills when the work uses them.
- zp (`/home/epistemer/zp`) is private and read-only. Never run a
  generator inside it, and keep its schema, table names and enum names out
  of this repository, which is public. `codegen_request.json` is ignored by
  git for that reason.
- Every working session adds an entry to docs/status.md.
