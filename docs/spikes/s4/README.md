# Spike S4: zp's enums from sqlc's request

A record, kept as a starting point. The decision it led to is genx's
[ADR 0006](https://github.com/epistemer/genx/blob/main/docs/adr/0006-foreign-file.md).
Nothing loads the files here: a metafile only counts in the directory of a
Go package, and this directory has none.

| File | What it is |
|------|------------|
| [genx.meta.yaml](genx.meta.yaml) | Two trivial generators over `sqlc.request/v1`: `enums` writes Go, `enum-labels` emits a `labels/v1` file |
| [enums.go.tmpl](enums.go.tmpl) | The template of `enums`, in the shape of zp's `models/enums.go` |

## What it shows

- sqlc v1.29.0 writes its codegen request as JSON, with no database, when
  the config has a `gen: json:` block. For zp it is 5.6 MB: 28 enums, 89
  tables and 794 columns in the `public` schema, and 391 queries.
- Enums are at `.catalog.schemas[].enums[]`, each with `name` and `vals`.
  Tables are at `.catalog.schemas[].tables[]`, each with `rel.name`, a
  `comment`, and `columns[]` with `name`, `type`, `not_null` and `comment`.
- `enums` reproduces zp's hand-written generator: 931 lines, identical to
  zp's `models/enums.go` after gofmt, apart from the header line.
- `enum-labels` writes 112 labels, one per enum value, with the key
  `enum.<type>.<value>` and the value as the English text.
- A run takes 0.36 s on a 2015 laptop, of which the go command is half
  (genx's README, "Timings").

Checked again on 2026-11-14 with genx `v0.1.0`, against a request made from
zp's schema that day.

## Running it again

zp's schema is private, so no request is committed here (`.gitignore` keeps
`codegen_request.json` out). To repeat the check:

1. Copy zp's `schema/`, `sql/`, `go.mod` and `go.sum` to a scratch folder,
   never working in zp itself. Add a `config/sqlc.yaml` that has only:

   ```yaml
   version: "2"
   sql:
     - engine: "postgresql"
       queries: "../sql/*.sql"
       schema: "../schema/*.sql"
       gen:
         json:
           out: "../models"
           filename: "codegen_request.json"
   ```

   Then `go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f config/sqlc.yaml`
   writes `models/codegen_request.json`.
2. In a second scratch module, put the two files of this directory in a
   package directory `gen/` (with a `gen.go` that says `package gen`), copy
   the request to `models/`, and write this `genx.yaml`:

   ```yaml
   version: 1
   producers:
     - name: models
       run: [go, version]   # a stand-in: the request is already there
       emits:
         - path: models/codegen_request.json
           domain: sqlc.request/v1
   use:
     - example.com/s4/gen: [enums, enum-labels]
   ```

3. `go run github.com/epistemer/genx/cmd/genx@v0.1.0`, then compare
   `models/zz_enums.go` with `gofmt`'s output for zp's `models/enums.go`.

## What it doesn't show

- **Names sqlc changes.** `rename:` and the initialisms list of sqlc's Go
  options aren't in the request, and sqlc resolves colliding names with
  state a query doesn't have. zp uses none of these. An app that does gets
  names that differ from sqlc's (ADR 0006).
- **The Go package.** The request doesn't name the Go package or its
  directory: the two generators hard-code `models`. The request's
  `plugin_options` holds only the JSON plugin's own `out` and `filename`.
- **Table and column labels, and views.** Not tried. zp's generators for
  them are described in [the brief](../../brief.md).
- **The enum variables' names.** The spike kept zp's `<Type><Value>Key`.
  neozp's convention for a label's variable is `L10n<GoName>`.
