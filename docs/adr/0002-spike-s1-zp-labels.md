# ADR 0002: Every zp label from a trivial generator (spike S1)

| | |
|---|---|
| Status | Accepted |
| Date | 2026-11-14 |
| Roadmap item | M1.1 |
| Deciders | Claude, as owner. Claude ran the spike on a copy of zp's schema and queries in a scratch folder |

## Context

ADR 0001 keeps zp's keys exactly (D4) and names variables as zp does (D6),
with jq doing the naming. Before the interface freezes, S1 had to show that
a query gets every zp key and variable right, not only the enums S4 did,
and that the generated Go compiles with sqlc's models. Its kill criterion:
if a key can't be matched in jq, a producer program writes the labels.

## What was tried

The files are in [spikes/s1](../spikes/s1/README.md). zp's `schema/` and
`sql/` were copied to a scratch module. sqlc v1.29.0 wrote its Go models
(zp's Go options without its type overrides, which point at zp's own
packages) and its request as JSON. genx `v0.1.0` then ran three generators:
S1's `labels` and `enums`, and i18n `v0.2.0`'s `labels`. The parameters
were hard-coded in `defs`, standing in for genx's.

## Results

| Check | Result |
|-------|--------|
| Keys | **995 of 995 equal** to the keys zp's `i18n.go` and `enums.go` declare: 112 enum values, 89 tables and views, 794 columns. None missing, none extra |
| Variables | Every table and column variable equals zp's. Every enum variable is `L10n` plus zp's name without `Key`, which is sqlc's constant |
| Uniqueness | No duplicate key, path or variable in zp's schema |
| English text | Equal for 994 labels. The one difference is zp's only table comment, a directive (D7, D10), which becomes the text |
| Compile | sqlc's models, the 28 `Label()` methods and i18n's generated Go pass `go vet` together |
| Time | A full run, sqlc included, about 3 s warm. The labels file is 200 KB; i18n's Go from it is 2,005 lines |

The audit ADR 0001 asked for:

- zp's `sqlc.yaml` sets no `rename` and no `initialisms`, so sqlc's default
  `[id]` applies, as D9 assumes. Its `overrides` map database and column
  types, never an enum.
- 13 views (15 `CREATE VIEW` statements across migrations); no
  materialized view. sqlc lists views with the tables, and zp's keys cover
  them.
- No column comments; one table comment, the directive above.
- The request is 5.6 MB, seven times sqlc's whole Go output for the same
  schema (about 0.8 MB). M3.2 measures what committing it costs (D12).

## What a query still can't see

Nothing new beyond ADR 0001's D9, which this narrows. sqlc drops non-ASCII
characters from enum values before naming constants, and the query does the
same, so values need no refusal. Table and column names need none either:
`zpPascal` upper-cases only ASCII, exactly as zp did, and every variable
starts with `L10n`, so it is exported whatever follows. Only an enum's own
name can't be matched, because it goes through sqlc's `StructName`, whose
`strings.Title` upper-cases any letter, and gojq has nothing but
`ascii_upcase`. An enum name with a letter outside ASCII stops the run.

## Decision

The kill criterion isn't met: trivial generators write every zp label. The
metafile in `spikes/s1` is the starting point for M2.1 and M2.2. The
refusal of non-ASCII names in D9 applies to enum names only.

The spike has none of ADR 0001's checks (unique keys, paths and variables,
an empty name after `strip`, Go names starting with `L10n`, the enum names
above, a second schema or request) and trims comments of ASCII white space
only. Those are M2.1's.

## Consequences

- M1.3 (parameters in genx) is the only unknown left before the
  generators; M2.1 adds the checks.
- zp's move keeps all 995 keys. Its code renames the enum `Key` variables
  (unused outside its models package), and its table comment goes.
