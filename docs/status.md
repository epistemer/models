# Status log

Newest first. Every working session adds an entry: what landed, the decisions
taken for you, what's next, and what's waiting on you.

## 2026-11-24 — A5: Confab iteration 1 (v0.5.0, M7)

**Locked** by epistemer on 2026-11-24 with every default (A5.1a–c,
A5.3a–b), from Confab's findings F5, F8, F15, F19 and F28 and epistemer's
answers on the board. ADR 0008 for A5.1. **Released** as v0.5.0 (#16), checked through proxy.golang.org.

**Landed:**
- F2 (#15, docs only): `table_prefix`'s trade-off, that sqlc keeps the
  prefix in its Go names, in interface §3 and the README.
- A5.1: `column.Field`; every getter but an enum's calls it, so one views
  file serves sqlc's defaults and `emit_pointers_for_null_types` with
  overrides. `gentest/shopptr` is the shop's sqlc Go with pointers and
  zp-shaped overrides (decimal, uuid, time.Time, time.Duration, an app's
  int64 type). It reuses shop's generated files with the package renamed,
  and a PostgreSQL test reads the same rows through both configurations
  to identical values. `shop_kinds_view` gains `timestamptz` and
  `interval`.
- A5.2 text methods on enums, A5.3 the `pgx` generator with
  `RegisterTypes` (tested against PostgreSQL: `[]OrderStatus` fails before,
  works after), A5.4 `option.Static`, A5.5 enum text made readable.
- The example uses pointers, time overrides and `RegisterTypes`.

**Who did what:** M7.2 (`column.Field`, `option.Static`) by Gemini 4
Hydrogen (Medium), about 72k tokens; M7.3 (text methods, readable values,
`pgx`) by Gemini 4 Hydrogen (High), about 199k tokens. Both ran in
worktrees outside the repo, and Claude reviewed them. Argon was cooling
down until about 2026-11-28. Claude did M7.4 and M7.5.

**Review** ([0014](reviews/0014-a5.md)): Gemini 4 Hydrogen (High), about
96k tokens. Twelve findings; seven changed the code or docs:
- `enums` refuses a query constant that hides `fmt` and the file's other
  names;
- a Valuer's named basic result is read;
- a pointer to an interface is refused at init;
- pgtype fields are read without reflection per row.

These are recorded as clarifications in A5.

**Decided for you:** the board's M1–M5 are A5.1–A5.5 here, since M1–M7
name this repository's milestones. The `db_type` spellings sqlc v1.29.0
matches are documented (it silently ignores `pg_catalog.timestamptz` and
`timestamp`).

**For Confab** (not edited here):
- It can delete `internal/pg`, its hand-written `LoadTypes` in
  `AfterConnect` (use `models.RegisterTypes`), and `dataLabel` (use
  `option.Static`).
- It should switch its sqlc config to `emit_pointers_for_null_types: true`
  with overrides, including `date` and `timestamptz` to `time.Time`, which
  stay pgtype otherwise. Then add `pgx` to models' generators.

**After the tag:** v0.5.0's docs gave pgx v5.6.0 as `pgx`'s minimum;
pgx's changelog puts `LoadTypes` in v5.7.0, and v5.6.0 lacks it. Docs on
`main`, README, CHANGELOG and the skill now say v5.7.0, recorded in A5;
the tag's code is unaffected.

**Second review** ([0015](reviews/0015-a5-docs-and-fixes.md)): Gemini 4
Hydrogen (Medium), about 52k tokens, of what 0014 didn't see: F2, A5's
interface and ADR text, and the fix commit. Ten of its 13 findings were
valid: two edge bugs in `column.Field` and contract wording. **Released**
as v0.5.1 (#17), checked through proxy.golang.org.

**Waiting on you:** nothing.

## 2026-11-22 — A4: column sets for views (v0.4.0)

**Accepted** by epistemer with the defaults of M1–M5 (ADR 0007).
**Released** as v0.4.0 (#14), checked through proxy.golang.org.

**Landed** (Claude):
- package `column` (`Kind`, `Value`, `Column[R]`, `UUID`);
- `views` writes `<Type>Columns` beside each `List<Type>`. The jq maps
  each SQL type, nullable or not, to its kind and getter as sqlc's pgx/v5
  types it; a type no kind takes is left out and named.
- The fixture gains `shop_kinds_view`, a column of every kind, nullable
  and not, and an array and jsonb that are left out. A test reads every
  getter, set and null, with its label and options.
- interface §4.3 and the new §4.6, §2's needs, §6's two refusals, §7's
  two limits, the ADR's correction (a need, not a refusal), CHANGELOG.

**Review** ([0013](reviews/0013-column-sets.md)): Gemini 4 Argon
(Medium), about 27k tokens. Ten findings; six changed the code or docs:
- numerics read from their exact digits, NaN and ±Inf included;
- infinite dates and times marked, not read as year 1;
- the clashes of `column` and `<Type>Columns` refused;
- a type qualified with the default schema found.

Rejected, after checking sqlc v1.29.0's mapping: `citext` and `name` as
`interface{}`.

## 2026-11-16: interfaces for UI packages (M6, v0.3.0)

- **You said** (2026-11-16): moving zp onto neozp is out of neozp's scope,
  and what zp's model layer got from its table and form packages becomes
  interfaces models owns. The roadmap no longer plans or waits for zp's
  move; `v1.0.0` is simply your call. Both rules are in your global
  CLAUDE.md, and the skill's "Moving zp over" section is gone.
- **You decided:** `option.SelectOption` (`Value()`, `Label()`) in a new
  package; `All<Type>Options()` as a function; `List<Type>` takes a
  `listing.Source`, which `listing.Query` implements; `listing.Labeler`
  declared now. Recorded as ADR 0006 and amendment A3.
- **Landed:** package `option`; the enums file gains `Value()` and
  `All<Type>Options()`; `listing.Source` and `listing.Labeler`; a refusal
  for an option list named like an enum's type or constant; the example
  lists options; README, CHANGELOG and the skill.
- **Review 0012** by Gemini 4 Argon (Medium), about 46k tokens: seven
  findings, four accepted (two in part), three rejected. Its worry that
  `Labeler` makes a query per row was wrong, as you pointed out: a table
  calls it once per link filter, and rows get their text from the views'
  joins. The docs now say so.
- **Tagged:** `v0.3.0`.
- **Waiting on you:** nothing.

## 2026-11-15: v0.2.0, listings of views (M5.2–M5.6)

- **You approved** amendment A2; it is in interface.md (§4.3, §4.4, §6,
  §7, §11).
- **Landed:** package `listing` ([#10](https://github.com/epistemer/models/pull/10)),
  the `views` generator ([#11](https://github.com/epistemer/models/pull/11)),
  and the example's listing, the docs and the skill (M5.5, M5.6). A new CI
  job runs the generated methods against PostgreSQL 16 on Linux.
- **zp on a copy** (ADR 0004's addendum): its 13 views get the same 13
  types with the same columns, in the same order, as zp's `gen/view`;
  nothing refused; the models package compiles; warm run 1.9 s.
- **Reviews:** 0009 by Gemini 4 Argon (Medium), about 20k tokens (a typed
  nil got past the nil check); 0010 by Argon (Medium), about 80k tokens
  (sqlc makes `oxen_view` `OxView`; a query named `Max` hid a builtin);
  0011 by Gemini 3.8 Flash (High), about 178k tokens, on the docs and the
  skill (the skill's example used columns its view lacks; `Querier` needs
  `emit_interface`). The PostgreSQL test also caught a bug of mine: past
  the last page, the total was the offset.
- **Decided for you,** as clarifications under A2 (interface §11): a query
  named `Listings` isn't refused (it declares no type); the total's floor
  holds only when rows came back; `oxen` names follow sqlc; control
  characters are refused; a query constant hiding any name the views file
  uses is refused, not only `listing`. Each refuses only apps whose
  generated code couldn't compile. You can overturn them.
- **Tagged:** `v0.2.0` on 2026-11-16, after #12 merged; it resolves
  through the module proxy.
- **Next:** v1.0.0 still waits for zp's business code to run on models.
  For views, zp's `view.View` must produce a `listing.Query`, which is the
  future admin-table package's job. The backlog: `rename` and
  `initialisms` parameters, several requests or schemas, option types with
  the forms package, and `IN`/`IS NULL` filters if an app needs them.
- **Waiting on you:** nothing.

## 2026-11-15: listings of views, the design (M5.1)

- **Why now:** v0.1.0 finished the roadmap. zp's business code still
  needs one thing from models before it can move: a replacement for its
  view builder (`gen/view`, 38 files call it). The other backlog items
  wait for a forms package or aren't needed by zp (its `sqlc.yaml` sets no
  `rename` or `initialisms`, and it uses one schema).
- **You decided** (2026-11-15): the query is models' own `listing.Query`,
  a method returns `([]T, int64, error)`, views are the relations whose
  names end in `_view`, and pgx/v5 only.
- **Spike S2** (scratch folder, PostgreSQL 16 in a throwaway container):
  sqlc's request doesn't mark views; `_view` names aren't made singular;
  pgx v5 encodes `any` filter values by column type (enums, numeric,
  dates); `CAST(col AS text) ILIKE` with escaping works for every type.
- **Landed:** [ADR 0005](adr/0005-view-listings.md), roadmap M5, and the
  draft of amendment A2 (now [interface §11](interface.md#a2-listings-of-views-2026-11-15-for-v020)).
- **Review 0008** by Gemini 4 Argon (High), about 49k tokens: twelve
  findings, nine accepted (two in part). Three of its claims about sqlc
  were wrong when checked, though two still became a refusal and a limit.
  The real catches: the search's `OR` needed parentheses, or a filter
  would apply to the last search column only; a query named `Listing`
  hides the package; a `nil` filter value is now an error.
- **Waiting on you:** approving A2's text (M5.2): the names `views`,
  `List<Type>`, `Listings`, and package `listing`'s API.

## 2026-11-14: v0.1.0 (M4)

- **Landed:** the README with a quickstart (followed in an empty module,
  from `go mod init` to `go tool genx check`), the CHANGELOG, and the
  `models` skill in `~/.claude/skills/models/`.
- **Review 0007** of the release by Gemini 4 Hydrogen (High), about 161k
  tokens: seven findings, five accepted, one in part. Two changed code:
  `dir` is now checked against genx's path rules on every run, and two
  enums that sqlc names alike (`user1`, `user_1`) are refused.
- **Decided for you:** interface amendment A1, that refusal. sqlc's own Go
  breaks there too, so it changes no working app; listed in interface §11
  so you can overturn it.
- **Tagged:** `v0.1.0`.
- **Next:** v1.0.0 waits until zp's business code runs on models (the
  roadmap's rule). The backlog: sqlc's `rename` and `initialisms` as
  parameters, several requests or schemas, option types with the forms
  package, views' query builder.
- **Waiting on you:** nothing.

## 2026-11-14: the example (M3.1)

- **Landed:** [examples/shop](../examples/shop/README.md), interface §2 as
  an app writes it: sqlc and genx as its tools, `table_prefix: shop_`, a
  `go:generate` that runs genx, `genxtest.AssertGenerated`, and a French
  translation shown through `Label()`, a variable and `L10n(path)`. CI
  builds it, runs `go tool genx check` and its tests on all three OSes.
- **Found:** an app's first run needs `go mod tidy -e`, because `main.go`
  imports the models package before sqlc has written it. The quickstart
  (M4.1) says so.
- **Review 0006** by Gemini 3.8 Flash (High), about 278k tokens: seven
  findings, five accepted, one in part.
- **Next:** M4, the quickstart, the skill and `v0.1.0`.

## 2026-11-14: the generators (M2), and zp on them (M3.2)

- **Landed:** `genx.meta.yaml` and `enums.go.tmpl` at the module root, with
  every check of interface §6. [ADR 0003](adr/0003-test-module.md): a
  nested module, `gentest`, tests them: package `shop` is generated from
  the shop fixture by its own `genx.yaml` (sqlc through the pinned tool,
  then models and i18n) and checked by `genxtest.AssertGenerated`; golden
  files cover the defaults, other `dir` and `package` values and the schema
  without enums, each rendered twice; 26 refusal cases run against both
  generators; `i18ntest` reads translations back through `Label()`, the
  variables and `L10n(path)`.
- **The fixture** gained an enum with no values, so its `Label()` is
  compiled too.
- **Review 0005** by Gemini 4 Hydrogen (High), about 79k tokens: six
  findings, three accepted. The real catch: the `L10n` check would have
  refused a table like `él10n`, because it read letters outside ASCII the
  way `struct_name` does, not the way sqlc does.
- **M3.2, zp on a copy** ([ADR 0004](adr/0004-zp-on-models.md)): the
  merged generators, with zp's prefix as a parameter, give all 995 of zp's
  keys and variables, refuse nothing, and compile. A warm run takes 1.9 s,
  of which sqlc is 1.2 s. Committing the 5.6 MB request costs about 344 KB
  of git objects, and a one-column change about 4 KB, so it stays
  committed.
- **Next:** M3.1, an example app; then M4, the quickstart, the skill and
  `v0.1.0`.
- **Waiting on you:** nothing.

## 2026-11-14: the interface, drafted and locked (M1.5)

- **Landed:** [interface.md](interface.md), the draft for your freeze: what
  an app writes (§2), the three parameters (§3), the two generators and
  their files (§4), every naming rule with sqlc's numbering spelled out
  (§5), the refusals (§6), the limits (§7), and what a later change to it
  would mean (§10).
- **Review 0004** by Gemini 4 Hydrogen (High), about 55k tokens (Argon was
  cooling down): fifteen findings, twelve accepted (three in part). The
  main changes: refusals for names that can't become Go identifiers and for
  a missing default schema; `package` checked by an exact pattern; sqlc's
  duplicate numbering stated as an algorithm (checked: sqlc itself breaks
  on `value_1` then `!!`); white space listed by character; §7 corrected
  row by row.
- **Decided for you, in the draft:** paths use SQL names (ADR 0001 D5);
  both generators make every check; an enum with no values gets a `Label()`
  that returns nil; message wording isn't part of the contract.
- **Locked by you** on 2026-11-14, as drafted. M1 is done.
- **Next:** M2, the generators.

## 2026-11-14: the fixture (M1.2)

- **Landed:** [testdata](../testdata/README.md): two public toy schemas
  (`shop`, which covers every naming case of ADR 0001's D4, D7 and D9, and
  `plain`, with no enums), the requests sqlc v1.29.0 makes from them, and a
  nested module pinning sqlc as a `tool`. A new CI job regenerates both and
  fails on any change, on all three OSes.
- **Checked:** sqlc's own Go for `shop` and the spike's generators agree on
  all 19 enum constants (numbered duplicates, an empty value, `café`, the
  `id` initialism, a leading digit), and `go vet` passes with i18n's
  generated Go.
- **Found:** sqlc v1.29.0 parses SQL in WebAssembly, so it needs no cgo,
  which ADR 0001 had assumed. Each request is 2.1 MB but about 37 KB
  gzipped, which is roughly what git stores.
- **Review 0003** by Gemini 3.8 Flash (High), about 260k tokens: five
  findings, three accepted.
- **Next:** M1.5, the interface draft for your freeze.
- **Waiting on you:** nothing yet; the freeze will be.

## 2026-11-14: genx v0.2.0 has parameters (M1.3, M1.4)

- **Landed in genx:** amendment A4 as you approved it
  ([genx PR #15](https://github.com/epistemer/genx/pull/15)), tagged
  `v0.2.0`. A metafile declares `params:`, an app sets them in the map form
  of its `use` entry, and every query reads them as `$params`;
  `genxtest.RenderParams` tests with other values. genx's e2e script
  `params.txtar` runs this package's shape (`dir`, `package` defaulting to
  null, `table_prefix`) end to end.
- **Decided for you, in genx:** values belong to the `use` entry; `$params`
  is `{}` without declarations; values aren't checked beyond their JSON
  type. genx's status log has the rest.
- **Next:** M2 on genx `v0.2.0`: the spike's `def params:` becomes the
  metafile's `params:` and `$params`.
- **Waiting on you:** nothing.

## 2026-11-14: ownership, ADR 0001, the roadmap, spike S1 (M0, M1.1)

- **Landed:**
  - Ownership delegated to Claude until `v1.0.0`, on genx's terms
    ([ownership.md](ownership.md)).
  - [ADR 0001](adr/0001-models-on-sqlc.md): two trivial generators,
    `labels` (a `labels/v1` file for enum values, tables, views and
    columns) and `enums` (`Label()` per enum type), over sqlc's request,
    with parameters from genx.
  - [roadmap.md](roadmap.md): M0 to M4 and `v0.1.0`.
  - [Review 0001](reviews/0001-design.md) of both, by Gemini 4 Hydrogen
    (High), about 59k tokens (Argon was cooling down for 21 hours): twelve
    findings, ten accepted (two of them in part), two rejected. The
    biggest changes: uniqueness checks on every key, path and variable,
    the fixture moved before the freeze, and S1 widened to every label.
- **Decided by you (2026-11-14), each as recommended:** parameters in genx
  `use` entries (D2); zp's keys exactly, and `L10n` variables (D4, D6);
  `Label() i18n.TKer` as the only enum method (D10); enums and labels in
  `v0.1.0`, the view builder later (D13).
- **Decided for you:** everything else in ADR 0001: paths are the names in
  SQL (D5); comments are the English text (D7); one schema and one request
  (D8); sqlc's naming defaults only (D9); the root package holds only the
  metafile (D11); the request is committed (D12).
- **M1.1, spike S1** ([ADR 0002](adr/0002-spike-s1-zp-labels.md), PR #2):
  on a copy of zp's schema in a scratch folder, all 995 of zp's keys and
  every variable came out of the spike's trivial generators, and the
  `Label()` methods compile with sqlc's models and i18n's Go. Review 0002
  by Gemini 3.8 Flash (High), about 247k tokens: seven findings, five
  accepted. One corrected me: every variable starts with `L10n`, so only
  enum names need the non-ASCII refusal.
- **Decided by you for genx amendment A4 (M1.3), each as recommended:**
  parameters are declared per metafile (`params:` with `doc` and
  `default`); an app sets them in a map form of its `use` entry
  (`{generators, params}`), the list form staying valid; an unknown
  parameter, or a value whose JSON type differs from a non-null default's,
  is an error; queries (`defs`, `units`, `out`) read `$params`, templates
  see only the unit; `genx list` shows the values; genxtest gains
  `RenderParams(t, dir, generator, params map[string]any, inputs ...Input)`.
- **Next:** M1.3, A4 in genx, in its own session; then M1.4 (genx
  `v0.2.0`). M1.2, the fixture, can go alongside.
- **Waiting on you:** nothing until the freeze (M1.5).

## 2026-11-14: the repository is prepared; no work has started

Done by the session that released genx `v0.1.0`, at epistemer's request, so
another session can start the package.

- **Landed:**
  - The module, the MIT license, and CI with the same 12 jobs as cfgx, i18n
    and genx.
  - [brief.md](brief.md): what the package is for, what zp does today, what
    is proven, what is already decided, and eight open questions, each with
    a default.
  - [spikes/s4](spikes/s4/README.md): the two trivial generators of genx's
    spike S4, rebuilt and checked again with genx `v0.1.0`. `enums`
    reproduces zp's `enums.go` exactly; `enum-labels` writes its 112 labels.
  - On GitHub, merged branches are deleted automatically.
- **Decided for you:** nothing. The brief's defaults are suggestions.
- **Next:** the first working session. The brief's §7 has the order the
  other packages followed.
- **Waiting on you:**
  - Delegating ownership, in the session that starts the work.
  - The brief's open questions. The first one decides the rest: how a
    generator learns an app's Go package and table prefix, which genx
    `v0.1.0` can't pass to it.
