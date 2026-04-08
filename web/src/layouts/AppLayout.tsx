import { Outlet, NavLink } from 'react-router-dom'
import { LayoutDashboard, Layers, History, BookOpen, Settings, Zap } from 'lucide-react'
import { ThemeToggle } from '../components/ThemeToggle'
import { cn } from '../lib/utils'

const NAV = [
  { to: '/',        label: 'Dashboard', icon: LayoutDashboard, end: true },
  { to: '/features', label: 'Features',  icon: Layers },
  { to: '/records',  label: 'Records',   icon: History },
  { to: '/lessons',  label: 'Lessons',   icon: BookOpen },
  { to: '/settings', label: 'Settings',  icon: Settings },
]

export function AppLayout() {
  return (
    <div className="flex h-screen overflow-hidden bg-background text-foreground">
      {/* Sidebar */}
      <aside className="flex w-56 flex-col border-r border-border bg-card">
        {/* Logo */}
        <div className="flex h-14 items-center gap-2 px-4 border-b border-border">
          <Zap className="h-5 w-5 text-accent" />
          <span className="font-semibold tracking-tight">ZCode</span>
        </div>

        {/* Nav */}
        <nav className="flex-1 space-y-0.5 p-2 py-3">
          {NAV.map(({ to, label, icon: Icon, end }) => (
            <NavLink
              key={to}
              to={to}
              end={end}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-2.5 rounded-md px-3 py-2 text-sm font-medium transition-colors cursor-pointer',
                  isActive
                    ? 'bg-accent/10 text-accent'
                    : 'text-muted-foreground hover:bg-muted hover:text-foreground',
                )
              }
            >
              <Icon className="h-4 w-4 shrink-0" />
              {label}
            </NavLink>
          ))}
        </nav>

        {/* Footer */}
        <div className="border-t border-border p-3">
          <ThemeToggle />
        </div>
      </aside>

      {/* Main */}
      <main className="flex-1 overflow-y-auto">
        <Outlet />
      </main>
    </div>
  )
}
