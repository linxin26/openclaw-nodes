import { ShieldCheck, ShieldOff } from "lucide-react"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Switch } from "@/components/ui/switch"
import { copy } from "@/lib/copy"
import type { CapabilityInfo } from "@/lib/wails"

export function CapabilityCard({
  capability,
  onToggle,
  disabled,
}: {
  capability: CapabilityInfo
  onToggle: (enabled: boolean) => void
  disabled?: boolean
}) {
  return (
    <Card className="h-full">
      <CardHeader className="flex-row items-start justify-between space-y-0">
        <div>
          <CardTitle>{capability.name}</CardTitle>
          <CardDescription className="mt-1">{capability.description}</CardDescription>
        </div>
        {capability.healthy ? (
          <ShieldCheck className="h-5 w-5 text-emerald-600" />
        ) : (
          <ShieldOff className="h-5 w-5 text-rose-600" />
        )}
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">{copy.pages.capabilities.enabled}</span>
          <Switch checked={capability.enabled} onCheckedChange={onToggle} disabled={disabled} />
        </div>
        <div className="flex items-center justify-between rounded-2xl bg-slate-50 px-3 py-2 text-sm">
          <span className="text-muted-foreground">健康状态</span>
          <span className="font-medium">
            {capability.healthy ? copy.pages.capabilities.healthy : copy.pages.capabilities.unhealthy}
          </span>
        </div>
        <div className="space-y-2">
          <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
            {copy.pages.capabilities.commands}
          </p>
          <div className="flex flex-wrap gap-2">
            {capability.commands.map((command) => (
              <span key={command} className="rounded-full bg-slate-100 px-2.5 py-1 text-xs text-slate-700">
                {command}
              </span>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
