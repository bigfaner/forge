import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { ArrowRight, CheckCircle2, AlertCircle, Clock } from 'lucide-react'
import { api } from '../lib/api'
import { FeatureCard } from '../components/FeatureCard'

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
    <div className="p-6 space-y-6">
      <div>
        <h1 className="text-xl font-semibold">Dashboard</h1>
        <p className="text-sm text-muted-foreground mt-0.5">Overview of all features and tasks</p>
      </div>

      {/* Stats grid */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-3">
        <StatCard label="Total Tasks" value={totals.total} icon={<Clock className="h-4 w-4 text-muted-foreground" />} />
        <StatCard label="Completed" value={totals.completed} icon={<CheckCircle2 className="h-4 w-4 text-green-400" />} accent="text-green-400" />
        <StatCard label="In Progress" value={totals.in_progress} icon={<ArrowRight className="h-4 w-4 text-blue-400" />} accent="text-blue-400" />
        <StatCard label="Blocked" value={totals.blocked} icon={<AlertCircle className="h-4 w-4 text-red-400" />} accent="text-red-400" />
      </div>

      {/* Active features */}
      <div>
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-sm font-semibold">Active Features</h2>
          <Link to="/features" className="text-xs text-accent hover:underline">
            View all →
          </Link>
        </div>
        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="h-24 rounded-lg bg-muted animate-pulse" />
            ))}
          </div>
        ) : activeFeatures.length === 0 ? (
          <EmptyState message="No active features" />
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
            {activeFeatures.map(f => <FeatureCard key={f.slug} feature={f} />)}
          </div>
        )}
      </div>

      {/* All features */}
      {features.length > 0 && activeFeatures.length < features.length && (
        <div>
          <h2 className="text-sm font-semibold mb-3">All Features</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3">
            {features
              .filter(f => f.stats.in_progress === 0 && f.stats.pending === 0)
              .map(f => <FeatureCard key={f.slug} feature={f} />)}
          </div>
        </div>
      )}
    </div>
  )
}

function StatCard({ label, value, icon, accent = 'text-foreground' }: {
  label: string; value: number; icon: React.ReactNode; accent?: string
}) {
  return (
    <div className="rounded-lg border border-border bg-card p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-xs text-muted-foreground">{label}</span>
        {icon}
      </div>
      <span className={`text-2xl font-bold ${accent}`}>{value}</span>
    </div>
  )
}

function EmptyState({ message }: { message: string }) {
  return (
    <div className="rounded-lg border border-dashed border-border p-8 text-center">
      <p className="text-sm text-muted-foreground">{message}</p>
    </div>
  )
}
