package option_test

import (
	"context"
	"testing"

	"github.com/epistemer/models/option"
)

func TestStatic(t *testing.T) {
	opt := option.Static("ev-1", "Annual Gala")
	if got := opt.Value(); got != "ev-1" {
		t.Errorf("Value() = %q, want %q", got, "ev-1")
	}
	lbl := opt.Label()
	if got := lbl.K(); got != "" {
		t.Errorf("Label().K() = %q, want empty", got)
	}
	if got := lbl.T(context.Background(), "ignored", 42); got != "Annual Gala" {
		t.Errorf("Label().T(...) = %q, want %q", got, "Annual Gala")
	}
}
