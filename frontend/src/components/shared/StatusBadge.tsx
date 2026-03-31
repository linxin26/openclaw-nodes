import { Badge } from "@/components/ui/badge"
import { copy } from "@/lib/copy"
import type { ConnectionStatus } from "@/lib/wails"

const statusMap: Record<ConnectionStatus["status"], { className: string }> = {
  connected: { className: "bg-emerald-100 text-emerald-700" },
  connecting: { className: "bg-amber-100 text-amber-700" },
  offline: { className: "bg-slate-200 text-slate-700" },
  error: { className: "bg-rose-100 text-rose-700" },
}

export function StatusBadge({
  status,
  showGateway = false,
}: {
  status: ConnectionStatus
  showGateway?: boolean
}) {
  const meta = statusMap[status.status]
  return (
    <div className="flex flex-wrap items-center gap-2">
      <Badge className={meta.className}>{copy.common.status[status.status]}</Badge>
      {showGateway && status.gateway ? (
        <span className="text-xs text-muted-foreground">{status.gateway}</span>
      ) : null}
    </div>
  )
}
