import { Outlet, NavLink } from 'react-router-dom'
import { LayoutDashboard, Layers, History, BookOpen, Settings, Zap } from 'lucide-react'

const NAV = [
  { to: '/',         label: 'Dashboard', icon: LayoutDashboard, end: true },
  { to: '/features', label: 'Features',  icon: Layers },
  { to: '/records',  label: 'Records',   icon: History },
  { to: '/lessons',  label: 'Lessons',   icon: BookOpen },
  { to: '/settings', label: 'Settings',  icon: Settings },
]

export function AppLayout() {
  return (
    <div style={{ display: 'flex', height: '100vh', overflow: 'hidden', background: '#07080a', color: '#f9f9f9' }}>
      {/* Sidebar */}
      <aside style={{ display: 'flex', flexDirection: 'column', width: 208, flexShrink: 0, background: '#07080a', borderRight: '1px solid rgba(255,255,255,0.06)' }}>
        {/* Logo */}
        <div style={{ display: 'flex', height: 56, alignItems: 'center', gap: 10, padding: '0 16px', borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
          <div style={{ display: 'flex', height: 24, width: 24, alignItems: 'center', justifyContent: 'center', borderRadius: 6, background: 'linear-gradient(135deg, #FF6363 0%, #ff4040 100%)', boxShadow: 'rgba(255,99,99,0.3) 0px 2px 8px' }}>
            <Zap size={14} color="white" />
          </div>
          <span style={{ fontSize: 14, fontWeight: 600, letterSpacing: '0.2px', color: '#f9f9f9' }}>
            ZCode
          </span>
        </div>

        {/* Nav */}
        <nav style={{ flex: 1, padding: '10px 8px', display: 'flex', flexDirection: 'column', gap: 2 }}>
          {NAV.map(({ to, label, icon: Icon, end }) => (
            <NavLink
              key={to}
              to={to}
              end={end}
              style={({ isActive }) => ({
                display: 'flex',
                alignItems: 'center',
                gap: 10,
                padding: '7px 10px',
                borderRadius: 8,
                fontSize: 14,
                fontWeight: 500,
                letterSpacing: '0.2px',
                textDecoration: 'none',
                cursor: 'pointer',
                transition: 'all 0.15s',
                color: isActive ? '#f9f9f9' : '#6a6b6c',
                background: isActive ? 'rgba(255,255,255,0.06)' : 'transparent',
              })}
            >
              <Icon size={15} style={{ flexShrink: 0 }} />
              {label}
            </NavLink>
          ))}
        </nav>

        {/* Footer */}
        <div style={{ padding: '10px 16px', borderTop: '1px solid rgba(255,255,255,0.06)' }}>
          <p style={{ fontSize: 11, fontFamily: 'monospace', color: '#434345', letterSpacing: '0.3px', margin: 0 }}>
            v0.1.0
          </p>
        </div>
      </aside>

      {/* Main */}
      <main style={{ flex: 1, overflowY: 'auto', background: '#07080a' }}>
        <Outlet />
      </main>
    </div>
  )
}
