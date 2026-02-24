package apiserver

import (
	"net/http"
	"time"

	"github.com/geminik12/krag/internal/apiserver/biz"
	"github.com/geminik12/krag/internal/apiserver/handler"
	"github.com/geminik12/krag/internal/apiserver/pkg/validation"
	"github.com/geminik12/krag/internal/apiserver/store"
	"github.com/geminik12/krag/internal/pkg/core"
	"github.com/geminik12/krag/internal/pkg/errorsx"
	"github.com/geminik12/krag/internal/pkg/known"
	"github.com/gin-gonic/gin"

	mw "github.com/geminik12/krag/internal/pkg/middleware"
	genericoptions "github.com/geminik12/krag/pkg/options"
	"github.com/geminik12/krag/pkg/token"
)

type HttpServerConfig struct {
	MySQLOptions *genericoptions.MySQLOptions
	Addr         string
	JWTKey       string
	Expiration   time.Duration
}

func (cfg *HttpServerConfig) NewGinServer() (Server, error) {
	token.Init(cfg.JWTKey, known.XUserID, cfg.Expiration)

	engine := gin.New()

	mws := []gin.HandlerFunc{gin.Recovery(), mw.NoCache, mw.Cors, mw.RequestID(), mw.Logger()}
	engine.Use(mws...)

	db, err := cfg.MySQLOptions.NewDB()
	if err != nil {
		return nil, err
	}
	store := store.NewStore(db)

	cfg.InstallRESTAPI(engine, store)

	return &ginServer{
		cfg: cfg,
		srv: &http.Server{
			Addr:    cfg.Addr,
			Handler: engine,
		},
	}, nil
}

func (cfg *HttpServerConfig) InstallRESTAPI(engine *gin.Engine, store store.IStore) {
	engine.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, nil, errorsx.ErrNotFound.WithMessage("Page not found."))
	})

	engine.GET("/health", func(c *gin.Context) {
		core.WriteResponse(c, map[string]string{"status": "ok"}, nil)
	})

	// 创建核心业务处理器
	handler := handler.NewHandler(biz.NewBiz(store), validation.NewValidator(store))

	// 注册用户登录和令牌刷新接口。这2个接口比较简单，所以没有 API 版本
	engine.POST("/login", handler.Login)
	// 注意：认证中间件要在 handler.RefreshToken 之前加载
	engine.PUT("/refresh-token", mw.Authn(), handler.RefreshToken)

	authMiddlewares := []gin.HandlerFunc{mw.Authn()}

	v1 := engine.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", handler.CreateUser)
			userv1.Use(authMiddlewares...)
			userv1.PUT(":userID", handler.UpdateUser)
			userv1.PUT(":userID/change-password", handler.ChangePassword)
			userv1.DELETE(":userID", handler.DeleteUser)
			userv1.GET(":userID", handler.GetUser)
		}
	}
}
