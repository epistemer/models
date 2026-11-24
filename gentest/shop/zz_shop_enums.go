// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package shop

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
	}
}

// Label returns the label of t, or nil when t isn't a value of PaymentMethod.
func (t PaymentMethod) Label() i18n.TKer {
	switch t {
	case PaymentMethodCreditCard:
		return L10nPaymentMethodCreditCard
	case PaymentMethodBankTransfer:
		return L10nPaymentMethodBankTransfer
	case PaymentMethodCashOnDelivery:
		return L10nPaymentMethodCashOnDelivery
	case PaymentMethodPayPal:
		return L10nPaymentMethodPayPal
	case PaymentMethodValue4:
		return L10nPaymentMethodValue4
	case PaymentMethodValue5:
		return L10nPaymentMethodValue5
	case PaymentMethodValue6:
		return L10nPaymentMethodValue6
	case PaymentMethodCaf:
		return L10nPaymentMethodCaf
	case PaymentMethodValue8:
		return L10nPaymentMethodValue8
	}
	return nil
}

// Value returns t as a string, for option.SelectOption.
func (t PaymentMethod) Value() string { return string(t) }

// MarshalText returns t's value, or an error when t isn't a value of
// PaymentMethod.
func (t PaymentMethod) MarshalText() ([]byte, error) {
	if t.Label() == nil {
		return nil, fmt.Errorf("models: %q is not a value of PaymentMethod", string(t))
	}
	return []byte(t), nil
}

// UnmarshalText sets t to text, or returns an error when text isn't a
// value of PaymentMethod.
func (t *PaymentMethod) UnmarshalText(text []byte) error {
	if PaymentMethod(text).Label() == nil {
		return fmt.Errorf("models: %q is not a value of PaymentMethod", string(text))
	}
	*t = PaymentMethod(text)
	return nil
}

// AllPaymentMethodOptions returns every value of PaymentMethod, in the
// schema's order, as select options.
func AllPaymentMethodOptions() []option.SelectOption {
	return []option.SelectOption{
		PaymentMethodCreditCard,
		PaymentMethodBankTransfer,
		PaymentMethodCashOnDelivery,
		PaymentMethodPayPal,
		PaymentMethodValue4,
		PaymentMethodValue5,
		PaymentMethodValue6,
		PaymentMethodCaf,
		PaymentMethodValue8,
	}
}

// Label returns the label of t, or nil when t isn't a value of IDKind.
func (t IDKind) Label() i18n.TKer {
	switch t {
	case IDKindID:
		return L10nIDKindID
	case IDKindUserID:
		return L10nIDKindUserID
	case IDKindUuid:
		return L10nIDKindUuid
	}
	return nil
}

// Value returns t as a string, for option.SelectOption.
func (t IDKind) Value() string { return string(t) }

// MarshalText returns t's value, or an error when t isn't a value of
// IDKind.
func (t IDKind) MarshalText() ([]byte, error) {
	if t.Label() == nil {
		return nil, fmt.Errorf("models: %q is not a value of IDKind", string(t))
	}
	return []byte(t), nil
}

// UnmarshalText sets t to text, or returns an error when text isn't a
// value of IDKind.
func (t *IDKind) UnmarshalText(text []byte) error {
	if IDKind(text).Label() == nil {
		return fmt.Errorf("models: %q is not a value of IDKind", string(text))
	}
	*t = IDKind(text)
	return nil
}

// AllIDKindOptions returns every value of IDKind, in the
// schema's order, as select options.
func AllIDKindOptions() []option.SelectOption {
	return []option.SelectOption{
		IDKindID,
		IDKindUserID,
		IDKindUuid,
	}
}

// Label returns the label of t, or nil when t isn't a value of NoValues.
func (t NoValues) Label() i18n.TKer {
	switch t {
	}
	return nil
}

// Value returns t as a string, for option.SelectOption.
func (t NoValues) Value() string { return string(t) }

// MarshalText returns t's value, or an error when t isn't a value of
// NoValues.
func (t NoValues) MarshalText() ([]byte, error) {
	if t.Label() == nil {
		return nil, fmt.Errorf("models: %q is not a value of NoValues", string(t))
	}
	return []byte(t), nil
}

// UnmarshalText sets t to text, or returns an error when text isn't a
// value of NoValues.
func (t *NoValues) UnmarshalText(text []byte) error {
	if NoValues(text).Label() == nil {
		return fmt.Errorf("models: %q is not a value of NoValues", string(text))
	}
	*t = NoValues(text)
	return nil
}

// AllNoValuesOptions returns every value of NoValues, in the
// schema's order, as select options.
func AllNoValuesOptions() []option.SelectOption {
	return []option.SelectOption{}
}

// Label returns the label of t, or nil when t isn't a value of _2faMethod.
func (t _2faMethod) Label() i18n.TKer {
	switch t {
	case _2faMethodTotp:
		return L10n_2faMethodTotp
	case _2faMethodSms:
		return L10n_2faMethodSms
	case _2faMethodU2f:
		return L10n_2faMethodU2f
	case _2faMethodLevel3:
		return L10n_2faMethodLevel3
	}
	return nil
}

// Value returns t as a string, for option.SelectOption.
func (t _2faMethod) Value() string { return string(t) }

// MarshalText returns t's value, or an error when t isn't a value of
// _2faMethod.
func (t _2faMethod) MarshalText() ([]byte, error) {
	if t.Label() == nil {
		return nil, fmt.Errorf("models: %q is not a value of _2faMethod", string(t))
	}
	return []byte(t), nil
}

// UnmarshalText sets t to text, or returns an error when text isn't a
// value of _2faMethod.
func (t *_2faMethod) UnmarshalText(text []byte) error {
	if _2faMethod(text).Label() == nil {
		return fmt.Errorf("models: %q is not a value of _2faMethod", string(text))
	}
	*t = _2faMethod(text)
	return nil
}

// All_2faMethodOptions returns every value of _2faMethod, in the
// schema's order, as select options.
func All_2faMethodOptions() []option.SelectOption {
	return []option.SelectOption{
		_2faMethodTotp,
		_2faMethodSms,
		_2faMethodU2f,
		_2faMethodLevel3,
	}
}
