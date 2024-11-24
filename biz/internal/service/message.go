package service

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/li1553770945/sheepim-connect-service/kitex_gen/message"
	"github.com/li1553770945/sheepim-online-service/kitex_gen/online"
	"github.com/li1553770945/sheepim-push-proxy-service/kitex_gen/push_proxy"
	"github.com/li1553770945/sheepim-room-service/kitex_gen/room"
)

func (s *MessageHandlerService) handler(ctx context.Context, keyBytes []byte, valueBytes []byte) {
	value := string(valueBytes)
	var messageObj push_proxy.PushMessageReq
	err := json.Unmarshal([]byte(value), &messageObj)
	if err != nil {
		klog.CtxErrorf(ctx, "反序列化失败: %v", err)
		return
	}
	roomRpcResp, err := s.RoomClient.GetRoomMembers(ctx, &room.GetRoomMembersReq{RoomId: messageObj.RoomId})
	if err != nil {
		klog.CtxErrorf(ctx, "调用room服务失败: %v", err)
		return
	}
	if roomRpcResp.BaseResp.Code != 0 {
		klog.CtxErrorf(ctx, "反序列化失败: %s", roomRpcResp.BaseResp.Message)
		return
	}

	members := roomRpcResp.Members

	onlineRpcResp, err := s.OnlineClient.GetOnlineMemberEndpoint(ctx, &online.GetOnlineMemberEndpointReq{
		ClientIdList: members,
	})
	if err != nil {
		klog.CtxErrorf(ctx, "调用online服务失败: %v", err)
		return
	}
	if onlineRpcResp.BaseResp.Code != 0 {
		klog.CtxErrorf(ctx, "反序列化失败: %s", onlineRpcResp.BaseResp.Message)
		return
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
			klog.CtxErrorf(ctx, "消息发送失败%v", err)
			continue
		}
		if sendMessageResp == nil {
			klog.CtxErrorf(ctx, "未收到回复resp")
			continue
		}

		if sendMessageResp.BaseResp.Code != 0 {
			klog.CtxErrorf(ctx, "发送失败：%s", sendMessageResp.BaseResp.Message)
			continue
		}

	}

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
		s.handler(ctx, key, value)

	}
}
