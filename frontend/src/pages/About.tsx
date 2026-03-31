import { Copy, FolderOpen } from "lucide-react"
import { useEffect, useState } from "react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { copy } from "@/lib/copy"
import { api, type AboutInfo } from "@/lib/wails"

export function About() {
  const [info, setInfo] = useState<AboutInfo | null>(null)

  useEffect(() => {
    api.about.get().then(setInfo).catch(() => undefined)
  }, [])

  if (!info) return null

  return (
    <div className="grid gap-6 xl:grid-cols-[0.9fr_1.1fr]">
      <Card>
        <CardHeader>
          <CardTitle>{copy.pages.about.runtimeIdentity}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4 text-sm">
          <CopyRow label={copy.pages.about.labels.deviceId} value={info.deviceId} />
          <CopyRow label={copy.pages.about.labels.publicKey} value={info.publicKey} />
          <CopyRow label={copy.pages.about.labels.dataDir} value={info.dataDir} />
          <Button variant="secondary" onClick={() => api.about.openPath(info.dataDir)}>
            <FolderOpen className="mr-2 h-4 w-4" />
            {copy.pages.about.buttons.openDataDirectory}
          </Button>
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>{copy.pages.about.buildMetadata}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3 text-sm">
          <InfoRow label={copy.pages.about.labels.version} value={info.version} />
          <InfoRow label={copy.pages.about.labels.platform} value={info.platform} />
          <InfoRow label={copy.pages.about.labels.hostname} value={info.hostname} />
          <InfoRow label={copy.pages.about.labels.goVersion} value={info.goVersion} />
          <InfoRow label={copy.pages.about.labels.arch} value={info.arch} />
          <InfoRow label={copy.pages.about.labels.protocolVersion} value={`${info.protocolVersion}`} />
        </CardContent>
      </Card>
    </div>
  )
}

function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-center justify-between rounded-2xl bg-slate-50 px-4 py-3">
      <span className="text-muted-foreground">{label}</span>
      <span className="font-medium">{value}</span>
    </div>
  )
}

function CopyRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-2xl bg-slate-50 p-4">
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="text-muted-foreground">{label}</p>
          <p className="mt-1 break-all font-medium">{value}</p>
        </div>
        <Button
          variant="ghost"
          size="icon"
          onClick={async () => {
            await navigator.clipboard.writeText(value)
            toast.success(copy.pages.about.toasts.copied(label))
          }}
        >
          <Copy className="h-4 w-4" />
        </Button>
      </div>
    </div>
  )
}
