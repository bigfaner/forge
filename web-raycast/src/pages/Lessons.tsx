import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../lib/api'
import { MarkdownViewer } from '../components/MarkdownViewer'
import type { LessonCategory } from '../lib/types'

const CATEGORIES: { value: LessonCategory | 'all'; label: string }[] = [
  { value: 'all', label: 'All' },
  { value: 'debug', label: 'Debug' },
  { value: 'arch', label: 'Arch' },
  { value: 'tool', label: 'Tool' },
  { value: 'pattern', label: 'Pattern' },
  { value: 'gotcha', label: 'Gotcha' },
]

// Raycast-style category colors
const CAT_COLOR: Record<LessonCategory | 'all', { bg: string; text: string }> = {
  all:     { bg: 'rgba(255,255,255,0.06)', text: '#9c9c9d' },
  debug:   { bg: 'rgba(255,99,99,0.12)',   text: '#FF6363' },
  arch:    { bg: 'rgba(180,100,255,0.12)', text: '#c084fc' },
  tool:    { bg: 'rgba(85,179,255,0.12)',  text: 'hsl(202,100%,67%)' },
  pattern: { bg: 'rgba(34,211,238,0.12)',  text: '#22d3ee' },
  gotcha:  { bg: 'rgba(255,188,51,0.12)',  text: 'hsl(43,100%,60%)' },
  other:   { bg: 'rgba(255,255,255,0.06)', text: '#9c9c9d' },
}

export function Lessons() {
  const [category, setCategory] = useState<LessonCategory | 'all'>('all')
  const [search, setSearch] = useState('')
  const [expanded, setExpanded] = useState<string | null>(null)

  const { data } = useQuery({
    queryKey: ['lessons'],
    queryFn: () => api.lessons.list(),
    staleTime: 60_000,
  })

  const { data: lessonContent } = useQuery({
    queryKey: ['lessons', expanded],
    queryFn: () => api.lessons.get(expanded!),
    enabled: !!expanded,
  })

  const lessons = (data?.lessons ?? []).filter(l => {
    const matchCat = category === 'all' || l.category === category
    const matchSearch = !search ||
      l.title.toLowerCase().includes(search.toLowerCase()) ||
      l.excerpt.toLowerCase().includes(search.toLowerCase())
    return matchCat && matchSearch
  })

  return (
    <div style={{ padding: '32px 40px', maxWidth: 900 }}>
      <div style={{ marginBottom: 28 }}>
        <h1 style={{ fontSize: 24, fontWeight: 500, color: '#f9f9f9', letterSpacing: '0.2px', margin: 0 }}>
          Lessons
        </h1>
        <p style={{ fontSize: 14, color: '#6a6b6c', marginTop: 4, letterSpacing: '0.2px' }}>
          Reusable knowledge extracted from past work
        </p>
      </div>

      {/* Search */}
      <input
        type="text"
        placeholder="Search lessons…"
        value={search}
        onChange={e => setSearch(e.target.value)}
        style={{
          width: '100%',
          maxWidth: 360,
          borderRadius: 8,
          border: '1px solid rgba(255,255,255,0.08)',
          background: '#07080a',
          color: '#f9f9f9',
          fontSize: 14,
          fontWeight: 500,
          letterSpacing: '0.2px',
          padding: '8px 12px',
          outline: 'none',
          marginBottom: 20,
          fontFamily: 'inherit',
        }}
        onFocus={e => {
          e.currentTarget.style.borderColor = 'rgba(255,255,255,0.16)'
          e.currentTarget.style.boxShadow = 'hsla(202,100%,67%,0.15) 0px 0px 0px 3px'
        }}
        onBlur={e => {
          e.currentTarget.style.borderColor = 'rgba(255,255,255,0.08)'
          e.currentTarget.style.boxShadow = 'none'
        }}
      />

      {/* Category pills */}
      <div style={{ display: 'flex', gap: 6, flexWrap: 'wrap', marginBottom: 24 }}>
        {CATEGORIES.map(c => {
          const isActive = category === c.value
          const colors = CAT_COLOR[c.value]
          return (
            <button
              key={c.value}
              onClick={() => setCategory(c.value)}
              style={{
                padding: '4px 12px',
                borderRadius: 86,
                fontSize: 12,
                fontWeight: 600,
                letterSpacing: '0.2px',
                cursor: 'pointer',
                border: `1px solid ${isActive ? colors.text + '40' : 'rgba(255,255,255,0.08)'}`,
                background: isActive ? colors.bg : 'transparent',
                color: isActive ? colors.text : '#6a6b6c',
                transition: 'all 0.15s',
              }}
            >
              {c.label}
            </button>
          )
        })}
      </div>

      {/* Lesson list */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
        {lessons.length === 0 ? (
          <p style={{ fontSize: 14, color: '#6a6b6c', textAlign: 'center', padding: '48px 0', letterSpacing: '0.2px' }}>
            No lessons found
          </p>
        ) : (
          lessons.map(l => {
            const catColors = CAT_COLOR[l.category]
            const isOpen = expanded === l.name
            return (
              <div
                key={l.name}
                style={{
                  background: '#101111',
                  border: '1px solid rgba(255,255,255,0.06)',
                  boxShadow: 'rgb(27,28,30) 0px 0px 0px 1px, rgb(7,8,10) 0px 0px 0px 1px inset',
                  borderRadius: 10,
                  overflow: 'hidden',
                }}
              >
                <button
                  onClick={() => setExpanded(isOpen ? null : l.name)}
                  style={{
                    width: '100%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    gap: 12,
                    padding: '12px 16px',
                    textAlign: 'left',
                    cursor: 'pointer',
                    border: 'none',
                    background: 'transparent',
                    color: 'inherit',
                    fontFamily: 'inherit',
                  }}
                >
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8, minWidth: 0 }}>
                    <span style={{
                      flexShrink: 0,
                      fontSize: 11,
                      fontWeight: 600,
                      letterSpacing: '0.2px',
                      padding: '2px 8px',
                      borderRadius: 4,
                      background: catColors.bg,
                      color: catColors.text,
                    }}>
                      {l.category}
                    </span>
                    <span style={{ fontSize: 14, fontWeight: 500, color: '#f9f9f9', letterSpacing: '0.2px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                      {l.title}
                    </span>
                  </div>
                  <span style={{ fontSize: 12, color: '#434345', flexShrink: 0 }}>
                    {isOpen ? '▲' : '▼'}
                  </span>
                </button>
                {!isOpen && (
                  <p style={{ padding: '0 16px 12px', fontSize: 13, color: '#6a6b6c', letterSpacing: '0.2px', lineHeight: 1.5 }}>
                    {l.excerpt}
                  </p>
                )}
                {isOpen && lessonContent && (
                  <div style={{ padding: '0 16px 16px', borderTop: '1px solid rgba(255,255,255,0.06)', paddingTop: 12 }}>
                    <MarkdownViewer content={lessonContent} />
                  </div>
                )}
              </div>
            )
          })
        )}
      </div>
    </div>
  )
}
