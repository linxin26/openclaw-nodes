import type { CapabilityInfo, ConnectionStatus } from "@/lib/wails"

// 当前仅提供中文展示文案，内部标识保持英文稳定。
// 这样后续接入 i18n 时，可以直接沿用这套语义分组与键结构，无需重排页面代码。
export const copy: any = {
  layout: {
    sidebar: {
      brand: "OpenClaw",
      title: "节点控制台",
      description: "Wails 桌面面板，用于配置、遥测和手动设备操作。",
    },
    header: {
      controlPanel: "控制面板",
      protocol: "协议",
      tlsWith: "通过 TLS",
      tlsWithout: "不使用 TLS",
    },
  },
  navigation: {
    dashboard: "仪表盘",
    connection: "连接设置",
    capabilities: "能力管理",
    operations: "设备操作",
    logs: "日志",
    about: "关于",
  },
  common: {
    loading: "加载中",
    save: "保存",
    test: "测试",
    connect: "连接",
    disconnect: "断开连接",
    copy: "复制",
    retry: "重试",
    notConfigured: "未配置",
    unavailable: "不可用",
    idle: "空闲",
    enabled: "已启用",
    disabled: "已禁用",
    noActivity: "暂无活动记录。",
    status: {
      connected: "已连接",
      connecting: "连接中",
      offline: "未连接",
      error: "异常",
    },
  },
  pages: {
    dashboard: {
      title: "仪表盘",
      description: "连接概览、能力状态和最近活动。",
      metrics: {
        gateway: "网关",
        uptime: "运行时长",
        enabledCapabilities: "已启用能力",
      },
      nodeState: {
        title: "节点状态",
        description: "当前连接概览和本地身份快照。",
      },
      rows: {
        protocol: "协议版本",
        retryCount: "重试次数",
        tls: "TLS",
        deviceId: "设备 ID",
        hostname: "主机名",
        platform: "平台",
        version: "版本",
      },
      activity: {
        title: "最近活动",
        description: "GUI 应用发出的最近生命周期和操作事件。",
        empty: "暂无活动记录。",
      },
    },
    connection: {
      title: "连接设置",
      description: "网关设置与会话控制。",
      settings: {
        title: "网关设置",
        description: "更新此节点如何连接网关并上报自身信息。",
        labels: {
          gateway: "网关地址",
          port: "端口",
          discovery: "发现方式",
          token: "网关令牌",
          tls: "TLS",
        },
        placeholders: {
          gateway: "gateway.example:443",
          token: "粘贴网关认证令牌",
          selectMode: "选择模式",
        },
        discovery: {
          auto: "自动（mDNS）",
          manual: "手动",
        },
        tlsDescription: "使用安全的 WebSocket 传输。",
      },
      actions: {
        save: "保存配置",
        test: "测试连接",
        connect: "连接",
        disconnect: "断开连接",
      },
      session: {
        title: "会话状态",
        description: "基于当前 Wails 应用绑定的实时状态。",
        labels: {
          status: "状态",
          gateway: "网关",
          tls: "TLS",
          testLatency: "测试延迟",
          retryCount: "重试次数",
          retryDelay: "重试延迟",
        },
      },
      toasts: {
        saved: "配置已保存",
        testSuccess: (latencyMs: number) => `连接测试成功，耗时 ${latencyMs}ms`,
        testFailed: "连接测试失败",
      },
    },
    capabilities: {
      loading: "正在加载能力",
      enabled: "已启用",
      disabled: "已禁用",
      healthy: "健康",
      unhealthy: "异常",
      commands: "命令",
    },
    operations: {
      tabs: {
        camera: "摄像头",
        screen: "屏幕",
        location: "位置",
        photos: "照片",
        notifications: "通知",
        motion: "运动",
        sms: "短信",
        calendar: "日历",
      },
      camera: {
        title: "摄像头",
        description: "列出摄像头、拍摄静态帧或录制短视频。",
        actions: {
          list: "列出摄像头",
          snapshot: "拍摄快照",
          clip: "录制 5 秒视频",
        },
      },
      screen: {
        title: "屏幕",
        description: "捕获当前桌面。",
        actions: {
          capture: "捕获屏幕",
        },
      },
      location: {
        title: "位置",
        description: "解析大致设备坐标。",
        actions: {
          get: "获取位置",
        },
      },
      photos: {
        title: "照片",
        description: "获取最新照片元数据。",
        actions: {
          latest: "最新照片",
        },
      },
      notifications: {
        title: "通知",
        description: "读取或触发桌面通知。",
        labels: {
          body: "通知内容",
        },
        actions: {
          list: "列出通知",
          trigger: "触发通知",
        },
      },
      motion: {
        title: "运动",
        description: "读取运动状态和计步器估计值。",
        actions: {
          activity: "运动活动",
          pedometer: "计步器",
        },
      },
      sms: {
        title: "短信",
        description: "通过节点接口发送或搜索短信数据。",
        labels: {
          to: "收件人",
          body: "短信内容",
        },
        actions: {
          send: "发送短信",
          search: "搜索短信",
        },
      },
      calendar: {
        title: "日历",
        description: "查看事件或创建新的本地事件文件条目。",
        labels: {
          title: "标题",
          start: "开始时间",
        },
        actions: {
          list: "列出事件",
          add: "添加事件",
        },
      },
      result: {
        title: "操作结果",
        completedIn: (durationMs: number) => `已完成，耗时 ${durationMs}ms`,
        copied: "结果已复制",
        copy: "复制",
        savePayload: "保存载荷",
        retry: "重试",
      },
    },
    logs: {
      filters: {
        info: "信息",
        warn: "警告",
        error: "错误",
      },
      searchPlaceholder: "搜索日志正文",
      copied: "日志已复制",
      empty: "当前筛选条件下没有匹配的日志。",
      copyNote: "日志正文、级别和值保持原始技术内容，不强制翻译，便于和后端输出、事件名与排障信息逐项对照。",
    },
    about: {
      runtimeIdentity: "运行标识",
      buildMetadata: "构建信息",
      labels: {
        deviceId: "设备 ID",
        publicKey: "公钥",
        dataDir: "数据目录",
        version: "版本",
        platform: "平台",
        hostname: "主机名",
        goVersion: "Go 版本",
        arch: "架构",
        protocolVersion: "协议版本",
      },
      buttons: {
        openDataDirectory: "打开数据目录",
      },
      toasts: {
        copied: (label: string) => `${label} 已复制`,
      },
    },
  },
}

