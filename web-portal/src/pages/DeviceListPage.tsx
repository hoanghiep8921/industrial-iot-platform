import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Box, Typography } from '@mui/material'
import { DataGrid, type GridColDef, type GridRowParams } from '@mui/x-data-grid'
import DeviceFilters from '../components/devices/DeviceFilters'
import DeviceStatusBadge from '../components/devices/DeviceStatusBadge'
import DeviceFormDialog from '../components/devices/DeviceFormDialog'
import ErrorAlert from '../components/common/ErrorAlert'
import { useDeviceStore } from '../stores/useDeviceStore'
import type { Device, DeviceStatus, CreateDeviceRequest, UpdateDeviceRequest } from '../types/device'
import { formatDate } from '../utils/formatters'

const columns: GridColDef<Device>[] = [
  { field: 'name', headerName: 'Name', flex: 2, minWidth: 150 },
  { field: 'serialNumber', headerName: 'Serial', flex: 1.5, minWidth: 120 },
  { field: 'factoryId', headerName: 'Factory', width: 80 },
  { field: 'area', headerName: 'Area', width: 100, valueGetter: (_, row) => row.area ?? '—' },
  {
    field: 'status', headerName: 'Status', width: 130,
    renderCell: (params) => <DeviceStatusBadge status={params.value as DeviceStatus} />,
  },
  {
    field: 'lastSeenAt', headerName: 'Last Seen', width: 160,
    valueGetter: (_, row) => formatDate(row.lastSeenAt),
  },
]

export default function DeviceListPage() {
  const navigate = useNavigate()
  const devices = useDeviceStore((s) => s.devices)
  const totalCount = useDeviceStore((s) => s.totalCount)
  const filters = useDeviceStore((s) => s.filters)
  const loading = useDeviceStore((s) => s.loading)
  const error = useDeviceStore((s) => s.error)
  const fetchDevices = useDeviceStore((s) => s.fetchDevices)
  const createDevice = useDeviceStore((s) => s.createDevice)
  const updateDevice = useDeviceStore((s) => s.updateDevice)
  const clearError = useDeviceStore((s) => s.clearError)

  const [search, setSearch] = useState('')
  const [statusFilter, setStatusFilter] = useState<DeviceStatus | ''>('')
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingDevice, setEditingDevice] = useState<Device | null>(null)

  useEffect(() => {
    fetchDevices({ page: 1, limit: 30 })
  }, [fetchDevices])

  useEffect(() => {
    const timeout = setTimeout(() => {
      fetchDevices({ search: search || undefined, status: statusFilter || undefined, page: 1 })
    }, 300)
    return () => clearTimeout(timeout)
  }, [search, statusFilter, fetchDevices])

  const handleRowClick = (params: GridRowParams<Device>) => {
    navigate(`/devices/${params.row.id}`)
  }

  const handleSubmit = (data: CreateDeviceRequest | UpdateDeviceRequest, id?: string) => {
    if (id) {
      updateDevice(id, data as UpdateDeviceRequest)
    } else {
      createDevice(data as CreateDeviceRequest)
    }
  }

  if (error && devices.length === 0) {
    return <ErrorAlert message={error} onRetry={() => { clearError(); fetchDevices(); }} />
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
      <Typography variant="h4">Devices</Typography>
      <DeviceFilters
        search={search}
        status={statusFilter}
        onSearchChange={setSearch}
        onStatusChange={setStatusFilter}
        onAddClick={() => { setEditingDevice(null); setDialogOpen(true); }}
      />
      <DataGrid
        rows={devices}
        columns={columns}
        loading={loading}
        onRowClick={handleRowClick}
        pageSizeOptions={[20, 30, 50]}
        rowCount={totalCount}
        paginationMode="server"
        paginationModel={{ page: (filters.page ?? 1) - 1, pageSize: filters.limit ?? 20 }}
        onPaginationModelChange={(m) => fetchDevices({ page: m.page + 1, limit: m.pageSize })}
        autoHeight
        sx={{ bgcolor: 'background.paper', '& .MuiDataGrid-row': { cursor: 'pointer' } }}
      />
      <DeviceFormDialog
        open={dialogOpen}
        device={editingDevice}
        onClose={() => setDialogOpen(false)}
        onSubmit={handleSubmit}
      />
    </Box>
  )
}
