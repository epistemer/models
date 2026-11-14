// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package models

import "github.com/epistemer/i18n"

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
