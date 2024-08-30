package config

import (
	"os"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
)

// Config 是全局配置的结构体
type Config struct {
	Logger  *logrus.Logger
	ARCH    string
	OS      string
	HomeDir string
	IsRoot  bool
}

var (
	cfg  *Config
	once sync.Once
)

// InitConfig 初始化配置
func InitConfig() {
	once.Do(func() {
		logger := logrus.New()
		// 这里可以根据需要配置日志格式、输出等
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		// 使用 runtime 包获取当前系统的架构和操作系统
		arch := runtime.GOARCH
		osType := runtime.GOOS

		// 获取当前用户的家目录
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Errorf("获取家目录失败: %s", err)
			os.Exit(1)
		}
		isRoot := os.Geteuid() == 0

		// 初始化全局配置对象
		cfg = &Config{
			Logger:  logger,
			ARCH:    arch,
			OS:      osType,
			HomeDir: homeDir,
			IsRoot:  isRoot,
		}
	})
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return cfg
}
