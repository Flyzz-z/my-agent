# 项目架构说明

## 概述
本项目是一个整合了多个功能模块的 Go 后端应用，采用领域驱动设计（DDD）架构。

## 三大主要功能

### 1. AI 搜索 (AI Search)
**路径**: `/api/v1/aisearch`

**说明**:
- AI 搜索是系统的核心功能之一，整合了 LLM（大语言模型）和 RAG（检索增强生成）能力
- 通过 RAG 技术从文档库中检索相关内容，结合 LLM 生成智能回答
- 使用 Eino 框架构建的 Graph 执行流程：检索 → 格式化 → 模板 → LLM 生成

**主要端点**:
- `POST /api/v1/aisearch/search` - 执行AI搜索
- `POST /api/v1/aisearch/search-stream` - 流式AI搜索
- `POST /api/v1/aisearch/document` - 添加文档到RAG索引

**技术栈**:
- **LLM**: Ollama (支持本地大模型)
- **RAG**: Redis Vector Search (向量检索)
- **Embedding**: 豆包 Embedding API
- **Framework**: Eino (CloudWeGo)

### 2. 秒杀系统 (Seckill)
**路径**: `/api/v1/seckill`

**说明**:
- 高性能秒杀系统，处理高并发抢购场景
- 使用 Redis 进行库存管理和分布式锁
- 通过消息队列异步处理订单

**主要端点**:
- `POST /api/v1/seckill` - 执行秒杀
- `GET /api/v1/seckill/coupon/:id` - 获取优惠券信息
- `POST /api/v1/seckill/init-stock` - 初始化库存

### 3. [第三个功能 - 待规划]
可以考虑添加：
- 工具调用系统 (Tool Calling)
- 多轮对话管理
- 知识图谱服务
- 等等...

## 项目结构

```
my-agent/
├── cmd/
│   └── server/          # HTTP 服务器入口
│       └── main.go
├── internal/
│   ├── domain/          # 领域层
│   │   ├── aisearch/    # AI搜索领域 (整合了LLM+RAG)
│   │   │   ├── model.go
│   │   │   └── service.go
│   │   └── seckill/     # 秒杀领域
│   │       ├── model.go
│   │       ├── repository.go
│   │       └── service.go
│   ├── infrastructure/  # 基础设施层
│   │   └── rag/         # RAG引擎实现
│   │       ├── rag_engine.go
│   │       ├── embedding.go
│   │       ├── indexer.go
│   │       ├── retriever.go
│   │       └── splitter.go
│   └── server/          # 服务器层
│       └── http/
│           ├── handler/ # HTTP处理器
│           ├── middleware/
│           └── router.go
├── pkg/                 # 公共包
│   ├── llm/            # LLM客户端
│   └── utils/          # 工具函数
├── config/             # 配置管理
├── docs/               # 文档
├── scripts/            # 脚本
├── main.go             # 简单测试入口
└── config.yaml         # 配置文件
```

## 关键设计决策

### 为什么将 LLM 和 RAG 整合到 AI Search？

1. **功能内聚**: LLM 和 RAG 都是为 AI 搜索服务的核心技术，它们不是独立的业务功能
2. **简化架构**: 避免过度拆分，三大主要功能更清晰：AI搜索、秒杀、[第三个功能]
3. **易于理解**: 用户视角看，他们使用的是"AI搜索"功能，而非直接使用"LLM"或"RAG"
4. **可扩展性**: AI搜索内部可以继续优化和扩展技术实现，不影响整体架构

### AI Search 内部架构

```
┌─────────────────────────────────────────────┐
│           AI Search Service                 │
├─────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │   RAG    │  │   LLM    │  │  Graph   │ │
│  │  Engine  │  │  Client  │  │  Runner  │ │
│  └──────────┘  └──────────┘  └──────────┘ │
│       │             │              │        │
│       └─────────────┴──────────────┘        │
│              Eino Framework                 │
└─────────────────────────────────────────────┘
         │
         ▼
    Search API
    /api/v1/aisearch/*
```

## 运行方式

### 快速测试 (main.go)
```bash
go run main.go
```
简单测试 AI 搜索功能，无需启动完整服务

### 完整 HTTP 服务
```bash
go run cmd/server/main.go
```
启动包含所有功能的 HTTP 服务器

## 配置说明

配置文件 `config.yaml` 包含：
- **Server**: HTTP 服务器配置
- **Redis**: 缓存和向量存储
- **LLM**: 大模型配置 (Ollama)
- **RAG**: 向量检索配置
- **Embedding**: 嵌入模型配置
- **Seckill**: 秒杀系统配置

## 技术栈

- **Go**: 1.21+
- **Web Framework**: Gin
- **AI Framework**: Eino (CloudWeGo)
- **LLM**: Ollama
- **Vector DB**: Redis Vector Search
- **Cache**: Redis
- **Config**: Viper
- **Message Queue**: RocketMQ (秒杀系统)

## API 示例

### AI 搜索
```bash
curl -X POST http://localhost:8080/api/v1/aisearch/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Kafka如何阻止重复消费",
    "session": "session-123"
  }'
```

### 添加文档
```bash
curl -X POST http://localhost:8080/api/v1/aisearch/document \
  -H "Content-Type: application/json" \
  -d '{
    "file_path": "/path/to/document.md"
  }'
```

### 秒杀
```bash
curl -X POST http://localhost:8080/api/v1/seckill \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 123,
    "coupon_id": 456
  }'
```

## 下一步计划

1. 完善 AI 搜索的流式响应
2. 实现秒杀系统的 MySQL 和 RocketMQ 集成
3. 规划和实现第三个主要功能
4. 添加更多工具调用能力到 AI 搜索
5. 完善测试覆盖率
6. 添加 API 文档 (Swagger)
