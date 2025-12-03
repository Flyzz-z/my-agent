# API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **Content-Type**: `application/json`

## 1. 秒杀系统 API

### 1.1 秒杀接口

**POST** `/seckill/`

**请求体**:
```json
{
  "user_id": 123,
  "coupon_id": 1
}
```

**响应**:
```json
{
  "success": true,
  "message": "秒杀成功，订单处理中",
  "order_id": 456
}
```

### 1.2 获取优惠券信息

**GET** `/seckill/coupon/:id`

**响应**:
```json
{
  "id": 1,
  "name": "双十一优惠券",
  "description": "满100减50",
  "total_stock": 1000,
  "remain_stock": 500,
  "start_time": "2024-11-01T00:00:00Z",
  "end_time": "2024-11-11T23:59:59Z",
  "status": 1
}
```

### 1.3 初始化库存

**POST** `/seckill/init-stock`

**请求体**:
```json
{
  "coupon_id": 1
}
```

## 2. AI Agent API

### 2.1 对话接口

**POST** `/agent/chat`

**请求体**:
```json
{
  "query": "Kafka如何阻止重复消费?",
  "use_rag": true,
  "session": "session-123"
}
```

**响应**:
```json
{
  "answer": "Kafka 可以通过以下方式阻止重复消费...",
  "documents": ["文档片段1", "文档片段2"],
  "session": "session-123"
}
```

### 2.2 流式对话接口

**POST** `/agent/chat-stream`

**请求体**:
```json
{
  "query": "解释一下 Go 的协程",
  "use_rag": false,
  "session": "session-456"
}
```

**响应**: Server-Sent Events (SSE)

## 3. 文档搜索 API

### 3.1 搜索文档

**POST** `/search/`

**请求体**:
```json
{
  "query": "Kafka 消费",
  "page": 1,
  "page_size": 10,
  "tags": ["技术", "消息队列"]
}
```

**响应**:
```json
{
  "total": 100,
  "documents": [
    {
      "id": "doc-1",
      "title": "Kafka 消费者指南",
      "content": "...",
      "author": "张三",
      "tags": ["技术", "消息队列"],
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "page_size": 10
}
```

### 3.2 索引文档

**POST** `/search/document`

**请求体**:
```json
{
  "id": "doc-1",
  "title": "文档标题",
  "content": "文档内容...",
  "author": "作者",
  "tags": ["标签1", "标签2"]
}
```

### 3.3 获取文档

**GET** `/search/document/:id`

**响应**:
```json
{
  "id": "doc-1",
  "title": "文档标题",
  "content": "文档内容...",
  "author": "作者",
  "tags": ["标签1", "标签2"],
  "created_at": "2024-01-01T00:00:00Z"
}
```

### 3.4 删除文档

**DELETE** `/search/document/:id`

**响应**:
```json
{
  "message": "文档删除成功"
}
```

## 4. 健康检查

**GET** `/health`

**响应**:
```json
{
  "status": "ok"
}
```

## 错误响应

所有 API 在发生错误时返回以下格式：

```json
{
  "error": "错误信息"
}
```

HTTP 状态码：
- `200`: 成功
- `400`: 请求参数错误
- `500`: 服务器内部错误
