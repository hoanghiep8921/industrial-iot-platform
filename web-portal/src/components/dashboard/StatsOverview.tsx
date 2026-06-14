import { Card, CardContent, Grid, Typography } from '@mui/material'
import type { DeviceStats } from '../../types/device'

const statCards = (stats: DeviceStats | null) => {
  const online = stats?.byStatus?.online ?? 0
  const offline = stats?.byStatus?.offline ?? 0
  const error = stats?.byStatus?.error ?? 0
  return [
    { label: 'Total Devices', value: stats?.total ?? 0, color: '#1976d2' },
    { label: 'Online', value: online, color: '#2e7d32' },
    { label: 'Offline', value: offline, color: '#757575' },
    { label: 'Errors', value: error, color: '#c62828' },
  ]
}

interface StatsOverviewProps {
  stats: DeviceStats | null
}

export default function StatsOverview({ stats }: StatsOverviewProps) {
  return (
    <Grid container spacing={2}>
      {statCards(stats).map((card) => (
        <Grid size={{ xs: 6, sm: 3 }} key={card.label}>
          <Card>
            <CardContent sx={{ textAlign: 'center', py: 3 }}>
              <Typography variant="h4" sx={{ fontWeight: 700, color: card.color }}>
                {card.value}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {card.label}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  )
}
