package service

import (
	"context"
	"github.com/li1553770945/sheepim-connect-service/kitex_gen/message/messageservice"
	"github.com/li1553770945/sheepim-online-service/kitex_gen/online/onlineservice"
	"github.com/li1553770945/sheepim-push-worker-service/biz/internal/repo"
	"github.com/li1553770945/sheepim-room-service/kitex_gen/room/roomservice"
)

type MessageHandlerService struct {
	Repo          repo.IRepository
	OnlineClient  onlineservice.Client
	RoomClient    roomservice.Client
	ConnectClient messageservice.Client
}

type IMessageHandlerService interface {
	HandleMessage()
	handler(context.Context, []byte, []byte) error
}

func NewMessageHandlerService(repo repo.IRepository,
	onlineClient onlineservice.Client,
	roomClient roomservice.Client,
	connectClient messageservice.Client,
) IMessageHandlerService {
	return &MessageHandlerService{
		Repo:          repo,
		OnlineClient:  onlineClient,
		RoomClient:    roomClient,
		ConnectClient: connectClient,
	}
}
