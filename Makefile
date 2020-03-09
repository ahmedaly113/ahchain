PACKAGES=$(shell GO111MODULE=on go list -mod=readonly ./...)

MODULES = argument backing category challenge expiration stake story

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=ahchaind \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

define \n


endef

benchmark:
	@go test -mod=readonly -bench=. $(PACKAGES)

buidl: build

build: build_cli build_daemon

download:
	go mod download

build_cli:
	@go build $(BUILD_FLAGS) -o bin/ahchaincli cmd/ahchaincli/*.go

build_daemon:
	@go build $(BUILD_FLAGS) -o bin/ahchaind cmd/ahchaind/*.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/ahchaind cmd/ahchaind/*.go
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/ahchaincli cmd/ahchaincli/*.go

doc:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/ahmedaly113/ahchain/"
	godoc -http=:6060

export:
	@bin/ahchaind export

create-wallet:
	bin/ahchaincli keys add validator --home ~/.octopus

init:
	rm -rf ~/.ahchaind
	bin/ahchaind init trunode $(shell bin/ahchaincli keys show validator -a --home ~/.octopus)
	bin/ahchaind add-genesis-account $(shell bin/ahchaincli keys show validator -a --home ~/.octopus) 10000000000utru
	bin/ahchaind gentx --name=validator --amount 10000000000utru --home-client ~/.octopus
	bin/ahchaind collect-gentxs

install:
	@go install $(BUILD_FLAGS) ./cmd/ahchaind
	@go install $(BUILD_FLAGS) ./cmd/ahchaincli
	@echo "Installed ahchaind and ahchaincli ..."
	@ahchaind version --long

reset:
	bin/ahchaind unsafe-reset-all

restart: build_daemon reset start

start:
	bin/ahchaind start --inv-check-period 10 --log_level "main:info,state:info,*:error,app:info,account:info,trubank:info,claim:info,community:info,truslashing:info,trustaking:info"

check:
	@echo "--> Running golangci"
	@golangci-lint run --tests=false --skip-files=\\btest_common.go

dep_graph: ; $(foreach dir, $(MODULES), godepgraph -s -novendor github.com/ahmedaly113/ahchain/x/$(dir) | dot -Tpng -o x/$(dir)/dep.png${\n})

install_tools_macos:
	brew install dep && brew upgrade dep
	brew install golangci/tap/golangci-lint
	brew upgrade golangci/tap/golangci-lint

go_test:
	@go test $(PACKAGES)

test: go_test

test_cover:
	@go test $(PACKAGES) -v -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic
	@go tool cover -html=coverage.txt

version:
	@bin/ahchaind version --long

########################################
### Local validator nodes using docker and docker-compose

build-docker-ahchaindnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/ahchaind/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/ahchaind:Z ahmedaly113/ahchaindnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

########################################

.PHONY: benchmark buidl build build_cli build_daemon check dep_graph test test_cover update_deps \
build-docker-ahchaindnode localnet-start localnet-stop
