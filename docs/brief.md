# Brief for the session that starts this package

Written on 2026-11-14 by the session that released genx `v0.1.0`, to hand
over what it knew. This is a starting point, not a design: nothing here is
decided unless it says who decided it. Replace this file with the package's
own design records as they appear.

## 1. What this package is for

neozp is the rebuild of zp's framework as separate libraries (cfgx, i18n,
genx so far). This is the fourth: support for database models that
[sqlc](https://sqlc.dev) generates.

It is the first item of genx's backlog
([roadmap](https://github.com/epistemer/genx/blob/main/docs/roadmap.md),
"After v0.1.0"): *"neozp's models package on sqlc: enums and table labels
as trivial generators, S4 made real."* It is also the second user genx's
formats need before genx can tag `v1.0.0`, so what this package finds hard
is feedback for genx.

## 2. Where the repository stands

| | |
|---|---|
| Module | `github.com/epistemer/models`, Go 1.26, MIT. `doc.go` is a placeholder: the package name isn't decided (§6, question 7) |
| CI | The same 12 jobs as cfgx, i18n and genx: Go 1.26 and 1.27 on Linux, macOS and Windows; gofmt, vet, staticcheck, tidy, tests with the race detector |
| GitHub | `main` is the default branch. Merged branches are deleted automatically, as in the other three |
| Ownership | **Not delegated yet.** epistemer grants it in the session that starts the work. The terms used so far are genx's [ownership.md](https://github.com/epistemer/genx/blob/main/docs/ownership.md): copy it here when it's granted |
| Not there yet | A design record, a roadmap, an interface, any code, a Claude Code skill |

## 3. What zp does today

zp's code is a record of what the framework must do, not a design to copy.
zp is private: keep its schema, table names and enum names out of this
repository. Counts and shapes are fine.

zp runs sqlc v1.29.0 (PostgreSQL, `pgx/v5`) into a `models` package, and
three generators of its own write into the same package:

| Generator | Reads | Writes |
|-----------|-------|--------|
| `gen/enum`, a sqlc process plugin | The request's enums | Per enum type: an i18n key per value (`enum.<type>.<value>`, the value as its English text), `Value_() string` and `Label_() i18n.Ter` methods, and `All<Type>Options`, a list of the values as select options |
| `gen/meta`, a sqlc process plugin | The request's tables in the `public` schema | An i18n key per table (`table.<table>`) and per column (`table.<table>.<Table>_<Column>`). The variable and the English text drop the app's table prefix: `shop_order_line` becomes `L10nOrderLine`, "Order Line". A table's SQL comment names its label column and id columns, for option lists |
| `gen/view`, a program | `models.go`, parsed as Go | For each struct whose name ends in `View`: a `View<Name>` method that counts and pages a SQL view through squirrel, with the caller's filters |

Sizes, for scale: 28 enum types with 112 values; 89 tables with 794
columns; 13 views; 391 queries. sqlc's request for it is 5.6 MB of JSON.

## 4. What is already proven

Spike S4, recorded in genx's
[ADR 0006](https://github.com/epistemer/genx/blob/main/docs/adr/0006-foreign-file.md),
with its files and how to run it again in [spikes/s4](spikes/s4/README.md):

- sqlc writes its request as JSON with no database (`gen: json:`). genx
  takes it as a **foreign file**: the manifest names its domain,
  `sqlc.request/v1`.
- A trivial generator (a jq query and a template, no program) reproduces
  zp's `enums.go` exactly, 931 lines.
- A second one emits the enum values' labels as a `labels/v1` file, which
  i18n's `labels` generator turns into keys.
- Both run in 0.36 s over the 5.6 MB request.

Checked again on 2026-11-14 with genx `v0.1.0`.

## 5. What is already decided

By epistemer, in the standing instructions (`~/.claude/CLAUDE.md`):

- **Generators publish labels as data.** This package writes `labels/v1`
  files and never Go that imports i18n. It is used only by neozp, so it
  writes labels always, with no switch.
- **Keys:** `table.<table>[.<Field>]` and `enum.<type>.<value>`. Key names
  are permanent: renaming one orphans its translations.
- **English text:** the source's description where it has one, else the
  name made readable; an enum value is itself.
- **Each label** may carry a `var`, named `L10n<GoName>`, and a `path` for
  generic code. The producer resolves the names; a consumer never derives
  one. `L10n`, `l10n` and every `var` stay free in the package.
- **Never** generate translations or locale files.
- **Tests:** the labels files pass genx's conformance files,
  `genxtest.AssertGenerated` keeps generated files current, and `i18ntest`
  shows a translation set on a generated key.
- **Go style:** no struct tags for metadata; libraries take small
  interfaces of getters; no generics whose purpose isn't obvious at the
  call site.

By genx's locked interface: what a manifest, a metafile, a query and a
template may do
([interface.md](https://github.com/epistemer/genx/blob/main/docs/interface.md)).

## 6. Open questions, each with a default

These are API design, so they are epistemer's. The defaults are this
brief's suggestions and nobody has agreed to them.

1. **How does a generator learn an app's choices?** The first real
   problem. sqlc's request doesn't name the Go package, its directory, or
   the app's table prefix, and genx `v0.1.0` has no parameters for a
   generator (interface §12 lists them as "not in this version"). The spike
   hard-codes `models`. The options:
   - a genx amendment that adds parameters to `use`, designed with this
     package as its first user;
   - fixed conventions (the package is `models` in `models/`, no prefix);
   - a sqlc plugin shipped by this module, as zp does, which gets options
     from `sqlc.yaml` but makes the package a program again.

   *Default:* propose the genx amendment. Parameters are already on genx's
   backlog, and this package is the user that shows what they must carry.
2. **What does `v0.1.0` cover?** *Default:* enums (methods and option
   lists) and the labels of enum values, tables and columns. Views are a
   query builder, a different problem, and come later.
3. **Trivial generators or a producer program?** *Default:* trivial
   generators, as ADR 0006 decided. A program only if an app needs sqlc's
   `rename` or initialisms, which a query can't see.
4. **Do neozp's keys match zp's exactly?** Deployments of zp hold
   translations under zp's keys. *Default:* yes, the same keys, so moving
   zp over loses none. That makes the table prefix part of question 1.
5. **Variable names.** zp names an enum value's label `<Type><Value>Key`;
   the convention is `L10n<GoName>`. *Default:* `L10n<Type><Value>`, and
   zp's `L10n<Table>` and `L10n<Table>_<Column>` for tables.
6. **Where do the option types live?** zp's `Value_`, `Label_` and its
   select-option interface are in zp's `types` package. *Default:* decide
   with the package that renders forms; until then, generate the methods
   and return `i18n.Ter`.
7. **The package's name.** Every app's own package is called `models`, so
   a package of this module with that name needs an alias wherever both are
   imported. *Default:* the module's root holds only the metafile and its
   templates; anything an app imports at run time goes in a subpackage with
   another name.
8. **Committing sqlc's request.** genx renders from committed domain files,
   and `genx check` fails when a declared file is missing. That puts a
   5.6 MB JSON file in the app's repository. *Default:* commit it and
   measure; raise it with genx if it hurts.

## 7. A way to start

How cfgx, i18n and genx started, in order:

1. epistemer delegates ownership; `docs/ownership.md` records the terms.
2. A design record (ADR 0001) and a roadmap with a done-when for each item.
3. Spikes for what's unknown, each with a kill criterion. Here: table and
   column labels as a trivial generator, and question 1.
4. An interface draft. **Its freeze is epistemer's.**
5. Milestones, each reviewed by the other model family before it merges,
   each logged in [status.md](status.md).
6. After the freeze, a Claude Code skill in `~/.claude/skills/models/`.

## 8. Where to read

| What | Where |
|------|-------|
| How to use genx, cfgx and i18n | The `genx`, `cfgx` and `i18n` skills (`~/.claude/skills/`) |
| genx's contract, guides and failure list | `/home/epistemer/genx/docs/`: `interface.md`, `guides/generator.md`, `guides/producer.md`, `failures.md` |
| A complete trivial generator with tests | `/home/epistemer/genx/examples/labels` |
| A generator in production, with a nested test module | `/home/epistemer/i18n`: `genx.meta.yaml`, `labels.go.tmpl`, `labels/` |
| What genx learned from its first users | `/home/epistemer/genx/docs/friction.md` |
| zp's generators | `/home/epistemer/zp`: `gen/enum`, `gen/meta`, `gen/view`, `config/sqlc.yaml`, `generate.go`. Read only: never run a generator inside zp |
| How the other packages recorded their work | Each repository's `docs/`: `ownership.md`, `roadmap.md`, `status.md`, `adr/`, `reviews/` |
