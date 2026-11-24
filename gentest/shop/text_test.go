package shop_test

import (
	"encoding/json"
	"testing"

	"github.com/epistemer/models/gentest/shop"
)

func TestEnumTextMethods(t *testing.T) {
	// Every value of PaymentMethod, including "", "café", "!!", round-trips.
	opts := shop.AllPaymentMethodOptions()
	seen := make(map[string]bool, len(opts))
	for _, o := range opts {
		v := o.(shop.PaymentMethod)
		seen[v.Value()] = true
		b, err := v.MarshalText()
		if err != nil {
			t.Fatalf("%q.MarshalText(): %v", v, err)
		}
		if string(b) != v.Value() {
			t.Errorf("%q.MarshalText() = %q, want %q", v, b, v.Value())
		}
		var got shop.PaymentMethod
		if err := got.UnmarshalText(b); err != nil {
			t.Fatalf("UnmarshalText(%q): %v", b, err)
		}
		if got != v {
			t.Errorf("round-trip %q gave %q", v, got)
		}
	}
	for _, want := range []string{"", "café", "!!"} {
		if !seen[want] {
			t.Errorf("AllPaymentMethodOptions missing %q", want)
		}
	}

	// UnmarshalText refuses unknown values and wrong case, leaving receiver unchanged.
	for _, bad := range []string{"bogus", "Paid"} {
		s := shop.OrderStatusShipped
		err := s.UnmarshalText([]byte(bad))
		want := `models: "` + bad + `" is not a value of OrderStatus`
		if err == nil || err.Error() != want {
			t.Errorf("UnmarshalText(%q): err = %v, want %q", bad, err, want)
		}
		if s != shop.OrderStatusShipped {
			t.Errorf("UnmarshalText(%q) modified receiver to %q", bad, s)
		}
	}

	// MarshalText of an out-of-enum value errors.
	if _, err := shop.OrderStatus("bogus").MarshalText(); err == nil || err.Error() != `models: "bogus" is not a value of OrderStatus` {
		t.Errorf("MarshalText(bogus): err = %v", err)
	}

	// encoding/json uses the text methods.
	type payload struct {
		Status shop.OrderStatus `json:"status"`
	}
	raw, err := json.Marshal(payload{Status: shop.OrderStatusPaid})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"status":"paid"}` {
		t.Errorf("json.Marshal = %s, want %s", raw, `{"status":"paid"}`)
	}
	var s shop.OrderStatus
	if err := json.Unmarshal([]byte(`"bogus"`), &s); err == nil {
		t.Error(`json.Unmarshal("bogus") succeeded, want error`)
	}
}
