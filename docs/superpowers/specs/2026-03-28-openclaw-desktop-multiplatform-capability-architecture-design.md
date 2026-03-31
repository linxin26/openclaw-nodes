# OpenClaw Desktop Node 多平台能力升级架构设计

**版本**: v0.1
**日期**: 2026-03-28
**状态**: Draft

---

## 1. 背景与目标

当前 `openclaw-node` 已具备桌面 Node 的基本协议接入与部分设备能力实现，但能力层仍明显偏向 Windows 单平台：

- `camera` 基于 `ffmpeg + dshow`
- `screen` 通过 PowerShell 脚本截图
- `photos` 直接依赖 `USERPROFILE/Pictures`
- `device.describe` / `device.status` / `device.permissions` 仍存在硬编码平台与能力状态

这套实现适合单平台原型验证，但不适合扩展到 `Windows`、`macOS`、`Linux` 三端统一支持。

本次方案的目标不是“给现有 Windows 代码补上 macOS/Linux 分支”，而是把 Node 升级为统一的跨平台桌面能力宿主：

- 采用统一能力内核，复用协议、配置、状态、UI 和运行时框架
- 为 `Windows`、`macOS`、`桌面 Linux` 建立平台适配层
- 优先支持拍照、短视频等媒体能力，同时覆盖当前仓库中的 8 类能力
- 尽量减少外部依赖，但允许少量平台适配依赖存在
- 对无头 Linux 不作为正式支持目标，只在文档中说明限制与后续方向

本次设计覆盖的能力范围为：

- `camera`
- `photos`
- `screen`
- `location`
- `motion`
- `notifications`
- `sms`
- `calendar`

---

## 2. 设计原则

### 2.1 统一能力内核，平台实现下沉

协议层、配置层、能力注册、权限模型、状态输出、Wails/Tray UI 应保持统一；平台差异仅下沉到能力 provider 层。

### 2.2 命令语义统一优先于底层实现统一

不同平台底层 API 可以不同，但 `camera.snap`、`camera.clip`、`screen.snapshot` 等命令的输入、输出、错误模型必须一致。

### 2.3 尽量少依赖，但允许少量适配依赖

不将能力架构建立在大量外部命令拼接之上。优先使用 Go 可集成的原生/轻量方案；必要时允许少量平台适配依赖存在，但必须可探测、可上报、可降级。

### 2.4 明确支持等级，拒绝伪等价实现

对桌面端天然不适合的能力，如 `sms`、`motion`，要在架构上保留接口与状态模型，但不强行伪造三端一致实现。

### 2.5 桌面 Linux 为正式目标，无头 Linux 明确受限

Linux 支持范围以带图形会话的桌面发行版为主。无头环境只记录限制和后续演进方向，不纳入首阶段正式支持承诺。

---

## 3. 总体架构

升级后的 Node 建议分为 5 层：

1. 协议与连接层
2. 能力编排层
3. 能力接口层
4. 平台适配层
5. 宿主与配置层

```text
+------------------------------------------------------------------+
|                         Desktop Node Host                        |
|                                                                  |
|  +--------------------+   +-----------------------------------+  |
|  | Protocol / Runtime |   | Host / Config / Wails / Tray UI   |  |
|  | WS / Auth / Invoke |   | Settings / Status / Diagnostics    |  |
|  +---------+----------+   +----------------+------------------+  |
|            |                               |                     |
|            v                               v                     |
|  +------------------------------------------------------------+  |
|  | Capability Registry & Orchestrator                        |  |
|  | - registration                                            |  |
|  | - permission checks                                       |  |
|  | - availability checks                                     |  |
|  | - timeout / cancellation                                  |  |
|  | - unified result / error mapping                          |  |
|  +-------------------------+----------------------------------+  |
|                            |                                     |
|                            v                                     |
|  +------------------------------------------------------------+  |
|  | Capability Interfaces                                      |  |
|  | camera / photos / screen / location / ...                 |  |
|  +-------------------------+----------------------------------+  |
|                            |                                     |
|                            v                                     |
|  +----------------+----------------+--------------------------+  |
|  | Windows        | macOS          | Linux (desktop)          |  |
|  | Providers      | Providers      | Providers                |  |
|  +----------------+----------------+--------------------------+  |
+------------------------------------------------------------------+
```

