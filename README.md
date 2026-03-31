# OpenClaw Node 桌面节点

> 版本说明：当前仓库仅在 OpenClaw `2026.3.8` 版本上完成过联调与验证，其他 OpenClaw 版本暂未测试，兼容性请自行评估。

OpenClaw Node 是一个基于 Go + Wails + React 的桌面节点应用，用于将桌面设备接入 OpenClaw 网关，并向网关暴露设备能力、运行状态与操作入口。

项目同时提供两种运行形态：

- `CLI` 模式：适合调试协议连接、参数注入和托盘行为
- `Wails` 桌面模式：提供图形界面，用于连接配置、能力开关、操作日志和运行信息查看

## 核心能力

当前仓库围绕“桌面节点接入网关”实现了以下能力：

- 连接 OpenClaw Gateway，支持手动指定网关地址、Token 和 TLS
- 在局域网内通过 mDNS 自动发现网关
- 注册并上报节点身份、能力集合、命令集合和权限信息
- 提供桌面端可控能力开关
- 提供桌面 GUI，用于查看连接状态、节点信息、最近活动和能力列表
- 支持托盘驻留和后台运行

当前能力模块主要包括：

- `camera`
- `location`
- `photos`
- `screen`
- `motion`
- `notifications`
- `sms`
- `calendar`

不同平台下的能力实现覆盖度不同，具体以 `internal/device/platform/` 下的实现为准。

## 技术栈

- 后端：Go 1.20
- 桌面框架：Wails v2
- 前端：React 18 + TypeScript + Vite + Tailwind CSS
- 通信：WebSocket + 自定义协议层

## 项目结构

```text
.
├─ main.go                 # Wails 桌面入口
├─ cmd/                    # CLI 入口
├─ internal/
│  ├─ protocol/            # 网关通信、连接与命令分发
│  ├─ device/              # 设备能力注册、平台能力实现
│  ├─ wails/               # Wails 绑定与桌面应用接口
│  ├─ tray/                # 托盘与桌面交互
│  ├─ config/              # 配置加载与保存
│  └─ discovery/           # mDNS 自动发现
├─ store/                  # 数据目录与持久化辅助
├─ frontend/               # React + Vite 前端
├─ docs/superpowers/       # 设计文档与实现计划
└─ wails.json              # Wails 构建配置
```

## 快速开始

### 1. 环境要求

- Go `1.20.x`
- Node.js `18+`
- npm
- Windows 桌面环境
- Wails CLI

### 2. 安装项目依赖

