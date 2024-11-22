package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func (s *MessageHandlerService) handler(keyBytes []byte, valueBytes []byte) error {
	key := string(keyBytes)
	value := string(valueBytes)
	fmt.Println(key, value)
	return nil
}
func (s *MessageHandlerService) HandleMessage() {

	for {
		// 从 Repository 获取消息
		ctx := context.Background()
		key, value, err := s.Repo.FetchMessage(ctx)
		if err != nil {
			klog.CtxErrorf(ctx, "Kafka 消费消息失败: %v", err)
			break
		}

		// 调用传入的消息处理函数
		err = s.handler(key, value)
		if err != nil {
			klog.CtxErrorf(ctx, "消息处理失败: %v\n", err)
		}
	}
}
