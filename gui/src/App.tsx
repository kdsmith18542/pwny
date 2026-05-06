import { BrowserRouter, Routes, Route } from "react-router-dom"
import AppShell from "./components/layout/AppShell"
import Dashboard from "./pages/Dashboard"
import Modules from "./pages/Modules"
import Sessions from "./pages/Sessions"
import Workspace from "./pages/Workspace"
import Payload from "./pages/Payload"
import Settings from "./pages/Settings"
import { Toaster } from "./components/ui/toaster"
import "./App.css"

function App() {
  return (
    <BrowserRouter>
      <AppShell>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/modules" element={<Modules />} />
          <Route path="/sessions" element={<Sessions />} />
          <Route path="/payload" element={<Payload />} />
          <Route path="/workspace" element={<Workspace />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </AppShell>
      <Toaster />
    </BrowserRouter>
  )
}

export default App
