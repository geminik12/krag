package apiserver

import (
	"net/http"
	"time"

	"github.com/geminik12/krag/internal/pkg/core"
	"github.com/geminik12/krag/internal/pkg/errorsx"
	"github.com/gin-gonic/gin"

	genericoptions "github.com/geminik12/krag/pkg/options"
)

type HttpServerConfig struct {
	MySQLOptions *genericoptions.MySQLOptions
	Addr         string
	JWTKey       string
	Expiration   time.Duration
}

func (cfg *HttpServerConfig) NewGinServer() (Server, error) {
	engine := gin.New()

	cfg.InstallRESTAPI(engine)

	return &ginServer{
		cfg: cfg,
		srv: &http.Server{
			Addr:    cfg.Addr,
			Handler: engine,
		},
	}, nil
}

func (cfg *HttpServerConfig) InstallRESTAPI(engine *gin.Engine) {
	engine.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, nil, errorsx.ErrNotFound.WithMessage("Page not found."))
	})

	engine.GET("/health", func(c *gin.Context) {
		core.WriteResponse(c, map[string]string{"status": "ok"}, nil)
	})
}
