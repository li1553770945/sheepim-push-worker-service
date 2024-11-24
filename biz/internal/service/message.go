package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/li1553770945/sheepim-connect-service/kitex_gen/message"
	"github.com/li1553770945/sheepim-online-service/kitex_gen/online"
	"github.com/li1553770945/sheepim-push-proxy-service/kitex_gen/push_proxy"
	"github.com/li1553770945/sheepim-room-service/kitex_gen/room"
	"log"
)

func (s *MessageHandlerService) handler(ctx context.Context, keyBytes []byte, valueBytes []byte) error {
	value := string(valueBytes)
	var returnError error
	returnError = nil
	var messageObj push_proxy.PushMessageReq
	err := json.Unmarshal([]byte(value), &messageObj)
	if err != nil {
		log.Fatalf("反序列化失败: %v", err)
	}
	roomRpcResp, err := s.RoomClient.GetRoomMembers(ctx, &room.GetRoomMembersReq{RoomId: messageObj.RoomId})
	if err != nil {
		return err
	}
	if roomRpcResp.BaseResp.Code != 0 {
		return errors.New(roomRpcResp.BaseResp.Message)
	}

	members := roomRpcResp.Members

	onlineRpcResp, err := s.OnlineClient.GetOnlineMemberEndpoint(ctx, &online.GetOnlineMemberEndpointReq{
		ClientIdList: members,
	})
	if err != nil {
		return err
	}
	if onlineRpcResp.BaseResp.Code != 0 {
		return errors.New(onlineRpcResp.BaseResp.Message)
	}
	klog.CtxInfof(ctx, "收到来自%s的消息：%s", messageObj.ClientId, messageObj.Message)
	onlineMembers := onlineRpcResp.Status
	for _, onlineMember := range onlineMembers {
		klog.CtxInfof(ctx, "找到客户端：%s", onlineMember.ClientId)

		if onlineMember.ClientId == messageObj.ClientId {
			continue
		}
		klog.CtxInfof(ctx, "发送到客户端：%s,%s", onlineMember.ClientId, onlineMember.ServerEndpoint)
		endpoint := onlineMember.ServerEndpoint
		sendMessageResp, err := s.ConnectClient.SendMessage(ctx, &message.SendMessageReq{
			ClientId: onlineMember.ClientId,
			Event:    messageObj.Event,
			Type:     messageObj.Type,
			Message:  messageObj.Message,
		}, callopt.WithHostPort(endpoint))
		if err != nil {
			returnError = err
			klog.CtxErrorf(ctx, "消息发送失败%v", err)
		}
		if sendMessageResp == nil {
			returnError = errors.New("未收到回复resp")
			klog.CtxErrorf(ctx, "未收到回复resp")
			continue
		}

		if sendMessageResp.BaseResp.Code != 0 {
			returnError = errors.New(fmt.Sprintf("发送失败：%s", sendMessageResp.BaseResp.Message))
			klog.CtxErrorf(ctx, "发送失败：%s", sendMessageResp.BaseResp.Message)
		}

	}

	return returnError
}
func (s *MessageHandlerService) HandleMessage() {

	klog.CtxInfof(context.Background(), "开始进行消息消费循环")
	for {
		// 从 Repository 获取消息
		ctx := context.Background()
		key, value, err := s.Repo.FetchMessage(ctx)
		if err != nil {
			klog.CtxErrorf(ctx, "Kafka 消费消息失败: %v", err)
			break
		}

		// 调用传入的消息处理函数
		err = s.handler(ctx, key, value)
		if err != nil {
			klog.CtxErrorf(ctx, "消息处理失败: %v\n", err)
		} else {
			klog.CtxInfof(ctx, "消息处理成功")
		}

	}
}
