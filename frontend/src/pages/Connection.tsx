import { useState } from "react"
import { Cable, Save, ServerCog, TestTube2 } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Switch } from "@/components/ui/switch"
import { useConfig } from "@/hooks/useConfig"
import { useConnection } from "@/hooks/useConnection"
import { copy } from "@/lib/copy"
import { api } from "@/lib/wails"

export function Connection() {
  const { config, setConfig, save, saving } = useConfig()
  const { status } = useConnection()
  const [testing, setTesting] = useState(false)
  const [connecting, setConnecting] = useState(false)
  const [latency, setLatency] = useState<string>("")

  const handleSave = async () => {
    await save(config)
    toast.success(copy.pages.connection.toasts.saved)
  }

  const handleTest = async () => {
    setTesting(true)
    try {
      const result = await api.connection.test()
      setLatency(result.success ? `${result.latencyMs}ms` : result.error)
      if (result.success) toast.success(copy.pages.connection.toasts.testSuccess(result.latencyMs))
      else toast.error(result.error || copy.pages.connection.toasts.testFailed)
    } finally {
      setTesting(false)
    }
  }

  const handleConnectToggle = async () => {
    setConnecting(true)
    try {
      if (status.status === "connected" || status.status === "connecting") {
        await api.connection.disconnect()
      } else {
        await api.connection.connect()
      }
    } finally {
      setConnecting(false)
    }
  }

  return (
    <div className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr]">
      <Card>
        <CardHeader>
          <CardTitle>{copy.pages.connection.settings.title}</CardTitle>
          <CardDescription>{copy.pages.connection.settings.description}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-5">
          <div className="space-y-2">
            <Label htmlFor="gateway">{copy.pages.connection.settings.labels.gateway}</Label>
            <Input
              id="gateway"
              value={config.gateway}
              onChange={(event) => setConfig({ ...config, gateway: event.target.value })}
              placeholder={copy.pages.connection.settings.placeholders.gateway}
            />
          </div>
          <div className="grid gap-5 md:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="port">{copy.pages.connection.settings.labels.port}</Label>
              <Input
                id="port"
                type="number"
                value={config.port}
                onChange={(event) =>
                  setConfig({ ...config, port: Number.parseInt(event.target.value || "0", 10) })
                }
              />
            </div>
            <div className="space-y-2">
              <Label>{copy.pages.connection.settings.labels.discovery}</Label>
              <Select value={config.discovery} onValueChange={(value) => setConfig({ ...config, discovery: value })}>
                <SelectTrigger>
                  <SelectValue placeholder={copy.pages.connection.settings.placeholders.selectMode} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="auto">{copy.pages.connection.settings.discovery.auto}</SelectItem>
                  <SelectItem value="manual">{copy.pages.connection.settings.discovery.manual}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label htmlFor="token">{copy.pages.connection.settings.labels.token}</Label>
              <Input
                id="token"
                type="password"
                value={config.token}
                onChange={(event) => setConfig({ ...config, token: event.target.value })}
                placeholder={copy.pages.connection.settings.placeholders.token}
              />
            </div>
          </div>
          <div className="flex items-center justify-between rounded-2xl bg-slate-50 p-4">
            <div>
              <p className="font-medium">{copy.pages.connection.settings.labels.tls}</p>
              <p className="text-sm text-muted-foreground">{copy.pages.connection.settings.tlsDescription}</p>
            </div>
            <Switch checked={config.tls} onCheckedChange={(checked) => setConfig({ ...config, tls: checked })} />
          </div>
          <div className="flex flex-wrap gap-3">
            <Button onClick={handleSave} disabled={saving}>
              <Save className="mr-2 h-4 w-4" />
              {copy.pages.connection.actions.save}
            </Button>
            <Button variant="secondary" onClick={handleTest} disabled={testing}>
              <TestTube2 className="mr-2 h-4 w-4" />
              {copy.pages.connection.actions.test}
            </Button>
            <Button variant="outline" onClick={handleConnectToggle} disabled={connecting}>
              <Cable className="mr-2 h-4 w-4" />
              {status.status === "connected" || status.status === "connecting"
                ? copy.pages.connection.actions.disconnect
                : copy.pages.connection.actions.connect}
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ServerCog className="h-5 w-5 text-primary" />
            {copy.pages.connection.session.title}
          </CardTitle>
          <CardDescription>{copy.pages.connection.session.description}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 text-sm">
          <StatusRow label={copy.pages.connection.session.labels.status} value={copy.common.status[status.status]} />
          <StatusRow label={copy.pages.connection.session.labels.gateway} value={status.gateway || copy.common.notConfigured} />
          <StatusRow label={copy.pages.connection.session.labels.tls} value={status.tls ? copy.common.enabled : copy.common.disabled} />
          <StatusRow label={copy.pages.connection.session.labels.testLatency} value={latency || copy.common.unavailable} />
          <StatusRow label={copy.pages.connection.session.labels.retryCount} value={`${status.retryCount}`} />
          <StatusRow label={copy.pages.connection.session.labels.retryDelay} value={`${status.retryDelayMs}ms`} />
        </CardContent>
      </Card>
    </div>
  )
}

function StatusRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between rounded-2xl bg-slate-50 px-4 py-3">
      <span className="text-muted-foreground">{label}</span>
      <span className="font-medium">{value}</span>
    </div>
  )
}
