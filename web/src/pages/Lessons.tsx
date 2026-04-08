import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../lib/api'
import { MarkdownViewer } from '../components/MarkdownViewer'
import { LESSON_CATEGORY_COLOR, cn } from '../lib/utils'
import type { LessonCategory } from '../lib/types'

const CATEGORIES: { value: LessonCategory | 'all'; label: string }[] = [
  { value: 'all', label: 'All' },
  { value: 'debug', label: 'Debug' },
  { value: 'arch', label: 'Arch' },
  { value: 'tool', label: 'Tool' },
  { value: 'pattern', label: 'Pattern' },
  { value: 'gotcha', label: 'Gotcha' },
]

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
    <div className="p-6 space-y-5">
      <div>
        <h1 className="text-xl font-semibold">Lessons</h1>
        <p className="text-sm text-muted-foreground mt-0.5">Reusable knowledge extracted from past work</p>
      </div>

      {/* Search */}
      <input
        type="text"
        placeholder="Search lessons…"
        value={search}
        onChange={e => setSearch(e.target.value)}
        className="w-full max-w-sm rounded-md border border-border bg-muted px-3 py-2 text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring"
      />

      {/* Category tabs */}
      <div className="flex gap-1 flex-wrap">
        {CATEGORIES.map(c => (
          <button
            key={c.value}
            onClick={() => setCategory(c.value)}
            className={cn(
              'px-3 py-1 rounded-full text-xs font-medium transition-colors cursor-pointer border',
              category === c.value
                ? 'bg-accent/10 text-accent border-accent/30'
                : 'border-border text-muted-foreground hover:text-foreground',
            )}
          >
            {c.label}
          </button>
        ))}
      </div>

      {/* List */}
      <div className="space-y-2">
        {lessons.length === 0 ? (
          <p className="text-sm text-muted-foreground py-8 text-center">No lessons found</p>
        ) : (
          lessons.map(l => (
            <div key={l.name} className="rounded-lg border border-border bg-card overflow-hidden">
              <button
                onClick={() => setExpanded(expanded === l.name ? null : l.name)}
                className="w-full flex items-center justify-between gap-3 px-4 py-3 text-left cursor-pointer hover:bg-muted/50 transition-colors"
              >
                <div className="flex items-center gap-2 min-w-0">
                  <span className={cn('shrink-0 text-xs px-1.5 py-0.5 rounded', LESSON_CATEGORY_COLOR[l.category])}>
                    {l.category}
                  </span>
                  <span className="text-sm font-medium truncate">{l.title}</span>
                </div>
                <span className="text-xs text-muted-foreground shrink-0">{expanded === l.name ? '▲' : '▼'}</span>
              </button>
              {expanded !== l.name && (
                <p className="px-4 pb-3 text-xs text-muted-foreground">{l.excerpt}</p>
              )}
              {expanded === l.name && lessonContent && (
                <div className="px-4 pb-4 border-t border-border pt-3">
                  <MarkdownViewer content={lessonContent} />
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  )
}
