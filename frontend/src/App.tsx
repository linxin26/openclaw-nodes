import { HashRouter, Navigate, Route, Routes } from "react-router-dom"

import { Header } from "@/components/layout/Header"
import { Sidebar } from "@/components/layout/Sidebar"
import { AppToaster } from "@/components/ui/sonner"
import { useConnection } from "@/hooks/useConnection"
import { About } from "@/pages/About"
import { Capabilities } from "@/pages/Capabilities"
import { Connection } from "@/pages/Connection"
import { Dashboard } from "@/pages/Dashboard"
import { Logs } from "@/pages/Logs"
import { Operations } from "@/pages/Operations"

export default function App() {
  const { status } = useConnection()

  return (
    <HashRouter>
      <div className="min-h-screen p-4 lg:p-6">
        <div className="mx-auto grid min-h-[calc(100vh-2rem)] max-w-[1440px] gap-4 lg:grid-cols-[280px_1fr]">
          <Sidebar />
          <div className="space-y-4">
            <Header status={status} />
            <main className="rounded-[30px] border border-white/70 bg-white/55 p-4 shadow-panel backdrop-blur md:p-6">
              <Routes>
                <Route path="/" element={<Dashboard />} />
                <Route path="/connection" element={<Connection />} />
                <Route path="/capabilities" element={<Capabilities />} />
                <Route path="/operations" element={<Operations />} />
                <Route path="/logs" element={<Logs />} />
                <Route path="/about" element={<About />} />
                <Route path="*" element={<Navigate to="/" replace />} />
              </Routes>
            </main>
          </div>
        </div>
      </div>
      <AppToaster />
    </HashRouter>
  )
}
