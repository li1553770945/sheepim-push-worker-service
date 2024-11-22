# sheepim用户服务

## 初始化项目
```bash
kitex -module "github.com/li1553770945/sheepim-push-worker-service" -service sheepim-push-worker-service idl/project.thrift
cd biz/infra/container
wire
```
## 配置文件示例

```yml
server:
  listen-address: 192.168.6.241:8897
  service-name: sheepim-push-worker-service

etcd:
  endpoint:
    - "xxx:2379"

open-telemetry:
  endpoint: "xxx:4417"


kafka:
  brokers:
    - "xxx:9092"
  topic: messages
  group-id: group1


```

## 开发环境

```bash
export ENV=development
```

## 生产环境

```bash
export ENV=production
```