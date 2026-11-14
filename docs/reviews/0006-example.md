# Review 0006: the example (M3.1)

| | |
|---|---|
| Date | 2026-11-14 |
| Reviewer | Gemini 3.8 Flash (High) through agy, adversarial review, about 278k tokens |
| Under review | The diff without generated files and `go.sum`, and interface §2, in a scratch folder |
| Triage | Claude, as owner |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Major | The example's tests check few labels, and `AssertGenerated` checks freshness, not correctness | **In part.** Correctness is gentest's job (goldens, 26 refusals, i18ntest); the example shows the workflow. One more assertion, on the column variable `main.go` prints |
| 2 | Major | A path (`table.shop_customer.email`) finds a label whose key differs; `L10n` returns nil on a typo | **Rejected.** Both are by design: keys are zp's, paths are SQL names (ADR 0001, D4 and D5), and `L10n`'s nil is i18n's contract |
| 3 | Major | CI never builds `main` | **Accepted.** The nested-module step runs `go build ./...` |
| 4 | Major | `sqlc.yaml` leaves out `sql_package: "pgx/v5"`, unlike interface §2 | **Accepted.** The example now matches §2 |
| 5 | Minor | `go build` leaves an untracked `shop` binary | **Accepted.** Ignored |
| 6 | Minor | sqlc compiles on six runners with no cache for nested modules | **Accepted.** `setup-go` caches on every `go.sum` |
| 7 | Minor | §2 says CI runs `go tool genx check`, CI only runs `AssertGenerated` | **Accepted.** CI runs both, in gentest and the example |
