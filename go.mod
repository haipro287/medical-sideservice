module github.com/sotanext-team/medical-chain/src/sideservice

go 1.16

require (
	cloud.google.com/go v0.81.0
	contrib.go.opencensus.io/exporter/stackdriver v0.13.8
	github.com/cosmos/cosmos-sdk v0.42.6
	github.com/ethereum/go-ethereum v1.10.9 // indirect
	github.com/gin-gonic/gin v1.7.4
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/sirupsen/logrus v1.8.1
	github.com/sonntuet1997/medical-chain-utils v0.1.9
	github.com/tendermint/spm v0.1.2
	github.com/urfave/cli/v2 v2.3.0
	go.opencensus.io v0.23.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/genproto v0.0.0-20210909211513-a8c4777a87af
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/cosmos/cosmos-sdk => github.com/haipro287/cosmos-sdk v0.42.10-0.20210924053243-5b97096b98c4

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
