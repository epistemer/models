# ADR 0007: a column set for every view

| | |
|---|---|
| Status | Accepted by epistemer on 2026-11-22, with the defaults of M1–M5 |
| Date | 2026-11-22 |
| Roadmap item | ui D7.4; models v0.4.0 |
| Deciders | epistemer decided the shape (ui's S2 and S3, 2026-11-22) and decides this text; Claude proposes it, from ui's design of tables (server-state-design.md §6.1, §10) and its spike (§12.1) |

## Context

ui's tables are declared from columns that models generates, so no column
name is typed as a string in an app (ui's principle 4). epistemer decided
on 2026-11-22:

- S2: models generates a column set per `_view` relation;
- S3: a column is its name, its label, a kind from its SQL type, an
  enum's options, and a getter returning a plain value with nulls
  unwrapped, so ui never imports the driver;
- S22: what a page decides about a column (searchable, sortable, a
  filter) lives in ui, not on models' column.

ui's spike (D7.2) declared exactly that shape locally and ran a table of
orders on it. One finding bears on models: a column carries its options,
not its Go enum type, so a filter's value reaches `listing` as the
option's string. §4.3 already says a filter on an enum column takes the
`Value()` string.

## Options

1. **A package `column` and a generated set per view** (this proposal).
2. **Reflection over sqlc's row struct in ui.** No generation, but ui
   then knows the driver's null types and the SQL names only by tags.
3. **Columns written by hand in each app.** What principle 4 rules out.

## Decision

Option 1. A new runtime package, `github.com/epistemer/models/column`,
depending on `option` and i18n as `option` does:

```go
package column

// Kind is a column's kind, from its SQL type.
type Kind uint8

const (
	Text    Kind = iota + 1 // text, varchar, char, citext, uuid
	Integer                 // int2, int4, int8
	Number                  // float4, float8, numeric
	Bool                    // bool
	Date                    // date
	Time                    // timestamp, timestamptz
	Enum                    // an enum of the schema
)

// Value is a cell's plain value: the field of its kind, unless Null.
type Value struct {
	Null bool
	Text string             // Text; a numeric's exact digits too
	Int  int64              // Integer
	Num  float64            // Number
	Bool bool               // Bool
	Time time.Time          // Date (UTC midnight), Time
	Enum option.SelectOption // Enum
}

// Column is one column of a relation whose rows are R.
type Column[R any] struct {
	Name    string                // its SQL name
	Label   i18n.TKer             // table.<relation>.<Field>, the label models already generates
	Kind    Kind
	Options []option.SelectOption // an enum's, from All<Type>Options
	Get     func(R) Value
}
```

The `views` generator writes, beside each `List<Type>`:

```go
// OrderSummaryViewColumns are order_summary_view's columns.
var OrderSummaryViewColumns = struct {
	ID, Status, Total, Placed, Note column.Column[OrderSummaryView]
}{
	ID:     column.Column[OrderSummaryView]{Name: "id", Label: L10nOrderSummaryView_Id, Kind: column.Integer, Get: func(r OrderSummaryView) column.Value { return column.Value{Int: r.ID} }},
	Status: column.Column[OrderSummaryView]{Name: "status", Label: …, Kind: column.Enum, Options: AllOrderStatusOptions(), Get: …},
	Note:   column.Column[OrderSummaryView]{Name: "note", …, Get: func(r OrderSummaryView) column.Value {
		if !r.Note.Valid { return column.Value{Null: true} }
		return column.Value{Text: r.Note.String}
	}},
	// …
}
```

- A column whose SQL type is none of the kinds above (json, bytea, an
  array, an interval) is left out of the set, with a comment naming it.
  It stays in the row and in `List<Type>`.
- `numeric` is `Number`, with `Num` and its exact digits in `Text`, since
  ui imports no decimal library.
- A column can be built by hand, for a test or rows that don't come from
  sqlc.

## Consequences

- A new package and new declarations in the views file: a minor version,
  `v0.4.0`, under v0's rules. An app on v0.3.0 changes nothing.
- The getters unwrap pgx's null types in generated code, so ui depends on
  `column` and `option` only.
- `IN`, `IS NULL` and "contains" in `listing`, which tables will want, are
  each an amendment of their own, later (ui's design §6.5).
- The labels file is unchanged: the columns reuse the labels' variables,
  so `views` needs `labels` and i18n's `labels`, and `enums` for a view
  with an enum column, as `enums` needs `labels` (interface §2). (The
  proposal said `views` would refuse without them; a generator can't see
  which others an app runs, so it is a need, as for `enums`: without them
  the views file doesn't compile.)
- `views` no longer works alone, which it did in v0.3.0. Under v0's rules,
  a minor version.

## Questions, answered with their defaults

| # | Question | Default |
|---|---|---|
| M1 | The package's name | `column` |
| M2 | The set's name | `<Type>Columns`, a struct of named fields |
| M3 | Types outside the kinds | left out of the set, named in a comment |
| M4 | `numeric` | `Number`, with the exact digits in `Text` too |
| M5 | `uuid` | `Text`, its canonical string |
