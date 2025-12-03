package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 结构体定义整个应用程序的配置
type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Redis       RedisConfig       `yaml:"redis"`
	MySQL       MySQLConfig       `yaml:"mysql"`
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	RocketMQ    RocketMQConfig    `yaml:"rocketmq"`
	LLM         LLMConfig         `yaml:"llm"`
	RAG         RAGConfig         `yaml:"rag"`
	Embedding   EmbeddingConfig   `yaml:"embedding"`
	Seckill     SeckillConfig     `yaml:"seckill"`
}

// RedisConfig Redis相关配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Protocol int    `yaml:"protocol"`
	UnstableResp3 bool `yaml:"unstable_resp3"`
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

// ServerConfig 服务器相关配置
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// MySQLConfig MySQL数据库配置
type MySQLConfig struct {
	DSN             string `yaml:"dsn"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// ElasticsearchConfig Elasticsearch配置
type ElasticsearchConfig struct {
	Addresses []string `yaml:"addresses"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	Index     string   `yaml:"index"`
}

// RocketMQConfig RocketMQ配置
type RocketMQConfig struct {
	NameServer   string `yaml:"name_server"`
	GroupName    string `yaml:"group_name"`
	Topic        string `yaml:"topic"`
	InstanceName string `yaml:"instance_name"`
}

// SeckillConfig 秒杀系统配置
type SeckillConfig struct {
	CachePrefix    string        `yaml:"cache_prefix"`
	LockPrefix     string        `yaml:"lock_prefix"`
	LockExpire     time.Duration `yaml:"lock_expire"`
	MaxRetry       int           `yaml:"max_retry"`
	OrderTimeout   time.Duration `yaml:"order_timeout"`
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

func GetConfig() *Config {
	return &GlobalConfig
}