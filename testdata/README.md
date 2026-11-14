# Fixtures

sqlc's codegen requests for two public toy schemas (roadmap M1.2). The
generators' tests read them as `sqlc.request/v1` files; nothing here comes
from zp.

| Directory | Schema | Request |
|-----------|--------|---------|
| [shop](shop/schema.sql) | Five enums (one with no values), four tables (three with the prefix `shop_`), a view, comments | `shop/request.json` |
| [plain](plain/schema.sql) | One table, no enums, no prefix | `plain/request.json` |
| [sqlc](sqlc/go.mod) | A nested module that pins sqlc v1.29.0 as a tool, so the root module doesn't depend on it | – |

Each request is 2.1 MB, most of it `pg_catalog` and `information_schema`,
and about 37 KB gzipped. The files are named `request.json` because
`.gitignore` keeps every `codegen_request.json` out.

## Regenerating

sqlc v1.29.0 parses SQL in WebAssembly, so it needs no cgo. From
`testdata/sqlc`:

```bash
go tool sqlc generate -f ../shop/sqlc.yaml && go tool sqlc generate -f ../plain/sqlc.yaml
```

CI's `fixtures` job does the same on Linux, macOS and Windows and fails if
a request changed, so sqlc writes the same bytes on all three. Two runs
give the same bytes too.

## What `shop` covers

Checked on 2026-11-14 against sqlc's own Go output for the schema, which
the spike's generators matched constant for constant (19 of 19) and which
passes `go vet` with them. With `table_prefix: shop_`:

Enum rows give sqlc's constant; the label's variable is `L10n` plus it.
Table rows give the label's key, variable, path or text, as named.

| Case | In the schema | Expected |
|------|---------------|----------|
| Plain values | `order_status`: `pending`, `paid`, `shipped` | `OrderStatusPending`, `OrderStatusPaid`, `OrderStatusShipped` |
| sqlc's replacement | `payment_method`: `credit-card`, `bank/transfer`, `cash:on:delivery` | `PaymentMethodCreditCard`, `…BankTransfer`, `…CashOnDelivery` |
| A duplicate after replacement | `pay-pal` (index 3), `pay_pal` (index 4) | `PaymentMethodPayPal`, `PaymentMethodValue4` |
| Nothing left after replacement | `!!` (index 5) | `PaymentMethodValue5` |
| A value that collides with a number | `value_5` (index 6) | `PaymentMethodValue6` |
| A letter outside ASCII in a value | `café` | `PaymentMethodCaf`; the key keeps `café` |
| An empty value | `''` (index 8) | `PaymentMethodValue8`; key `enum.payment_method.`, empty text |
| An enum with no values | `no_values` | type `NoValues`, no labels; its `Label()` always returns nil |
| The initialism `id` | type `id_kind`, values `id`, `user_id`, `uuid` | `IDKind`, `IDKindID`, `IDKindUserID`, `IDKindUuid` |
| A leading digit in a type | `2fa_method` | type `_2faMethod`, constant `_2faMethodTotp` |
| Digits inside names | values `u2f`, `level_3`; column `shop_order.address_line2` | `_2faMethodU2f`, `_2faMethodLevel3`; key `table.shop_order.Order_AddressLine2` |
| The prefix on tables | `shop_order_line` | key `table.shop_order_line`, variable `L10nOrderLine`, text `Order Line` |
| The prefix on a column | `shop_order.shop_ref` | key `table.shop_order.Order_Ref`, path `table.shop_order.shop_ref` |
| No prefix | `audit_log` | `L10nAuditLog`, `Audit Log` |
| zp's casing, not sqlc's | `shop_order.customer_id` | key `table.shop_order.Order_CustomerId`, variable `L10nOrder_CustomerId` (sqlc's field is `CustomerID`) |
| A table comment | `shop_customer` | text `Customer` |
| A column comment with spaces around it | `shop_customer.email` | text `Where receipts go` |
| A comment of white space only | `shop_order.shop_ref` | text `Ref` |
| A view, with a comment | `shop_order_summary` | key `table.shop_order_summary`, text `Orders at a glance` |

The refusals (another schema, two requests, duplicates, an empty name after
the prefix, a Go name starting with `L10n`, an enum name outside ASCII) are
tested on small requests written for each test (M2.1), not here.
