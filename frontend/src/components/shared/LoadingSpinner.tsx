import { copy } from "@/lib/copy"

export function LoadingSpinner({ label = copy.common.loading }: { label?: string }) {
  return (
    <div className="inline-flex items-center gap-2 text-sm text-muted-foreground">
      <span className="h-4 w-4 animate-spin rounded-full border-2 border-primary/30 border-t-primary" />
      <span>{label}</span>
    </div>
  )
}
