import { Card, CardContent, CardHeader, Typography, Chip, Table, TableHead, TableRow, TableCell, TableBody } from '@mui/material'
import type { Firmware, FirmwareUpdate } from '../../types/firmware'
import EmptyState from '../common/EmptyState'

interface FirmwareInfoCardProps {
  firmwareVersion: string | undefined
  firmwares: Firmware[]
  updates: FirmwareUpdate[]
}

const updateStatusColor: Record<string, 'default' | 'primary' | 'success' | 'error' | 'warning'> = {
  pending: 'default',
  downloading: 'primary',
  installing: 'warning',
  success: 'success',
  failed: 'error',
  rolled_back: 'warning',
}

export default function FirmwareInfoCard({ firmwareVersion, firmwares, updates }: FirmwareInfoCardProps) {
  return (
    <Card>
      <CardHeader
        title="Firmware"
        subheader={`Current version: ${firmwareVersion || 'Unknown'}`}
      />
      <CardContent>
        <Typography variant="subtitle2" sx={{ mb: 1 }}>Available Firmwares</Typography>
        {firmwares.length === 0 ? (
          <EmptyState message="No firmwares registered" />
        ) : (
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Version</TableCell>
                <TableCell>Model</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Size</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {firmwares.slice(0, 5).map((fw) => (
                <TableRow key={fw.id}>
                  <TableCell>{fw.version}</TableCell>
                  <TableCell>{fw.model ?? 'All'}</TableCell>
                  <TableCell><Chip size="small" label={fw.status} /></TableCell>
                  <TableCell>{fw.fileSize ? `${(fw.fileSize / 1024).toFixed(1)} KB` : '—'}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}

        <Typography variant="subtitle2" sx={{ mt: 2, mb: 1 }}>Update History</Typography>
        {updates.length === 0 ? (
          <EmptyState message="No update history" />
        ) : (
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>To Version</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Retries</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {updates.slice(0, 5).map((u) => (
                <TableRow key={u.id}>
                  <TableCell>{u.toVersion}</TableCell>
                  <TableCell>
                    <Chip size="small" color={updateStatusColor[u.status] ?? 'default'} label={u.status} />
                  </TableCell>
                  <TableCell>{u.retryCount}/{u.maxRetries}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  )
}
