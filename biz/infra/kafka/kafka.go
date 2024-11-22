package kafka

import (
	"context"
	"github.com/li1553770945/sheepim-push-worker-service/biz/infra/config"
	"github.com/segmentio/kafka-go"
	"log"
)

// KafkaClient 封装了 Kafka 的生产者和消费者
type KafkaClient struct {
	Producer *kafka.Writer
	Consumer *kafka.Reader
}

// 创建 Kafka Admin Client，用于管理主题
func ensureTopicExists(brokers []string, topic string, numPartitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", brokers[0]) // 使用第一个 Broker 进行管理
	if err != nil {
		return err
	}
	defer func(conn *kafka.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}

	// 遍历分区元数据，检查主题是否已存在
	for _, p := range partitions {
		if p.Topic == topic {
			log.Printf("Topic %s 已存在", topic)
			return nil // 主题已存在，无需创建
		}
	}

	// 创建主题
	log.Printf("创建新的 Topic: %s", topic)
	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	})
}

// NewKafkaClient 初始化 Kafka 生产者和消费者
func NewKafkaClient(cfg *config.Config) *KafkaClient {
	conf := cfg.KafkaConfig
	if len(conf.Brokers) == 0 || conf.Topic == "" {
		panic("Kafka 配置错误：Brokers 或 Topic 未设置")
	}
	err := ensureTopicExists(conf.Brokers, conf.Topic, 1, 1) // 设置分区数和副本因子
	if err != nil {
		panic("创建 Kafka Topic 失败: " + err.Error())
	}
	// 创建生产者
	producer := kafka.Writer{
		Addr:     kafka.TCP(conf.Brokers...),
		Topic:    conf.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	// 创建消费者
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  conf.Brokers,
		GroupID:  conf.GroupID,
		Topic:    conf.Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	// 返回封装的 KafkaClient
	return &KafkaClient{
		Producer: &producer,
		Consumer: consumer,
	}
}

// ProduceMessage 发送消息到 Kafka
func (kc *KafkaClient) ProduceMessage(ctx context.Context, key, value []byte) error {
	err := kc.Producer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		log.Printf("Kafka 生产消息失败: %v\n", err)
		return err
	}
	log.Println("消息生产成功")
	return nil
}

// ConsumeMessages 消费 Kafka 消息
func (kc *KafkaClient) ConsumeMessages(ctx context.Context, handler func(key, value []byte) error) {
	for {
		msg, err := kc.Consumer.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka 消费消息失败: %v\n", err)
			break
		}

		// 调用处理函数
		err = handler(msg.Key, msg.Value)
		if err != nil {
			log.Printf("消息处理失败: %v\n", err)
		}
	}
}
