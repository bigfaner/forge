import { Sun, Moon } from 'lucide-react'
import { useState, useEffect } from 'react'
import { getStoredTheme, setTheme } from '../lib/theme'
import { cn } from '../lib/utils'

export function ThemeToggle() {
  const [theme, setLocalTheme] = useState<'light' | 'dark'>(() => {
    const t = getStoredTheme()
    return t === 'light' ? 'light' : 'dark'
  })

  useEffect(() => {
    setTheme(theme)
  }, [theme])

  const toggle = () => setLocalTheme(t => t === 'dark' ? 'light' : 'dark')

  return (
    <button
      onClick={toggle}
      className={cn(
        'flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm',
        'text-muted-foreground hover:bg-muted hover:text-foreground transition-colors cursor-pointer',
      )}
      aria-label="Toggle theme"
    >
      {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
      {theme === 'dark' ? 'Light mode' : 'Dark mode'}
    </button>
  )
}
