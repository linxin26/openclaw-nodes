import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatUptime(uptimeMs: number): string {
  const hours = Math.floor(uptimeMs / 3600000)
  const minutes = Math.floor((uptimeMs % 3600000) / 60000)
  return `${hours} 小时 ${minutes} 分钟`
}

export function formatTimestamp(ts: number): string {
  return new Date(ts).toLocaleTimeString("zh-CN", { hour12: false })
}
