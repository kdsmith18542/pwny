import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"

export default function Settings() {
  return (
    <div className="h-full flex flex-col gap-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
        <p className="text-muted-foreground">Configure global framework options and UI preferences.</p>
      </div>

      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>API Connection</CardTitle>
            <CardDescription>Configure how the GUI connects to the Pwny server.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <label className="text-sm font-medium">Server Address</label>
                <Input defaultValue="127.0.0.1" />
              </div>
              <div className="space-y-1.5">
                <label className="text-sm font-medium">Server Port</label>
                <Input defaultValue="31337" />
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="outline" className="bg-green-500/10 text-green-500 border-green-500/20">Connected</Badge>
              <span className="text-xs text-muted-foreground">Latency: 2ms</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Appearance</CardTitle>
            <CardDescription>Customize the look and feel of the application.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex gap-2">
              <Button variant="outline" className="ring-2 ring-primary">Dark Mode</Button>
              <Button variant="outline">Light Mode</Button>
              <Button variant="outline">System</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
