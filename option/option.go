// Package option declares what a select input or a table's filter offers:
// a value and its label. models owns the interface so that neozp's form
// and table packages depend on models, not models on them. Every enum
// value of an app's models implements it (models' interface §4.5).
package option

import "github.com/epistemer/i18n"

// SelectOption is one choice of a select input or a table's filter: the
// value sent and stored, and its label for people.
type SelectOption interface {
	Value() string
	Label() i18n.TKer
}
