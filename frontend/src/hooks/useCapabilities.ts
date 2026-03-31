import { useCallback, useEffect, useState } from "react"

import { api, type CapabilityInfo } from "@/lib/wails"
import { EventsOn } from "../../wailsjs/runtime/runtime"

export function useCapabilities() {
  const [capabilities, setCapabilities] = useState<CapabilityInfo[]>([])
  const [loading, setLoading] = useState(true)

  const reload = useCallback(async () => {
    const next = await api.capabilities.get()
    setCapabilities(next)
  }, [])

  const setCapability = useCallback(async (key: string, enabled: boolean) => {
    await api.capabilities.set(key, enabled)
    setCapabilities((current) =>
      current.map((item) => (item.key === key ? { ...item, enabled } : item)),
    )
  }, [])

  useEffect(() => {
    reload().finally(() => setLoading(false))
    const dispose = EventsOn("capability:change", (payload) => {
      const capability = payload as CapabilityInfo | undefined
      if (!capability) return
      setCapabilities((current) =>
        current.map((item) => (item.key === capability.key ? capability : item)),
      )
    })
    return dispose
  }, [reload])

  return { capabilities, loading, reload, setCapability }
}
