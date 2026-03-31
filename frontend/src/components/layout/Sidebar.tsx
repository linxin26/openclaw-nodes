import { Cpu, Logs, PlugZap, Radar, Settings2, Sparkles } from "lucide-react"
import { NavLink } from "react-router-dom"

import { copy } from "@/lib/copy"
import { cn } from "@/lib/utils"

const items = [
  { to: "/", label: copy.navigation.dashboard, icon: Sparkles },
  { to: "/connection", label: copy.navigation.connection, icon: PlugZap },
  { to: "/capabilities", label: copy.navigation.capabilities, icon: Settings2 },
  { to: "/operations", label: copy.navigation.operations, icon: Radar },
  { to: "/logs", label: copy.navigation.logs, icon: Logs },
  { to: "/about", label: copy.navigation.about, icon: Cpu },
]

export function Sidebar() {
  return (
    <aside className="flex h-full flex-col rounded-[30px] border border-white/70 bg-slate-950 px-4 py-5 text-white shadow-panel">
      <div className="mb-8 px-3">
        <p className="text-xs uppercase tracking-[0.32em] text-sky-300">{copy.layout.sidebar.brand}</p>
        <h1 className="mt-3 text-2xl font-semibold">{copy.layout.sidebar.title}</h1>
      </div>
      <nav className="flex flex-1 flex-col gap-2">
        {items.map((item) => {
          const Icon = item.icon
          return (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === "/"}
              className={({ isActive }) =>
                cn(
                  "flex items-center gap-3 rounded-2xl px-3 py-3 text-sm transition-colors",
                  isActive ? "bg-white/12 text-white" : "text-slate-300 hover:bg-white/8",
                )
              }
            >
              <Icon className="h-4 w-4" />
              <span>{item.label}</span>
            </NavLink>
          )
        })}
      </nav>
      <div className="rounded-2xl border border-white/10 bg-white/5 p-4 text-xs text-slate-300">
        {copy.layout.sidebar.description}
      </div>
    </aside>
  )
}
