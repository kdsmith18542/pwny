import { useState } from "react"
import { useQuery, useQueryClient } from "@tanstack/react-query"
import { workspaceApi } from "@/lib/api/client"
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { 
  Loader2, 
  Monitor, 
  Server, 
  Laptop, 
  HelpCircle, 
  Plus,
  RefreshCcw,
  Search
} from "lucide-react"


export default function Workspace() {
  const queryClient = useQueryClient()
  const [activeWorkspaceId, setActiveWorkspaceId] = useState<string | null>(null)
  const [search, setSearch] = useState("")

  const { data: workspaces, isLoading: workspacesLoading } = useQuery({
    queryKey: ["workspaces"],
    queryFn: workspaceApi.list,
  })

  const { data: hosts, isLoading: hostsLoading } = useQuery({
    queryKey: ["hosts", activeWorkspaceId],
    queryFn: () => activeWorkspaceId ? workspaceApi.getHosts(activeWorkspaceId) : Promise.resolve([]),
    enabled: !!activeWorkspaceId,
  })

  // Select first workspace by default if none selected
  if (!activeWorkspaceId && workspaces && workspaces.length > 0) {
    setActiveWorkspaceId(workspaces[0].id)
  }

  const filteredHosts = hosts?.filter((h: any) => 
    h.address.toLowerCase().includes(search.toLowerCase()) ||
    h.os_name?.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="h-full flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Workspace</h1>
          <p className="text-muted-foreground">Track discovered hosts, services, and vulnerabilities in the current engagement.</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" className="gap-2" onClick={() => queryClient.invalidateQueries({ queryKey: ["workspaces"] })}>
            <RefreshCcw className="h-4 w-4" /> Refresh
          </Button>
          <Button size="sm" className="gap-2">
            <Plus className="h-4 w-4" /> New Workspace
          </Button>
        </div>
      </div>

      <div className="flex-1 flex gap-6 overflow-hidden">
        {/* Workspace Selector */}
        <div className="w-1/4 flex flex-col gap-4">
          <Card className="flex-1 overflow-hidden flex flex-col">
            <CardHeader className="py-4">
              <CardTitle className="text-sm">Active Workspace</CardTitle>
            </CardHeader>
            <CardContent className="p-2 space-y-1">
              {workspacesLoading ? (
                <div className="flex items-center justify-center p-4">
                  <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />
                </div>
              ) : workspaces?.length === 0 ? (
                <div className="text-center p-4 text-xs text-muted-foreground">
                  No workspaces found
                </div>
              ) : workspaces?.map((w: any) => (
                <button
                  key={w.id}
                  className={`w-full text-left px-3 py-2 rounded-md transition-colors text-sm ${
                    activeWorkspaceId === w.id 
                      ? "bg-primary text-primary-foreground" 
                      : "hover:bg-muted"
                  }`}
                  onClick={() => setActiveWorkspaceId(w.id)}
                >
                  <div className="font-medium truncate">{w.name}</div>
                  <div className={`text-xs truncate ${activeWorkspaceId === w.id ? "text-primary-foreground/70" : "text-muted-foreground"}`}>
                    {w.description || "No description"}
                  </div>
                </button>
              ))}
            </CardContent>
          </Card>
        </div>

        {/* Hosts Table */}
        <div className="flex-1 flex flex-col gap-4 overflow-hidden">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input 
              placeholder="Search hosts..." 
              className="pl-10"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
          
          <Card className="flex-1 overflow-hidden flex flex-col">
            <CardContent className="flex-1 overflow-hidden p-0">
              {hostsLoading ? (
                <div className="flex items-center justify-center p-12">
                  <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                </div>
              ) : !activeWorkspaceId ? (
                <div className="flex flex-col items-center justify-center p-12 text-muted-foreground">
                  <Plus className="h-12 w-12 mb-4 opacity-20" />
                  <p>Select a workspace to view hosts</p>
                </div>
              ) : filteredHosts?.length === 0 ? (
                <div className="flex flex-col items-center justify-center p-12 text-muted-foreground">
                  <Monitor className="h-12 w-12 mb-4 opacity-20" />
                  <p>No hosts discovered yet</p>
                </div>
              ) : (
                <Table>
                  <TableHeader className="bg-muted/50 sticky top-0 z-10">
                    <TableRow>
                      <TableHead>Address</TableHead>
                      <TableHead>OS</TableHead>
                      <TableHead>Services</TableHead>
                      <TableHead>Status</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredHosts?.map((host: any) => (
                      <TableRow key={host.id} className="cursor-pointer hover:bg-muted/50">
                        <TableCell className="font-medium">
                          <div className="flex items-center gap-2">
                            {getHostIcon(host.purpose)}
                            {host.address}
                          </div>
                        </TableCell>
                        <TableCell className="text-xs">
                          {host.os_name || "Unknown"} {host.os_sp}
                          <div className="text-[10px] text-muted-foreground">{host.os_flavor}</div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline" className="font-mono text-[10px]">
                            {host.services_count || 0} Open
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <Badge variant={host.state === 'alive' ? 'default' : 'secondary'} className="capitalize text-[10px]">
                            {host.state}
                          </Badge>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}

function getHostIcon(purpose: string) {
  switch (purpose) {
    case 'server': return <Server className="h-4 w-4 text-primary" />
    case 'client': return <Laptop className="h-4 w-4 text-blue-500" />
    case 'device': return <Monitor className="h-4 w-4 text-green-500" />
    default: return <HelpCircle className="h-4 w-4 text-muted-foreground" />
  }
}
