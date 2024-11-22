package repo

import (
	"context"
	"github.com/li1553770945/sheepim-push-worker-service/biz/infra/kafka"
)

type IRepository interface {
	FetchMessage(ctx context.Context) ([]byte, []byte, error) // 返回 Key、Value 和错误

}

type Repository struct {
	KafkaClient *kafka.KafkaClient
}

func NewRepository(kafkaClient *kafka.KafkaClient) IRepository {
	return &Repository{
		KafkaClient: kafkaClient,
	}
}
