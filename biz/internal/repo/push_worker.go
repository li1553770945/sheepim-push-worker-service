package repo

import (
	"context"
)

func (r *Repository) FetchMessage(ctx context.Context) ([]byte, []byte, error) {
	// 调用 Kafka 消费者读取消息
	msg, err := r.KafkaClient.Consumer.ReadMessage(ctx)
	if err != nil {
		return nil, nil, err
	}
	return msg.Key, msg.Value, nil // 返回消息的 Key 和 Value
}