各层职责如下：

- 协议与连接层：负责 Gateway 连接、鉴权、命令收发、重连和结果回传
- 能力编排层：负责命令到能力的路由、统一错误模型、权限检查、超时控制
- 能力接口层：定义 8 类能力的公共接口、结果结构和元数据描述
- 平台适配层：分别实现 `windows` / `darwin` / `linux` 的 provider
- 宿主与配置层：统一负责配置、运行状态、诊断、UI 展示与能力开关

---

## 4. 能力分组与架构边界

8 类能力不应以完全相同的方式抽象，建议按差异来源分为 4 组：

### 4.1 媒体采集组

- `camera`
- `photos`
- `screen`

这一组受平台采集 API、权限、媒体格式和编码链路影响最大，是多平台升级的核心。

### 4.2 系统感知组

- `location`
- `motion`

这一组主要受系统服务可用性和设备类型差异影响。

### 4.3 系统集成组

- `notifications`
- `calendar`

这一组主要涉及桌面平台通知中心、数据源适配和权限边界。

### 4.4 受限占位组

- `sms`

桌面端通常缺乏统一、稳定的原生 SMS 能力，应在架构上保留接口和状态位，但不在首阶段承诺真实实现。

---

## 5. 统一能力接口设计

每类能力建议拆成两层接口：

- `CapabilityProvider`
- `CommandHandler`

前者负责能力本身，后者负责具体命令执行。

### 5.1 CapabilityProvider

负责：

- 能力元数据描述
- 平台支持判断
- 当前可用性检查
- 权限状态检查
- 健康检查

建议统一暴露以下能力级接口：

- `Descriptor()`
- `PermissionStatus(ctx)`
- `Availability(ctx)`
- `Commands()`

### 5.2 CommandHandler

负责：

- 参数校验
- 命令执行
- 结果模型填充
- 平台错误归一化

命令处理不应直接关心协议层帧结构，也不应直接依赖某个平台脚本命令细节。

### 5.3 统一结果模型

为保证 Gateway 与上层调用方无需感知平台差异，建议统一以下结果类型：

- 图片结果：`format`、`width`、`height`、`size`、`timestamp`、`payload`
- 视频结果：`format`、`durationMs`、`size`、`timestamp`、`includeAudio`、`payload`
- 列表结果：`items`、`total`、`cursor/after`
- 状态结果：`enabled`、`available`、`reason`

其中 `payload` 可根据后续协议约束选择内嵌 Base64 或文件引用，但对三端必须保持一致策略。

---

## 6. 能力责任矩阵

### 6.1 Tier 1: 首要实现能力

- `camera`
- `photos`
- `screen`

这三类能力必须在 `Windows`、`macOS`、`桌面 Linux` 都有明确实现路径。

### 6.2 Tier 2: 统一接口但允许平台差异

- `location`
- `notifications`
- `calendar`

这三类能力要求有统一接口和状态模型，但允许不同平台的能力级别不同。

### 6.3 Tier 3: 协议保留与状态占位

- `motion`
- `sms`

这两类能力应保留能力接口和状态声明，但不强行进入首阶段真实实现范围。

---

## 7. 代码组织与模块重构建议

当前 `internal/device` 下的大部分文件同时承担了命令注册、平台调用和结果拼装，后续扩展会迅速耦合失控。建议调整为“能力内核 + 能力编排 + 平台 provider + 运行时装配”的结构。

建议目录方向如下：

```text
internal/device/
  core/
    capability.go
    registry.go
    permissions.go
    errors.go
    result.go
  runtime/
    runtime.go
    platform_windows.go
    platform_darwin.go
    platform_linux.go
  capabilities/
    camera/
    photos/
    screen/
    location/
    motion/
    notifications/
    sms/
    calendar/
  platform/
    windows/
    darwin/
    linux/
```

职责建议如下：

