package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 结构体定义整个应用程序的配置
type Config struct {
	Redis  RedisConfig  `yaml:"redis"`
	LLM    LLMConfig    `yaml:"llm"`
	RAG    RAGConfig    `yaml:"rag"`
	Embedding EmbeddingConfig `yaml:"embedding"`
}

// RedisConfig Redis相关配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// LLMConfig 大语言模型相关配置
type LLMConfig struct {
	BaseURL string        `yaml:"base_url"`
	Model   string        `yaml:"model"`
	Timeout time.Duration `yaml:"timeout"`
}

// RAGConfig RAG相关配置
type RAGConfig struct {
	IndexName    string   `yaml:"index_name"`
	Prefix       string   `yaml:"prefix"`
	Dimension    int64    `yaml:"dimension"`
	VectorField  string   `yaml:"vector_field"`
	TopK         int      `yaml:"top_k"`
	Dialect      int      `yaml:"dialect"`
	ReturnFields []string `yaml:"return_fields"`
}

type EmbeddingConfig struct {
	Model string `yaml:"model"`
	APIKey string `yaml:"api_key"`
}

var (
	DefaultConfigPath = "/home/flyzz/agent/config.yaml"
	GlobalConfig      Config
)

// LoadConfig 从YAML文件加载配置到全局配置
func LoadConfig(configPath string) error {
	// 确保路径是绝对路径
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	// 读取配置文件
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	// 解析YAML配置到全局结构体
	if err := yaml.Unmarshal(data, &GlobalConfig); err != nil {
		return err
	}

	return nil
}

// LoadDefaultConfig 初始化全局配置为默认配置
func LoadDefaultConfig() {
	GlobalConfig = Config{
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		LLM: LLMConfig{
			BaseURL: "http://localhost:11434",
			Model:   "qwen:4b",
			Timeout: 60 * time.Second,
		},
		RAG: RAGConfig{
			IndexName:    "rag_index:",
			Prefix:       "rag_prefix:",
			Dimension:    1536, // 常见嵌入模型的维度
			VectorField:  "vector_content",
			TopK:         1,
			Dialect:      2,
			ReturnFields: []string{"vector_content", "content"},
		},
	}
}

func GetConfig() *Config {
	return &GlobalConfig
}