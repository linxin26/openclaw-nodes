import { useEffect, useState } from "react"
import { Activity, Gauge, Server, Shield } from "lucide-react"

import { LoadingSpinner } from "@/components/shared/LoadingSpinner"
import { StatusBadge } from "@/components/shared/StatusBadge"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { useConnection } from "@/hooks/useConnection"
import { copy, getLevelLabel } from "@/lib/copy"
import { api, type ActivityEntry, type CapabilityInfo, type DeviceInfo } from "@/lib/wails"
import { formatTimestamp, formatUptime } from "@/lib/utils"

export function Dashboard() {
  const { status, loading } = useConnection()
  const [device, setDevice] = useState<DeviceInfo | null>(null)
  const [capabilities, setCapabilities] = useState<CapabilityInfo[]>([])
  const [activity, setActivity] = useState<ActivityEntry[]>([])

  useEffect(() => {
    api.device.getInfo().then(setDevice).catch(() => undefined)
    api.capabilities.get().then(setCapabilities).catch(() => undefined)
    api.logs.getRecentActivity().then(setActivity).catch(() => undefined)
  }, [])

  if (loading) return <LoadingSpinner />

  const enabledCount = capabilities.filter((item) => item.enabled).length

  return (
    <div className="grid gap-6 xl:grid-cols-[1.15fr_0.85fr]">
      <div className="space-y-6">
        <Card className="overflow-hidden">
          <CardContent className="grid gap-6 p-6 md:grid-cols-3">
            <MetricCard
              icon={Server}
              label={copy.pages.dashboard.metrics.gateway}
              value={status.gateway || copy.common.notConfigured}
            />
            <MetricCard
              icon={Gauge}
              label={copy.pages.dashboard.metrics.uptime}
              value={status.uptimeMs > 0 ? formatUptime(status.uptimeMs) : copy.common.idle}
            />
            <MetricCard
              icon={Shield}
              label={copy.pages.dashboard.metrics.enabledCapabilities}
              value={`${enabledCount}/8`}
            />
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>{copy.pages.dashboard.nodeState.title}</CardTitle>
            <CardDescription>{copy.pages.dashboard.nodeState.description}</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4 md:grid-cols-2">
            <div className="rounded-2xl bg-slate-50 p-4">
              <StatusBadge status={status} showGateway />
              <dl className="mt-4 space-y-2 text-sm">
                <InfoRow label={copy.pages.dashboard.rows.protocol} value={`v${status.protocolVersion}`} />
                <InfoRow label={copy.pages.dashboard.rows.retryCount} value={`${status.retryCount}`} />
                <InfoRow label={copy.pages.dashboard.rows.tls} value={status.tls ? copy.common.enabled : copy.common.disabled} />
              </dl>
            </div>
            <div className="rounded-2xl bg-slate-50 p-4">
              <dl className="space-y-2 text-sm">
                <InfoRow label={copy.pages.dashboard.rows.deviceId} value={device?.deviceId ?? copy.common.unavailable} />
                <InfoRow label={copy.pages.dashboard.rows.hostname} value={device?.hostname ?? copy.common.unavailable} />
                <InfoRow label={copy.pages.dashboard.rows.platform} value={device?.platform ?? copy.common.unavailable} />
                <InfoRow label={copy.pages.dashboard.rows.version} value={device?.version ?? copy.common.unavailable} />
              </dl>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Activity className="h-5 w-5 text-primary" />
            {copy.pages.dashboard.activity.title}
          </CardTitle>
          <CardDescription>{copy.pages.dashboard.activity.description}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {activity.length === 0 ? (
              <p className="text-sm text-muted-foreground">{copy.pages.dashboard.activity.empty}</p>
            ) : (
              activity.map((entry) => (
                <div key={`${entry.timestamp}-${entry.event}`} className="rounded-2xl bg-slate-50 p-4">
                  <div className="flex items-center justify-between gap-4">
                    {/* 事件名保留原文，便于和后端事件日志、排障输出逐项对照。 */}
                    <p className="font-medium">{entry.event}</p>
                    <span className="text-xs text-muted-foreground">{formatTimestamp(entry.timestamp)}</span>
                  </div>
                  <p className="mt-1 text-sm text-muted-foreground">{getLevelLabel(entry.level)}</p>
                </div>
              ))
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

function MetricCard({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ComponentType<{ className?: string }>
  label: string
  value: string
}) {
  return (
    <div className="rounded-3xl bg-slate-950 p-5 text-white">
      <Icon className="h-5 w-5 text-sky-300" />
      <p className="mt-6 text-xs uppercase tracking-[0.24em] text-slate-400">{label}</p>
      <p className="mt-2 text-xl font-semibold">{value}</p>
    </div>
  )
}

function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between gap-4">
      <dt className="text-muted-foreground">{label}</dt>
      <dd className="max-w-[60%] truncate font-medium">{value}</dd>
    </div>
  )
}

