package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogConfig     *LogConfig     `yaml:"log_config"`
	HunyuanConfig *HunyuanConfig `yaml:"hunyuan_config"`
}

type HunyuanConfig struct {
	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别 debug info warn error
	Filename   string `yaml:"filename"`    // 日志文件位置
	MaxSize    int32  `yaml:"max_size"`    // 进行切割之前,日志文件的最大大小(MB为单位)
	MaxAge     int32  `yaml:"max_age"`     // 保留旧文件的最大天数
	MaxBackups int32  `yaml:"max_backups"` // 保留旧文件的最大个数
}

var config *Config

func InitConfig(path string) {
	byteFile, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("无法读取 yaml 文件: %s", path))
	}
	if err = yaml.Unmarshal(byteFile, &config); err != nil {
		panic(fmt.Errorf("无法解析 yaml 文件: %s", string(byteFile)))
	}
}

func GetLogConfig() *LogConfig {
	return config.LogConfig
}

func GetHunyuanConfig() *HunyuanConfig {
	return config.HunyuanConfig
}
