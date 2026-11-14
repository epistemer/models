package main

import (
	"testing"

	"github.com/epistemer/genx/genxtest"
	"github.com/epistemer/i18n"
	"github.com/epistemer/i18n/i18ntest"
	"github.com/epistemer/models/examples/shop/models"
)

// TestGenerated fails when a generated file is stale: sqlc's models and
// request (through sqlc diff), models' labels and enums, i18n's Go.
func TestGenerated(t *testing.T) { genxtest.AssertGenerated(t) }

// TestTranslated: a translation typed into the running app shows through
// the generated labels.
func TestTranslated(t *testing.T) {
	tr := i18ntest.New(t, i18n.Locale{Code: "en", NativeName: "English"})
	if err := tr.CreateLocale(i18n.Locale{Code: "fr", NativeName: "Français"}); err != nil {
		t.Fatal(err)
	}
	for key, text := range map[string]string{
		"enum.order_status.paid":             "payée",
		"table.shop_order":                   "Commande",
		"table.shop_customer.Customer_Email": "Adresse e-mail",
	} {
		if err := tr.Set("fr", key, text); err != nil {
			t.Fatal(err)
		}
	}
	fr := i18ntest.Context(t, tr, "fr")
	for _, c := range []struct{ got, want string }{
		{models.OrderStatusPaid.Label().T(fr), "payée"},
		{models.OrderStatusShipped.Label().T(fr), "shipped"}, // not translated: the code's text
		{models.L10nOrder.T(fr), "Commande"},
		{models.L10nOrder_Status.T(fr), "Status"},
		{models.L10n("table.shop_customer.email").T(fr), "Adresse e-mail"},
	} {
		if c.got != c.want {
			t.Errorf("got %q, want %q", c.got, c.want)
		}
	}
}
