import { request } from './client'
import { snakeToCamel, camelToSnake } from '../utils/transformKeys'
import type { Device, DeviceFilter, DeviceStats, DeviceListResponse, CreateDeviceRequest, UpdateDeviceRequest, UpdateStatusRequest } from '../types/device'
import type { DeviceShadow } from '../types/shadow'
import type { Firmware, FirmwareUpdate, CreateUpdateRequest } from '../types/firmware'

// ========== Device Registry ==========

export async function fetchDevices(filter: Partial<DeviceFilter>): Promise<DeviceListResponse> {
  const params = new URLSearchParams()
  if (filter.factoryId) params.set('factory_id', filter.factoryId)
  if (filter.area) params.set('area', filter.area)
  if (filter.status) params.set('status', filter.status)
  if (filter.protocol) params.set('protocol', filter.protocol)
  if (filter.search) params.set('search', filter.search)
  params.set('page', String(filter.page ?? 1))
  params.set('limit', String(filter.limit ?? 20))
  const data = await request<{ count: number; page: number; limit: number; devices: unknown[] }>(
    `/api/v1/devices?${params}`,
  )
  return { ...data, devices: data.devices.map(snakeToCamel as (obj: unknown) => Device) }
}

export async function fetchDevice(id: string): Promise<Device> {
  return (snakeToCamel as (obj: unknown) => Device)(await request(`/api/v1/devices/${id}`))
}

export async function createDevice(data: CreateDeviceRequest): Promise<Device> {
  return (snakeToCamel as (obj: unknown) => Device)(
    await request('/api/v1/devices', {
      method: 'POST',
      body: JSON.stringify(camelToSnake(data)),
    }),
  )
}

export async function updateDevice(id: string, data: UpdateDeviceRequest): Promise<Device> {
  return (snakeToCamel as (obj: unknown) => Device)(
    await request(`/api/v1/devices/${id}`, {
      method: 'PUT',
      body: JSON.stringify(camelToSnake(data)),
    }),
  )
}

export async function deleteDevice(id: string): Promise<void> {
  return request(`/api/v1/devices/${id}`, { method: 'DELETE' })
}

export async function updateDeviceStatus(id: string, data: UpdateStatusRequest): Promise<Device> {
  return (snakeToCamel as (obj: unknown) => Device)(
    await request(`/api/v1/devices/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify(camelToSnake(data)),
    }),
  )
}

export async function fetchDeviceStats(): Promise<DeviceStats> {
  return (snakeToCamel as (obj: unknown) => DeviceStats)(await request('/api/v1/stats'))
}

// ========== Device Shadow ==========

export async function fetchShadow(deviceId: string): Promise<DeviceShadow> {
  return (snakeToCamel as (obj: unknown) => DeviceShadow)(await request(`/api/v1/devices/${deviceId}/shadow`))
}

export async function updateShadowDesired(deviceId: string, state: Record<string, unknown>): Promise<DeviceShadow> {
  return (snakeToCamel as (obj: unknown) => DeviceShadow)(
    await request(`/api/v1/devices/${deviceId}/shadow/desired`, {
      method: 'PUT',
      body: JSON.stringify(state),
    }),
  )
}

// ========== Firmware ==========

export async function fetchFirmwares(status?: string, model?: string): Promise<{ count: number; firmwares: Firmware[] }> {
  const params = new URLSearchParams()
  if (status) params.set('status', status)
  if (model) params.set('model', model)
  const data = await request<{ count: number; firmwares: unknown[] }>(
    `/api/v1/firmwares${params.toString() ? '?' + params : ''}`,
  )
  return { count: data.count, firmwares: data.firmwares.map(snakeToCamel as (obj: unknown) => Firmware) }
}

// ========== Firmware Updates ==========

export async function fetchFirmwareUpdates(
  deviceId?: string,
  status?: string,
  campaignName?: string,
): Promise<{ count: number; updates: FirmwareUpdate[] }> {
  const params = new URLSearchParams()
  if (deviceId) params.set('device_id', deviceId)
  if (status) params.set('status', status)
  if (campaignName) params.set('campaign_name', campaignName)
  const data = await request<{ count: number; updates: unknown[] }>(
    `/api/v1/firmware-updates?${params}`,
  )
  return { count: data.count, updates: data.updates.map(snakeToCamel as (obj: unknown) => FirmwareUpdate) }
}

export async function createFirmwareUpdate(data: CreateUpdateRequest): Promise<unknown> {
  return snakeToCamel(
    await request(
      '/api/v1/firmware-updates',
      { method: 'POST', body: JSON.stringify(camelToSnake(data)) },
    ),
  )
}
