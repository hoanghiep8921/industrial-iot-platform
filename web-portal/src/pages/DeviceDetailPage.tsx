import { useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { Box } from '@mui/material'
import DeviceInfoCard from '../components/device-detail/DeviceInfoCard'
import TelemetryHistoryChart from '../components/device-detail/TelemetryHistoryChart'
import ShadowStateViewer from '../components/device-detail/ShadowStateViewer'
import FirmwareInfoCard from '../components/device-detail/FirmwareInfoCard'
import LoadingSpinner from '../components/common/LoadingSpinner'
import ErrorAlert from '../components/common/ErrorAlert'
import { useDeviceStore } from '../stores/useDeviceStore'
import { fetchDeviceTelemetry } from '../api/telemetryApi'
import { useState } from 'react'
import type { TelemetryRecord } from '../types/telemetry'

export default function DeviceDetailPage() {
  const { deviceId } = useParams<{ deviceId: string }>()
  const selectedDevice = useDeviceStore((s) => s.selectedDevice)
  const shadow = useDeviceStore((s) => s.shadow)
  const firmwares = useDeviceStore((s) => s.firmwares)
  const firmwareUpdates = useDeviceStore((s) => s.firmwareUpdates)
  const loading = useDeviceStore((s) => s.loading)
  const error = useDeviceStore((s) => s.error)
  const fetchDeviceById = useDeviceStore((s) => s.fetchDeviceById)
  const fetchShadow = useDeviceStore((s) => s.fetchShadow)
  const fetchFirmwares = useDeviceStore((s) => s.fetchFirmwares)
  const fetchFirmwareUpdates = useDeviceStore((s) => s.fetchFirmwareUpdates)

  const [telemetryRecords, setTelemetryRecords] = useState<TelemetryRecord[]>([])

  useEffect(() => {
    if (deviceId) {
      fetchDeviceById(deviceId)
      fetchShadow(deviceId)
      fetchFirmwares()
      fetchFirmwareUpdates(deviceId)
    }
  }, [deviceId, fetchDeviceById, fetchShadow, fetchFirmwares, fetchFirmwareUpdates])

  const handleFetchTelemetry = async (from: string, to: string) => {
    if (!deviceId) return
    try {
      const data = await fetchDeviceTelemetry(deviceId, from, to)
      setTelemetryRecords(data.records)
    } catch { /* ignore */ }
  }

  if (loading && !selectedDevice) return <LoadingSpinner message="Loading device..." />
  if (error && !selectedDevice) return <ErrorAlert message={error} onRetry={() => deviceId && fetchDeviceById(deviceId)} />
  if (!selectedDevice) return <ErrorAlert message="Device not found" />

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
      <DeviceInfoCard device={selectedDevice} />
      <TelemetryHistoryChart
        deviceId={deviceId!}
        records={telemetryRecords}
        onFetch={handleFetchTelemetry}
      />
      <ShadowStateViewer shadow={shadow} />
      <FirmwareInfoCard
        firmwareVersion={selectedDevice.firmwareVersion}
        firmwares={firmwares}
        updates={firmwareUpdates}
      />
    </Box>
  )
}
