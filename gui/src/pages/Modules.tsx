import { useState } from "react"
import { useModules, useModule } from "@/hooks/useModules"
import { Input } from "@/components/ui/input"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Search, Loader2, Box } from "lucide-react"
import { modulesApi } from "@/lib/api/client"

export default function Modules() {
  const [search, setSearch] = useState("")
  const [selectedModule, setSelectedModule] = useState<string | null>(null)
  
  const { data: modules, isLoading } = useModules()
  const { data: moduleDetails, isLoading: isDetailsLoading } = useModule(selectedModule)

  const filteredModules = modules?.filter(m => 
    m.name.toLowerCase().includes(search.toLowerCase()) ||
    m.description.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="h-full flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Modules</h1>
          <p className="text-muted-foreground">Browse and run exploits, auxiliary, and post-exploitation modules.</p>
        </div>
      </div>

      <div className="flex-1 flex gap-6 overflow-hidden">
        {/* Module List */}
        <div className="w-1/3 flex flex-col gap-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input 
              placeholder="Search modules..." 
              className="pl-10"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
          <Card className="flex-1 overflow-hidden flex flex-col">
            <ScrollArea className="flex-1">
              <div className="p-2 space-y-1">
                {isLoading ? (
                  <div className="flex items-center justify-center p-8">
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                  </div>
                ) : filteredModules?.map((m) => (
                  <button
                    key={m.name}
                    className={`w-full text-left px-3 py-2 rounded-md transition-colors text-sm ${
                      selectedModule === m.name 
                        ? "bg-primary text-primary-foreground" 
                        : "hover:bg-muted"
                    }`}
                    onClick={() => setSelectedModule(m.name)}
                  >
                    <div className="font-medium truncate">{m.name}</div>
                    <div className={`text-xs truncate ${selectedModule === m.name ? "text-primary-foreground/70" : "text-muted-foreground"}`}>
                      {m.description}
                    </div>
                  </button>
                ))}
              </div>
            </ScrollArea>
          </Card>
        </div>

        {/* Module Details */}
        <div className="flex-1 overflow-hidden flex flex-col">
          {selectedModule ? (
            <Card className="flex-1 overflow-hidden flex flex-col">
              {isDetailsLoading ? (
                <div className="flex-1 flex items-center justify-center">
                  <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                </div>
              ) : moduleDetails && (
                <>
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <div className="space-y-1">
                        <CardTitle className="text-2xl">{moduleDetails.name}</CardTitle>
                        <CardDescription>{moduleDetails.description}</CardDescription>
                      </div>
                      <Badge variant="outline" className="capitalize">
                        {moduleDetails.type}
                      </Badge>
                    </div>
                  </CardHeader>
                  <ScrollArea className="flex-1">
                    <CardContent className="space-y-6">
                      <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-1">
                          <div className="text-sm font-medium">Platforms</div>
                          <div className="flex flex-wrap gap-1">
                            {moduleDetails.platforms?.map((p: string) => (
                              <Badge key={p} variant="secondary">{p}</Badge>
                            ))}
                          </div>
                        </div>
                        <div className="space-y-1">
                          <div className="text-sm font-medium">Architectures</div>
                          <div className="flex flex-wrap gap-1">
                            {moduleDetails.arch?.map((a: string) => (
                              <Badge key={a} variant="secondary">{a}</Badge>
                            ))}
                          </div>
                        </div>
                      </div>

                      <div className="space-y-4">
                        <div className="text-lg font-semibold">Options</div>
                        <div className="space-y-4">
                          {Object.entries(moduleDetails.options || {}).map(([name, opt]: [string, any]) => (
                            <div key={name} className="space-y-1.5">
                              <div className="flex items-center gap-2">
                                <span className="text-sm font-medium">{name}</span>
                                {opt.required && <Badge variant="destructive" className="h-4 px-1 text-[10px]">Required</Badge>}
                              </div>
                              <Input 
                                defaultValue={opt.default} 
                                placeholder={opt.description}
                              />
                              <p className="text-xs text-muted-foreground">{opt.description}</p>
                            </div>
                          ))}
                        </div>
                      </div>
                    </CardContent>
                  </ScrollArea>
                  <div className="p-6 border-t flex justify-end gap-3">
                    <Button variant="outline" onClick={async () => {
                      if (selectedModule) {
                        try {
                          const res = await modulesApi.validate(selectedModule, {})
                          alert(res.valid ? "Validation successful" : "Validation failed: " + res.error)
                        } catch (err) {
                          alert("Error validating module")
                        }
                      }
                    }}>Validate</Button>
                    <Button onClick={async () => {
                      if (selectedModule) {
                        try {
                          const res = await modulesApi.run(selectedModule, {})
                          alert("Job started: " + res.job_id)
                        } catch (err) {
                          alert("Error starting job")
                        }
                      }
                    }}>Run Module</Button>
                  </div>
                </>
              )}
            </Card>
          ) : (
            <div className="flex-1 flex flex-col items-center justify-center text-muted-foreground bg-muted/20 rounded-lg border-2 border-dashed">
              <Box className="h-12 w-12 mb-4 opacity-20" />
              <p>Select a module to view details</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
