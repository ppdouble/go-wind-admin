# Kratos v3 服务端迁移

## 分支与更新

`main` 保持与 `tx7do/go-wind-admin:main` 一致，继续接收上游功能。`codex/kratos-v3` 从 `main` 建立，只维护服务端迁移。上游更新后，先更新 fork 的 `main`，再将 `main` 合入迁移分支，解决冲突并运行验证。不要把迁移分支合回 `main`，否则 `main` 不再是上游镜像。

前端源码、业务 protobuf、HTTP 路径和 JSON 请求响应契约没有迁移到新框架。上游新增 API 时，使用本分支的 v3 生成配置重新生成服务端 HTTP 与错误代码。

## v2 到 v3 的变化

| 范围 | v2 | v3 分支 |
|---|---|---|
| 框架模块 | `github.com/go-kratos/kratos/v2` | `github.com/go-kratos/kratos/v3`；服务端依赖图不得含 v2 |
| 日志 | Kratos `log.Logger` | 核心改用 `log/slog`，现有日志提供器由 `pkg/logbridge` 接入 |
| HTTP 与错误生成 | v2 插件 | v3 插件重新生成；`buf.kratos-v3.gen.yaml` 固定错误 proto 的 Go 导入映射 |
| JSON | v2 默认 codec | v3 codec 拆分，`pkg/jsoncompat` 保持现有 protojson 行为 |
| 链路追踪 | Kratos 核心 middleware | `kratos/contrib/otel/v3/tracing` |
| gRPC 客户端 | `DialInsecure` | `NewClient`，根据配置选择 TLS 凭据 |
| 限流与校验 | tx7do 扩展依赖 v2 类型 | 本地 v3 适配模块；保留 Aegis BBR 行为 |
| 认证 | v2 transport/context 读取 token | 从 v3 transport 请求头读取 token |

[Kratos 官方 v2→v3 迁移指南](https://go-kratos.com/docs/migration/v2-to-v3/)说明了框架层接口变化。

当前发布的 tx7do 扩展模块仍依赖 Kratos v2。本分支将服务端用到的模块放在 `backend/third_party/tx7do` 并保留许可证，通过 `go.mod` 的 `replace` 使用 v3 适配版本。`pkg/v2logcompat` 只包含旧日志接口的兼容源码，运行时不导入 Kratos v2。若这些模块正式支持 v3，可逐个移除本地适配。

## 生成、验证与 CI

在 `backend/` 下执行：

```sh
make api_v3
go build ./app/admin/service/cmd/server
go list -deps -f '{{.ImportPath}}' ./app/admin/service/cmd/server | grep '^github.com/go-kratos/kratos/v2/'
go test ./app/admin/service/internal/server ./app/admin/service/internal/service \
  ./pkg/middleware/auth ./pkg/middleware/ent ./pkg/middleware/logging \
  ./pkg/metadata ./pkg/jsoncompat ./pkg/logbridge ./pkg/jwt
```

第三条命令应无输出。`.github/workflows/kratos-v3.yml` 在 Linux 构建服务端、检查依赖图并运行上述定向测试；构建产物只保留 3 天。为限制 Actions 用量，迁移分支普通代码提交不自动运行该工作流。首次加入或修改工作流文件时运行一次，之后可推送 `kratos-v3-ci-*` 标签按需触发。`workflow_dispatch` 仅在工作流文件存在于默认分支时可从 GitHub 网页运行。上游原有工作流保留其触发条件；涉及其监听路径的迁移提交仍可能运行一次。

完整 `go test ./...` 仍有待解决的测试失败。迁移前的 `main` 提交 `534f6852` 与此前 v3 验证均复现：`internal/data` 的六个 SQLite `*_ExcludeExpired` 测试返回过期关联，`pkg/sseevent/TestPublishedEventsReachAllFrontends` 查找的前端 `event.ts` 路径不存在。相关实现和测试在迁移中没有修改，因此这些失败尚不能归因于 Kratos v3；其根因仍需单独确认。定向 CI 通过不等于全量测试通过。
