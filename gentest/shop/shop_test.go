package shop_test

import (
	"testing"

	"github.com/epistemer/i18n"
	"github.com/epistemer/i18n/i18ntest"
	"github.com/epistemer/models/gentest/shop"
)

// TestTranslations sets translations on generated keys and reads them back
// through Label(), the variables and L10n(path) (roadmap M2.3).
func TestTranslations(t *testing.T) {
	tr := i18ntest.New(t, i18n.Locale{Code: "en", NativeName: "English"})
	if err := tr.CreateLocale(i18n.Locale{Code: "ar", NativeName: "Arabic", RTL: true}); err != nil {
		t.Fatal(err)
	}
	for key, text := range map[string]string{
		"enum.order_status.paid":      "مدفوع",
		"enum.payment_method.pay_pal": "باي بال",
		"table.shop_order_line":       "سطر الطلب",
		"table.shop_order_summary":    "ملخص الطلبات",
		"table.shop_order.Order_Ref":  "المرجع",
	} {
		if err := tr.Set("ar", key, text); err != nil {
			t.Fatalf("Set(%q): %v", key, err)
		}
	}
	ar := i18ntest.Context(t, tr, "ar")
	en := i18ntest.Context(t, tr, "en")

	for _, c := range []struct {
		name   string
		label  i18n.TKer
		key    string
		ar, en string
	}{
		{"enum value through Label", shop.OrderStatusPaid.Label(), "enum.order_status.paid", "مدفوع", "paid"},
		{"numbered duplicate", shop.PaymentMethodValue4.Label(), "enum.payment_method.pay_pal", "باي بال", "pay_pal"},
		{"the value it duplicates", shop.PaymentMethodPayPal.Label(), "enum.payment_method.pay-pal", "pay-pal", "pay-pal"},
		{"empty value", shop.PaymentMethodValue8.Label(), "enum.payment_method.", "", ""},
		{"table variable", shop.L10nOrderLine, "table.shop_order_line", "سطر الطلب", "Order Line"},
		{"view by path", shop.L10n("table.shop_order_summary"), "table.shop_order_summary", "ملخص الطلبات", "Orders at a glance"},
		{"column by path", shop.L10n("table.shop_order.shop_ref"), "table.shop_order.Order_Ref", "المرجع", "Ref"},
		{"column variable", shop.L10nOrder_Ref, "table.shop_order.Order_Ref", "المرجع", "Ref"},
		{"commented column", shop.L10nCustomer_Email, "table.shop_customer.Customer_Email", "Where receipts go", "Where receipts go"},
	} {
		if c.label == nil {
			t.Errorf("%s: nil label", c.name)
			continue
		}
		if got := c.label.K(); got != c.key {
			t.Errorf("%s: key %q, want %q", c.name, got, c.key)
		}
		if got := c.label.T(ar); got != c.ar {
			t.Errorf("%s: ar %q, want %q", c.name, got, c.ar)
		}
		if got := c.label.T(en); got != c.en {
			t.Errorf("%s: en %q, want %q", c.name, got, c.en)
		}
	}

	if l := shop.OrderStatus("refunded").Label(); l != nil {
		t.Errorf("Label of a value OrderStatus doesn't have: %v, want nil", l.K())
	}
	if l := shop.NoValues("x").Label(); l != nil {
		t.Errorf("Label of an enum with no values: %v, want nil", l.K())
	}
	if l := shop.L10n("table.no_such_table"); l != nil {
		t.Errorf("L10n of an unknown path: %v, want nil", l.K())
	}
}