copy.capabilities = {
  enabled: copy.pages.capabilities.enabled,
  disabled: copy.pages.capabilities.disabled,
  healthy: copy.pages.capabilities.healthy,
  unhealthy: copy.pages.capabilities.unhealthy,
  commands: copy.pages.capabilities.commands,
  unavailableDescription: "当前未提供额外说明。",
}

copy.operations = {
  result: {
    title: copy.pages.operations.result.title,
    success: copy.pages.operations.result.completedIn,
    savePayload: copy.pages.operations.result.savePayload,
  },
}

copy.logs = {
  searchPlaceholder: copy.pages.logs.searchPlaceholder,
  copy: "复制日志",
  clear: "清空日志",
  technicalHint: "日志正文保留原始技术内容，避免翻译后影响故障定位和协议排查。",
  empty: copy.pages.logs.empty,
  toast: {
    copied: copy.pages.logs.copied,
  },
}

copy.about = {
  identity: copy.pages.about.runtimeIdentity,
  metadata: copy.pages.about.buildMetadata,
  deviceId: copy.pages.about.labels.deviceId,
  publicKey: copy.pages.about.labels.publicKey,
  dataDir: copy.pages.about.labels.dataDir,
  openDataDir: copy.pages.about.buttons.openDataDirectory,
  version: copy.pages.about.labels.version,
  platform: copy.pages.about.labels.platform,
  hostname: copy.pages.about.labels.hostname,
  goVersion: copy.pages.about.labels.goVersion,
  architecture: copy.pages.about.labels.arch,
  protocolVersion: copy.pages.about.labels.protocolVersion,
  toast: {
    copied: copy.pages.about.toasts.copied,
  },
}

copy.actions = {
  resultCopied: copy.pages.operations.result.copied,
  savedFile: (filename: string) => `已保存 ${filename}`,
  invokeCompleted: (method: string) => `${method} 执行完成`,
  invokeFailed: (method: string) => `${method} 执行失败`,
}

copy.capabilitiesMap = {
  camera: { name: "摄像头", description: "访问摄像头设备，执行拍照、枚举或录制相关命令。" },
  location: { name: "位置", description: "获取设备当前位置或近似地理位置信息。" },
  photos: { name: "照片", description: "浏览本地照片资源并读取最近照片元数据。" },
  screen: { name: "屏幕", description: "采集当前桌面画面或屏幕截图。" },
  motion: { name: "运动", description: "读取运动状态、活动类型和计步相关数据。" },
  notifications: { name: "通知", description: "读取桌面通知或触发本地通知。" },
  sms: { name: "短信", description: "发送短信并检索本地短信记录。" },
  calendar: { name: "日历", description: "查看或创建本地日历事件。" },
}

export function getStatusLabel(status: ConnectionStatus["status"]) {
  return copy.common.status[status]
}

export function getLevelLabel(level: string) {
  return copy.pages.logs.filters[level as keyof typeof copy.pages.logs.filters] ?? level.toUpperCase()
}

export function getCapabilityDisplay(capability: Pick<CapabilityInfo, "key" | "name" | "description">) {
  const mapped = copy.capabilitiesMap[capability.key]
  return {
    name: mapped?.name ?? capability.name,
    description: mapped?.description ?? (capability.description || copy.capabilities.unavailableDescription),
  }
}
