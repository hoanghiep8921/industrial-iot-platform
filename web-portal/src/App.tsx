import { Routes, Route } from 'react-router-dom'
import AppLayout from './components/layout/AppLayout'
import DashboardPage from './pages/DashboardPage'
import DeviceListPage from './pages/DeviceListPage'
import DeviceDetailPage from './pages/DeviceDetailPage'
import AlarmPage from './pages/AlarmPage'
import ErrorBoundary from './components/common/ErrorBoundary'

export default function App() {
  return (
    <ErrorBoundary>
      <Routes>
        <Route element={<AppLayout />}>
          <Route index element={<ErrorBoundary><DashboardPage /></ErrorBoundary>} />
          <Route path="devices" element={<ErrorBoundary><DeviceListPage /></ErrorBoundary>} />
          <Route path="devices/:deviceId" element={<ErrorBoundary><DeviceDetailPage /></ErrorBoundary>} />
          <Route path="alarms" element={<ErrorBoundary><AlarmPage /></ErrorBoundary>} />
        </Route>
      </Routes>
    </ErrorBoundary>
  )
}
