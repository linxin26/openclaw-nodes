import {
  Connect,
  Disconnect,
  GetAbout,
  GetCapabilities,
  GetConfig,
  GetDeviceInfo,
  GetLogs,
  GetRecentActivity,
  GetStatus,
  InvokeCommand,
  OpenPath,
  SaveConfig,
  SaveFileToDisk,
  SetCapability,
  TestConnection,
  type AboutInfo,
  type ActivityEntry,
  type CapabilityInfo,
  type Config,
  type ConnectionStatus,
  type DeviceInfo,
  type InvokeResult,
  type LogEntry,
  type LogFilter,
  type TestResult,
} from "../../wailsjs/go/app/App"

export type {
  AboutInfo,
  ActivityEntry,
  CapabilityInfo,
  Config,
  ConnectionStatus,
  DeviceInfo,
  InvokeResult,
  LogEntry,
  LogFilter,
  TestResult,
}

export const api = {
  status: {
    get: GetStatus,
  },
  device: {
    getInfo: GetDeviceInfo,
  },
  config: {
    get: GetConfig,
    save: SaveConfig,
  },
  connection: {
    connect: Connect,
    disconnect: Disconnect,
    test: TestConnection,
  },
  capabilities: {
    get: GetCapabilities,
    set: SetCapability,
  },
  operations: {
    invoke: InvokeCommand,
    saveFileToDisk: SaveFileToDisk,
  },
  logs: {
    get: GetLogs,
    getRecentActivity: GetRecentActivity,
  },
  about: {
    get: GetAbout,
    openPath: OpenPath,
  },
} as const
