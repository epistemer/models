// Package models holds neozp's generators for database models that sqlc
// generates: genx.meta.yaml and enums.go.tmpl, run by genx. It has no API at
// run time. An app switches the generators on in its genx.yaml and imports
// this package only blank, under the genx build tag, so the module stays
// required. The contract is docs/interface.md.
package models
