# Status log

Newest first. Every working session adds an entry: what landed, the decisions
taken for you, what's next, and what's waiting on you.

## 2026-11-14: ownership, ADR 0001, the roadmap (M0)

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
- **Next:** M1.1, spike S1 against a copy of zp; M1.2, the fixture; M1.3,
  parameters in genx.
- **Waiting on you:** nothing yet. M1.3's amendment text will be.

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
