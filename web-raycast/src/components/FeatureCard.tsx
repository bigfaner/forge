import { Link } from 'react-router-dom'
import type { FeatureSummary } from '../lib/types'
import { progressPercent, formatDate } from '../lib/utils'

interface Props {
  feature: FeatureSummary
}

export function FeatureCard({ feature }: Props) {
  const pct = progressPercent(feature.stats)
  const { stats } = feature

  return (
    <Link
      to={`/features/${feature.slug}`}
      style={{
        display: 'block',
        background: '#101111',
        border: '1px solid rgba(255,255,255,0.06)',
        boxShadow: 'rgb(27,28,30) 0px 0px 0px 1px, rgb(7,8,10) 0px 0px 0px 1px inset',
        borderRadius: 12,
        padding: '16px 20px',
        textDecoration: 'none',
        color: 'inherit',
        transition: 'border-color 0.15s, box-shadow 0.15s',
      }}
      onMouseEnter={e => {
        const el = e.currentTarget as HTMLElement
        el.style.borderColor = 'rgba(255,255,255,0.12)'
        el.style.boxShadow = 'rgb(40,41,43) 0px 0px 0px 1px, rgb(7,8,10) 0px 0px 0px 1px inset'
      }}
      onMouseLeave={e => {
        const el = e.currentTarget as HTMLElement
        el.style.borderColor = 'rgba(255,255,255,0.06)'
        el.style.boxShadow = 'rgb(27,28,30) 0px 0px 0px 1px, rgb(7,8,10) 0px 0px 0px 1px inset'
      }}
    >
      <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', gap: 8, marginBottom: 12 }}>
        <div style={{ minWidth: 0 }}>
          <p style={{ fontFamily: 'monospace', fontSize: 11, color: '#434345', marginBottom: 3, letterSpacing: '0.3px' }}>
            {feature.slug}
          </p>
          <h3 style={{ fontSize: 14, fontWeight: 600, color: '#f9f9f9', letterSpacing: '0.2px', lineHeight: 1.3, margin: 0 }}>
            {feature.title}
          </h3>
        </div>
        <span style={{ flexShrink: 0, fontSize: 13, fontWeight: 600, color: '#FF6363', letterSpacing: '0.1px' }}>
          {pct}%
        </span>
      </div>

      {/* Progress bar */}
      <div style={{ height: 2, width: '100%', background: 'rgba(255,255,255,0.06)', borderRadius: 1, overflow: 'hidden', marginBottom: 12 }}>
        <div
          style={{
            height: '100%',
            borderRadius: 1,
            width: `${pct}%`,
            background: pct === 100 ? 'hsl(151,59%,59%)' : '#FF6363',
            transition: 'width 0.3s ease',
          }}
        />
      </div>

      {/* Stats row */}
      <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
        <span style={{ fontSize: 12, color: 'hsl(202,100%,67%)', letterSpacing: '0.2px' }}>
          {stats.in_progress} active
        </span>
        <span style={{ fontSize: 12, color: 'hsl(151,59%,59%)', letterSpacing: '0.2px' }}>
          {stats.completed}/{stats.total} done
        </span>
        {stats.blocked > 0 && (
          <span style={{ fontSize: 12, color: '#FF6363', letterSpacing: '0.2px' }}>
            {stats.blocked} blocked
          </span>
        )}
        <span style={{ marginLeft: 'auto', fontSize: 11, color: '#434345', letterSpacing: '0.2px' }}>
          {formatDate(feature.lastUpdated)}
        </span>
      </div>
    </Link>
  )
}