- `core/`: 能力接口、统一错误、权限状态、结果模型、注册中心
- `runtime/`: 当前平台识别、provider 装配、命令注册、运行时状态聚合
- `capabilities/<name>/`: 各能力的公共模型、参数定义、命令编排逻辑
- `platform/<os>/`: 各平台的具体 provider 实现

---

## 8. 注册机制与运行时装配

当前基于 `init()` 的 `protocol.RegisterHandler(...)` 方式不适合继续扩展，因为：

- 命令注册时机隐式
- 平台差异难以控制
- 能力开关与可用性状态难以动态表达
- `device.describe`、`device.status`、`device.permissions` 只能依赖硬编码

建议改为显式注册流程：

1. 启动时识别当前平台
2. 创建对应平台 provider 集合
3. 将 provider 注入对应能力编排器
4. 由能力编排器向协议层注册命令
5. 由能力注册中心动态生成 `device.describe` / `device.status` / `device.permissions`

这样可以实现：

- 命令注册与平台适配解耦
- 能力启停与状态输出统一
- `device.info` 不再写死 `windows`
- UI、配置和协议层都从同一份能力元数据读取信息

---

## 9. 权限模型与可用性状态

多平台能力设计中，权限和支持状态必须先于具体实现统一。

建议所有能力统一采用以下状态：

- `granted`
- `denied`
- `restricted`
- `not_supported`
- `not_applicable`

状态含义如下：

- `granted`: 已授权且允许使用
- `denied`: 平台支持，但用户或系统明确拒绝
- `restricted`: 平台理论支持，但当前运行环境不满足
- `not_supported`: 当前平台不支持该能力
- `not_applicable`: 当前设备类型上该能力无意义

每个能力都应按以下顺序判断：

1. 平台是否支持
2. 当前运行环境是否满足
3. 权限是否已授予

只有三者都满足，命令才进入实际执行阶段。

---

## 10. 降级策略与受限能力处理

多平台设计不要求所有能力绝对等价，但要求每项能力的“为何可用/为何不可用”可解释、可上报。

### 10.1 camera

- 优先真实摄像头枚举与采集
- 无设备时返回不可用状态
- 不允许伪造拍照/录像结果

### 10.2 photos

- 首阶段优先采用统一受控目录模型
- 后续再逐步引入系统媒体库接入

### 10.3 screen

- 桌面会话中必须可实现截图
- 无头 Linux 明确标记为 `restricted`

### 10.4 location

- 优先系统位置服务
- 缺失时允许降级为网络定位
- 返回结果需标记来源与精度级别

### 10.5 motion

- 桌面端默认允许长期 `not_supported`
- 不建议构造伪传感器实现

### 10.6 notifications

- 至少支持本地通知能力
- “读取系统通知”应视作增强能力，不纳入首阶段统一承诺

### 10.7 calendar

- 优先支持本地日历数据源或标准化导入目录
- 允许不同平台具备不同访问级别

### 10.8 sms

- 默认 `not_supported`
- 保留接口和扩展点，不承诺首阶段真实能力

---

## 11. 平台适配策略

### 11.1 Windows

Windows 是当前基线平台，应从“脚本与命令驱动实现”升级为“原生能力优先，少量依赖兜底”。

建议方向：

- `camera`: 优先原生媒体采集路径，必要时保留轻量视频编码兜底
- `photos`: 默认用户图片目录或应用受控目录
- `screen`: 优先系统截图 API，不再依赖 PowerShell 作为正式主路径
- `location`: 优先系统位置服务
- `notifications`: 对接系统通知能力
- `calendar`: 对接本地日历数据源或受控导入目录
- `motion` / `sms`: 默认受限

### 11.2 macOS

macOS 的关键点在于权限建模和媒体/屏幕采集限制。

建议方向：

- `camera`: 支持拍照与短视频，摄像头权限作为一等公民处理
- `photos`: 区分受控目录与系统照片库两种模式
- `screen`: 截图能力可支持，屏幕捕获权限需单独建模
- `location`: 支持系统位置服务权限状态回传
- `notifications`: 支持系统通知集成
- `calendar`: 支持本地日历访问与权限判断
- `motion` / `sms`: 默认受限

