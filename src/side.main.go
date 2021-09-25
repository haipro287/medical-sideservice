package main

import (
	"cloud.google.com/go/profiler"
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/sirupsen/logrus"
	"github.com/sonntuet1997/medical-chain-utils/common"
	"github.com/sonntuet1997/medical-chain-utils/common_service"
	pb2 "github.com/sonntuet1997/medical-chain-utils/common_service/pb"
	api2 "github.com/sotanext-team/medical-chain/src/sideservice/src/api"
	"github.com/sotanext-team/medical-chain/src/sideservice/src/services"
	"github.com/tendermint/spm/cosmoscmd"
	"github.com/urfave/cli/v2"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var logger *logrus.Logger

func runMain(appCtx *cli.Context) error {
	logger = common.InitLogger(appCtx)

	if appCtx.Bool("disable-tracing") {
		logger.Info("Tracing disabled.")
	} else {
		logger.Info("Tracing enabled. But not implement")
	}
	if appCtx.Bool("disable-profiler") {
		logger.Info("Profiling disabled.")
	} else {
		logger.Info("Profiling enabled. But not implement")
		//go initProfiling(serviceName, appCtx.String("runtime-version"))
	}

	var wg sync.WaitGroup

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	// Start gRPC server
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", appCtx.Int("grpc-port")))
	if err != nil {
		return err
	}
	defer func() { _ = grpcListener.Close() }()

	var srv *grpc.Server
	log.Printf("%v", appCtx.Bool("allow-kill"))
	commonServer := common_service.NewCommonServiceServer(logger, appCtx.Bool("allow-kill"))
	wg.Add(1)
	go func() {
		defer wg.Done()
		if appCtx.Bool("disable-stats") {
			logger.Info("Stats disabled.")
			srv = grpc.NewServer()
		} else {
			logger.Info("Stats enabled.")
			srv = grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
		}
		healthpb.RegisterHealthServer(srv, commonServer)
		pb2.RegisterCommonServiceServer(srv, commonServer)
		reflection.Register(srv)
		logger.WithField("port", appCtx.Int("grpc-port")).Info("listening for gRPC connections")
		if err := srv.Serve(grpcListener); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err != nil {
			logger.Errorf("err: %v", err)
		}

		cosmoscmd.SetPrefixes("medichain")
		unsafeKeyring := keyring.NewUnsafe(keyring.NewInMemory())
		cosmosService := services.NewCosmosService(ctx, unsafeKeyring, appCtx.String("chain-id"), appCtx.String("cosmos-endpoint"))
		_, err := cosmosService.AddAccountFromMnemonic("admin", appCtx.String("mnemonic"))
		if err != nil {
			logger.Errorf("error while adding admin mnemonic: %v", err)
		}

		app, err := api2.InitApiService(ctx, api2.ApiServiceOptions{
			Addr:   fmt.Sprintf(":%d", appCtx.Int("port")),
			Logger: logger,
			CosmosService :cosmosService,
		})
		logger.WithField("port", appCtx.Int("port")).Info("listening for HTTP connections")
		err = app.ListenAndServe()
		if err != nil {
			logger.Errorf("err: %v", err)
		}
	}()

	// Start pprof server
	pprofListener, err := net.Listen("tcp", fmt.Sprintf(":%d", appCtx.Int("pprof-port")))
	if err != nil {
		return err
	}
	defer func() { _ = pprofListener.Close() }()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.WithField("port", appCtx.Int("pprof-port")).Info("listening for pprof requests")
		srv := new(http.Server)
		_ = srv.Serve(pprofListener)
	}()
	// Start signal watcher
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGHUP)
		select {
		case s := <-sigCh:
			logger.WithField("signal", s.String()).Infof("shutting down due to signal")
			//_ = grpcListener.Close()
			_ = pprofListener.Close()
			cancelFn()
			logger.WithField("signal", s.String()).Infof("shutdown success")
		case <-ctx.Done():
		}
	}()

	// Keep running until we receive a signal
	wg.Wait()
	return nil
}

func NewSideService(context context.Context, logger *logrus.Logger) (SideService, error) {
	return nil, xerrors.Errorf("unsupported DB URI scheme")
}

func initTracing() {
	// initJaegerTracing()
	initStackdriverTracing()
}

// func initJaegerTracing() {
// 	svcAddr := os.Getenv("JAEGER_SERVICE_ADDR")
// 	if svcAddr == "" {
// 		logger.Info("jaeger initialization disabled.")
// 		return
// 	}

// 	// Register the Jaeger exporter to be able to retrieve
// 	// the collected spans.
// 	exporter, err := jaeger.NewExporter(jaeger.Options{
// 		CollectorEndpoint: fmt.Sprintf("http://%s", svcAddr),
// 		Process: jaeger.Process{
// 			ServiceName: serviceName,
// 		},
// 	})
// 	if err != nil {
// 		logger.Fatal(err)
// 	}
// 	trace.RegisterExporter(exporter)
// 	logger.Info("jaeger initialization completed.")
// }

func initStats(exporter *stackdriver.Exporter) {
	view.SetReportingPeriod(60 * time.Second)
	view.RegisterExporter(exporter)
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		logger.Warn("Error registering default server views")
	} else {
		logger.Info("Registered default server views")
	}
}

func initStackdriverTracing() {
	// TODO(ahmetb) this method is duplicated in other microservices using Go
	// since they are not sharing packages.
	for i := 1; i <= 3; i++ {
		exporter, err := stackdriver.NewExporter(stackdriver.Options{})
		if err != nil {
			logger.Infof("failed to initialize stackdriver exporter: %+v", err)
		} else {
			trace.RegisterExporter(exporter)
			logger.Info("registered Stackdriver tracing")

			// Register the views to collect server stats.
			initStats(exporter)
			return
		}
		d := time.Second * 10 * time.Duration(i)
		logger.Infof("sleeping %v to retry initializing Stackdriver exporter", d)
		time.Sleep(d)
	}
	logger.Warn("could not initialize Stackdriver exporter after retrying, giving up")
}

func initProfiling(service, version string) {
	// TODO(ahmetb) this method is duplicated in other microservices using Go
	// since they are not sharing packages.
	for i := 1; i <= 3; i++ {
		if err := profiler.Start(profiler.Config{
			Service:        service,
			ServiceVersion: version,
			// ProjectID must be set if not running on GCP.
			// ProjectID: "my-project",
		}); err != nil {
			logger.Warnf("failed to start profiler: %+v", err)
		} else {
			logger.Info("started Stackdriver profiler")
			return
		}
		d := time.Second * 10 * time.Duration(i)
		logger.Infof("sleeping %v to retry initializing Stackdriver profiler", d)
		time.Sleep(d)
	}
	logger.Warn("could not initialize Stackdriver profiler after retrying, giving up")
}
