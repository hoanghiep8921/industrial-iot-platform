import { Chip } from '@mui/material'
import type { DeviceStatus } from '../../types/device'

const statusConfig: Record<DeviceStatus, { color: 'success' | 'default' | 'warning' | 'error'; label: string }> = {
  online: { color: 'success', label: 'Online' },
  offline: { color: 'default', label: 'Offline' },
  maintenance: { color: 'warning', label: 'Maintenance' },
  error: { color: 'error', label: 'Error' },
}

interface DeviceStatusBadgeProps {
  status: DeviceStatus
}

export default function DeviceStatusBadge({ status }: DeviceStatusBadgeProps) {
  const config = statusConfig[status] ?? { color: 'default' as const, label: status }
  return <Chip size="small" color={config.color} label={config.label} />
}
