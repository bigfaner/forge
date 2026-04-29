import { useMutation, useQueryClient } from '@tanstack/react-query'
import { Zap } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'

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
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        gap: 6,
        padding: '7px 16px',
        borderRadius: 86,
        fontSize: 13,
        fontWeight: 600,
        letterSpacing: '0.3px',
        cursor: claim.isPending ? 'not-allowed' : 'pointer',
        border: '1px solid rgba(255,255,255,0.1)',
        background: 'rgba(255,255,255,0.07)',
        color: '#f9f9f9',
        transition: 'opacity 0.15s',
        opacity: claim.isPending ? 0.5 : 1,
        boxShadow: 'rgba(255,255,255,0.05) 0px 1px 0px 0px inset',
        fontFamily: 'inherit',
      }}
      onMouseEnter={e => { if (!claim.isPending) e.currentTarget.style.opacity = '0.6' }}
      onMouseLeave={e => { if (!claim.isPending) e.currentTarget.style.opacity = '1' }}
    >
      <Zap size={13} />
      {claim.isPending ? 'Claiming…' : 'Claim Task'}
    </button>
  )
}
