package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/li1553770945/sheepim-connect-service/kitex_gen/message/messageservice"
	"github.com/li1553770945/sheepim-push-worker-service/biz/infra/config"
)

func NewConnectClient(config *config.Config) messageservice.Client {
	r, err := etcd.NewEtcdResolver(config.EtcdConfig.Endpoint)
	Client, err := messageservice.NewClient(
		config.RpcConfig.ConnectServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.ServerConfig.ServiceName}),
	)
	if err != nil {
		panic("认证 RPC 客户端启动失败" + err.Error())
	}
	return Client
}
