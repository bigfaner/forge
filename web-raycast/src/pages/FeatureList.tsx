import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../lib/api'
import { FeatureCard } from '../components/FeatureCard'
import type { FeatureSummary } from '../lib/types'

type Filter = 'all' | 'active' | 'completed'

function filterFeatures(features: FeatureSummary[], filter: Filter) {
  if (filter === 'active') return features.filter(f => f.stats.in_progress > 0 || f.stats.pending > 0)
  if (filter === 'completed') return features.filter(f => f.stats.completed === f.stats.total && f.stats.total > 0)
  return features
}

const TABS: { value: Filter; label: string }[] = [
  { value: 'all', label: 'All' },
  { value: 'active', label: 'Active' },
  { value: 'completed', label: 'Completed' },
]

export function FeatureList() {
  const [filter, setFilter] = useState<Filter>('all')

  const { data, isLoading } = useQuery({
    queryKey: ['features'],
    queryFn: () => api.features.list(),
    refetchInterval: 30_000,
  })

  const features = filterFeatures(data?.features ?? [], filter)

  return (
    <div style={{ padding: '32px 40px', maxWidth: 1200 }}>
      {/* Header */}
      <div style={{ marginBottom: 28 }}>
        <h1 style={{ fontSize: 24, fontWeight: 500, color: '#f9f9f9', letterSpacing: '0.2px', margin: 0 }}>
          Features
        </h1>
        <p style={{ fontSize: 14, color: '#6a6b6c', marginTop: 4, letterSpacing: '0.2px' }}>
          {data?.features.length ?? 0} features total
        </p>
      </div>

      {/* Filter tabs — Raycast pill style */}
      <div
        style={{
          display: 'inline-flex',
          gap: 2,
          background: 'rgba(255,255,255,0.04)',
          border: '1px solid rgba(255,255,255,0.06)',
          borderRadius: 8,
          padding: 3,
          marginBottom: 24,
        }}
      >
        {TABS.map(t => (
          <button
            key={t.value}
            onClick={() => setFilter(t.value)}
            style={{
              padding: '5px 14px',
              borderRadius: 6,
              fontSize: 13,
              fontWeight: 500,
              letterSpacing: '0.2px',
              cursor: 'pointer',
              border: 'none',
              transition: 'all 0.15s',
              background: filter === t.value ? 'rgba(255,255,255,0.08)' : 'transparent',
              color: filter === t.value ? '#f9f9f9' : '#6a6b6c',
              boxShadow: filter === t.value ? 'rgb(27,28,30) 0px 0px 0px 1px' : 'none',
            }}
          >
            {t.label}
          </button>
        ))}
      </div>

      {/* Grid */}
      {isLoading ? (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: 12 }}>
          {[...Array(6)].map((_, i) => (
            <div key={i} style={{ height: 112, borderRadius: 12, background: 'rgba(255,255,255,0.03)' }} />
          ))}
        </div>
      ) : features.length === 0 ? (
        <div style={{ border: '1px dashed rgba(255,255,255,0.1)', borderRadius: 12, padding: '64px 32px', textAlign: 'center' }}>
          <p style={{ fontSize: 14, color: '#6a6b6c', letterSpacing: '0.2px' }}>No features found</p>
        </div>
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: 12 }}>
          {features.map(f => <FeatureCard key={f.slug} feature={f} />)}
        </div>
      )}
    </div>
  )
}
