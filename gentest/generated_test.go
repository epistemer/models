package gentest

import (
	"testing"

	"github.com/epistemer/genx/genxtest"
)

// TestGenerated fails when anything genx.yaml generates is stale, sqlc's
// request included (through sqlc diff).
func TestGenerated(t *testing.T) { genxtest.AssertGenerated(t) }
