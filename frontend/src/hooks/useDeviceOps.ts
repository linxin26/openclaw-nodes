import { useState } from "react"
import { toast } from "sonner"

import { copy } from "@/lib/copy"
import { api, type InvokeResult } from "@/lib/wails"

export function useDeviceOps() {
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<InvokeResult | null>(null)
  const [lastCall, setLastCall] = useState<{
    method: string
    params?: Record<string, unknown>
  } | null>(null)

  const invoke = async (method: string, params?: Record<string, unknown>) => {
    setLoading(true)
    setLastCall({ method, params })
    try {
      const next = await api.operations.invoke(method, params)
      setResult(next)
      if (next.success) toast.success(copy.actions.invokeCompleted(method))
      else toast.error(next.error || copy.actions.invokeFailed(method))
      return next
    } finally {
      setLoading(false)
    }
  }

  const retry = async () => {
    if (!lastCall) return null
    return invoke(lastCall.method, lastCall.params)
  }

  const saveFile = async (base64Data: string, filename: string) => {
    await api.operations.saveFileToDisk(base64Data, filename)
    toast.success(copy.actions.savedFile(filename))
  }

  return { loading, result, invoke, retry, saveFile, setResult }
}
