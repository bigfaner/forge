import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { ArrowRight, CheckCircle2, AlertCircle, Clock } from 'lucide-react'
import { api } from '../lib/api'
import { FeatureCard } from '../components/FeatureCard'

const rCard = {
  background: '#101111',
  border: '1px solid rgba(255,255,255,0.06)',
  boxShadow: 'rgb(27,28,30) 0px 0px 0px 1px, rgb(7,8,10) 0px 0px 0px 1px inset',
  borderRadius: 12,
  padding: '20px 24px',
} as const

export function Dashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ['features'],
    queryFn: () => api.features.list(),
    refetchInterval: 30_000,
  })

  const features = data?.features ?? []

  const totals = features.reduce(
    (acc, f) => ({
      total: acc.total + f.stats.total,
      completed: acc.completed + f.stats.completed,
      in_progress: acc.in_progress + f.stats.in_progress,
      blocked: acc.blocked + f.stats.blocked,
    }),
    { total: 0, completed: 0, in_progress: 0, blocked: 0 },
  )

  const activeFeatures = features.filter(f => f.stats.in_progress > 0 || f.stats.pending > 0)

  return (
    <div style={{ padding: '32px 40px', maxWidth: 1200 }}>
      {/* Header */}
      <div style={{ marginBottom: 32 }}>
        <h1 style={{ fontSize: 24, fontWeight: 500, color: '#f9f9f9', letterSpacing: '0.2px', margin: 0 }}>
          Dashboard
        </h1>
        <p style={{ fontSize: 14, color: '#6a6b6c', marginTop: 4, letterSpacing: '0.2px' }}>
          Overview of all features and tasks
        </p>
      </div>

      {/* Stats grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: 12, marginBottom: 40 }}>
        <StatCard
          label="Total Tasks"
          value={totals.total}
          icon={<Clock size={14} style={{ color: '#6a6b6c' }} />}
          valueColor="#f9f9f9"
        />
        <StatCard
          label="Completed"
          value={totals.completed}
          icon={<CheckCircle2 size={14} style={{ color: 'hsl(151,59%,59%)' }} />}
          valueColor="hsl(151,59%,59%)"
        />
        <StatCard
          label="In Progress"
          value={totals.in_progress}
          icon={<ArrowRight size={14} style={{ color: 'hsl(202,100%,67%)' }} />}
          valueColor="hsl(202,100%,67%)"
        />
        <StatCard
          label="Blocked"
          value={totals.blocked}
          icon={<AlertCircle size={14} style={{ color: '#FF6363' }} />}
          valueColor="#FF6363"
        />
      </div>

      {/* Active features */}
      <div style={{ marginBottom: 40 }}>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 16 }}>
          <h2 style={{ fontSize: 14, fontWeight: 600, color: '#f9f9f9', letterSpacing: '0.2px', margin: 0 }}>
            Active Features
          </h2>
          <Link
            to="/features"
            style={{ fontSize: 12, color: '#9c9c9d', textDecoration: 'none', letterSpacing: '0.3px', display: 'flex', alignItems: 'center', gap: 4 }}
            onMouseEnter={e => (e.currentTarget.style.color = '#f9f9f9')}
            onMouseLeave={e => (e.currentTarget.style.color = '#9c9c9d')}
          >
            View all <ArrowRight size={12} />
          </Link>
        </div>

        {isLoading ? (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 12 }}>
            {[...Array(3)].map((_, i) => (
              <div key={i} style={{ ...rCard, height: 96, background: 'rgba(255,255,255,0.03)' }} />
            ))}
          </div>
        ) : activeFeatures.length === 0 ? (
          <EmptyState message="No active features" />
        ) : (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: 12 }}>
            {activeFeatures.map(f => <FeatureCard key={f.slug} feature={f} />)}
          </div>
        )}
      </div>

      {/* Completed features */}
      {features.length > 0 && activeFeatures.length < features.length && (
        <div>
          <h2 style={{ fontSize: 14, fontWeight: 600, color: '#6a6b6c', letterSpacing: '0.2px', marginBottom: 16 }}>
            Completed Features
          </h2>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: 12 }}>
            {features
              .filter(f => f.stats.in_progress === 0 && f.stats.pending === 0)
              .map(f => <FeatureCard key={f.slug} feature={f} />)}
          </div>
        </div>
      )}
    </div>
  )
}

function StatCard({ label, value, icon, valueColor }: {
  label: string; value: number; icon: React.ReactNode; valueColor: string
}) {
  return (
    <div style={rCard}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
        <span style={{ fontSize: 12, color: '#6a6b6c', fontWeight: 500, letterSpacing: '0.2px' }}>{label}</span>
        {icon}
      </div>
      <span style={{ fontSize: 28, fontWeight: 600, color: valueColor, letterSpacing: '-0.5px' }}>
        {value}
      </span>
    </div>
  )
}

function EmptyState({ message }: { message: string }) {
  return (
    <div style={{ border: '1px dashed rgba(255,255,255,0.1)', borderRadius: 12, padding: '48px 32px', textAlign: 'center' }}>
      <p style={{ fontSize: 14, color: '#6a6b6c', letterSpacing: '0.2px' }}>{message}</p>
    </div>
  )
}
