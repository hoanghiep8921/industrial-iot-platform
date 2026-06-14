import { useEffect } from 'react'
import { Box } from '@mui/material'
import StatsOverview from '../components/dashboard/StatsOverview'
import TelemetryChart from '../components/dashboard/TelemetryChart'
import DeviceStatusCards from '../components/dashboard/DeviceStatusCards'
import LoadingSpinner from '../components/common/LoadingSpinner'
import ErrorAlert from '../components/common/ErrorAlert'
import { useDeviceStore } from '../stores/useDeviceStore'
import { useTelemetryStore } from '../stores/useTelemetryStore'
import { usePolling } from '../hooks/usePolling'

export default function DashboardPage() {
  const deviceStats = useDeviceStore((s) => s.deviceStats)
  const deviceLoading = useDeviceStore((s) => s.loading)
  const deviceError = useDeviceStore((s) => s.error)
  const fetchDeviceStats = useDeviceStore((s) => s.fetchDeviceStats)
  const fetchDevices = useDeviceStore((s) => s.fetchDevices)

  const records = useTelemetryStore((s) => s.records)
  const latestValues = useTelemetryStore((s) => s.latestValues)
  const telemetryLoading = useTelemetryStore((s) => s.loading)
  const telemetryError = useTelemetryStore((s) => s.error)
  const fetchRecent = useTelemetryStore((s) => s.fetchRecent)
  const fetchLatest = useTelemetryStore((s) => s.fetchLatest)

  useEffect(() => {
    fetchDeviceStats()
    fetchDevices({ limit: 20 })
  }, [fetchDeviceStats, fetchDevices])

  usePolling(() => {
    fetchRecent(100)
    fetchLatest()
  }, 10_000)

  const loading = (deviceLoading || telemetryLoading) && deviceStats === null
  const error = deviceError ?? telemetryError

  if (loading) return <LoadingSpinner message="Loading dashboard..." />
  if (error) return <ErrorAlert message={error} onRetry={() => { fetchDeviceStats(); fetchRecent(); }} />

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
      <StatsOverview stats={deviceStats} />
      <TelemetryChart records={records} />
      <DeviceStatusCards latestValues={latestValues} />
    </Box>
  )
}
