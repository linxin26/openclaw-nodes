import { useMemo } from "react"
import { Copy, Download, RotateCcw } from "lucide-react"
import { toast } from "sonner"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { copy } from "@/lib/copy"
import type { InvokeResult } from "@/lib/wails"

function findBase64Candidate(data: Record<string, unknown>) {
  for (const [key, value] of Object.entries(data)) {
    if (typeof value === "string" && value.length > 128 && /^[A-Za-z0-9+/=]+$/.test(value)) {
      return { key, value }
    }
  }
  return null
}

export function ResultPanel({
  title = copy.pages.operations.result.title,
  result,
  onRetry,
  onSave,
}: {
  title?: string
  result: InvokeResult | null
  onRetry?: () => void
  onSave?: (base64Data: string, filename: string) => Promise<void>
}) {
  const serialized = useMemo(
    () => (result ? JSON.stringify(result.success ? result.data : { error: result.error }, null, 2) : ""),
    [result],
  )
  const base64Candidate = result?.success ? findBase64Candidate(result.data) : null

  if (!result) {
    return null
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>
          {result.success ? copy.pages.operations.result.completedIn(result.durationMs) : result.error}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="rounded-2xl bg-slate-950 p-4 text-xs text-slate-100">
          <pre className="overflow-auto whitespace-pre-wrap break-all">{serialized}</pre>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={async () => {
              await navigator.clipboard.writeText(serialized)
              toast.success(copy.pages.operations.result.copied)
            }}
          >
            <Copy className="mr-2 h-4 w-4" />
            {copy.pages.operations.result.copy}
          </Button>
          {base64Candidate && onSave ? (
            <Button
              variant="secondary"
              size="sm"
              onClick={() => onSave(base64Candidate.value, `${base64Candidate.key}.bin`)}
            >
              <Download className="mr-2 h-4 w-4" />
              {copy.pages.operations.result.savePayload}
            </Button>
          ) : null}
          {onRetry ? (
            <Button variant="ghost" size="sm" onClick={onRetry}>
              <RotateCcw className="mr-2 h-4 w-4" />
              {copy.pages.operations.result.retry}
            </Button>
          ) : null}
        </div>
      </CardContent>
    </Card>
  )
}
