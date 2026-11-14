// Package gentest tests models' generators (genx.meta.yaml and
// enums.go.tmpl in the module's root directory) through genx. It is a module
// of its own, so models itself requires nothing. Package shop is a small app:
// sqlc writes its models from testdata/shop, and genx runs models' and
// i18n's generators over sqlc's request (genx.yaml).
package gentest
