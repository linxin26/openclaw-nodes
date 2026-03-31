import { useCallback, useEffect, useState } from "react"

import { api, type LogEntry, type LogFilter } from "@/lib/wails"
import { EventsOn } from "../../wailsjs/runtime/runtime"

export function useLogs(initialFilter?: LogFilter) {
  const [logs, setLogs] = useState<LogEntry[]>([])
  const [loading, setLoading] = useState(true)

  const reload = useCallback(async (filter = initialFilter) => {
    const next = await api.logs.get(filter)
    setLogs(next)
  }, [initialFilter])

  useEffect(() => {
    reload().finally(() => setLoading(false))
    const dispose = EventsOn("log", (payload) => {
      const entry = payload as LogEntry | undefined
      if (!entry) return
      setLogs((current) => [...current, entry])
    })
    return dispose
  }, [reload])

  return { logs, setLogs, loading, reload }
}
