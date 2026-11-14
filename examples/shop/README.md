# Example: a shop

An app whose models sqlc generates, with models' and i18n's generators run
by genx. It is interface §2 as an app writes it:

| File | What it is |
|------|------------|
| [sqlc.yaml](sqlc.yaml) | sqlc's config: Go models and the request as JSON, both into `models/` |
| [genx.yaml](genx.yaml) | sqlc as a producer, then `labels` and `enums` with `table_prefix: shop_`, then i18n's `labels` |
| [tools.go](tools.go) | The blank imports that keep models and i18n required |
| [main.go](main.go) | Prints labels; its `go:generate` runs genx |
| [main_test.go](main_test.go) | `genx check` from a test, and a translation shown through the labels |

```bash
go generate ./... && go run .
```

The module requires models through a `replace` to the repository's root,
so it always tests the current generators.
