# ADR 0006: Interfaces models owns for UI packages (v0.3.0)

| | |
|---|---|
| Status | Accepted |
| Date | 2026-11-16 |
| Roadmap item | M6.1 |
| Deciders | epistemer decided the rule and D1–D4 on 2026-11-16 (D1–D3 as recommended; D4 epistemer chose to declare now). Claude decided D5 and D6 as owner |

## Context

epistemer's rule (2026-11-16): whatever zp's model layer got from a UI
package (tables, forms) is, in neozp, an interface that models owns. The
UI package implements or consumes it later, so the dependency points from
the UI to models, never the other way. Moving zp onto neozp is out of
neozp's scope; zp's code is only a record of what was needed.

From zp's code, its model layer depended on UI packages in three places:

- `types.SelectOption` (`Value_()`, `Label_()`): zp's enums implemented it,
  `All<Type>Options` listed them, and zp's form select input and table
  enum columns consumed it.
- `Builder`: the interface, declared in zp's models package, through which
  zp's table (`view.View`) shaped a view's query. models v0.2.0 took a
  `listing.Query` value instead (ADR 0005, D1).
- A link column's `Labeler func(id string) (string, bool)`: zp's table
  turned a row's id into the text shown for it.

Checked on sqlc v1.29.0: sqlc declares `Value() (driver.Value, error)` on
`Null<Type>` only, with either driver; the enum type itself has `Scan`,
`Valid` and nothing else, so `Value() string` is free on it.

## Decisions

**D1. `option.SelectOption`** (epistemer). A package
`github.com/epistemer/models/option`:
`type SelectOption interface { Value() string; Label() i18n.TKer }`. Not in
the root package: an app's own sqlc package is usually named `models`
too, which would force an alias in generated code. The enums generator
adds `Value()` (`string(t)`) to every enum type, which already has
`Label()`, so every enum value is a `SelectOption`.

**D2. `All<Type>Options()`** (epistemer). The enums generator writes
`func All<Type>Options() []option.SelectOption`, every value in the
schema's order, next to sqlc's `All<Type>Values()`. A function, so no
caller can change a shared slice.

**D3. `listing.Source`** (epistemer). `List<Type>` takes
`src listing.Source`, `type Source interface { ListingQuery() Query }`.
`Query` implements it itself, so a call that passes a `Query` value
doesn't change; a table package implements `Source` with its own state.
A nil `src` panics, as a nil `io.Reader` does.

**D4. `listing.Labeler`** (epistemer chose to declare it now):
`type Labeler interface { LabelOf(ctx context.Context, id string) (string, bool, error) }`,
the text that stands for a row, such as a customer's name for its id in a
link column; `false` when there is no such row. Apps implement it; a table
package calls it. models generates no implementation yet.

**D5. Dependencies** (Claude). `option` imports i18n for `TKer`, so the
models module requires i18n `v0.2.0` or later; apps that use models already
do. `listing` keeps no dependency.

**D6. Release** (Claude). D1–D3 change what an app gets (the enums file
gains two declarations, `List<Type>`'s parameter type changes), so this is
`v0.3.0`, interface amendment A3. An app that declared its own `Value()` or
`All<Type>Options` for an enum gets a compile error, listed as a limit.

## Consequences

- interface.md gets A3 (§4.2, §4.3, §4.4, a new §4.5, §6, §7, §9).
- The roadmap gets M6, ending in `v0.3.0`. Its backlog loses option types.
- When neozp's table and form packages start, they import
  `github.com/epistemer/models/option` and `.../listing`; models imports
  neither.
