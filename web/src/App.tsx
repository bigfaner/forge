import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AppLayout } from './layouts/AppLayout'
import { Dashboard } from './pages/Dashboard'
import { FeatureList } from './pages/FeatureList'
import { FeatureDetail } from './pages/FeatureDetail'
import { TaskDetail } from './pages/TaskDetail'
import { Records } from './pages/Records'
import { Lessons } from './pages/Lessons'
import { Settings } from './pages/Settings'

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route index element={<Dashboard />} />
          <Route path="features" element={<FeatureList />} />
          <Route path="features/:slug" element={<FeatureDetail />} />
          <Route path="features/:slug/tasks/:id" element={<TaskDetail />} />
          <Route path="records" element={<Records />} />
          <Route path="lessons" element={<Lessons />} />
          <Route path="settings" element={<Settings />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
