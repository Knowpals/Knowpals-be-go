# Knowpals Backend (Go)

教学场景下的后端服务：教师上传视频、班级任务分发、学生观看与答题、学情统计；并通过 **Kafka** 与 Python 侧流水线协作完成分段、知识点与习题等异步处理。

## 技术栈

- **语言 / 运行时**：Go（`go.mod` 要求 **Go 1.25+**）
- **Web**：Gin
- **ORM**：GORM + MySQL
- **缓存 / 验证码等**：Redis
- **对象存储**：腾讯云 COS（视频上传）
- **消息队列**：Kafka（Sarama）
- **依赖注入**：Google Wire
- **API 文档**：Swagger 注解 + `docs/` 生成物（需自行挂载或导出 JSON/YAML）

## 项目结构（简要）

```
├── api/              # Kafka 消息体、HTTP DTO
├── config/           # 配置结构与示例配置
├── controller/       # HTTP 入口，参数校验与 DTO 映射
├── service/          # 业务编排，错误封装
├── repository/dao/   # 数据访问
├── repository/model/ # GORM 模型
├── domain/           # 领域对象（供 service / controller 使用）
├── events/           # Kafka 生产/消费、流水线 Worker
├── web/              # 路由注册
├── ioc/              # DB / Redis / Kafka / Logger 等初始化
├── middleware/       # JWT、日志、CORS 等
├── errors/           # 统一错误码与 errorx 包装
├── wire.go           # Wire 注入图（仅 wireinject 构建标签）
└── wire_gen.go       # Wire 生成代码（勿手改逻辑，改完 wire.go 后重新生成）
```

分层约定：**Controller → Service → DAO**；Controller 不直连数据库；Service 内返回业务错误时尽量通过 `errors` 包封装。

## 核心业务流程

### 1. 视频处理流水线（Go ↔ Python）

- Go 向 Kafka **`task`** 主题投递任务消息（`api/message` 中定义）。
- Python 消费任务，按阶段产出 **`result`** 主题的结果消息。
- Go 侧 **`PipelineWorker`** 消费 **`result`**，调用 `PipelineService.ProcessResult` 将分段、知识点映射、习题等写入 MySQL。

主题常量见 `events/topic/topic.go`：

- `TASK_TOPIC = "task"`
- `RESULT_TOPIC = "result"`

> 消费端在 `ioc/kafka.go` 中调整了 **会话超时、心跳、MaxProcessingTime** 等，避免业务处理较慢时反复再均衡导致消息重复消费。

### 2. 业务功能概览

- 用户注册 / 登录 / 忘记密码（邮箱验证码，依赖 Redis + SMTP）
- 班级创建、加入、学生列表
- 教师上传视频至 COS、将视频任务下发到班级
- 学生观看行为上报、进度更新、按班级查询进度
- 学生答题、按视频拉取习题
- 学情统计（个人 / 班级、总览等）

## 环境依赖

- **MySQL**：与 `config` 中 DSN 一致的数据库。
- **Redis**：验证码、会话相关能力。
- **Kafka**：`task` / `result` 主题需存在或由集群自动创建（视集群配置而定）。
- **SMTP**：发信（注册/找回密码等）。
- **腾讯云 COS**：视频上传。

## 配置说明

1. 复制示例配置并修改：

   ```bash
   cp config/config-example.yaml config/config.yaml
   ```

2. 启动时通过 **文件路径** 指定配置（二选一）：

   - 命令行：`go run . --config ./config/config.yaml`
   - 环境变量：`CONFIG_PATH=/path/to/config.yaml`

3. 配置字段与结构体对应关系见 `config/conf.go`。注意示例文件里个别键名可能与结构体 `yaml` 标签不一致，例如 **COS** 在代码中为 `secretID` / `secretKey` / `url`（见 `COSConf`），请与 `config-example.yaml` 对照修正。

主要段落：

