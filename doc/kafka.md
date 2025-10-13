# Kafka 总结
## Kafka 介绍
Kafka 是一个分布式流式处理平台。这到底是什么意思呢？流平台具有三个关键功能：
1. 消息队列：发布和订阅消息流，这个功能类似于消息队列，这也是 Kafka 也被归类为消息队列的原因。
2. 容错的持久方式存储记录消息流：Kafka 会把消息持久化到磁盘，有效避免了消息丢失的风险。
3. 流式处理平台： 在消息发布的时候进行处理，Kafka 提供了一个完整的流式处理类库。
Kafka 主要有两大应用场景：
1. 消息队列：建立实时流数据管道，以可靠地在系统或应用程序之间获取数据。
2. 数据处理： 构建实时的流数据处理程序来转换或处理数据流

## Kafka 核心概念
### 基本概念
1. Producer（生产者） : 产生消息的一方。
2. Consumer（消费者） : 消费消息的一方。
3. Broker（代理） : 可以看作是一个独立的 Kafka 实例。多个 Kafka Broker 组成一个 Kafka Cluster。
同时，你一定也注意到每个 Broker 中又包含了 Topic 以及 Partition 这两个重要的概念：
* Topic（主题） : Producer 将消息发送到特定的主题，Consumer 通过订阅特定的 Topic(主题) 来消费消息。
* Partition（分区） : Partition 属于 Topic 的一部分。一个 Topic 可以有多个 Partition ，并且同一 Topic 下的 Partition 可以分布在不同的 Broker 上，这也就表明一个 Topic 可以横跨多个 Broker 。这正如我上面所画的图一样。

### Kafka多副本
还有一点我觉得比较重要的是 Kafka 为分区（Partition）引入了多副本（Replica）机制。分区（Partition）中的多个副本之间会有一个叫做 leader 的家伙，其他副本称为 follower。我们发送的消息会被发送到 leader 副本，然后 follower 副本才能从 leader 副本中拉取消息进行同步。
生产者和消费者只与 leader 副本交互。你可以理解为其他副本只是 leader 副本的拷贝，它们的存在只是为了保证消息存储的安全性。当 leader 副本发生故障时会从 follower 中选举出一个 leader,但是 follower 中如果有和 leader 同步程度达不到要求的参加不了 leader 的竞选。
Kafka 的多分区（Partition）以及多副本（Replica）机制有什么好处呢？
1. Kafka 通过给特定 Topic 指定多个 Partition, 而各个 Partition 可以分布在不同的 Broker 上, 这样便能提供比较好的并发能力（负载均衡）。
2. Partition 可以指定对应的 Replica 数, 这也极大地提高了消息存储的安全性, 提高了容灾能力，不过也相应的增加了所需要的存储空间。

## Kafka 消费顺序、消息丢失和重复消费
### 消费顺序
同一 Partition 下的消息是有序的，但是不同 Partition 下的消息是无序的。 若要保证消息有序则可以指定 Producer 发送消息时指定 key, 并将 key 作为 partition key, 这样同一个 key 的消息就会被发送到同一个 Partition 中。

### 消息丢失
生产者：基于ACK机制确保消息发送成功, 即生产者发送消息后会等待 Broker 确认消息是否成功写入到 Partition 中。
消费者：基于offset机制确保消息消费成功, 并通过提交offset来标记消息已被消费。
Broker：基于多副本机制确保消息不丢失, 即消息会被复制到多个副本中, 确保消息的可靠性。

死信队列：当消息被消费失败时，重试仍旧失败后， 可以将消息发送到一个特殊的 Topic 中, 这个 Topic 就被称为死信队列（Dead Letter Queue），等待后续消费。

### 重复消费
* 消费消息服务做幂等校验，比如 Redis 的 set、MySQL 的主键等天然的幂等功能。这种方法最有效。
* 将 enable.auto.commit 参数设置为 false，关闭自动提交，开发者在代码中手动提交 offset。那么这里会有个问题：什么时候提交 offset 合适？
  。处理完消息再提交：依旧有消息重复消费的风险，和自动提交一样
	。拉取到消息即提交：会有消息丢失的风险。允许消息延时的场景，一般会采用这种方式。然后，通过定时任务在业务不繁忙（比如凌晨）的时候做数据兜底

------
著作权归JavaGuide(javaguide.cn)所有
基于MIT协议
原文链接：https://javaguide.cn/high-performance/message-queue/kafka-questions-01.html