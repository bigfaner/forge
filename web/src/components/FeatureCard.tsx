import { Link } from 'react-router-dom'
import type { FeatureSummary } from '../lib/types'
import { progressPercent, formatDate, cn } from '../lib/utils'

interface Props {
  feature: FeatureSummary
}

export function FeatureCard({ feature }: Props) {
  const pct = progressPercent(feature.stats)
  const { stats } = feature

  return (
    <Link
      to={`/features/${feature.slug}`}
      className={cn(
        'block rounded-lg border border-border bg-card p-4',
        'hover:border-accent/40 hover:bg-card/80 transition-colors',
      )}
    >
      <div className="flex items-start justify-between gap-2 mb-3">
        <div>
          <p className="text-xs text-muted-foreground font-mono mb-0.5">{feature.slug}</p>
          <h3 className="font-semibold text-sm leading-snug">{feature.title}</h3>
        </div>
        <span className="shrink-0 text-xs font-medium text-accent">{pct}%</span>
      </div>

      {/* Progress bar */}
      <div className="h-1.5 w-full rounded-full bg-muted overflow-hidden mb-3">
        <div
          className="h-full rounded-full bg-accent transition-all"
          style={{ width: `${pct}%` }}
        />
      </div>

      {/* Stats row */}
      <div className="flex gap-3 text-xs text-muted-foreground">
        <span className="text-blue-400">{stats.in_progress} active</span>
        <span className="text-green-400">{stats.completed}/{stats.total} done</span>
        {stats.blocked > 0 && <span className="text-red-400">{stats.blocked} blocked</span>}
        <span className="ml-auto">{formatDate(feature.lastUpdated)}</span>
      </div>
    </Link>
  )
}
