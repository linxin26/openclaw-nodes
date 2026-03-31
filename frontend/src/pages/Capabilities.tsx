import { toast } from "sonner"

import { CapabilityCard } from "@/components/shared/CapabilityCard"
import { LoadingSpinner } from "@/components/shared/LoadingSpinner"
import { useCapabilities } from "@/hooks/useCapabilities"
import { copy, getCapabilityDisplay } from "@/lib/copy"

export function Capabilities() {
  const { capabilities, loading, setCapability } = useCapabilities()

  if (loading) {
    return <LoadingSpinner label={copy.pages.capabilities.loading} />
  }

  return (
    <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
      {capabilities.map((capability) => {
        const display = getCapabilityDisplay(capability)
        return (
        <CapabilityCard
          key={capability.key}
          capability={capability}
          onToggle={async (enabled) => {
            await setCapability(capability.key, enabled)
            toast.success(`${display.name} ${enabled ? copy.common.enabled : copy.common.disabled}`)
          }}
        />
        )
      })}
    </div>
  )
}

