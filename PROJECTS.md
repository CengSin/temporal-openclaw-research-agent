# OpenClaw Workflow Usecase

## 项目目标

使用 Temporal 对 OpenClaw Agent 能力进行工作流编排，当前先实现一个最小可运行链路：

核心目标：
- 由 Workflow 触发本地 OpenClaw 的 `Chronicle` agent。
- 把“文章主题”交给 agent，并要求其调用 skills 生成一篇短文。
- 工作流执行失败时自动重试（基础 RetryPolicy）。
- 提供可直接运行的本地测试入口。

---

## 当前代码结构

```text
.
├── main.go
├── namespace/
│   └── namespace.go
├── worker/
│   └── worker.go
└── workflow/
    ├── workflow.go
    ├── workflow_integration_test.go
    └── activity/
        ├── activity.go
        └── model.go
```

说明：
- `namespace/namespace.go`：定义 Temporal Task Queue 名称。
- `main.go`：开发态一体化入口，会先启动内嵌 Worker，再提交 Workflow 并等待文章结果。
- `worker/worker.go`：独立 Worker 启动入口，适合与 Workflow 提交端分离部署。
- `workflow/workflow.go`：Temporal 工作流（执行 Chronicle 生成文章 Activity，并携带重试策略）。
- `workflow/activity/activity.go`：通过 `openclaw-go/chatcompletions` 调本地 OpenClaw。
- `workflow/workflow_integration_test.go`：本地集成测试（可选开关执行，避免 CI/本地无服务时误失败）。

---

## OpenClaw 文档对齐结论

- 使用 `chatcompletions.Client` 调用 `POST /v1/chat/completions`。
- 通过 `AgentID: "Chronicle"` 自动带 `x-openclaw-agent-id`，实现定向到本地 Chronicle agent。
- `Model` 使用 `openclaw:chronicle`，主题与“调用 skills 写文章”要求放在 user prompt 中。
- 支持 `OPENCLAW_BASE_URL` 与 `OPENCLAW_TOKEN` 环境变量。

---

## 当前实现状态（已完成）

- ✅ 修复原先无法编译的问题（未定义类型/函数、注册不一致、空入口）。
- ✅ 新建最小 Workflow：`ChronicleArticleWorkflow(topic string) -> article string`。
- ✅ 新建 Activity：`GenerateArticleWithChronicle(topic string)` 调用本地 OpenClaw。
- ✅ Worker 已注册新 Workflow 与 Activity。
- ✅ 根目录入口已支持内嵌启动 Worker，`go run . "<主题>"` 可直接本地联调。
- ✅ 新增集成测试：`workflow/workflow_integration_test.go`。
- ✅ `go test ./...` 可通过（默认跳过真正的 OpenClaw 集成调用）。

---

## 本地运行说明

### 1) 直接本地运行（一体化模式）

```bash
go run . "半导体行业周报"
```

### 2) 分离运行（可选）

```bash
go run ./worker
go run . "半导体行业周报"
```

### 3) 执行 OpenClaw 集成测试（可选）

```bash
OPENCLAW_INTEGRATION=1 OPENCLAW_TEST_TOPIC="消费电子行业观察" go test ./workflow -run TestChronicleArticleWorkflow_Integration -v
```

可选环境变量：
- `OPENCLAW_BASE_URL`（默认 `http://localhost:18789`）
- `OPENCLAW_TOKEN`（本地网关若需要鉴权时填写）

## 下一步建议

- 将“skills 执行结果”结构化输出（标题、摘要、正文、标签）而非纯文本。
- 增加 Workflow 输入参数对象（主题、字数、语气、目标读者）。
- 在 Activity 中补充超时与错误分类（鉴权错误/网络错误/模型错误）。
- 下一阶段再接入人工审核 signal 与 RAG 入库。

