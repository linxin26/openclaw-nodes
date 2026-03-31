import { useMemo, useState } from "react"
import { Copy, Trash2 } from "lucide-react"
import { toast } from "sonner"

import { VirtualLogList } from "@/components/shared/VirtualLogList"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { useLogs } from "@/hooks/useLogs"
import { copy } from "@/lib/copy"
import { cn } from "@/lib/utils"

export function Logs() {
  const { logs, setLogs } = useLogs()
  const [search, setSearch] = useState("")
  const [levels, setLevels] = useState<string[]>(["info", "warn", "error"])

  const filtered = useMemo(
    () =>
      logs.filter((log) => {
        if (!levels.includes(log.level)) return false
        if (search && !log.message.toLowerCase().includes(search.toLowerCase())) return false
        return true
      }),
    [logs, levels, search],
  )

  const copyLogs = async () => {
    // 复制时保留原始级别和正文，避免把排障上下文翻译掉。
    await navigator.clipboard.writeText(
      filtered.map((entry) => `${new Date(entry.timestamp).toISOString()} ${entry.level}: ${entry.message}`).join("\n"),
    )
    toast.success(copy.pages.logs.copied)
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-2">
        {["info", "warn", "error"].map((level) => (
          <Button
            key={level}
            variant={levels.includes(level) ? "default" : "outline"}
            size="sm"
            onClick={() =>
              setLevels((current) =>
                current.includes(level) ? current.filter((item) => item !== level) : [...current, level],
              )
            }
          >
            {copy.pages.logs.filters[level as keyof typeof copy.pages.logs.filters]}
          </Button>
        ))}
        <Input
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder={copy.pages.logs.searchPlaceholder}
          className="max-w-sm"
        />
        <Button variant="outline" size="icon" onClick={copyLogs} title={copy.common.copy}>
          <Copy className="h-4 w-4" />
        </Button>
        <Button variant="outline" size="icon" onClick={() => setLogs([])} title={copy.common.clear}>
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>

      <p className="text-sm text-muted-foreground">{copy.pages.logs.copyNote}</p>

      {filtered.length > 1000 ? (
        <VirtualLogList logs={filtered} />
      ) : (
        <div className="space-y-1 overflow-auto rounded-3xl bg-slate-950 p-4 font-mono text-xs text-slate-100">
          {filtered.map((log) => (
            <div key={`${log.timestamp}-${log.message}`} className="flex gap-3">
              <span className="text-slate-500">{new Date(log.timestamp).toLocaleTimeString()}</span>
              <span
                className={cn(
                  "w-12",
                  log.level === "error" && "text-rose-300",
                  log.level === "warn" && "text-amber-300",
                  log.level === "info" && "text-sky-300",
                )}
              >
                {log.level.toUpperCase()}
              </span>
              <span className="break-all">{log.message}</span>
            </div>
          ))}
          {filtered.length === 0 ? <p className="text-slate-400">{copy.pages.logs.empty}</p> : null}
        </div>
      )}
    </div>
  )
}

