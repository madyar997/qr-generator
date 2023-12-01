// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/madyar997/qr-generator/config"
	"github.com/madyar997/qr-generator/internal/controller/http/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/madyar997/qr-generator/docs"
	"github.com/madyar997/qr-generator/internal/usecase"
	"github.com/madyar997/qr-generator/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @description Using a translation service as an example
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.Translation, q usecase.Qr, cfg *config.Config) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	handler.Use(middleware.HTTPMetrics())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")
	{
		newTranslationRoutes(h, t, l, cfg)
		newQrRoutes(h, q, l, cfg)
	}

	handler.GET("/code-200", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	handler.GET("/code-400", func(c *gin.Context) {
		c.Status(http.StatusBadRequest)
	})
	handler.GET("/code-500", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})
}
