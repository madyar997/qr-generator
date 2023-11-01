package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/madyar997/qr-generator/config"
	"github.com/madyar997/qr-generator/internal/controller/http/middleware"
	"github.com/madyar997/qr-generator/internal/usecase"
	"github.com/madyar997/qr-generator/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type qrRoutes struct {
	q   usecase.Qr
	l   logger.Interface
	cfg *config.Config
}

func newQrRoutes(handler *gin.RouterGroup, q usecase.Qr, l logger.Interface, cfg *config.Config) {
	r := &qrRoutes{q, l, cfg}

	h := handler.Group("/qr")
	{
		h.GET("/me", middleware.JwtVerify(cfg), r.me)
	}
}

func (r *qrRoutes) me(ctx *gin.Context) {
	span := opentracing.StartSpan("qr-generator-service /me handler method")
	defer span.Finish()

	userID, ok := ctx.Get("user_id")
	if !ok {
		errorResponse(ctx, http.StatusInternalServerError, "can not parse user id ")

		return
	}

	context := opentracing.ContextWithSpan(ctx.Request.Context(), span)

	res, err := r.q.Me(context, int(userID.(float64)))
	if err != nil {
		r.l.Error(err, "http - v1 - me")
		errorResponse(ctx, http.StatusInternalServerError, err.Error())

		return
	}
	ctx.Data(http.StatusOK, "image/png", res)
}
