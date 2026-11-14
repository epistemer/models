# ADR 0001: Enums and labels from sqlc's request, as trivial generators

| | |
|---|---|
| Status | Accepted |
| Date | 2026-11-14 |
| Roadmap item | M0.2, reviewed in M0.3 ([review 0001](../reviews/0001-design.md)) |
| Deciders | epistemer decided D2, D4, D6, D10 and D13 on 2026-11-14 (the brief's questions 1, 2, 4, 5 and 6, each as recommended). Claude decided the rest as owner; each can be overturned |

## Context

neozp's fourth package supports database models that [sqlc](https://sqlc.dev)
generates. sqlc writes the Go models itself. What it doesn't write, and zp's
own generators do today, is:

- an i18n key for every enum value, table and column, so a form or an admin
  page can show a translated name;
- a way to get an enum value's label from the value.

zp does it with two sqlc process plugins and a third program for views
([brief](../brief.md), §3). genx's spike S4
([ADR 0006](https://github.com/epistemer/genx/blob/main/docs/adr/0006-foreign-file.md))
showed that sqlc's codegen request, written as JSON with no database
(`gen: json:`), can drive trivial generators, and that a jq query matches
sqlc's enum names when the app uses sqlc's naming defaults.

The standing rules of neozp apply: a generator publishes its labels as a
`labels/v1` file and never writes Go that registers i18n keys; keys are
permanent; labels carry a `var` named `L10n<GoName>` and a `path`.

Facts checked for this record, on sqlc v1.29.0 with a public toy schema:

- The request carries no Go option: not the package, its directory, `rename`
  or `initialisms`. `gen: json:` refuses keys of its own, so an app can't
  add any there.
- Views are listed with the tables, and a column's SQL comment is in the
  request. sqlc skips `pg_catalog` and `information_schema` and prefixes the
  Go names of enums from any schema other than the default one.
- An enum type's Go name is sqlc's `StructName(enum)`, and a value's
  constant is `StructName(enum + "_" + value)`, after replacing some
  characters and numbering duplicates (`value_<i>`). `StructName` upper-cases
  a part found in `initialisms`, which defaults to `["id"]`.
- A request from a one-table schema is 2.1 MB. About 0.9 MB of it is
  `pg_catalog` and `information_schema`, the same for every app on a given
  sqlc version.

## Decisions

**D1. Trivial generators only.** The module ships a metafile and templates
in its root package, and no program. sqlc's request is a foreign file the
app's manifest names as `sqlc.request/v1`. A producer program waits until an
app needs something a query can't see, such as `rename` (D9).

**D2. An app's choices come through genx parameters** (decided by
epistemer). genx `v0.1.0` can't pass anything to a generator, so models
needs a genx amendment: a metafile declares parameters with defaults, an app
sets them on its `use` entry, every query (`units` and `out`) reads them, and
an unknown parameter is an error. The generators check the values: `dir`
must be a clean relative path and `package` a Go identifier, so `dir: .`
needs an explicit `package`. This module is the
first user (roadmap M1.3, M1.4). The parameters:

| Parameter | Default | Meaning |
|-----------|---------|---------|
| `dir` | `models` | The Go package's directory, where sqlc writes its models |
| `package` | the last element of `dir` | The Go package's name |
| `table_prefix` | `""` | Stripped from table and column names in variables, English text and column keys |

Rejected: fixed conventions (they can't carry zp's prefix, which zp's keys
contain); a sqlc plugin shipped by this module (sqlc's `process` runs a
binary by path, so every app builds one, and the package becomes a program
again).

**D3. Two generators.** An app switches on both, and i18n's `labels`:

```yaml
use:
  - github.com/epistemer/models:
      generators: [labels, enums]   # the shape of use entries is M1.3's
      params: {table_prefix: app_}
  - github.com/epistemer/i18n: [labels]
```

| Generator | Writes | Holds |
|-----------|--------|-------|
| `labels` | `<dir>/zz_<package>.labels.json` (`labels/v1`) | A label for every enum value, table, view and column of the default schema |
| `enums` | `<dir>/zz_<package>_enums.go` | `func (t <Type>) Label() i18n.TKer` for every enum type |

i18n turns the labels file into `<dir>/zz_<package>_i18n.go`, in the same
package as sqlc's models.

**D4. Keys are zp's** (decided by epistemer), so moving zp over loses no
translation. With `strip` removing `table_prefix` from the front of a name
and `zpPascal` splitting on `_` and upper-casing the first letter of each
part:

| Label | Key |
|-------|-----|
| An enum value | `enum.<enum>.<value>` |
| A table or view | `table.<table>` |
| A column | `table.<table>.<T>_<C>`, where `<T>` is `zpPascal(strip(table))` and `<C>` is `zpPascal(strip(column))` |

`<enum>`, `<value>` and `<table>` are the names in SQL, prefix included.
`zpPascal` isn't sqlc's `StructName`: it knows no initialisms, so `user_id`
gives `UserId`. It is what zp's keys were made with, and keys are permanent.
zp strips the prefix from column names as well as table names, so D4 does
too.

`zpPascal` and `strip` aren't one-to-one: `user1` and `user_1` both give
`User1`, and with prefix `app_` the columns `app_name` and `name` both give
`Name`. So the `labels` generator checks that every key, path and variable
it writes is unique, and that no `strip` leaves an empty name, and stops the
run naming both sources when one isn't. Such a schema needs a column
renamed; it never gets a key that silently means two things.

**D5. Paths are the names in SQL:** `enum.<enum>.<value>`, `table.<table>`
and `table.<table>.<column>`. Generic code, such as an admin page, knows a
table and a column by those names, and i18n's generated `L10n(path)` finds
the label.

**D6. Variables** (decided by epistemer): `L10n<T>` for a table,
`L10n<T>_<C>` for a column, as zp names them (zp's code refers to about two
hundred), and `L10n<Const>` for an enum value, where `<Const>` is the name
of sqlc's constant for it. zp's `<Type><Value>Key` names go: nothing outside
zp's models package uses them. Any two labels with one variable stop the run
(D4). So does a table, view or enum whose name in Go starts with `L10n`,
which could clash with sqlc's own declarations in the package.

**D7. English text** is the SQL comment, trimmed, where that isn't empty, else
`strip(name)` made readable (`app_order_line` with prefix `app_` →
`Order Line`). An enum value's text is the value. zp ignores comments and
used its only table comment as a directive (D10); that comment goes when zp
moves over.

**D8. One schema.** Labels and enums come from the request's default schema
(`.catalog.default_schema`). `pg_catalog`, `information_schema` and
`pg_temp` are skipped. A table or enum in any other schema stops the run
with a message, so nothing goes without a label unnoticed. One request per
app, too: a second `sqlc.request/v1` file stops the run. Lifting either
limit needs its own key namespace for the other schema or database, designed
then; the default schema's keys stay as they are.

**D9. sqlc's naming defaults.** The generators name each enum type and
constant as sqlc v1.29.0 does with its default options: `StructName` with
`initialisms: [id]`, after sqlc's character replacement and its numbering
of duplicate values. interface.md (M1.5) spells the rules out, and the
fixture tests each one. Not followed, and documented as limits:

- `rename` and a non-default `initialisms`, which can change a type's or a
  constant's name. A type mismatch is a compile error naming the type in
  `zz_<package>_enums.go`; a constant mismatch only changes a variable's
  name, never a key. Parameters that mirror the options are on the
  backlog.
- An `overrides` entry that maps an enum's type to another Go type, which
  leaves `Label()` on a type sqlc doesn't declare: a compile error.
- Names with letters outside ASCII, where jq's upper-casing and Go's
  differ: the run stops. [ADR 0002](0002-spike-s1-zp-labels.md) narrows
  this to enum names: sqlc drops such letters from enum values itself, and
  table and column names follow zp's ASCII-only casing.

**D10. `Label()` and nothing else for enums** (decided by epistemer).
`Label()` returns the value's `L10n` variable, or nil for a value the enum
doesn't have. The generated file imports i18n for the return type only; the
keys still reach i18n through the labels file. zp's `Value_()` is
`string(t)`; its `All<Type>Options` lists give way to sqlc's
`All<Type>Values()` (`emit_all_enum_values`). A select-option type waits for
the package that renders forms. zp's table-comment directive for option
lists (a label column and id columns) is never used by zp's template, and is
dropped.

**D11. The root package holds the metafile.** `github.com/epistemer/models`
is `package models`, with the metafile, its templates and a `doc.go`. It
has no API at run time, so the clash with every app's own `models` package
is harmless: an app imports it only blank, under `//go:build genx`.

**D12. The request is committed.** genx renders from committed domain files.
M3.2 measures what that costs zp's repository; a problem goes to genx as
friction, not into a workaround here.

**D13. Not in v0.1.0** (decided by epistemer): views' query builder (zp's `gen/view`, a runtime
library on squirrel), option lists, labels for enum types themselves, and
schemas other than the default.

## Consequences

- models can't ship before genx has parameters: M1.3 and M1.4 are in genx,
  under its terms, and the amendment's exact text comes back to epistemer.
- Every key is a function of names in SQL and one parameter. Changing
  `table_prefix` changes column keys, so the interface calls it permanent
  once translations exist.
- The two generators share their naming functions through the metafile's
  `defs`.
- zp's move needs mechanical renames (`.Label_()` → `.Label()`,
  `All<Type>Options` → sqlc's values, the enum `Key` variables), no key
  changes, and its one table comment removed.
- Tests run without sqlc: a request made from a public toy schema is a
  committed fixture (`testdata/`), and CI regenerates it with sqlc pinned
  in a nested module, on all three OSes (sqlc v1.29.0 needs no cgo).
