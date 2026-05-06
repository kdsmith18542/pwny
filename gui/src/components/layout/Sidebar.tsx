import { NavLink } from "react-router-dom"
import { 
  LayoutDashboard, 
  Box, 
  Terminal, 
  Shield, 
  Database, 
  Settings,
  Skull
} from "lucide-react"
import { cn } from "@/lib/utils"

const navItems = [
  { icon: LayoutDashboard, label: "Dashboard", to: "/" },
  { icon: Box, label: "Modules", to: "/modules" },
  { icon: Terminal, label: "Sessions", to: "/sessions" },
  { icon: Shield, label: "Payload", to: "/payload" },
  { icon: Database, label: "Workspace", to: "/workspace" },
  { icon: Settings, label: "Settings", to: "/settings" },
]

export default function Sidebar() {
  return (
    <aside className="w-64 border-r bg-muted/30 flex flex-col">
      <div className="p-6 flex items-center gap-2 mb-4">
        <Skull className="h-8 w-8 text-primary" />
        <span className="text-2xl font-bold tracking-tight">Pwny</span>
      </div>
      <nav className="flex-1 px-4 space-y-1">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              cn(
                "flex items-center gap-3 px-3 py-2 rounded-md transition-colors text-sm font-medium",
                isActive 
                  ? "bg-primary text-primary-foreground" 
                  : "text-muted-foreground hover:bg-muted hover:text-foreground"
              )
            }
          >
            <item.icon className="h-4 w-4" />
            {item.label}
          </NavLink>
        ))}
      </nav>
      <div className="p-4 border-t text-xs text-muted-foreground">
        v0.1.0-alpha
      </div>
    </aside>
  )
}
