# Status log

Newest first. Every working session adds an entry: what landed, the decisions
taken for you, what's next, and what's waiting on you.

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
