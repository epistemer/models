package gentest

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestPointerPackage: package shopptr is sqlc's Go for the shop with
// emit_pointers_for_null_types and overrides (sqlc.yaml's second block),
// and models' files for it are shop's with the package renamed. sqlc's
// request carries none of its Go options, so one set of generated files
// must serve both configurations (A5, ADR 0008); shopptr's tests then
// compile and run them against the other types. -update rewrites the
// copies.
func TestPointerPackage(t *testing.T) {
	files, err := filepath.Glob("shop/zz_shop_*.go")
	if err != nil || len(files) == 0 {
		t.Fatalf("no generated files in shop: %v", err)
	}
	for _, src := range files {
		b, err := os.ReadFile(src)
		if err != nil {
			t.Fatal(err)
		}
		want := bytes.Replace(b, []byte("\npackage shop\n"), []byte("\npackage shopptr\n"), 1)
		if bytes.Equal(want, b) {
			t.Fatalf("%s has no package clause to rename", src)
		}
		dst := filepath.Join("shopptr", strings.Replace(filepath.Base(src), "zz_shop_", "zz_shopptr_", 1))
		if *update {
			if err := os.WriteFile(dst, want, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		got, err := os.ReadFile(dst)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s isn't %s with the package renamed; run go test -update -run TestPointerPackage", dst, src)
		}
	}
}
