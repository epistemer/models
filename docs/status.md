# Status log

Newest first. Every working session adds an entry: what landed, the decisions
taken for you, what's next, and what's waiting on you.

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
