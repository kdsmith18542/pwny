import { useEffect, useRef, useState } from "react"
import { Terminal as XTerm } from "xterm"
import { FitAddon } from "@xterm/addon-fit"
import "xterm/css/xterm.css"
import { useWebSocket } from "@/hooks/useWebSocket"
import { useQuery } from "@tanstack/react-query"
import { sessionsApi } from "@/lib/api/client"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Card } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Terminal as TerminalIcon, Loader2, X } from "lucide-react"

export default function Sessions() {
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null)
  const { data: sessions, isLoading } = useQuery({
    queryKey: ["sessions"],
    queryFn: sessionsApi.list,
    refetchInterval: 5000,
  })

  return (
    <div className="h-full flex flex-col gap-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Sessions</h1>
        <p className="text-muted-foreground">Manage and interact with active compromised sessions.</p>
      </div>

      <div className="flex-1 flex gap-6 overflow-hidden">
        {/* Session List */}
        <div className="w-1/4">
          <Card className="h-full overflow-hidden flex flex-col">
            <ScrollArea className="flex-1">
              <div className="p-2 space-y-1">
                {isLoading ? (
                  <div className="flex items-center justify-center p-8">
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                  </div>
                ) : sessions?.length === 0 ? (
                  <div className="text-center p-8 text-sm text-muted-foreground">
                    No active sessions
                  </div>
                ) : sessions?.map((s) => (
                  <button
                    key={s.id}
                    className={`w-full text-left px-3 py-2 rounded-md transition-colors text-sm ${
                      activeSessionId === s.id 
                        ? "bg-primary text-primary-foreground" 
                        : "hover:bg-muted"
                    }`}
                    onClick={() => setActiveSessionId(s.id)}
                  >
                    <div className="flex items-center justify-between gap-2">
                      <span className="font-medium truncate">Session {s.id.slice(0, 8)}</span>
                      <Badge variant="outline" className="text-[10px] px-1 h-4 border-primary-foreground/30">
                        {s.type}
                      </Badge>
                    </div>
                    <div className={`text-xs truncate ${activeSessionId === s.id ? "text-primary-foreground/70" : "text-muted-foreground"}`}>
                      {s.target} • {s.platform}
                    </div>
                  </button>
                ))}
              </div>
            </ScrollArea>
          </Card>
        </div>

        {/* Console Area */}
        <div className="flex-1 flex flex-col gap-4 overflow-hidden">
          {activeSessionId ? (
            <SessionConsole sessionId={activeSessionId} onClose={() => setActiveSessionId(null)} />
          ) : (
            <div className="flex-1 flex flex-col items-center justify-center text-muted-foreground bg-muted/20 rounded-lg border-2 border-dashed">
              <TerminalIcon className="h-12 w-12 mb-4 opacity-20" />
              <p>Select a session to interact</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

function SessionConsole({ sessionId, onClose }: { sessionId: string; onClose: () => void }) {
  const terminalRef = useRef<HTMLDivElement>(null)
  const xtermRef = useRef<XTerm | null>(null)
  const { lastMessage, sendMessage, isConnected } = useWebSocket(`ws://127.0.0.1:31337/api/v1/sessions/${sessionId}/ws`)

  useEffect(() => {
    if (!terminalRef.current) return

    const term = new XTerm({
      theme: {
        background: "transparent",
        foreground: "#ffffff",
      },
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'Menlo, Monaco, "Courier New", monospace',
    })
    const fitAddon = new FitAddon()
    term.loadAddon(fitAddon)
    term.open(terminalRef.current)
    fitAddon.fit()

    term.onData((data) => {
      sendMessage({ type: "session:input", data })
    })

    xtermRef.current = term

    return () => {
      term.dispose()
    }
  }, [sessionId, sendMessage])

  useEffect(() => {
    if (lastMessage && lastMessage.type === "session:output") {
      xtermRef.current?.write(lastMessage.data)
    }
  }, [lastMessage])

  return (
    <Card className="flex-1 bg-black text-white flex flex-col overflow-hidden relative">
      <div className="bg-zinc-900 px-4 py-2 flex items-center justify-between text-xs border-b border-zinc-800">
        <div className="flex items-center gap-2">
          <div className={`w-2 h-2 rounded-full ${isConnected ? "bg-green-500" : "bg-red-500"}`} />
          <span>Interactive Console — {sessionId}</span>
        </div>
        <button onClick={onClose} className="hover:text-red-400">
          <X className="h-3 w-3" />
        </button>
      </div>
      <div ref={terminalRef} className="flex-1 p-2" />
    </Card>
  )
}
