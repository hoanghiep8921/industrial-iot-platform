export type DeviceStatus = 'online' | 'offline' | 'maintenance' | 'error'

export interface Device {
  id: string
  serialNumber: string
  name: string
  model?: string
  vendor?: string
  factoryId: string
  area?: string
  line?: string
  protocol: string
  capabilities: string[]
  status: DeviceStatus
  firmwareVersion?: string
  metadata?: Record<string, unknown>
  lastSeenAt?: string | null
  createdAt: string
  updatedAt: string
}

export interface DeviceFilter {
  factoryId?: string
  area?: string
  status?: DeviceStatus
  protocol?: string
  search?: string
  page: number
  limit: number
}

export interface DeviceStats {
  total: number
  byFactory: Record<string, number>
  byStatus: Record<string, number>
  byModel: Record<string, number>
}

export interface DeviceListResponse {
  count: number
  page: number
  limit: number
  devices: Device[]
}

export interface CreateDeviceRequest {
  serialNumber: string
  name: string
  factoryId: string
  model?: string
  vendor?: string
  area?: string
  line?: string
  protocol?: string
  capabilities?: string[]
  firmwareVersion?: string
  metadata?: Record<string, unknown>
}

export interface UpdateDeviceRequest {
  name?: string
  model?: string
  vendor?: string
  factoryId?: string
  area?: string
  line?: string
  protocol?: string
  capabilities?: string[]
  firmwareVersion?: string
  metadata?: Record<string, unknown>
}

export interface UpdateStatusRequest {
  status: DeviceStatus
}
