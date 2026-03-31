type EventCallback = (payload?: unknown) => void

type RuntimeWindow = Window & {
  runtime?: {
    EventsOn?: (event: string, cb: EventCallback) => (() => void) | void
    EventsOff?: (event: string) => void
    EventsEmit?: (event: string, data?: unknown) => void
  }
}

export function EventsOn(event: string, callback: EventCallback) {
  const runtime = (window as RuntimeWindow).runtime
  if (!runtime?.EventsOn) {
    return () => undefined
  }
  const dispose = runtime.EventsOn(event, callback)
  if (typeof dispose === "function") {
    return dispose
  }
  return () => runtime.EventsOff?.(event)
}

export function EventsOff(event: string) {
  ;(window as RuntimeWindow).runtime?.EventsOff?.(event)
}

export function EventsEmit(event: string, data?: unknown) {
  ;(window as RuntimeWindow).runtime?.EventsEmit?.(event, data)
}
