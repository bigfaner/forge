import type {
  ClaimResult,
  FeatureDetail,
  FeatureSummary,
  HealthInfo,
  LessonMeta,
  RecordEntry,
  Task,
  TaskDetail,
  TaskStatus,
} from './types'

const BASE = '/api'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...init,
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `HTTP ${res.status}`)
  }
  // Text endpoints (prd, design, lessons)
  const ct = res.headers.get('content-type') ?? ''
  if (ct.includes('text/plain') || ct.includes('text/markdown')) {
    return res.text() as unknown as T
  }
  return res.json()
}

export const api = {
  features: {
    list: (): Promise<{ features: FeatureSummary[] }> =>
      request('/features'),

    get: (slug: string): Promise<FeatureDetail> =>
      request(`/features/${slug}`),

    tasks: (slug: string): Promise<{ tasks: Task[] }> =>
      request(`/features/${slug}/tasks`),

    task: (slug: string, id: string): Promise<TaskDetail> =>
      request(`/features/${slug}/tasks/${encodeURIComponent(id)}`),

    prd: (slug: string): Promise<string> =>
      request(`/features/${slug}/prd`),

    design: (slug: string): Promise<string> =>
      request(`/features/${slug}/design`),

    records: (slug: string): Promise<{ records: RecordEntry[] }> =>
      request(`/features/${slug}/records`),
  },

  tasks: {
    claim: (): Promise<ClaimResult> =>
      request('/tasks/claim', { method: 'POST' }),

    setStatus: (slug: string, id: string, status: TaskStatus): Promise<void> =>
      request(`/features/${slug}/tasks/${encodeURIComponent(id)}/status`, {
        method: 'POST',
        body: JSON.stringify({ status }),
      }),
  },

  lessons: {
    list: (): Promise<{ lessons: LessonMeta[] }> =>
      request('/lessons'),

    get: (name: string): Promise<string> =>
      request(`/lessons/${name}`),
  },

  health: (): Promise<HealthInfo> =>
    request('/health'),

  records: {
    all: (): Promise<{ records: RecordEntry[] }> =>
      request('/records'),
  },
}
