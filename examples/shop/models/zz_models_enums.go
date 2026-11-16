// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import (
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
