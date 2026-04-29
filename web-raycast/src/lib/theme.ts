type Theme = 'light' | 'dark' | 'system'

function prefersDark(): boolean {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

export function applyTheme(theme: Theme): void {
  const root = document.documentElement
  if (theme === 'dark' || (theme === 'system' && prefersDark())) {
    root.classList.add('dark')
  } else {
    root.classList.remove('dark')
  }
}

export function getStoredTheme(): Theme {
  return (localStorage.getItem('theme') as Theme) ?? 'dark'
}

export function setTheme(theme: Theme): void {
  localStorage.setItem('theme', theme)
  applyTheme(theme)
}

export function initTheme(): void {
  applyTheme(getStoredTheme())
}
