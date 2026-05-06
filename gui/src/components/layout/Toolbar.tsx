import { Bell, Search, User } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"

export default function Toolbar() {
  return (
    <header className="h-16 border-b bg-background px-6 flex items-center justify-between">
      <div className="relative w-96">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input 
          placeholder="Fuzzy find modules, sessions, hosts..." 
          className="pl-10 bg-muted/50 border-none focus-visible:ring-1"
        />
      </div>
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon">
          <Bell className="h-4 w-4" />
        </Button>
        <div className="h-8 w-[1px] bg-border" />
        <Button variant="ghost" size="icon" className="rounded-full">
          <User className="h-4 w-4" />
        </Button>
      </div>
    </header>
  )
}
