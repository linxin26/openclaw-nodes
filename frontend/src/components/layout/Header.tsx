import { useLocation } from "react-router-dom"

import { StatusBadge } from "@/components/shared/StatusBadge"
import { copy } from "@/lib/copy"
import type { ConnectionStatus } from "@/lib/wails"

const titles: Record<string, { title: string; description: string }> = {
  "/": copy.pages.dashboard,
  "/connection": copy.pages.connection,
  "/capabilities": {
    title: copy.navigation.capabilities,
    description: "切换此节点对外暴露的命令。",
  },
  "/operations": {
    title: copy.navigation.operations,
    description: "触发设备命令并查看结果。",
  },
  "/logs": {
    title: copy.navigation.logs,
    description: "查看实时事件流和本地操作日志。",
  },
  "/about": {
    title: copy.navigation.about,
    description: "身份、存储与运行时元数据。",
  },
}

export function Header({ status }: { status: ConnectionStatus }) {
  const location = useLocation()
  const meta = titles[location.pathname] ?? titles["/"]

  return (
    <header className="flex flex-col gap-4 rounded-[30px] border border-white/70 bg-white/75 p-6 shadow-panel backdrop-blur md:flex-row md:items-center md:justify-between">
      <div>
        <p className="text-sm uppercase tracking-[0.24em] text-primary/75">{copy.layout.header.controlPanel}</p>
        <h2 className="mt-2 text-3xl font-semibold tracking-tight">{meta.title}</h2>
        <p className="mt-2 text-sm text-muted-foreground">{meta.description}</p>
      </div>
      <div className="space-y-2 rounded-2xl bg-slate-50/90 px-4 py-3">
        <StatusBadge status={status} showGateway />
        <p className="text-xs text-muted-foreground">
          {copy.layout.header.protocol} v{status.protocolVersion} {status.tls ? copy.layout.header.tlsWith : copy.layout.header.tlsWithout}
        </p>
      </div>
    </header>
  )
}