### 11.3 Linux（桌面）

Linux 以桌面图形环境为正式支持目标，不承诺无头环境。

建议方向：

- `camera`: 支持桌面 Linux 上的枚举、拍照、短视频
- `photos`: 优先受控目录，再逐步适配不同桌面媒体目录
- `screen`: 仅在图形会话中支持，适配 Wayland/X11 差异
- `location`: 优先网络定位或桌面环境可用位置服务
- `notifications`: 支持桌面通知能力
- `calendar`: 优先标准化数据源或导入目录
- `motion` / `sms`: 默认受限

### 11.4 Linux（无头）

无头 Linux 不属于首阶段正式目标。

应在架构中明确：

- `camera`、`screen`、`notifications` 常处于 `restricted`
- 允许保留能力注册，但状态应明确说明运行环境限制
- 后续如要支持，应单独形成子方案，而非挤入当前桌面方案

---

## 12. 外部依赖分级

为落实“尽量不要依赖”的目标，建议把能力实现依赖分为 3 级：

- `Level A`: Go 内部库或随应用集成的轻量绑定，优先使用
- `Level B`: 平台原生命令/API 桥接，可接受，但必须统一封装
- `Level C`: 大型外部工具链，仅作兜底，不作为正式架构主路径

目标约束：

- `Tier 1` 能力尽量落在 `Level A/B`
- `Level C` 仅保留给少数平台差异最难抹平的场景
- 所有外部依赖都必须可检测、可诊断、可降级

---

## 13. 对现有代码的直接影响

以下文件在升级中应被视为重构入口，而不是继续叠加平台分支：

- `internal/device/camera.go`
- `internal/device/screen.go`
- `internal/device/photos.go`
- `internal/protocol/commands_device.go`
- `internal/config/config.go`

重构方向：

- `camera.go` 从 Windows 相机实现升级为 `camera` 能力编排入口
- `screen.go` 从 PowerShell 截图入口升级为跨平台 `screen` 编排层
- `photos.go` 从硬编码 Windows 图片目录升级为统一媒体目录抽象
- `commands_device.go` 从静态能力描述升级为动态能力注册中心输出
- `config.go` 从简单布尔开关扩展为支持平台能力参数、依赖探测和策略配置

---

## 14. 风险与约束

### 14.1 平台权限差异

同一能力在三端上的授权方式和运行限制不同，若权限模型不先统一，后续实现会出现行为漂移。

### 14.2 媒体链路复杂度

`camera.snap` 与 `camera.clip` 对设备枚举、采集、编码、权限、文件大小控制都更敏感，是首要风险项。

### 14.3 Linux 图形协议差异

桌面 Linux 在 `X11` 与 `Wayland` 上行为差异明显，必须通过 provider 层统一，不应泄漏到上层命令语义。

### 14.4 协议输出一致性

如果 `device.describe`、`device.status`、`device.permissions` 继续使用硬编码，能力层再怎么抽象也无法对外表现一致。

### 14.5 不现实能力承诺

`motion`、`sms` 在桌面端天然受限，若在方案阶段承诺首批真实落地，会拖累整体多平台升级进度。

---

## 15. GUI 受影响点说明

本次升级的中心是能力层和运行时装配，不是重做 GUI。但由于 GUI 当前展示的数据来自静态配置和静态能力假设，能力模型升级后，GUI 必然需要跟随调整数据来源和展示结构。

### 15.1 受影响范围

当前受影响的 GUI/宿主代码主要包括：

- `internal/wails/app.go`
- `internal/wails/types.go`
- `internal/wails/app_test.go`
- `internal/tray/menu.go`
- `internal/tray/dialog.go`

这些模块后续不应继续维护一套独立的能力常量、平台判断或固定状态，而应统一读取运行时能力注册中心输出的元数据。

### 15.2 数据来源变化

升级前，GUI 更接近“读取配置 + 展示预设能力”的方式。

升级后，GUI 应改为读取运行时聚合结果，至少包括：

- 能力是否启用
- 能力当前是否可用
- 能力权限状态
- 不可用或受限原因
- 当前平台暴露的命令与支持等级

