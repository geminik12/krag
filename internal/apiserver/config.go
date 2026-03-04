package apiserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/geminik12/krag/internal/apiserver/biz"
	"github.com/geminik12/krag/internal/apiserver/handler"
	"github.com/geminik12/krag/internal/apiserver/pkg/validation"
	"github.com/geminik12/krag/internal/apiserver/store"
	"github.com/geminik12/krag/internal/pkg/core"
	"github.com/geminik12/krag/internal/pkg/errorsx"
	"github.com/geminik12/krag/internal/pkg/known"
	"github.com/geminik12/krag/internal/pkg/storage"
	"github.com/geminik12/krag/internal/pkg/worker"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"

	mw "github.com/geminik12/krag/internal/pkg/middleware"
	genericoptions "github.com/geminik12/krag/pkg/options"
	"github.com/geminik12/krag/pkg/token"
)

type HttpServerConfig struct {
	MySQLOptions  *genericoptions.MySQLOptions
	Addr          string
	JWTKey        string
	Expiration    time.Duration
	OllamaOptions *genericoptions.OllamaOptions
	MinIOOptions  *genericoptions.MinIOOptions
	QdrantOptions *genericoptions.QdrantOptions
	RedisOptions  *genericoptions.RedisOptions
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
	client := cfg.OllamaOptions.NewClient()

	// 初始化 MinIO 存储
	minioStore, err := storage.NewMinio(
		cfg.MinIOOptions.Endpoint,
		cfg.MinIOOptions.AccessKeyID,
		cfg.MinIOOptions.SecretAccessKey,
		cfg.MinIOOptions.Bucket,
		cfg.MinIOOptions.UseSSL,
	)
	if err != nil {
		return nil, err
	}

	// Init Asynq Client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.RedisOptions.Addr,
		Password: cfg.RedisOptions.Password,
		DB:       cfg.RedisOptions.DB,
	})

	// Create Biz
	bizInstance := biz.NewBiz(store, client, minioStore, cfg.QdrantOptions, asynqClient)

	// Start Worker
	w := worker.NewWorker(cfg.RedisOptions.Addr, cfg.RedisOptions.Password, cfg.RedisOptions.DB, bizInstance.KnowledgeV1())
	go func() {
		if err := w.Run(); err != nil {
			slog.Error("Failed to start worker", "err", err)
		}
	}()

	cfg.InstallRESTAPI(engine, bizInstance, store)

	// Monitor
	h := asynqmon.New(asynqmon.Options{
		RootPath: "/admin/queues",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     cfg.RedisOptions.Addr,
			Password: cfg.RedisOptions.Password,
			DB:       cfg.RedisOptions.DB,
		},
	})
	engine.Any("/admin/queues/*filepath", gin.WrapH(h))

	return &ginServer{
		cfg: cfg,
		srv: &http.Server{
			Addr:    cfg.Addr,
			Handler: engine,
		},
	}, nil
}

func (cfg *HttpServerConfig) InstallRESTAPI(engine *gin.Engine, bizInstance biz.IBiz, store store.IStore) {
	engine.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, nil, errorsx.ErrNotFound.WithMessage("Page not found."))
	})

	engine.GET("/health", func(c *gin.Context) {
		core.WriteResponse(c, map[string]string{"status": "ok"}, nil)
	})

	// 创建核心业务处理器
	handler := handler.NewHandler(bizInstance, validation.NewValidator(store))

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

		conversationv1 := v1.Group("/conversations")
		{
			conversationv1.Use(authMiddlewares...)
			conversationv1.POST("", handler.CreateConversation)
			conversationv1.GET("", handler.ListConversation)
			conversationv1.GET(":id", handler.GetConversation)
			conversationv1.DELETE(":id", handler.DeleteConversation)
			conversationv1.GET(":id/messages", handler.ListMessage)
		}

		chatv1 := v1.Group("/chat")
		{
			chatv1.Use(authMiddlewares...)
			chatv1.POST("", handler.Chat)
		}

		knowledgev1 := v1.Group("/knowledge")
		{
			knowledgev1.Use(authMiddlewares...)
			knowledgev1.POST("/upload", handler.UploadKnowledge)
			knowledgev1.POST("/search", handler.SearchKnowledge)
		}
	}
}
