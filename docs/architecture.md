# 项目架构文档

## 系统架构

本项目是一个基于 Go 语言开发的多功能服务器应用，整合了以下三大核心功能：

1. **优惠券秒杀系统** - 基于 Redis、MySQL、RocketMQ
2. **AI Agent 服务** - 基于 Eino 框架和 RAG 技术
3. **文档搜索服务** - 基于 Elasticsearch

## 技术栈

### 核心框架
- **Eino**: 字节跳动开源的 AI Agent 框架
- **Gin**: HTTP Web 框架
- **Go 1.24**: 编程语言

### 存储层
- **Redis**: 缓存、分布式锁、向量数据库
- **MySQL**: 关系型数据库
- **Elasticsearch**: 全文搜索引擎

### 消息队列
- **RocketMQ**: 异步消息处理

### AI 相关
- **Ollama**: 本地 LLM 部署 (qwen3:4b)
- **Ark Embedding**: 文本嵌入模型

## 目录结构

```
my-agent/
├── cmd/server/              # 服务器启动入口
├── internal/                # 私有应用代码
│   ├── server/http/        # HTTP 服务器
│   ├── domain/             # 业务领域层
│   │   ├── seckill/        # 秒杀业务
│   │   ├── agent/          # AI Agent
│   │   └── search/         # 文档搜索
│   └── infrastructure/     # 基础设施层
│       ├── repository/     # 数据仓库
│       ├── mq/             # 消息队列
│       └── rag/            # RAG 组件
├── pkg/                    # 公共库
├── config/                 # 配置
├── api/                    # API 定义
├── docs/                   # 文档
└── scripts/                # 脚本
```

## 架构设计原则

### 分层架构
1. **Server 层**: 处理 HTTP 请求和响应
2. **Domain 层**: 业务逻辑和领域模型
3. **Infrastructure 层**: 技术实现细节

### 依赖注入
通过接口抽象，实现依赖倒置，便于测试和替换实现。

### 配置管理
统一的配置文件管理，支持多环境配置。

## 核心功能

### 1. 秒杀系统

**流程设计**:
1. 用户发起秒杀请求
2. 获取分布式锁
3. 检查 Redis 中的库存并扣减
4. 发送订单消息到 RocketMQ
5. 异步处理订单，写入 MySQL

**关键技术**:
- Redis 分布式锁防止超卖
- Redis 缓存库存提升性能
- RocketMQ 异步处理削峰填谷

### 2. AI Agent

**RAG 流程**:
1. 接收用户查询
2. 使用 Embedding 模型向量化查询
3. 从 Redis Vector Store 检索相关文档
4. 将文档和查询传递给 LLM
5. 生成回答

**Eino Graph 节点**:
- Retriever: 检索文档
- Format: 格式化文档
- ChatTemplate: 构建提示词
- ChatModel: 生成回答

### 3. 文档搜索

**倒排索引搜索**:
- 基于 Elasticsearch
- 支持全文搜索
- 支持标签过滤
- 分页查询

## 数据流

```
┌─────────┐      ┌─────────┐      ┌──────────┐
│ Client  │─────▶│   Gin   │─────▶│ Handler  │
└─────────┘      └─────────┘      └──────────┘
                                         │
                                         ▼
                                   ┌──────────┐
                                   │ Service  │
                                   └──────────┘
                                         │
                    ┌────────────────────┼────────────────────┐
                    ▼                    ▼                    ▼
              ┌──────────┐         ┌─────────┐         ┌─────────┐
              │  MySQL   │         │  Redis  │         │   MQ    │
              └──────────┘         └─────────┘         └─────────┘
```

## 配置说明

配置文件位于 `config.yaml`，包含以下部分：
- `server`: HTTP 服务器配置
- `redis`: Redis 连接配置
- `mysql`: MySQL 连接配置
- `elasticsearch`: ES 连接配置
- `rocketmq`: RocketMQ 配置
- `llm`: 大语言模型配置
- `rag`: RAG 系统配置
- `embedding`: Embedding 模型配置
- `seckill`: 秒杀系统配置

## 后续优化

1. 实现基础设施层的具体实现 (MySQL、Elasticsearch、RocketMQ)
2. 添加单元测试和集成测试
3. 实现流式对话接口
4. 添加监控和日志
5. 实现用户认证和鉴权
6. 性能优化和压测
