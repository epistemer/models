package shop_test

import (
	"slices"
	"testing"

	"github.com/epistemer/models/gentest/shop"
	"github.com/epistemer/models/option"
)

// Every enum value is an option.SelectOption (models' interface §4.2).
var _ option.SelectOption = shop.OrderStatusPaid

func TestOptions(t *testing.T) {
	opts := shop.AllPaymentMethodOptions()
	var values []string
	for _, o := range opts {
		values = append(values, o.Value())
		if o.Label() == nil {
			t.Errorf("%q has no label", o.Value())
		}
	}
	want := []string{"credit-card", "bank/transfer", "cash:on:delivery", "pay-pal", "pay_pal", "!!", "value_5", "café", ""}
	if !slices.Equal(values, want) {
		t.Errorf("values %q, want %q", values, want)
	}
	opts[0] = shop.PaymentMethodCashOnDelivery
	if shop.AllPaymentMethodOptions()[0].Value() != "credit-card" {
		t.Error("a caller changed the next call's options")
	}
	if got := shop.AllNoValuesOptions(); got == nil || len(got) != 0 {
		t.Errorf("an enum with no values: %#v, want an empty slice", got)
	}
}
