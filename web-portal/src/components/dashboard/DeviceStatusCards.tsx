import { Card, CardContent, CardHeader, Grid, Typography } from '@mui/material'
import type { TelemetryRecord } from '../../types/telemetry'
import EmptyState from '../common/EmptyState'

interface DeviceStatusCardsProps {
  latestValues: TelemetryRecord[]
}

export default function DeviceStatusCards({ latestValues }: DeviceStatusCardsProps) {
  if (latestValues.length === 0) {
    return (
      <Card>
        <CardHeader title="Device Status" />
        <CardContent><EmptyState message="No device data available" /></CardContent>
      </Card>
    )
  }

  // Group latest values by device
  const deviceMap = new Map<string, TelemetryRecord[]>()
  for (const r of latestValues) {
    const existing = deviceMap.get(r.deviceId) ?? []
    existing.push(r)
    deviceMap.set(r.deviceId, existing)
  }

  return (
    <Card>
      <CardHeader title="Device Status" />
      <CardContent>
        <Grid container spacing={2}>
          {[...deviceMap.entries()].map(([deviceId, records]) => (
            <Grid size={{ xs: 12, sm: 6, md: 4 }} key={deviceId}>
              <Card variant="outlined">
                <CardContent>
                  <Typography variant="subtitle2" noWrap>{deviceId.slice(0, 12)}</Typography>
                  {records.slice(0, 3).map((r) => (
                    <Typography key={r.metricName} variant="body2" color="text.secondary">
                      {r.metricName}: {typeof r.value === 'number' ? r.value.toFixed(1) : String(r.value ?? '—')} {r.unit ?? ''}
                    </Typography>
                  ))}
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </CardContent>
    </Card>
  )
}
