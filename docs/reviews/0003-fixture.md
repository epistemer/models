# Review 0003: the fixture (M1.2)

| | |
|---|---|
| Date | 2026-11-14 |
| Reviewer | Gemini 3.8 Flash (High) through agy, adversarial review, about 260k tokens |
| Under review | The PR's diff without the two requests and `go.sum`, plus the `public` schema of `shop/request.json`, in a scratch folder |
| Triage | Claude, as owner |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Critical | `go.sum` and the requests aren't committed, so CI would fail on a clean checkout | **Rejected.** They are committed; they were left out of the reviewed diff for size |
| 2 | Major | The roadmap promises an empty enum value, and `'!!'` (empty after replacement) isn't one | **Accepted.** `payment_method` gains `''`; sqlc names it `PaymentMethodValue8`, and the spike's generators still match all 19 constants |
| 3 | Major | `git status -- ..` also blames `go.mod` changes on the fixtures, and `git diff` hides untracked files | **Accepted.** The check runs from the repository root on named paths and prints `git status`'s list |
| 4 | Minor | README rows: `order_status` missing, a leading-digit row listing names without one, a column row with no key or variable | **Accepted.** Fixed |
| 5 | Minor | Three OSes are redundant: sqlc runs in WebAssembly and writes LF | **Rejected.** The repository is public, so the minutes are free, and the job turns "the same bytes on every OS" from a claim into a check, which apps whose `genx check` runs on Windows rely on |
