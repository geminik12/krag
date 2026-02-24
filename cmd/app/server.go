package app

import (
	"io"
	"log/slog"
	"os"

	"github.com/geminik12/krag/cmd/app/options"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

func NewKragCommand() *cobra.Command {
	opts := options.NewServerOptions()

	cmd := &cobra.Command{
		Use:   "krag",
		Short: "Krag is a rag api server.",
		Long:  `Krag is a rag api server that provides a RESTful interface for managing rag resources.`,
		// 命令出错时，不打印帮助信息。设置为 true 可以确保命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
		// 设置命令运行时的参数检查，不需要指定命令行参数。
		Args: cobra.NoArgs,
	}

	// 初始化配置函数，在每个命令运行时调用
	cobra.OnInitialize(onInitialize)

	// cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	// 推荐使用配置文件来配置应用，便于管理配置项
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the krag-apiserver configuration file.")

	//TODO: 支持版本号注入方式

	return cmd
}

func run(opts *options.ServerOptions) error {
	initLog()

	if err := viper.Unmarshal(opts); err != nil {
		return err
	}

	if err := opts.Validate(); err != nil {
		return err
	}

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	server, err := cfg.NewGinServer()
	if err != nil {
		return err
	}

	return server.RunOrDie()
}

func initLog() {
	// 获取日志配置
	format := viper.GetString("log.format") // 日志格式，支持：json、text
	level := viper.GetString("log.level")   // 日志级别，支持：debug, info, warn, error
	output := viper.GetString("log.output") // 日志输出路径，支持：标准输出stdout和文件

	// 转换日志级别
	var slevel slog.Level
	switch level {
	case "debug":
		slevel = slog.LevelDebug
	case "info":
		slevel = slog.LevelInfo
	case "warn":
		slevel = slog.LevelWarn
	case "error":
		slevel = slog.LevelError
	default:
		slevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: slevel}

	var w io.Writer
	var err error
	// 转换日志输出路径
	switch output {
	case "":
		w = os.Stdout
	case "stdout":
		w = os.Stdout
	default:
		w, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
	}

	// 转换日志格式
	if err != nil {
		return
	}
	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	case "text":
		handler = slog.NewTextHandler(w, opts)
	default:
		handler = slog.NewJSONHandler(w, opts)

	}

	// 设置全局的日志实例为自定义的日志实例
	slog.SetDefault(slog.New(handler))
}
