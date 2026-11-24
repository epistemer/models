// Package option declares what a select input or a table's filter offers:
// a value and its label. models owns the interface so that neozp's form
// and table packages depend on models, not models on them. Every enum
// value of an app's models implements it (models' interface §4.5).
package option

import (
	"context"

	"github.com/epistemer/i18n"
)

// SelectOption is one choice of a select input or a table's filter: the
// value sent and stored, and its label for people.
type SelectOption interface {
	Value() string
	Label() i18n.TKer
}

// Static returns an option whose label is text from data, such as an
// event's name: its Label renders text in every locale, and its K is ""
// (A5).
func Static(value, text string) SelectOption {
	return staticOption{value: value, label: staticLabel(text)}
}

type staticOption struct {
	value string
	label staticLabel
}

func (o staticOption) Value() string    { return o.value }
func (o staticOption) Label() i18n.TKer { return o.label }

type staticLabel string

func (s staticLabel) T(_ context.Context, _ ...any) string { return string(s) }
func (s staticLabel) K() string                            { return "" }
