import { useCallback, useEffect, useState } from "react"

import { api, type Config } from "@/lib/wails"
import { EventsOn } from "../../wailsjs/runtime/runtime"

const fallbackConfig: Config = {
  gateway: "",
  port: 18789,
  token: "",
  tls: false,
  discovery: "auto",
  capabilities: {},
}

export function useConfig() {
  const [config, setConfig] = useState<Config>(fallbackConfig)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  const reload = useCallback(async () => {
    const next = await api.config.get()
    setConfig(next)
    return next
  }, [])

  const save = useCallback(async (next: Config) => {
    setSaving(true)
    try {
      await api.config.save(next)
      setConfig(next)
    } finally {
      setSaving(false)
    }
  }, [])

  useEffect(() => {
    reload().finally(() => setLoading(false))
    const dispose = EventsOn("config:change", (payload) => {
      if (payload) setConfig(payload as Config)
    })
    return dispose
  }, [reload])

  return { config, setConfig, save, reload, loading, saving }
}
