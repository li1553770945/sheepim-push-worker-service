package container

import (
	"github.com/li1553770945/sheepim-online-service/biz/infra/config"
	"github.com/li1553770945/sheepim-online-service/biz/infra/log"
	"github.com/li1553770945/sheepim-online-service/biz/infra/trace"
	"github.com/li1553770945/sheepim-online-service/biz/internal/service"
	"sync"
)

type Container struct {
	Trace          *trace.TraceStruct
	Logger         *log.TraceLogger
	Config         *config.Config
	ProjectService project.IProjectService
}

var APP *Container
var once sync.Once

func GetGlobalContainer() *Container {
	if APP == nil {
		panic("APP在使用前未初始化")
	}
	return APP
}

func InitGlobalContainer(env string) {
	once.Do(func() {
		APP = GetContainer(env)
	})
}

func NewContainer(config *config.Config,
	logger *log.TraceLogger,
	traceStruct *trace.TraceStruct,

	projectService project.IProjectService,
) *Container {
	return &Container{
		Config:         config,
		Logger:         logger,
		ProjectService: projectService,
		Trace:          traceStruct,
	}

}
