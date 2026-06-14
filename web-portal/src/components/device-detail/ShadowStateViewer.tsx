import { Card, CardContent, CardHeader, Typography, Box } from '@mui/material'
import type { DeviceShadow } from '../../types/shadow'
import EmptyState from '../common/EmptyState'

interface ShadowStateViewerProps {
  shadow: DeviceShadow | null
}

export default function ShadowStateViewer({ shadow }: ShadowStateViewerProps) {
  if (!shadow) {
    return (
      <Card>
        <CardHeader title="Device Shadow" />
        <CardContent><EmptyState message="No shadow data" /></CardContent>
      </Card>
    )
  }

  const hasDelta = Object.keys(shadow.delta).length > 0

  return (
    <Card>
      <CardHeader title="Device Shadow" subheader={`Version: ${shadow.version}`} />
      <CardContent>
        <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 2 }}>
          <Box>
            <Typography variant="subtitle2" color="primary.main">Desired</Typography>
            <Box component="pre" sx={{ bgcolor: 'grey.100', p: 1, borderRadius: 1, fontSize: 12, overflow: 'auto', maxHeight: 200 }}>
              {JSON.stringify(shadow.desiredState, null, 2) || '{}'}
            </Box>
          </Box>
          <Box>
            <Typography variant="subtitle2" color={hasDelta ? 'warning.main' : 'success.main'}>
              Reported {hasDelta && '⚠'}
            </Typography>
            <Box component="pre" sx={{ bgcolor: hasDelta ? '#fff3e0' : 'grey.100', p: 1, borderRadius: 1, fontSize: 12, overflow: 'auto', maxHeight: 200 }}>
              {JSON.stringify(shadow.reportedState, null, 2) || '{}'}
            </Box>
          </Box>
        </Box>
        {hasDelta && (
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle2" color="error.main">Delta (Desired ≠ Reported)</Typography>
            <Box component="pre" sx={{ bgcolor: '#ffebee', p: 1, borderRadius: 1, fontSize: 12, overflow: 'auto', maxHeight: 150 }}>
              {JSON.stringify(shadow.delta, null, 2)}
            </Box>
          </Box>
        )}
      </CardContent>
    </Card>
  )
}
