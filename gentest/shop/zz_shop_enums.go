// Code generated from sqlc's request by github.com/epistemer/models. DO NOT EDIT.

package shop

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
	}
	return nil
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

// Label returns the label of t, or nil when t isn't a value of NoValues.
func (t NoValues) Label() i18n.TKer {
	switch t {
	}
	return nil
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
