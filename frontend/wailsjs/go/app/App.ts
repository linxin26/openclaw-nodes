export type ConnectionStatus = {
  status: "offline" | "connecting" | "connected" | "error"
  gateway: string
  tls: boolean
  uptimeMs: number
  retryCount: number
  retryDelayMs: number
  protocolVersion: number
}

export type DeviceInfo = {
  deviceId: string
  platform: string
  hostname: string
  mode: string
  version: string
}

export type Config = {
  gateway: string
  port: number
  token: string
  tls: boolean
  discovery: string
  capabilities: Record<string, boolean>
}

export type TestResult = {
  success: boolean
  latencyMs: number
  error: string
}

export type CapabilityInfo = {
  key: string
  name: string
  description: string
  enabled: boolean
  commands: string[]
  dependencies: string[]
  healthy: boolean
}

export type InvokeResult = {
  success: boolean
  data: Record<string, unknown>
  error: string
  durationMs: number
}

export type AboutInfo = {
  deviceId: string
  publicKey: string
  version: string
  platform: string
  hostname: string
  goVersion: string
  arch: string
  dataDir: string
  protocolVersion: number
}

export type LogEntry = {
  timestamp: number
  level: string
  message: string
}

export type ActivityEntry = {
  timestamp: number
  event: string
  level: string
}

export type LogFilter = {
  levels: string[]
  search: string
  limit: number
  offset: number
}

type AppBindings = {
  GetStatus: () => Promise<ConnectionStatus>
  GetDeviceInfo: () => Promise<DeviceInfo>
  GetConfig: () => Promise<Config>
  SaveConfig: (config: Config) => Promise<void>
  Connect: () => Promise<void>
  Disconnect: () => Promise<void>
  TestConnection: () => Promise<TestResult>
  GetCapabilities: () => Promise<CapabilityInfo[]>
  SetCapability: (key: string, enabled: boolean) => Promise<void>
  InvokeCommand: (
    method: string,
    params?: Record<string, unknown>,
  ) => Promise<InvokeResult>
  GetLogs: (filter?: LogFilter) => Promise<LogEntry[]>
  GetRecentActivity: () => Promise<ActivityEntry[]>
  GetAbout: () => Promise<AboutInfo>
  OpenPath: (path: string) => Promise<void>
  SaveFileToDisk: (base64Data: string, filename: string) => Promise<void>
}

function getBindings(): Partial<AppBindings> | undefined {
  const win = window as Window & {
    go?: { wails?: { App?: Partial<AppBindings> } }
  }
  return win.go?.wails?.App
}

function invokeBinding<T>(name: keyof AppBindings, ...args: unknown[]): Promise<T> {
  const fn = getBindings()?.[name]
  if (typeof fn !== "function") {
    return Promise.reject(new Error(`Wails binding not generated yet: ${String(name)}`))
  }
  return (fn as (...items: unknown[]) => Promise<T>)(...args)
}

export function GetStatus() {
  return invokeBinding<ConnectionStatus>("GetStatus")
}
export function GetDeviceInfo() {
  return invokeBinding<DeviceInfo>("GetDeviceInfo")
}
export function GetConfig() {
  return invokeBinding<Config>("GetConfig")
}
export function SaveConfig(config: Config) {
  return invokeBinding<void>("SaveConfig", config)
}
export function Connect() {
  return invokeBinding<void>("Connect")
}
export function Disconnect() {
  return invokeBinding<void>("Disconnect")
}
export function TestConnection() {
  return invokeBinding<TestResult>("TestConnection")
}
export function GetCapabilities() {
  return invokeBinding<CapabilityInfo[]>("GetCapabilities")
}
export function SetCapability(key: string, enabled: boolean) {
  return invokeBinding<void>("SetCapability", key, enabled)
}
export function InvokeCommand(
  method: string,
  params?: Record<string, unknown>,
) {
  return invokeBinding<InvokeResult>("InvokeCommand", method, params)
}
export function GetLogs(filter?: LogFilter) {
  return invokeBinding<LogEntry[]>("GetLogs", filter)
}
export function GetRecentActivity() {
  return invokeBinding<ActivityEntry[]>("GetRecentActivity")
}
export function GetAbout() {
  return invokeBinding<AboutInfo>("GetAbout")
}
export function OpenPath(path: string) {
  return invokeBinding<void>("OpenPath", path)
}
export function SaveFileToDisk(base64Data: string, filename: string) {
  return invokeBinding<void>("SaveFileToDisk", base64Data, filename)
}
