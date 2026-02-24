package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultHomeDir = ".krag"

	defaultConfigName = "krag-apiserver.yaml"
)

func onInitialize() {
	if len(configFile) != 0 {
		viper.SetConfigFile(configFile)
	} else {
		for _, dir := range searchDirs() {
			viper.AddConfigPath(dir)
		}
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigName)
	}

	setupEnvironmentVariables()

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果读取配置文件失败，则打印错误信息并退出程序
		cobra.CheckErr(err)
	}
}

// searchDirs 返回默认的配置文件搜索目录.
func searchDirs() []string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	// 如果获取用户主目录失败，则打印错误信息并退出程序
	cobra.CheckErr(err)
	return []string{filepath.Join(homeDir, defaultHomeDir), "."}
}

// filePath 获取默认配置文件的完整路径.
func filePath() string {
	home, err := os.UserHomeDir()
	// 如果不能获取用户主目录，则记录错误并返回空路径
	cobra.CheckErr(err)
	return filepath.Join(home, defaultHomeDir, defaultConfigName)
}

// setupEnvironmentVariables 配置环境变量规则.
func setupEnvironmentVariables() {
	// 允许 viper 自动匹配环境变量
	viper.AutomaticEnv()
	// 设置环境变量前缀
	viper.SetEnvPrefix("FASTGO")
	// 替换环境变量 key 中的分隔符 '.' 和 '-' 为 '_'
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
}
