# Review 0005: the generators (M2)

| | |
|---|---|
| Date | 2026-11-14 |
| Reviewer | Gemini 4 Hydrogen (High) through agy, adversarial review, about 79k tokens. Argon was cooling down |
| Under review | The diff without generated files, goldens and the `shop` package; interface §3–§6; the generated enums file; in a scratch folder |
| Triage | Claude, as owner |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Major | No `i18ntest` test in the diff | **Rejected.** It is `gentest/shop/shop_test.go`, which the reviewed diff left out with the generated `shop` files |
| 2 | Major | The failure wrapper returns from `Fatalf` and doesn't catch `Fatal` or `Error` | **Rejected.** It is genx's documented pattern (its generator guide), and `RenderParams` returns after `Fatalf`. Were genxtest to call something else, `msg` stays empty and the case fails with "no refusal", so nothing passes silently |
| 3 | Major | No test that parameters are checked without a request, that `package` rescues a `dir` like `my-models`, or that `information_schema` and `pg_temp` are skipped | **Accepted.** Three tests added |
| 4 | Minor | The `L10n` check read `él10n` as `L10n`, because `struct_name` turns `é` into `_`; sqlc keeps it | **Accepted.** The check uses `sqlc_start`, which keeps letters outside ASCII; a test covers `él10n` |
| 5 | Minor | The path check can't fire; a missing request array; two newlines at the template's end | **Rejected.** The path check stays as a guard for §6's "one path"; genx always passes an array (its §5); gofmt normalizes the end |
| 6 | Minor | `tools.go` never compiled; edge cases missing (`dir: ""`, `a/./b`, a schema with only an enum, an enum with no values compiled); assertions on prefixes only | **Accepted.** CI vets with `-tags genx`; the cases are added; `testdata/shop` gains `no_values`, so its `Label()` compiles in `gentest/shop`; the refusal messages are asserted in full. Non-string `dir` and `table_prefix` are genx's to refuse (A4) |
