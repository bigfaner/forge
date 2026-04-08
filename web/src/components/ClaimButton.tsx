import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Zap } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { cn } from '../lib/utils'

interface Props {
  featureSlug: string
}

export function ClaimButton({ featureSlug }: Props) {
  const qc = useQueryClient()
  const navigate = useNavigate()

  const claim = useMutation({
    mutationFn: () => api.tasks.claim(),
    onSuccess: (result) => {
      qc.invalidateQueries({ queryKey: ['features', featureSlug, 'tasks'] })
      navigate(`/features/${featureSlug}/tasks/${result.taskId}`)
    },
  })

  return (
    <button
      onClick={() => claim.mutate()}
      disabled={claim.isPending}
      className={cn(
        'flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium',
        'bg-accent text-accent-foreground hover:bg-accent/90 transition-colors cursor-pointer',
        'disabled:opacity-50 disabled:cursor-not-allowed',
      )}
    >
      <Zap className="h-3.5 w-3.5" />
      {claim.isPending ? 'Claiming…' : 'Claim Task'}
    </button>
  )
}
