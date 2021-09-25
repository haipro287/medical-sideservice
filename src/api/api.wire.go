//go:build wireinject
// +build wireinject

package api

import (
	"context"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/sotanext-team/medical-chain/src/sideservice/src/services"
)

type ApiServiceOptions struct {
	Addr           string
	Logger         *logrus.Logger
	CosmosService *services.CosmosService
}

func InitApiService(ctx context.Context, opts ApiServiceOptions) (*Service, error) {
	wire.Build(
		wire.FieldsOf(&opts, "Addr", "Logger", "CosmosService"),
		NewIndexApi,
		NewCosmosApi,
		NewApiService,
	)
	return &Service{}, nil
}

//func InitIndexApi(ctx context.Context) (*IndexApi, error) {
//	wire.Build(
//		NewIndexApi,
//	)
//	return &IndexApi{}, nil
//}
//
