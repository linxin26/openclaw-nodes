# Repository Guidelines

## 项目结构与模块组织
`cmd/` 存放 CLI 与 Wails 构建入口。核心 Go 代码位于 `internal/`：`internal/protocol` 负责网关通信，`internal/device` 负责设备能力处理，`internal/wails` 提供桌面绑定，`internal/tray` 处理托盘 UI。持久化相关工具在 `store/`。前端位于 `frontend/src/`，按 `components/`、`hooks/`、`lib/` 与 `pages/` 组织。架构说明和实现计划集中在 `docs/superpowers/`。

## 构建、测试与开发命令
在仓库根目录运行 `go test ./...` 执行全部 Go 测试。调试协议逻辑时可使用 `go test -v ./internal/protocol/`。通过 `go build -o openclaw-node.exe ./cmd` 构建 Windows 可执行文件。进入 `frontend/` 后，`npm run dev` 启动 Vite 开发监听，`npm run build` 负责 TypeScript 检查与前端打包。若本机已安装 Wails CLI，可使用 `wails dev` 启动桌面开发环境，使用 `wails build` 按 `wails.json` 打包应用。

## 编码风格与命名约定
Go 代码遵循 `gofmt` 默认格式，使用制表符缩进，包名保持小写，导出标识符使用 `PascalCase`。React 与 TypeScript 文件使用 2 空格缩进，组件名使用 `PascalCase`，hooks 与工具函数使用 `camelCase`，例如 `useConnection`、`normalizeGatewayAddress`。前端内部导入优先使用 `@/` 别名。Go lint 使用 `golangci-lint run --no-config --disable-all -E errcheck ./...`。

## 测试指南
测试文件与实现代码同目录存放，统一使用 Go 自带的 `testing` 包。测试文件命名为 `*_test.go`，测试函数保持聚焦，例如 `TestClientCanReconnectAfterDisconnect`。修改协议流程、配置读写、Wails 绑定或设备命令注册时，应同步补充或更新测试。当前仓库未配置独立的前端测试框架，因此 UI 变更至少要确认 `npm run build` 成功。

## 提交与 Pull Request 规范
最近提交历史采用 Conventional Commits，例如 `fix: ...`、`docs: ...`、`merge: ...`。继续使用简洁的 `<type>: <description>` 主题格式。Pull Request 需要说明用户可见影响、列出已执行的验证命令，并关联对应 issue 或计划文档。涉及 `frontend/` 或 Wails 界面改动时，应附带截图。

## 安全与配置提示
不要提交生成的二进制文件、运行日志或机器本地密钥。`identity.json`、`config.yaml`、网关 token 以及 `%APPDATA%\\OpenClaw` 下的内容都应视为本地专用数据。
