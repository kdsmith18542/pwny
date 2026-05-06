import { useQuery } from "@tanstack/react-query"
import { modulesApi } from "@/lib/api/client"

export function useModules(type?: string) {
  return useQuery({
    queryKey: ["modules", type],
    queryFn: () => modulesApi.list(type),
  })
}

export function useModule(name: string | null) {
  return useQuery({
    queryKey: ["module", name],
    queryFn: () => (name ? modulesApi.get(name) : null),
    enabled: !!name,
  })
}
