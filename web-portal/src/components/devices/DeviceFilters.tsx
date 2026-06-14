import { Box, TextField, Button, ToggleButtonGroup, ToggleButton } from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import type { DeviceStatus } from '../../types/device'

interface DeviceFiltersProps {
  search: string
  status: DeviceStatus | ''
  onSearchChange: (v: string) => void
  onStatusChange: (v: DeviceStatus | '') => void
  onAddClick: () => void
}

export default function DeviceFilters({ search, status, onSearchChange, onStatusChange, onAddClick }: DeviceFiltersProps) {
  return (
    <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap', alignItems: 'center' }}>
      <TextField
        size="small"
        placeholder="Search by name or serial..."
        value={search}
        onChange={(e) => onSearchChange(e.target.value)}
        sx={{ minWidth: 250 }}
      />
      <ToggleButtonGroup
        size="small"
        value={status}
        exclusive
        onChange={(_, v) => onStatusChange(v ?? '')}
      >
        <ToggleButton value="">All</ToggleButton>
        <ToggleButton value="online">Online</ToggleButton>
        <ToggleButton value="offline">Offline</ToggleButton>
        <ToggleButton value="maintenance">Maintenance</ToggleButton>
        <ToggleButton value="error">Error</ToggleButton>
      </ToggleButtonGroup>
      <Button variant="contained" startIcon={<AddIcon />} onClick={onAddClick}>
        Add Device
      </Button>
    </Box>
  )
}