也就是说，GUI 的职责从“描述系统应该具备什么能力”，转为“展示当前平台和当前运行时真实具备什么能力”。

### 15.3 对现有界面的具体影响

#### 状态页

状态页不应再只展示 `connected`、`gatewayConnected` 以及简单 capability 布尔值，而应展示每项能力的真实状态，例如：

- `enabled=true, available=false, permission=denied`
- `enabled=true, available=false, permission=restricted`
- `enabled=false, available=false, reason=disabled_by_config`

这样用户才能分辨“能力关闭”、“平台不支持”、“没有授权”、“当前环境受限”之间的差别。

#### 配置页

配置页仍然保留能力开关，但其作用应限定为“是否启用能力”，而不是暗示“启用后一定可用”。

后续允许少量新增配置项，例如：

- 受控目录路径
- provider 策略选择
- 平台兜底方案开关

但不应在首阶段引入大量平台专属配置，避免 GUI 复杂度失控。

#### 托盘菜单

托盘菜单不应继续写死能力列表和状态说明。后续建议从运行时 capability metadata 动态生成菜单项或状态摘要，否则一旦 `Windows / macOS / Linux` 的能力级别不同，托盘层会迅速与真实运行状态脱节。

#### 诊断与错误提示

GUI 后续应能展示限制原因，而不仅是“失败”。

例如：

- 摄像头权限未授权
- 当前 Linux 会话为无头环境
- 当前平台不支持短信能力
- 当前 provider 缺失必要依赖

这类信息不应由 GUI 自己拼接，而应来自运行时与 provider 层统一输出的 reason/status 字段。

### 15.4 对当前 GUI 的影响等级

本次升级对当前 GUI 的影响属于“数据模型和展示逻辑升级”，不属于“整套界面重做”。

明确来说：

- 会影响 GUI 的状态展示模型
- 会影响设置页的数据绑定字段
- 会影响托盘菜单的数据来源
- 会影响 Wails 返回给前端的数据结构和测试

但不会要求本次架构升级同步完成：

- 全新视觉设计
- 大规模前端交互重构
- 跨平台 GUI 风格统一改版
- 完整的 GUI 信息架构重做

### 15.5 推荐处理方式

GUI 相关改造应遵循以下顺序：

1. 先完成能力注册中心和运行时状态聚合
2. 再让 Wails App 层暴露新的 capability status 模型
3. 最后调整托盘和配置/状态界面的展示逻辑

不建议反过来先改 GUI，因为在运行时状态模型稳定之前，GUI 改造会反复返工。

### 15.6 结论

当前 GUI 会受到本次升级影响，但影响点主要集中在“读什么数据、怎么解释状态”，而不是“界面本身长什么样”。

因此，GUI 在本方案中的正确定位是：

- 作为运行时能力状态的展示层
- 作为能力启停和少量参数配置的入口
- 不承载平台能力判断逻辑
- 不维护独立于能力注册中心之外的状态体系

---

## 16. 结论

本次升级应采用“统一能力内核 + 平台适配层”的方案，而不是继续沿用单平台能力实现叠加条件分支。

总体结论如下：

- 统一 Node 宿主、协议层、配置层、运行时和 UI
- 将 8 类能力拆分为公共能力接口与平台 provider 实现
- 以 `camera`、`photos`、`screen` 为多平台首要能力
- 对 `location`、`notifications`、`calendar` 保持统一接口但接受平台差异
- 对 `motion`、`sms` 保留协议与状态位，不强行首阶段落地
- 以桌面 `Windows / macOS / Linux` 为正式目标，明确无头 Linux 的限制边界

该设计文档适合作为后续实施计划文档的输入基线。下一步应基于本设计继续拆分实施阶段、文件改造顺序、测试矩阵和验收标准。
# Implementation Addendum

- The repository keeps a minimal tracked `frontend/dist/index.html` placeholder so desktop entrypoints with embedded assets compile under `go test ./...` in clean worktrees.
- Platform skeleton providers for `darwin` and `linux` are emitted as one file per capability to preserve explicit package boundaries and keep future real implementations isolated.
