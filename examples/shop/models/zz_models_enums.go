// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
	"fmt"

	"github.com/epistemer/i18n"
	"github.com/epistemer/models/option"
)

// Label returns the label of t, or nil when t isn't a value of OrderStatus.
func (t OrderStatus) Label() i18n.TKer {
	switch t {
	case OrderStatusPending:
		return L10nOrderStatusPending
	case OrderStatusPaid:
		return L10nOrderStatusPaid
	case OrderStatusShipped:
		return L10nOrderStatusShipped
	case OrderStatusCancelled:
		return L10nOrderStatusCancelled
	}
	return nil
}

// Value returns t as a string, for option.SelectOption.
func (t OrderStatus) Value() string { return string(t) }

// MarshalText returns t's value, or an error when t isn't a value of
// OrderStatus.
func (t OrderStatus) MarshalText() ([]byte, error) {
	if t.Label() == nil {
		return nil, fmt.Errorf("models: %q is not a value of OrderStatus", string(t))
	}
	return []byte(t), nil
}

// UnmarshalText sets t to text, or returns an error when text isn't a
// value of OrderStatus.
func (t *OrderStatus) UnmarshalText(text []byte) error {
	if OrderStatus(text).Label() == nil {
		return fmt.Errorf("models: %q is not a value of OrderStatus", string(text))
	}
	*t = OrderStatus(text)
	return nil
}

// AllOrderStatusOptions returns every value of OrderStatus, in the
// schema's order, as select options.
func AllOrderStatusOptions() []option.SelectOption {
	return []option.SelectOption{
		OrderStatusPending,
		OrderStatusPaid,
		OrderStatusShipped,
		OrderStatusCancelled,
	}
}
