import { useQuery } from "@tanstack/react-query"
import { sessionsApi, modulesApi } from "@/lib/api/client"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Badge } from "@/components/ui/badge"
import { 
  Terminal, 
  Box, 
  Activity, 
  Monitor,
  Loader2
} from "lucide-react"

export default function Dashboard() {
  const { data: sessions, isLoading: sessionsLoading } = useQuery({
    queryKey: ["sessions"],
    queryFn: sessionsApi.list,
    refetchInterval: 5000,
  })

  const { data: modules, isLoading: modulesLoading } = useQuery({
    queryKey: ["modules"],
    queryFn: () => modulesApi.list(),
  })

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
        <Badge variant="outline" className="bg-primary/10 text-primary border-primary/20">
          Server: 127.0.0.1:31337
        </Badge>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatsCard 
          title="Active Sessions" 
          value={sessions?.length || 0} 
          icon={<Terminal className="h-4 w-4" />}
          loading={sessionsLoading}
        />
        <StatsCard 
          title="Total Modules" 
          value={modules?.length || 0} 
          icon={<Box className="h-4 w-4" />}
          loading={modulesLoading}
        />
        <StatsCard 
          title="Running Jobs" 
          value={0} 
          icon={<Activity className="h-4 w-4" />}
        />
        <StatsCard 
          title="Hosts in Scope" 
          value={0} 
          icon={<Monitor className="h-4 w-4" />}
        />
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
        <Card className="col-span-4">
          <CardHeader>
            <CardTitle>Recent Sessions</CardTitle>
          </CardHeader>
          <CardContent>
            <ScrollArea className="h-[300px]">
              {sessionsLoading ? (
                <div className="flex items-center justify-center h-full">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : sessions?.length === 0 ? (
                <div className="text-center py-10 text-muted-foreground">
                  No active sessions
                </div>
              ) : (
                <div className="space-y-4">
                  {sessions?.map((s) => (
                    <div key={s.id} className="flex items-center justify-between p-3 rounded-lg border bg-muted/30">
                      <div className="flex items-center gap-3">
                        <Terminal className="h-4 w-4 text-primary" />
                        <div>
                          <div className="text-sm font-medium">Session {s.id.slice(0, 8)}</div>
                          <div className="text-xs text-muted-foreground">{s.target} • {s.platform}</div>
                        </div>
                      </div>
                      <Badge variant="secondary">{s.type}</Badge>
                    </div>
                  ))}
                </div>
              )}
            </ScrollArea>
          </CardContent>
        </Card>

        <Card className="col-span-3">
          <CardHeader>
            <CardTitle>Job Queue</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-center py-10 text-muted-foreground text-sm">
              No active jobs
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

function StatsCard({ title, value, icon, loading }: { title: string, value: number | string, icon: React.ReactNode, loading?: boolean }) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <div className="text-muted-foreground">{icon}</div>
      </CardHeader>
      <CardContent>
        {loading ? (
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
        ) : (
          <div className="text-2xl font-bold">{value}</div>
        )}
      </CardContent>
    </Card>
  )
}
