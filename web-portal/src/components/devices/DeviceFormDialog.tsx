import { useState, useEffect } from 'react'
import {
  Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button,
  FormControl, InputLabel, Select, MenuItem,
} from '@mui/material'
import type { Device, CreateDeviceRequest, UpdateDeviceRequest } from '../../types/device'

interface DeviceFormDialogProps {
  open: boolean
  device?: Device | null
  onClose: () => void
  onSubmit: (data: CreateDeviceRequest | UpdateDeviceRequest, id?: string) => void
}

const defaultForm = {
  serialNumber: '', name: '', factoryId: '', model: '', vendor: '',
  area: '', line: '', protocol: 'mqtt',
}

export default function DeviceFormDialog({ open, device, onClose, onSubmit }: DeviceFormDialogProps) {
  const [form, setForm] = useState(defaultForm)

  useEffect(() => {
    if (device) {
      setForm({
        serialNumber: device.serialNumber,
        name: device.name,
        factoryId: device.factoryId,
        model: device.model ?? '',
        vendor: device.vendor ?? '',
        area: device.area ?? '',
        line: device.line ?? '',
        protocol: device.protocol,
      })
    } else {
      setForm(defaultForm)
    }
  }, [device, open])

  const handleSubmit = () => {
    const data: CreateDeviceRequest = {
      serialNumber: form.serialNumber,
      name: form.name,
      factoryId: form.factoryId,
      model: form.model || undefined,
      vendor: form.vendor || undefined,
      area: form.area || undefined,
      line: form.line || undefined,
      protocol: form.protocol,
    }
    onSubmit(data, device?.id)
    onClose()
  }

  const set = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm({ ...form, [field]: e.target.value })

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{device ? 'Edit Device' : 'Add Device'}</DialogTitle>
      <DialogContent>
        <TextField label="Serial Number *" value={form.serialNumber} onChange={set('serialNumber')} fullWidth margin="dense" disabled={!!device} />
        <TextField label="Name *" value={form.name} onChange={set('name')} fullWidth margin="dense" />
        <TextField label="Factory ID *" value={form.factoryId} onChange={set('factoryId')} fullWidth margin="dense" />
        <TextField label="Model" value={form.model} onChange={set('model')} fullWidth margin="dense" />
        <TextField label="Vendor" value={form.vendor} onChange={set('vendor')} fullWidth margin="dense" />
        <TextField label="Area" value={form.area} onChange={set('area')} fullWidth margin="dense" />
        <TextField label="Line" value={form.line} onChange={set('line')} fullWidth margin="dense" />
        <FormControl fullWidth margin="dense">
          <InputLabel>Protocol</InputLabel>
          <Select value={form.protocol} label="Protocol" onChange={(e) => setForm({ ...form, protocol: e.target.value })}>
            <MenuItem value="mqtt">MQTT</MenuItem>
            <MenuItem value="modbus">Modbus</MenuItem>
            <MenuItem value="opcua">OPC-UA</MenuItem>
          </Select>
        </FormControl>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button variant="contained" onClick={handleSubmit} disabled={!form.serialNumber || !form.name || !form.factoryId}>
          {device ? 'Update' : 'Create'}
        </Button>
      </DialogActions>
    </Dialog>
  )
}
