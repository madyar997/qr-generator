// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
	"github.com/madyar997/qr-generator/config"
	v1 "github.com/madyar997/qr-generator/internal/controller/http/v1"
	"github.com/madyar997/qr-generator/internal/usecase"
	"github.com/madyar997/qr-generator/internal/usecase/repo"
	"github.com/madyar997/qr-generator/internal/usecase/webapi"
	"github.com/madyar997/qr-generator/pkg/httpserver"
	"github.com/madyar997/qr-generator/pkg/jaeger"
	"github.com/madyar997/qr-generator/pkg/logger"
	"github.com/madyar997/qr-generator/pkg/postgres"
	ssoClient "github.com/madyar997/user-client/client/grpc"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	//tracing
	tracer, closer, _ := jaeger.InitJaeger()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	ssoCli, err := ssoClient.NewClient(ssoClient.WithAddress(cfg.Grpc.Port), ssoClient.WithInsecure())
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - ssoClient : %w", err))
	}

	// Use case
	translationUseCase := usecase.New(
		repo.New(pg),
		webapi.New(),
	)

	qrUseCase := usecase.NewQrUseCase(http.DefaultClient, ssoCli)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, translationUseCase, qrUseCase, cfg)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
