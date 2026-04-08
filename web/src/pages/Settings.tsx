import { useQuery } from '@tanstack/react-query'
import { api } from '../lib/api'

export function Settings() {
  const { data } = useQuery({
    queryKey: ['health'],
    queryFn: () => api.health(),
    refetchInterval: 60_000,
  })

  return (
    <div className="p-6 max-w-lg space-y-6">
      <div>
        <h1 className="text-xl font-semibold">Settings</h1>
        <p className="text-sm text-muted-foreground mt-0.5">Server info and configuration</p>
      </div>

      <div className="rounded-lg border border-border bg-card divide-y divide-border">
        <InfoRow label="task-cli version" value={data?.version ?? '—'} />
        <InfoRow label="Project root" value={data?.projectRoot ?? '—'} mono />
        <InfoRow label="Current feature" value={data?.currentFeature || 'none'} mono />
        <InfoRow label="Server port" value="7300" />
        <InfoRow label="Poll interval" value="30s" />
      </div>
    </div>
  )
}

function InfoRow({ label, value, mono = false }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="flex items-center justify-between gap-4 px-4 py-3">
      <span className="text-sm text-muted-foreground">{label}</span>
      <span className={`text-sm font-medium ${mono ? 'font-mono' : ''}`}>{value}</span>
    </div>
  )
}
