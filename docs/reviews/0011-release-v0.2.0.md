# Review 0011: the release, v0.2.0 (M5.6)

| | |
|---|---|
| Date | 2026-11-15 |
| Reviewer | Gemini 3.8 Flash (High) through agy, adversarial review, about 178k tokens |
| Under review | The README and CHANGELOG diff, the example's diff and the `models` skill, against interface §4.3–§7, in a scratch folder |
| Triage | Claude, as owner |

| # | Severity | Finding | Outcome |
|---|----------|---------|---------|
| 1 | Critical | The skill's example names columns its view doesn't have, so it fails with `ErrQuery` | **Accepted.** The example names its view's columns |
| 2 | Major | The skill still says two generators and no runtime API | **Accepted** |
| 3 | Major | `type Q interface { Querier; Listings }` needs sqlc's `emit_interface: true` | **Accepted** in the README and the skill |
| 4 | Major | A non-nil value that encodes as NULL silently matches nothing; undocumented outside §4.4 | **Accepted** |
| 5 | Major | The README gives neither `listing`'s import path nor the `ErrQuery`/driver error split | **Accepted** |
| 6 | Major | Paging needs a `Sort` ending on a unique column; only §4.4 says so | **Accepted** |
| 7 | Minor | The skill's limits miss `omit_unused_structs` for views | **Accepted** |
| 8 | Minor | The skill's description of when the count runs is incomplete | **Accepted** |
| 9 | Minor | The search's leading `%` can't use a B-tree index | **Accepted** as a note in the skill |
| 10 | Minor | A stale comment in the example's test | **Accepted** |
