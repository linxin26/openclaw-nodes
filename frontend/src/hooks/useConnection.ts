import { useEffect, useState } from "react"

import { api, type ConnectionStatus } from "@/lib/wails"
import { EventsOn } from "../../wailsjs/runtime/runtime"

const fallbackStatus: ConnectionStatus = {
  status: "offline",
  gateway: "",
  tls: false,
  uptimeMs: 0,
  retryCount: 0,
  retryDelayMs: 0,
  protocolVersion: 3,
}

export function useConnection() {
  const [status, setStatus] = useState<ConnectionStatus>(fallbackStatus)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true
    api.status
      .get()
      .then((next) => {
        if (mounted) setStatus(next)
      })
      .finally(() => {
        if (mounted) setLoading(false)
      })

    const dispose = EventsOn("status:change", (payload) => {
      if (payload) setStatus(payload as ConnectionStatus)
    })
    return () => {
      mounted = false
      dispose()
    }
  }, [])

  return { status, loading }
}
