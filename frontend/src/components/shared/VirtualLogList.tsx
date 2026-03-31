import { FixedSizeList as List, type ListChildComponentProps } from "react-window"

import { getLevelLabel } from "@/lib/copy"
import type { LogEntry } from "@/lib/wails"

function Row({ index, style, data }: ListChildComponentProps<LogEntry[]>) {
  const log = data[index]
  return (
    <div style={style} className="flex gap-3 px-4 font-mono text-xs leading-6 text-slate-100">
      <span className="text-slate-500">{new Date(log.timestamp).toLocaleTimeString("zh-CN", { hour12: false })}</span>
      <span className="w-12 text-sky-300">{getLevelLabel(log.level)}</span>
      <span className="flex-1 break-all">{log.message}</span>
    </div>
  )
}

export function VirtualLogList({ logs }: { logs: LogEntry[] }) {
  return (
    <div className="overflow-hidden rounded-3xl bg-slate-950">
      <List height={420} itemCount={logs.length} itemData={logs} itemSize={28} width="100%">
        {Row}
      </List>
    </div>
  )
}
