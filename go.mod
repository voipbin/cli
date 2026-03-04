module github.com/voipbin/vn-cli

go 1.23.2

require (
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/cobra v1.10.2
	github.com/voipbin/voipbin-go v0.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/oapi-codegen/runtime v1.1.2 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)

replace github.com/voipbin/voipbin-go => ../voipbin-go
