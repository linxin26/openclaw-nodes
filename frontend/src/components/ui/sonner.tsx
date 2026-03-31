import { Toaster } from "sonner"

export function AppToaster() {
  return (
    <Toaster
      position="top-right"
      richColors
      toastOptions={{
        className: "!rounded-2xl !border !border-white/70 !bg-white/90 !text-slate-900",
      }}
    />
  )
}
