# 部署文档

## 本地开发环境部署

### 1. 前置要求

- Go 1.24+
- Docker & Docker Compose
- Ollama (用于本地 LLM)

### 2. 启动基础设施

使用 Docker Compose 启动所有依赖服务：

```bash
cd scripts
docker-compose up -d
```

这将启动以下服务：
- MySQL (端口 3306)
- Redis (端口 6379, 8001)
- Elasticsearch (端口 9200, 9300)
- RocketMQ NameServer (端口 9876)
- RocketMQ Broker (端口 10909, 10911, 10912)

### 3. 初始化数据库

数据库会在容器启动时自动初始化。如需手动初始化：

```bash
mysql -h localhost -u root -p < scripts/init_db.sql
# 密码: password
```

### 4. 配置 Ollama

启动 Ollama 并下载模型：

```bash
# 启动 Ollama
ollama serve

# 下载模型
ollama pull qwen3:4b
```

### 5. 配置应用

编辑 `config.yaml`，确保所有配置项正确：

```yaml
server:
  host: "0.0.0.0"
  port: 8080

redis:
  addr: "localhost:6379"

mysql:
  dsn: "root:password@tcp(localhost:3306)/seckill?charset=utf8mb4&parseTime=True&loc=Local"

# ... 其他配置
```

### 6. 安装依赖

```bash
go mod tidy
```

### 7. 运行服务器

```bash
go run cmd/server/main.go
```

服务器将在 `http://localhost:8080` 启动。

### 8. 验证部署

访问健康检查接口：

```bash
curl http://localhost:8080/health
```

预期响应：
```json
{"status":"ok"}
```

## 生产环境部署

### 1. 构建可执行文件

```bash
# 构建 Linux 二进制文件
GOOS=linux GOARCH=amd64 go build -o bin/server cmd/server/main.go
```

### 2. 使用 Makefile

创建 `Makefile`:

```makefile
.PHONY: build run test clean

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

deps:
	go mod tidy
	go mod download

docker-up:
	cd scripts && docker-compose up -d

docker-down:
	cd scripts && docker-compose down

install: deps build
```

使用方式：
```bash
make install  # 安装依赖并构建
make run      # 运行服务
make test     # 运行测试
```

### 3. 系统服务配置 (systemd)

创建 `/etc/systemd/system/seckill-agent.service`:

```ini
[Unit]
Description=Seckill Agent Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=app
WorkingDirectory=/opt/seckill-agent
ExecStart=/opt/seckill-agent/bin/server
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable seckill-agent
sudo systemctl start seckill-agent
sudo systemctl status seckill-agent
```

### 4. Nginx 反向代理

配置文件 `/etc/nginx/sites-available/seckill-agent`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

启用站点：
```bash
sudo ln -s /etc/nginx/sites-available/seckill-agent /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 监控和日志

### 日志查看

应用日志：
```bash
sudo journalctl -u seckill-agent -f
```

Docker 服务日志：
```bash
docker-compose logs -f mysql
docker-compose logs -f redis
docker-compose logs -f elasticsearch
```

### 性能监控

建议安装：
- Prometheus (指标收集)
- Grafana (可视化)
- ELK Stack (日志聚合)

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 MySQL 是否运行: `docker ps`
   - 验证 DSN 配置是否正确

2. **Redis 连接失败**
   - 检查 Redis 是否运行
   - 验证端口是否开放

3. **Ollama 连接失败**
   - 确保 Ollama 服务正在运行
   - 检查 base_url 配置

4. **RocketMQ 消息发送失败**
   - 检查 NameServer 和 Broker 状态
   - 验证网络配置

## 备份策略

### MySQL 备份
```bash
# 每日备份
mysqldump -u root -p seckill > backup_$(date +%Y%m%d).sql
```

### Redis 备份
```bash
# 触发 RDB 快照
redis-cli SAVE
```

### Elasticsearch 备份
使用 Elasticsearch 快照功能定期备份索引。
