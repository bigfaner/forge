import { useMemo } from 'react'
import { ReactFlow, Background, Controls, type Node, type Edge, Position } from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import type { Task } from '../../lib/types'
import { STATUS_LABEL } from '../../lib/utils'

const STATUS_BG: Record<string, string> = {
  pending:     '#1e293b',
  in_progress: '#1e3a5f',
  completed:   '#14532d',
  blocked:     '#7f1d1d',
  skipped:     '#713f12',
}
const STATUS_BORDER: Record<string, string> = {
  pending:     '#475569',
  in_progress: '#3b82f6',
  completed:   '#22c55e',
  blocked:     '#ef4444',
  skipped:     '#eab308',
}

interface Props {
  tasks: Task[]
  featureSlug: string
}

function buildLayout(tasks: Task[]): { nodes: Node[]; edges: Edge[] } {
  // Group tasks by phase for hierarchical layout
  const phases = [...new Set(tasks.map(t => t.phase))].sort((a, b) => a - b)
  const phaseMap: Record<number, Task[]> = {}
  tasks.forEach(t => {
    if (!phaseMap[t.phase]) phaseMap[t.phase] = []
    phaseMap[t.phase].push(t)
  })

  const NODE_W = 180
  const NODE_H = 60
  const H_GAP = 60
  const V_GAP = 100

  const nodes: Node[] = []
  const edges: Edge[] = []

  phases.forEach((phase, pi) => {
    const phaseTasks = phaseMap[phase]
    phaseTasks.forEach((task, ti) => {
      const totalInPhase = phaseTasks.length
      const x = pi * (NODE_W + H_GAP)
      const y = (ti - (totalInPhase - 1) / 2) * (NODE_H + V_GAP)

      nodes.push({
        id: task.id,
        position: { x, y },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
        data: { label: `${task.id}: ${task.title}`, status: task.status },
        style: {
          background: STATUS_BG[task.status] ?? '#1e293b',
          border: `1px solid ${STATUS_BORDER[task.status] ?? '#475569'}`,
          borderRadius: 6,
          color: '#f8fafc',
          fontSize: 11,
          width: NODE_W,
          padding: '6px 10px',
        },
      })

      task.dependencies.forEach(dep => {
        const depId = dep.endsWith('x')
          ? tasks.find(t2 => t2.id.startsWith(dep.replace('.x', '.')))?.id
          : dep
        if (depId) {
          edges.push({
            id: `${depId}->${task.id}`,
            source: depId,
            target: task.id,
            style: { stroke: '#475569', strokeWidth: 1 },
            animated: task.status === 'in_progress',
          })
        }
      })
    })
  })

  return { nodes, edges }
}

export function DagView({ tasks }: Props) {
  const { nodes, edges } = useMemo(() => buildLayout(tasks), [tasks])

  if (tasks.length === 0) {
    return <p className="text-sm text-muted-foreground text-center py-12">No tasks</p>
  }

  return (
    <div className="h-[500px] rounded-lg border border-border overflow-hidden bg-background">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
        proOptions={{ hideAttribution: true }}
      >
        <Background color="#334155" gap={20} size={1} />
        <Controls showInteractive={false} />
      </ReactFlow>
      {/* Legend */}
      <div className="absolute bottom-4 right-4 flex gap-2 flex-wrap">
        {Object.entries(STATUS_LABEL).map(([s, label]) => (
          <div key={s} className="flex items-center gap-1 text-xs">
            <div className="h-2 w-2 rounded-sm" style={{ background: STATUS_BORDER[s] }} />
            <span className="text-muted-foreground">{label}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
