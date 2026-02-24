package apiserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server interface {
	RunOrDie() error
}

type ginServer struct {
	cfg *HttpServerConfig
	srv *http.Server
}

func (s *ginServer) RunOrDie() error {
	slog.Info("Start to listening the incoming requests on http address", "addr", s.cfg.Addr)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	// 创建一个 os.Signal 类型的 channel，用于接收系统信号
	quit := make(chan os.Signal, 1)
	// 当执行 kill 命令时（不带参数），默认会发送 syscall.SIGTERM 信号
	// 使用 kill -2 命令会发送 syscall.SIGINT 信号（例如按 CTRL+C 触发）
	// 使用 kill -9 命令会发送 syscall.SIGKILL 信号，但 SIGKILL 信号无法被捕获，因此无需监听和处理
	signal.Notify(quit, os.Interrupt)
	// 阻塞程序，等待从 quit channel 中接收到信号
	<-quit

	slog.Info("Shutdown the http server")

	// 优雅关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("Insecure Server forced to shutdown", "err", err)
		return err
	}

	slog.Info("Server exiting")

	return nil
}
