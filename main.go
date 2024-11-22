// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/li1553770945/sheepim-push-worker-service/biz/infra/container"
	"os"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	container.InitGlobalContainer(env)
	App := container.GetGlobalContainer()

	defer func(p provider.OtelProvider, ctx context.Context) {
		err := p.Shutdown(ctx)
		if err != nil {
			klog.Fatalf("server stopped with error:%s", err)
		}
	}(App.Trace.Provider, context.Background())

	App.MessageHandlerService.HandleMessage()

}
