module github.com/ahmedaly113/ahchain

go 1.13

require (
	github.com/ahmedaly113/ahchain v0.3.3
	github.com/cosmos/cosmos-sdk v0.34.4-0.20191029195223-3099b42aa1a9
	github.com/gorilla/mux v1.7.3
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/tendermint v0.32.7
	github.com/tendermint/tm-db v0.2.0
	github.com/tendermint/tmlibs v0.9.0
)

replace github.com/cosmos/cosmos-sdk => github.com/ahmedaly113/cosmos-sdk v0.34.4-0.20191114003118-2268a8498fdd
