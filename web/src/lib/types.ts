// --- Feature ---
export interface FeatureStats {
  total: number
  pending: number
  in_progress: number
  completed: number
  blocked: number
  skipped: number
}

export interface FeatureSummary {
  slug: string
  title: string
  stats: FeatureStats
  lastUpdated: string
}

export interface FeatureDetail extends FeatureSummary {
  tasks: Task[]
}

// --- Task ---
export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'blocked' | 'skipped'
export type TaskPriority = 'P0' | 'P1' | 'P2'

export interface TaskRecord {
  summary: string
  filesCreated: string[]
  filesModified: string[]
  decisions: string[]
  testResults: string
  coverage: string
  commitHash: string
  completedAt: string
}

export interface Task {
  id: string
  title: string
  description: string
  phase: number
  priority: TaskPriority
  status: TaskStatus
  estimatedTime: string
  dependencies: string[]
  files: string[]
  record?: string | TaskRecord
}

export type TaskDetail = Omit<Task, 'record'> & {
  record?: TaskRecord
}

export interface ClaimResult {
  taskId: string
  key: string
  title: string
  file: string
}

// --- Lesson ---
export type LessonCategory = 'debug' | 'arch' | 'tool' | 'pattern' | 'gotcha' | 'other'

export interface LessonMeta {
  name: string
  category: LessonCategory
  title: string
  excerpt: string
}

// --- Health ---
export interface HealthInfo {
  version: string
  projectRoot: string
  currentFeature: string
}

// --- Record timeline ---
export interface RecordEntry {
  featureSlug: string
  taskId: string
  taskTitle: string
  coverage: string
  filesChanged: number
  commitHash: string
  completedAt: string
}
