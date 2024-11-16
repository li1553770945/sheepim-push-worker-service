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
  listen-address: 0.0.0.0:8890
  service-name: sheepim-push-worker-service

etcd:
  endpoint:
    - 127.0.0.1:2379

open-telemetry:
  endpoint: 127.0.0.1:4317

database:
  username: xxx
  password: xxx
  database: xxx
  address: xxx
  port: xxx

```

## 开发环境

```bash
export ENV=development
```

## 生产环境

```bash
export ENV=production
```