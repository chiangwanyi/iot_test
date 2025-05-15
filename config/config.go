package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Config struct {
	Debug  bool   `json:"debug"`
	Port   int    `json:"port"`
	DBPath string `json:"db_path"`
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig() {
	once.Do(func() {
		data, err := os.ReadFile("config.json")
		if err != nil {
			log.Fatalf("读取配置文件失败: %v", err)
		}

		instance = &Config{}
		if err := json.Unmarshal(data, instance); err != nil {
			log.Fatalf("解析配置文件失败: %v", err)
		}
	})
}

func GetConfig() *Config {
	if instance == nil {
		LoadConfig()
	}
	return instance
}
