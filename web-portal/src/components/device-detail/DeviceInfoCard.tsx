import { Card, CardContent, CardHeader, Typography, Grid } from '@mui/material'
import type { Device } from '../../types/device'
import { formatDate } from '../../utils/formatters'
import DeviceStatusBadge from '../devices/DeviceStatusBadge'

interface DeviceInfoCardProps {
  device: Device
}

export default function DeviceInfoCard({ device }: DeviceInfoCardProps) {
  const fields: [string, string | undefined][] = [
    ['Name', device.name],
    ['Model', device.model],
    ['Vendor', device.vendor],
    ['Serial Number', device.serialNumber],
    ['Protocol', device.protocol],
    ['Factory', device.factoryId],
    ['Area', device.area],
    ['Line', device.line],
    ['Firmware Version', device.firmwareVersion],
    ['Capabilities', device.capabilities?.join(', ')],
    ['Created', formatDate(device.createdAt)],
    ['Last Seen', formatDate(device.lastSeenAt)],
  ]

  return (
    <Card>
      <CardHeader
        title={device.name}
        action={<DeviceStatusBadge status={device.status} />}
      />
      <CardContent>
        <Grid container spacing={1}>
          {fields.filter(([, v]) => v).map(([label, value]) => (
            <Grid size={{ xs: 6 }} key={label}>
              <Typography variant="caption" color="text.secondary">{label}</Typography>
              <Typography variant="body2">{value}</Typography>
            </Grid>
          ))}
        </Grid>
      </CardContent>
    </Card>
  )
}
