import { useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Copy, Download, Loader2, CheckCircle2 } from "lucide-react"
import { payloadApi } from "@/lib/api/client"
import { useToast } from "@/components/ui/use-toast"

const PLATFORMS = ["windows", "linux", "macos"]
const ARCHS = ["x64", "x86"]
const PAYLOADS = [
  { name: "windows/x64/tcp", label: "Windows x64 Reverse TCP", platform: "windows", arch: "x64" },
  { name: "windows/x86/tcp", label: "Windows x86 Reverse TCP", platform: "windows", arch: "x86" },
  { name: "linux/x64/tcp", label: "Linux x64 Reverse TCP", platform: "linux", arch: "x64" },
  { name: "multi/http", label: "Multi HTTP Reverse Stager", platform: "any", arch: "any" },
]
const FORMATS = ["raw", "exe", "python", "powershell", "c"]
const ENCODERS = ["none", "x86/shikata_ga_nai"]

export default function Payload() {
  const { toast } = useToast()
  const [loading, setLoading] = useState(false)
  const [platform, setPlatform] = useState("windows")
  const [arch, setArch] = useState("x64")
  const [selectedPayload, setSelectedPayload] = useState("windows/x64/tcp")
  const [lhost, setLhost] = useState("127.0.0.1")
  const [lport, setLport] = useState("4444")
  const [format, setFormat] = useState("raw")
  const [encoder, setEncoder] = useState("none")
  const [iterations, setIterations] = useState(1)
  
  const [result, setResult] = useState<{ payload: string; size: number } | null>(null)

  const handleGenerate = async () => {
    setLoading(true)
    try {
      const data = await payloadApi.generate({
        name: selectedPayload,
        platform,
        arch,
        lhost,
        lport: parseInt(lport),
        format,
        encoder: encoder === "none" ? "" : encoder,
        iterations,
      })
      setResult(data)
      toast({
        title: "Payload generated",
        description: `Successfully generated ${data.size} bytes of shellcode.`,
      })
    } catch (err) {
      toast({
        title: "Generation failed",
        description: "An error occurred while generating the payload.",
        variant: "destructive",
      })
    } finally {
      setLoading(false)
    }
  }

  const copyToClipboard = () => {
    if (result) {
      navigator.clipboard.writeText(result.payload)
      toast({ title: "Copied to clipboard" })
    }
  }

  return (
    <div className="h-full flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Payload Generator</h1>
          <p className="text-muted-foreground">Configure and generate custom stagers and stages.</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-1 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Target Configuration</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-1.5">
                <label className="text-sm font-medium">Platform</label>
                <Select value={platform} onValueChange={setPlatform}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {PLATFORMS.map(p => <SelectItem key={p} value={p}>{p}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <label className="text-sm font-medium">Architecture</label>
                <Select value={arch} onValueChange={setArch}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {ARCHS.map(a => <SelectItem key={a} value={a}>{a}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <label className="text-sm font-medium">Payload</label>
                <Select value={selectedPayload} onValueChange={setSelectedPayload}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {PAYLOADS.filter(p => p.platform === platform || p.platform === "any").map(p => (
                      <SelectItem key={p.name} value={p.name}>{p.label}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Network Options</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-1.5">
                <label className="text-sm font-medium">LHOST</label>
                <Input value={lhost} onChange={e => setLhost(e.target.value)} />
              </div>
              <div className="space-y-1.5">
                <label className="text-sm font-medium">LPORT</label>
                <Input value={lport} onChange={e => setLport(e.target.value)} />
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="lg:col-span-1 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Output & Encoding</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-1.5">
                <label className="text-sm font-medium">Format</label>
                <Select value={format} onValueChange={setFormat}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {FORMATS.map(f => <SelectItem key={f} value={f}>{f}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-1.5">
                <label className="text-sm font-medium">Encoder</label>
                <Select value={encoder} onValueChange={setEncoder}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {ENCODERS.map(e => <SelectItem key={e} value={e}>{e}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              {encoder !== "none" && (
                <div className="space-y-1.5">
                  <label className="text-sm font-medium">Iterations</label>
                  <Input type="number" value={iterations} onChange={e => setIterations(parseInt(e.target.value))} />
                </div>
              )}
              <Button className="w-full" onClick={handleGenerate} disabled={loading}>
                {loading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : "Generate Payload"}
              </Button>
            </CardContent>
          </Card>
        </div>

        <div className="lg:col-span-1 flex flex-col h-full">
          <Card className="flex-1 flex flex-col">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Result</CardTitle>
                  <CardDescription>Generated shellcode output</CardDescription>
                </div>
                {result && (
                  <div className="flex gap-2">
                    <Button variant="ghost" size="icon" onClick={copyToClipboard}>
                      <Copy className="h-4 w-4" />
                    </Button>
                  </div>
                )}
              </div>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col">
              {result ? (
                <div className="flex-1 flex flex-col gap-4">
                  <div className="flex items-center gap-2 text-green-500 text-sm font-medium">
                    <CheckCircle2 className="h-4 w-4" />
                    Success! {result.size} bytes generated.
                  </div>
                  <div className="flex-1 bg-muted p-4 rounded-md font-mono text-[10px] break-all overflow-auto max-h-[400px]">
                    {result.payload}
                  </div>
                  <Button variant="outline" className="w-full gap-2">
                    <Download className="h-4 w-4" /> Download Binary
                  </Button>
                </div>
              ) : (
                <div className="flex-1 flex flex-col items-center justify-center text-muted-foreground opacity-30">
                  <Loader2 className="h-12 w-12 mb-4" />
                  <p>Awaiting generation...</p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