除了 `go.mod` 中声明的 Go 模块依赖外，本项目还依赖 Wails CLI 作为桌面开发和打包工具。建议在开始开发前统一安装：

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@v2.7.1
cd frontend
npm install
cd ..
```

说明：

- 当前仓库 `go.mod` 使用 `Go 1.20`
- 为了避免新版 Wails CLI 与本机 Go 版本不兼容，建议优先使用 `v2.7.1`
- `npm install` 用于安装前端依赖

如果你希望在任意目录直接执行 `wails`，还需要确保 Go 的 `bin` 目录已加入 `PATH`。

### 3. 运行 Go 测试

```powershell
go test ./...
```

### 4. 构建 CLI 可执行文件

```powershell
go build -o openclaw-node.exe ./cmd
```

生成产物位于仓库根目录：

- `openclaw-node.exe`

## 桌面开发与打包

### Wails 开发模式

```powershell
wails dev
```

说明：

- Wails 桌面模式的正式入口是仓库根目录 `main.go`
- `wails dev` / `wails build` 会使用该入口，并打包 `frontend/dist` 中的前端资源
- 不建议把 `go run -tags wails ./cmd` 当作正式 GUI 启动方式

如果 `wails` 没有加入 `PATH`，也可以直接使用绝对路径，例如：

```powershell
& "$env:USERPROFILE\\go\\bin\\wails.exe" dev
```

### Wails 生产打包

```powershell
wails build
```

同样，如果 `PATH` 未配置：

```powershell
& "$env:USERPROFILE\\go\\bin\\wails.exe" build
```

默认打包产物位于：

- `build/bin/openclaw-node.exe`

说明：

- 该产物包含 Go 后端、Wails 绑定和前端构建资源
- 正常情况下，最终用户可直接打开 `build/bin/openclaw-node.exe` 使用
- 如果你的目标是桌面分发，应始终以 `wails build` 产物为准

## 发布产物说明

仓库里常见的可执行文件主要有两类，不建议混用理解：

### 1. CLI 构建产物

通过下面命令生成：

```powershell
go build -o openclaw-node.exe ./cmd
```

产物位置：

- 仓库根目录 `openclaw-node.exe`

特点：

- 使用 CLI 入口 `cmd/main.go`
- 更适合协议调试、参数注入、日志观察和命令行启动
- 可直接配合 `-gateway`、`-token`、`-tls`、`-no-mdns` 使用
- 不包含桌面 GUI 窗口

适用场景：

- 本地调试协议连接
- 验证 CLI 启动逻辑
- 做自动化或脚本化测试

### 2. Wails 桌面打包产物

通过下面命令生成：

```powershell
wails build
```

产物位置：

- `build/bin/openclaw-node.exe`

特点：

- 使用 Wails 桌面入口 `main.go`
- 包含前端 GUI
- 适合最终用户直接运行
- 配置、连接、能力开关和日志查看都通过桌面界面完成

适用场景：

- 内部验收
- 桌面分发
- 给非开发用户使用

### 应该分发哪个产物

一般建议：

- 给开发者或调试环境：优先使用根目录的 CLI 产物
- 给实际使用者：优先使用 `build/bin/openclaw-node.exe`

如果你的目标是“交付一个可直接打开并配置的桌面应用”，应以 `wails build` 的产物为准。

## CLI 模式

仓库默认 `cmd/main.go` 提供 CLI 启动入口，可通过命令行参数覆盖连接配置：

```powershell
.\openclaw-node.exe -gateway 127.0.0.1:18789 -token your-token -tls
```

可用参数：

- `-gateway`：网关地址，格式为 `host:port`
- `-token`：网关鉴权 Token
- `-tls`：启用 TLS
- `-no-mdns`：关闭 mDNS 自动发现

说明：

- 如果指定 `-gateway`，程序会切换到手动发现模式
- 如果未配置网关地址，程序可依赖配置文件或 mDNS 自动发现

## 前端开发

```powershell
cd frontend
npm run dev
```

前端构建命令：

```powershell
cd frontend
npm run build
```

当前前端主要页面包括：

- `Dashboard`：连接摘要、节点状态和最近活动
- `Connection`：网关连接配置
- `Capabilities`：能力开关管理
- `Operations`：设备操作入口
- `Logs`：运行日志
- `About`：设备标识、构建元数据和数据目录信息

## 配置与数据目录

程序默认将本地运行数据保存在：

```text
%APPDATA%\OpenClaw
```

主要文件包括：

- `identity.json`：设备身份与密钥材料
- `config.yaml`：本地配置

### GUI 配置如何存储

GUI 与 CLI 共用同一份本地数据目录，默认位置为：

```text
%APPDATA%\OpenClaw
```

其中：

- `config.yaml`：保存 GUI 中的连接配置、发现模式、能力开关和能力参数
- `identity.json`：保存节点身份、设备 ID 和密钥材料

在 GUI 中执行以下操作时，配置会被直接写入 `config.yaml`：

- 在 `Connection` 页面点击“保存”
- 在 `Capabilities` 页面切换能力开关

配置保存后，GUI 会同步更新内存中的运行配置；部分设置会立即影响后续连接行为。

### GUI 配置与 CLI 参数的关系

两者的关系可以理解为：

- `config.yaml` 提供默认配置
- CLI 参数提供“当前进程级覆盖”

也就是说：

- GUI 保存的内容会持久化到磁盘
- CLI 参数只影响当前这次启动，不会自动回写到 `config.yaml`

常见覆盖关系如下：

- `-gateway`：覆盖本次启动使用的网关地址
- `-tls`：覆盖本次启动使用的 TLS 选项
- `-token`：作为本次连接使用的 Token
- `-no-mdns`：强制本次启动关闭 mDNS 自动发现

如果你既使用 GUI 又使用 CLI，建议把长期配置保存在 GUI 或 `config.yaml` 中，把临时调试参数放在 CLI 中。

### `config.yaml` 示例

下面是一份可直接参考的配置示例：

```yaml
gateway: gateway.example.com
port: 18789
token: your-gateway-token
tls: true
discovery: manual
capabilities:
  camera: true
  location: true
  photos: true
  screen: true
  motion: false
  notifications: true
  sms: false
  calendar: false
capabilityOptions: {}
```

字段说明：

- `gateway`：网关地址，可以写主机名、`host:port`，也可以写完整协议地址
- `port`：当 `gateway` 未显式包含端口时使用
- `token`：网关鉴权 Token
- `tls`：为 `true` 时优先使用安全 WebSocket 连接
- `discovery`：`auto` 表示启用 mDNS 自动发现，`manual` 表示使用手动配置
- `capabilities`：控制节点启动时默认启用的能力
- `capabilityOptions`：能力的扩展参数，目前保留给平台相关实现使用

### 常见连接写法

以下几种写法都可以：

```yaml
# 写主机名，由 port 补齐端口
gateway: gateway.example.com
port: 18789
tls: true
```

```yaml
# 直接写 host:port
gateway: 192.168.1.50:18789
port: 18789
tls: false
```

```yaml
# 直接写完整 URL
gateway: https://gateway.example.com
port: 443
tls: true
```

连接规则说明：

- 如果 `gateway` 不带端口，程序会用 `port` 自动补齐
- 如果 `gateway` 已经写成 `host:port`，则优先使用该端口
- 如果 `gateway` 写成 `http://` 或 `https://`，程序会在内部转换成 `ws://` 或 `wss://`
- `gateway` 为空时，GUI 连接测试默认回落到 `ws://localhost:18789`

### CLI 覆盖配置示例

CLI 参数会覆盖部分本地配置，常见用法如下：

