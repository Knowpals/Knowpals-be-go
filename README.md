# Knowpals-be-go

数据类型：
用户注册表：
id username phone_number password role

## OTel + Grafana

1. 复制配置：
   `cp config/config-example.yaml config/config.yaml`
2. 启动观测组件：
   `docker compose -f docker/docker-compose.yaml up -d`
3. 启动应用：
   `go run .`
4. 发起请求后查看：
   Grafana: `http://localhost:3000`
   账号密码: `admin / admin`
5. 在 Grafana 中：
   进入 `Explore`
   选择 `Tempo` 查看链路
   选择 `Prometheus` 查看指标
