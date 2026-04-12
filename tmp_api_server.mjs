#!/usr/bin/env node
// API server for ZCode Dashboard — reads docs/features/ filesystem
import http from 'http'
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const PROJECT_ROOT = __dirname  // /Users/nasuki/zcode
const FEATURES_DIR = path.join(PROJECT_ROOT, 'docs', 'features')
const LESSONS_DIR = path.join(PROJECT_ROOT, 'docs', 'lessons')
const PORT = 7300

function json(res, data, status = 200) {
  res.writeHead(status, { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' })
  res.end(JSON.stringify(data))
}

function text(res, data, status = 200) {
  res.writeHead(status, { 'Content-Type': 'text/plain; charset=utf-8', 'Access-Control-Allow-Origin': '*' })
  res.end(data)
}

function err(res, msg, status = 404) {
  res.writeHead(status, { 'Content-Type': 'text/plain', 'Access-Control-Allow-Origin': '*' })
  res.end(msg)
}

// Parse index.json into FeatureSummary
function loadFeature(slug) {
  const indexPath = path.join(FEATURES_DIR, slug, 'index.json')
  if (!fs.existsSync(indexPath)) return null
  const index = JSON.parse(fs.readFileSync(indexPath, 'utf8'))
  const tasks = Object.values(index.tasks || {})
  const stats = {
    total: tasks.length,
    pending: tasks.filter(t => t.status === 'pending').length,
    in_progress: tasks.filter(t => t.status === 'in_progress').length,
    completed: tasks.filter(t => t.status === 'completed').length,
    blocked: tasks.filter(t => t.status === 'blocked').length,
    skipped: tasks.filter(t => t.status === 'skipped').length,
  }
  // lastUpdated: newest record mtime, or index mtime
  const recordsDir = path.join(FEATURES_DIR, slug, 'records')
  let lastUpdated = fs.statSync(indexPath).mtime.toISOString()
  if (fs.existsSync(recordsDir)) {
    for (const f of fs.readdirSync(recordsDir)) {
      const mt = fs.statSync(path.join(recordsDir, f)).mtime.toISOString()
      if (mt > lastUpdated) lastUpdated = mt
    }
  }
  // title: feature field or slug
  const title = index.feature
    ? index.feature.replace(/-/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
    : slug
  return { slug, title, stats, lastUpdated, _index: index }
}

// Parse a task entry into Task shape expected by frontend
function taskFromEntry(key, entry, slug) {
  // Load record if exists
  let record = undefined
  if (entry.record) {
    const recPath = path.join(FEATURES_DIR, slug, entry.record.replace(/\.md$/, '.json'))
    const recPathMd = path.join(FEATURES_DIR, slug, entry.record)
    if (fs.existsSync(recPath)) {
      try {
        const r = JSON.parse(fs.readFileSync(recPath, 'utf8'))
        record = {
          summary: r.summary || '',
          filesCreated: r.filesCreated || [],
          filesModified: r.filesModified || [],
          decisions: r.keyDecisions || [],
          testResults: r.testsPassed !== undefined ? `${r.testsPassed} passed` : '',
          coverage: r.coverage ? `${r.coverage}%` : '',
          commitHash: r.commitHash || '',
          completedAt: r.completedAt || '',
        }
      } catch {}
    }
  }
  // Load task file for description/phase/files
  let description = ''
  let phase = 1
  let files = []
  const taskFile = path.join(FEATURES_DIR, slug, entry.file || `tasks/${key}.md`)
  if (fs.existsSync(taskFile)) {
    const content = fs.readFileSync(taskFile, 'utf8')
    const descMatch = content.match(/## Description\s+([\s\S]*?)(?=##|$)/)
    if (descMatch) description = descMatch[1].trim()
    const phaseMatch = content.match(/Phase[:\s]+(\d+)/)
    if (phaseMatch) phase = parseInt(phaseMatch[1])
    const filesMatch = content.match(/## Files\s+([\s\S]*?)(?=##|$)/)
    if (filesMatch) {
      files = filesMatch[1].trim().split('\n').map(l => l.replace(/^[-*]\s*/, '').trim()).filter(Boolean)
    }
  }

  return {
    id: entry.id || key,
    title: entry.title,
    description,
    phase,
    priority: entry.priority || 'P2',
    status: entry.status,
    estimatedTime: entry.estimatedTime || '',
    dependencies: entry.dependencies || [],
    files,
    record,
  }
}

const server = http.createServer((req, res) => {
  if (req.method === 'OPTIONS') {
    res.writeHead(204, { 'Access-Control-Allow-Origin': '*', 'Access-Control-Allow-Methods': 'GET,POST', 'Access-Control-Allow-Headers': 'Content-Type' })
    return res.end()
  }

  const url = new URL(req.url, `http://localhost:${PORT}`)
  const p = url.pathname.replace(/^\/api/, '')

  // GET /features
  if (p === '/features' && req.method === 'GET') {
    const slugs = fs.existsSync(FEATURES_DIR)
      ? fs.readdirSync(FEATURES_DIR).filter(d => fs.statSync(path.join(FEATURES_DIR, d)).isDirectory())
      : []
    const features = slugs.map(loadFeature).filter(Boolean).map(({ _index, ...f }) => f)
    return json(res, { features })
  }

  // GET /features/:slug
  const featureMatch = p.match(/^\/features\/([^/]+)$/)
  if (featureMatch && req.method === 'GET') {
    const slug = featureMatch[1]
    const f = loadFeature(slug)
    if (!f) return err(res, 'feature not found')
    const { _index, ...summary } = f
    const tasks = Object.entries(_index.tasks || {}).map(([k, v]) => taskFromEntry(k, v, slug))
    return json(res, { ...summary, tasks })
  }

  // GET /features/:slug/tasks
  const tasksMatch = p.match(/^\/features\/([^/]+)\/tasks$/)
  if (tasksMatch && req.method === 'GET') {
    const slug = tasksMatch[1]
    const f = loadFeature(slug)
    if (!f) return err(res, 'feature not found')
    const tasks = Object.entries(f._index.tasks || {}).map(([k, v]) => taskFromEntry(k, v, slug))
    return json(res, { tasks })
  }

  // GET /features/:slug/tasks/:id
  const taskMatch = p.match(/^\/features\/([^/]+)\/tasks\/(.+)$/)
  if (taskMatch && req.method === 'GET') {
    const [, slug, id] = taskMatch
    const f = loadFeature(slug)
    if (!f) return err(res, 'feature not found')
    const entry = Object.entries(f._index.tasks || {}).find(([k, v]) => v.id === id || k === id)
    if (!entry) return err(res, 'task not found')
    const task = taskFromEntry(entry[0], entry[1], slug)
    return json(res, task)
  }

  // POST /features/:slug/tasks/:id/status
  const statusMatch = p.match(/^\/features\/([^/]+)\/tasks\/(.+)\/status$/)
  if (statusMatch && req.method === 'POST') {
    const [, slug, id] = statusMatch
    let body = ''
    req.on('data', c => body += c)
    req.on('end', () => {
      try {
        const { status } = JSON.parse(body)
        const indexPath = path.join(FEATURES_DIR, slug, 'index.json')
        const index = JSON.parse(fs.readFileSync(indexPath, 'utf8'))
        const key = Object.keys(index.tasks).find(k => index.tasks[k].id === id || k === id)
        if (!key) return err(res, 'task not found')
        index.tasks[key].status = status
        fs.writeFileSync(indexPath, JSON.stringify(index, null, 2) + '\n')
        json(res, { ok: true })
      } catch (e) { err(res, e.message, 400) }
    })
    return
  }

  // GET /features/:slug/prd
  const prdMatch = p.match(/^\/features\/([^/]+)\/prd$/)
  if (prdMatch && req.method === 'GET') {
    const slug = prdMatch[1]
    const f = path.join(FEATURES_DIR, slug, 'prd.md')
    if (!fs.existsSync(f)) return err(res, 'prd not found')
    return text(res, fs.readFileSync(f, 'utf8'))
  }

  // GET /features/:slug/design
  const designMatch = p.match(/^\/features\/([^/]+)\/design$/)
  if (designMatch && req.method === 'GET') {
    const slug = designMatch[1]
    const f = path.join(FEATURES_DIR, slug, 'design.md')
    if (!fs.existsSync(f)) return err(res, 'design not found')
    return text(res, fs.readFileSync(f, 'utf8'))
  }

  // GET /records
  if (p === '/records' && req.method === 'GET') {
    const records = []
    if (fs.existsSync(FEATURES_DIR)) {
      for (const slug of fs.readdirSync(FEATURES_DIR)) {
        const recDir = path.join(FEATURES_DIR, slug, 'records')
        if (!fs.existsSync(recDir)) continue
        const indexPath = path.join(FEATURES_DIR, slug, 'index.json')
        const index = fs.existsSync(indexPath) ? JSON.parse(fs.readFileSync(indexPath, 'utf8')) : {}
        for (const file of fs.readdirSync(recDir).filter(f => f.endsWith('.json'))) {
          try {
            const r = JSON.parse(fs.readFileSync(path.join(recDir, file), 'utf8'))
            const taskKey = r.taskId || file.replace('.json', '')
            const taskEntry = Object.values(index.tasks || {}).find(t => t.id === taskKey)
            records.push({
              featureSlug: slug,
              taskId: taskKey,
              taskTitle: taskEntry?.title || taskKey,
              coverage: r.coverage ? `${r.coverage}%` : '',
              filesChanged: (r.filesCreated?.length || 0) + (r.filesModified?.length || 0),
              commitHash: r.commitHash || '',
              completedAt: r.completedAt || '',
            })
          } catch {}
        }
      }
    }
    records.sort((a, b) => (b.completedAt || '').localeCompare(a.completedAt || ''))
    return json(res, { records })
  }

  // GET /lessons
  if (p === '/lessons' && req.method === 'GET') {
    const lessons = []
    const CATS = ['debug', 'arch', 'tool', 'pattern', 'gotcha']
    if (fs.existsSync(LESSONS_DIR)) {
      for (const file of fs.readdirSync(LESSONS_DIR).filter(f => f.endsWith('.md'))) {
        const content = fs.readFileSync(path.join(LESSONS_DIR, file), 'utf8')
        const name = file.replace('.md', '')
        // Parse frontmatter
        const fmMatch = content.match(/^---\s*([\s\S]*?)\s*---/)
        let category = 'other', title = name
        if (fmMatch) {
          const fm = fmMatch[1]
          const catM = fm.match(/category:\s*(\S+)/)
          const titleM = fm.match(/title:\s*(.+)/)
          if (catM && CATS.includes(catM[1])) category = catM[1]
          if (titleM) title = titleM[1].trim()
        }
        // Excerpt: first non-empty line after frontmatter
        const body = content.replace(/^---[\s\S]*?---\s*/, '')
        const excerpt = body.split('\n').map(l => l.trim()).find(l => l && !l.startsWith('#')) || ''
        lessons.push({ name, category, title, excerpt: excerpt.slice(0, 120) })
      }
    }
    return json(res, { lessons })
  }

  // GET /lessons/:name
  const lessonMatch = p.match(/^\/lessons\/(.+)$/)
  if (lessonMatch && req.method === 'GET') {
    const name = lessonMatch[1]
    const f = path.join(LESSONS_DIR, `${name}.md`)
    if (!fs.existsSync(f)) return err(res, 'lesson not found')
    return text(res, fs.readFileSync(f, 'utf8'))
  }

  // POST /tasks/claim
  if (p === '/tasks/claim' && req.method === 'POST') {
    // Find first pending task across all features
    if (fs.existsSync(FEATURES_DIR)) {
      for (const slug of fs.readdirSync(FEATURES_DIR)) {
        const indexPath = path.join(FEATURES_DIR, slug, 'index.json')
        if (!fs.existsSync(indexPath)) continue
        const index = JSON.parse(fs.readFileSync(indexPath, 'utf8'))
        for (const [key, task] of Object.entries(index.tasks || {})) {
          if (task.status === 'pending') {
            index.tasks[key].status = 'in_progress'
            fs.writeFileSync(indexPath, JSON.stringify(index, null, 2) + '\n')
            return json(res, { taskId: task.id || key, key, title: task.title, file: task.file || '' })
          }
        }
      }
    }
    return err(res, 'no pending tasks', 404)
  }

  // GET /health
  if (p === '/health' && req.method === 'GET') {
    return json(res, {
      version: '0.1.0',
      projectRoot: PROJECT_ROOT,
      currentFeature: 'raycast-ui',
    })
  }

  err(res, `not found: ${p}`, 404)
})

server.listen(PORT, () => {
  console.log(`API server running at http://localhost:${PORT}`)
  console.log(`Project root: ${PROJECT_ROOT}`)
  console.log(`Features dir: ${FEATURES_DIR}`)
})
