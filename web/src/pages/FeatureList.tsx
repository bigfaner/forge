import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../lib/api'
import { FeatureCard } from '../components/FeatureCard'
import type { FeatureSummary } from '../lib/types'
import { cn } from '../lib/utils'

type Filter = 'all' | 'active' | 'completed'

function filterFeatures(features: FeatureSummary[], filter: Filter) {
  if (filter === 'active') return features.filter(f => f.stats.in_progress > 0 || f.stats.pending > 0)
  if (filter === 'completed') return features.filter(f => f.stats.completed === f.stats.total && f.stats.total > 0)
  return features
}

export function FeatureList() {
  const [filter, setFilter] = useState<Filter>('all')

  const { data, isLoading } = useQuery({
    queryKey: ['features'],
    queryFn: () => api.features.list(),
    refetchInterval: 30_000,
  })

  const features = filterFeatures(data?.features ?? [], filter)

  const TABS: { value: Filter; label: string }[] = [
    { value: 'all', label: 'All' },
    { value: 'active', label: 'Active' },
    { value: 'completed', label: 'Completed' },
  ]

  return (
    <div className="p-6 space-y-5">
      <div>
        <h1 className="text-xl font-semibold">Features</h1>
        <p className="text-sm text-muted-foreground mt-0.5">{data?.features.length ?? 0} features total</p>
      </div>

      {/* Filter tabs */}
      <div className="flex gap-1 rounded-lg bg-muted p-1 w-fit">
        {TABS.map(t => (
          <button
            key={t.value}
            onClick={() => setFilter(t.value)}
            className={cn(
              'px-3 py-1.5 rounded-md text-sm font-medium transition-colors cursor-pointer',
              filter === t.value
                ? 'bg-card text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground',
            )}
          >
            {t.label}
          </button>
        ))}
      </div>

      {/* Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
          {[...Array(6)].map((_, i) => <div key={i} className="h-28 rounded-lg bg-muted animate-pulse" />)}
        </div>
      ) : features.length === 0 ? (
        <div className="rounded-lg border border-dashed border-border p-12 text-center">
          <p className="text-sm text-muted-foreground">No features found</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
          {features.map(f => <FeatureCard key={f.slug} feature={f} />)}
        </div>
      )}
    </div>
  )
}