| 段落 | 用途 |
|------|------|
| `mysql` | 数据库 DSN |
| `jwt` | 签名与过期时间 |
| `redis` | 缓存 |
| `smtp` | 邮件 |
| `cos` | 对象存储 |
| `kafka` | Broker 地址与消费者组 ID |
| `otel` | 可选链路/指标导出 |

## 本地运行

```bash
# 使用 Go 1.25+（或与 go.mod 一致的 toolchain）
go mod download

# 确保 config/config.yaml 已就绪
go run . --config ./config/config.yaml
```

默认 HTTP 地址为 Gin 默认 **`0.0.0.0:8080`**（未在代码中显式修改时）。

### 重新生成 Wire（修改 wire.go 或构造函数签名后）

```bash
go generate ./...
```

若环境无法拉取 toolchain，可在本机安装匹配版本的 Go 后再执行。

## HTTP API 前缀

所有业务接口挂在 **`/api/v1`** 下（见 `web/web.go`）。

### 用户 `/api/v1/user`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/register` | 注册 |
| POST | `/sendCode` | 发验证码 |
| POST | `/loginByPassword` | 密码登录 |
| POST | `/loginByCode` | 验证码登录 |
| POST | `/forgotPassword` | 忘记密码 |
| GET | `/getUser/:id` | 用户信息（需登录） |

### 班级 `/api/v1/class`（需登录）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/create` | 创建班级 |
| POST | `/join` | 加入班级 |
| POST | `/quit/:class_id` | 退出班级 |
| GET | `/info/:class_id` | 班级信息 |
| GET | `/my-created` | 我创建的班级 |
| GET | `/my-joined` | 我加入的班级 |
| GET | `/students/:class_id` | 班级学生 |

### 视频 `/api/v1/video`（需登录）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/upload` | 上传视频 |
| GET | `/getDetail/:video_id` | 视频详情（分段、知识点、题目等） |
| POST | `/post-to-class` | 下发到班级 |
| GET | `/getTasks/:class_id` | 班级视频任务列表 |
| POST | `/task/process` | 查询上传/流水线任务进度 |
| GET | `/my-uploaded` | 教师本人上传的视频列表 |

### 题目 `/api/v1/question`（需登录）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/answer` | 学生答题 |
| GET | `/generate/:video_id` | 获取视频相关习题 |

### 行为与进度 `/api/v1/behavior`（需登录）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/record` | 观看行为（暂停/回放等） |
| POST | `/update-progress` | 更新观看进度 |
| GET | `/class-progress/:class_id/:status` | 班级内视频进度筛选 |
| GET | `/my/unfinished` | 学生未完成任务（含班级信息） |

### 统计 `/api/v1/stat`（需登录）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/student/:video_id` | 学生单视频学情 |
| GET | `/student/overview` | 学生总览（总观看时长、完成数/总数、总正确率等） |
| GET | `/class` | 班级学情（请求体参数见 Swagger / `api/http/statistic`） |

更完整的请求/响应模型可参考 `api/http/` 下各包及 `docs/swagger.json`。

## Kafka 与 Worker

进程启动时，除 Gin HTTP 外，会在后台启动 **流水线结果消费者**（`main.go` 中 `App.Run`），订阅 **`result`** 主题并将结果持久化。

请保证：

- `config.kafka.addrs` 可达；
- `consumerGroup` 在环境中唯一或符合运维规范；
- Python 与 Go 对 **`task` / `result`** 的 JSON 字段约定一致（见 `api/message/message.go` 注释）。

## 开发与规范建议

- 新增接口时保持 **Controller → Service → DAO**，并在 `errors` 中为 Service 层补充可映射到 HTTP 的错误构造器。
- 修改构造函数依赖后运行 **`go generate ./...`** 更新 `wire_gen.go`。
- 数据库批量写入注意 MySQL + GORM 的 upsert / 关联字段命名陷阱（历史问题可参考 `KnowledgeSegmentMapping` 等模型的字段命名）。

## 许可证

若仓库未包含 `LICENSE` 文件，使用前请与项目维护者确认授权范围。
