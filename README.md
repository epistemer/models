# models

neozp's generators for database models that [sqlc](https://sqlc.dev)
generates. sqlc writes the Go models; models adds a translatable label for
every enum value, table, view and column, and a `Label()` method on every
enum type. It is two trivial generators run by
[genx](https://github.com/epistemer/genx), with no program and no API at run
time; [i18n](https://github.com/epistemer/i18n) turns the labels into Go.

```go
status.Label().T(ctx)                              // "paid", or its translation
models.L10nOrder_Status.T(ctx)                     // a column's label
models.L10n("table.shop_order.status").T(ctx)      // the same, by SQL name
```

## Quickstart: an app on sqlc

You need genx v0.2.0 or later, i18n v0.2.0 or later and sqlc v1.29.0.
[examples/shop](examples/shop/README.md) is the result.

**1. Add the tools and the libraries.** In the app's module root:

```bash
go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0
go get -tool github.com/epistemer/genx/cmd/genx@latest
go get github.com/epistemer/models@latest github.com/epistemer/i18n@latest
```

**2. Have sqlc write its request.** Next to `gen: go:` in `sqlc.yaml`, a
`gen: json:` block writing into the same directory:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "db/schema.sql"
    queries: "db/query.sql"
    gen:
      go:
        package: "models"
        out: "models"
        sql_package: "pgx/v5"
        emit_all_enum_values: true
      json:
        out: "models"
        filename: "codegen_request.json"
```

**3. Write the manifest**, `genx.yaml`, at the module root:

```yaml
version: 1
producers:
  - name: sqlc
    run:   [go, tool, sqlc, generate, -f, sqlc.yaml]
    check: [go, tool, sqlc, diff, -f, sqlc.yaml]
    emits:
      - path: models/codegen_request.json
        domain: sqlc.request/v1
use:
  - github.com/epistemer/models:
      generators: [labels, enums]
      params:
        table_prefix: shop_   # if your tables carry one; else use [labels, enums]
  - github.com/epistemer/i18n: [labels]
```

`table_prefix` is removed from table and column names in variables,
English text and column keys, so set it before translations exist. If
sqlc's `out` isn't `models`, set `dir` (and `package`, if the directory's
name isn't the package's).

**4. Keep the generators required.** A file in any package, with that
package's name:

```go
//go:build genx

package main

import (
	_ "github.com/epistemer/i18n"
	_ "github.com/epistemer/models"
)
```

**5. Generate.** If your code already imports the models package and sqlc
hasn't written it yet, `go mod tidy` fails until it has; `-e` gets past it:

```bash
go mod tidy -e
go tool genx
go mod tidy
```

genx runs sqlc, then models' `labels` and `enums`, then i18n's `labels`,
and writes `models/zz_models.labels.json`, `models/zz_models_enums.go` and
`models/zz_models_i18n.go`. Commit them with `codegen_request.json` and
`genx.lock`. To regenerate from `go generate`, put
`//go:generate go tool genx` in a file of any package genx doesn't write
into, and not in the `//go:build genx` file, which `go generate` skips.

**6. Check in CI.** `go tool genx check` fails when anything is stale, the
request included (through `sqlc diff`). From a test, with
`github.com/epistemer/genx/genxtest`:

```go
func TestGenerated(t *testing.T) { genxtest.AssertGenerated(t) }
```

People then translate the labels in the running app; see i18n.

## What you get

| SQL | Key | Variable | Path |
|-----|-----|----------|------|
| enum `order_status`, value `paid` | `enum.order_status.paid` | `L10nOrderStatusPaid` | `enum.order_status.paid` |
| table `shop_order` | `table.shop_order` | `L10nOrder` | `table.shop_order` |
| its column `customer_id` | `table.shop_order.Order_CustomerId` | `L10nOrder_CustomerId` | `table.shop_order.customer_id` |

The English text is the SQL comment, else the name made readable. Keys are
permanent, and changing `table_prefix` changes column keys. Names that
would give a key two meanings, or a variable that can't compile, stop the
run with a message. [docs/interface.md](docs/interface.md) has every rule,
refusal and limit.

## Documentation

| Document | Answers |
|----------|---------|
| [docs/interface.md](docs/interface.md) | The contract: what an app writes and gets back (locked) |
| [docs/adr/](docs/adr/) | Why: the design (0001), the zp spike (0002), the test module (0003), zp on v0 (0004) |
| [docs/roadmap.md](docs/roadmap.md) | The plan, with a done-when for each item |
| [docs/status.md](docs/status.md) | The log of working sessions |
| [CHANGELOG.md](CHANGELOG.md) | Releases |

Part of neozp, alongside [cfgx](https://github.com/epistemer/cfgx),
[i18n](https://github.com/epistemer/i18n) and
[genx](https://github.com/epistemer/genx). License: [MIT](LICENSE).
