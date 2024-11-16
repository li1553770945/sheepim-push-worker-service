package trace

import (
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/li1553770945/sheepim-online-service/biz/infra/config"
)

type TraceStruct struct {
	Provider provider.OtelProvider
}

func InitTrace(config *config.Config) *TraceStruct {
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.ServerConfig.ServiceName),
		provider.WithExportEndpoint(config.OpenTelemetryConfig.Endpoint),
		provider.WithInsecure(),
	)
	return &TraceStruct{Provider: p}

}