```powershell
.\openclaw-node.exe -gateway gateway.example.com:18789 -token your-token -tls
```

```powershell
.\openclaw-node.exe -gateway 192.168.1.50:18789 -no-mdns
```

说明：

- `-gateway` 会覆盖 `config.yaml` 中的 `gateway`
- `-tls` 会覆盖 `config.yaml` 中的 `tls`
- `-token` 作为本次进程连接使用的 Token
- 指定 `-gateway` 或 `-no-mdns` 时，发现模式会转为手动模式

默认配置重点如下：

- 默认端口：`18789`
- 默认发现模式：`auto`
- 默认开启：`camera`、`location`、`photos`、`screen`、`notifications`
- 默认关闭：`motion`、`sms`、`calendar`

这些文件属于本地私有数据，不应提交到版本库。

## 测试与质量检查

执行全部 Go 测试：

```powershell
go test ./...
```

调试协议层时可单独运行：

```powershell
go test -v ./internal/protocol/
```

前端变更后，至少应确认：

```powershell
cd frontend
npm run build
```

可选的 Go lint 命令：

```powershell
golangci-lint run --no-config --disable-all -E errcheck ./...
```

## 开发约定

- Go 代码使用 `gofmt`
- 前端组件使用 `PascalCase`
- hooks 与工具函数使用 `camelCase`
- 前端内部导入优先使用 `@/` 别名
- 提交信息遵循 Conventional Commits，例如 `fix: ...`、`docs: ...`

## 安全说明

不要提交以下内容：

- 生成的二进制文件
- 运行日志
- `identity.json`
- `config.yaml`
- 网关 Token
- `%APPDATA%\\OpenClaw` 下的本地数据

## 故障排查

### 1. `wails` 命令不存在

如果执行 `wails dev` 或 `wails build` 时提示：

```text
wails : The term 'wails' is not recognized ...
```

说明本机没有安装 Wails CLI，或者 `PATH` 没包含 Go 的 `bin` 目录。

先按前文“安装项目依赖”中的命令安装 Wails CLI；如果只想单独补装，也可以执行：

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@v2.7.1
```

然后直接用绝对路径运行：

```powershell
& "$env:USERPROFILE\\go\\bin\\wails.exe" build
```

### 2. 点击“连接”后提示 `gateway is required`

这说明当前 GUI 配置里没有可用的网关地址。处理方式：

1. 打开 `Connection` 页面
2. 填写“网关地址”
3. 点击“保存”
4. 再次点击“连接”

如果你依赖 mDNS 自动发现，也建议先确认：

- “发现方式”是否为“自动（mDNS）”
- 网关和当前节点是否在同一局域网

### 3. `Test` 连接测试失败

常见原因包括：

- “网关地址”填写错误
- “端口”配置不正确
- `TLS` 开关与网关实际协议不匹配
- 网关未启动
- 本机到网关的网络不可达

建议按下面顺序检查：

1. 先确认网关地址是否正确
2. 确认端口是否开放
3. 如果网关是 HTTPS / WSS，打开 `TLS`
4. 如果是本地或内网纯 WS，关闭 `TLS`
5. 再次执行“测试连接”

补充说明：

- `http://` 会被内部转换成 `ws://`
- `https://` 会被内部转换成 `wss://`
- 当 `gateway` 为空时，GUI 测试默认尝试 `ws://localhost:18789`

### 4. GUI 保存了配置，但 CLI 启动结果不一致

这是因为：

- GUI 保存的是默认配置，写入 `config.yaml`
- CLI 参数只覆盖当前进程，不会自动写回 `config.yaml`

例如你执行：

```powershell
.\openclaw-node.exe -gateway 192.168.1.50:18789 -tls
```

这只会影响这一次启动；下次打开 GUI 时，仍然会读取 `config.yaml` 中保存的值。

### 5. GUI 中切换能力后，下次启动状态变化了

这是预期行为。`Capabilities` 页面中的开关会直接写入 `config.yaml`，因此它既影响当前运行时，也影响下次启动时的默认能力状态。

如果你想恢复默认值，可以：

- 在 GUI 中手动重新切换
- 或直接编辑 `%APPDATA%\\OpenClaw\\config.yaml`

### 6. CLI 能运行，但桌面打包失败

通常应先分开确认两部分：

1. Go 是否能编译

```powershell
go build -o openclaw-node.exe ./cmd
```

2. 前端是否能构建

```powershell
cd frontend
npm run build
```

如果这两步都成功，再执行：

```powershell
wails build
```

这样可以快速判断问题出在 Go、前端，还是 Wails CLI 本身。

## 文档

设计和实现相关文档集中在：

- `docs/superpowers/specs/`
- `docs/superpowers/plans/`

如果你是第一次接手该项目，建议按以下顺序阅读：

1. `README.md`
2. `docs/superpowers/specs/2026-03-21-openclaw-node-architecture.md`
3. `docs/superpowers/specs/2026-03-21-openclaw-windows-node-design.md`
4. `docs/superpowers/plans/` 下最近的实现计划
