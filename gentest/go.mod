module github.com/epistemer/models/gentest

go 1.26.0

replace github.com/epistemer/models => ../

require (
	github.com/epistemer/genx v0.2.0
	github.com/epistemer/i18n v0.2.0
	github.com/epistemer/models v0.0.0-20261114012404-fbeb5e9d81be
)

require (
	github.com/itchyny/gojq v0.12.20 // indirect
	github.com/itchyny/timefmt-go v0.1.9 // indirect
	go.etcd.io/bbolt v1.5.0 // indirect
	go.yaml.in/yaml/v4 v4.0.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

tool github.com/epistemer/genx/cmd/genx
