import { create } from 'zustand'
import {
  fetchDevices, fetchDevice, createDevice, updateDevice, deleteDevice,
  updateDeviceStatus, fetchDeviceStats, fetchShadow, updateShadowDesired,
  fetchFirmwares, fetchFirmwareUpdates, createFirmwareUpdate,
} from '../api/deviceApi'
import type { Device, DeviceFilter, DeviceStats, CreateDeviceRequest, UpdateDeviceRequest, UpdateStatusRequest } from '../types/device'
import type { DeviceShadow } from '../types/shadow'
import type { Firmware, FirmwareUpdate, CreateUpdateRequest } from '../types/firmware'

interface DeviceState {
  devices: Device[]
  totalCount: number
  filters: DeviceFilter
  deviceStats: DeviceStats | null
  selectedDevice: Device | null
  shadow: DeviceShadow | null
  firmwares: Firmware[]
  firmwareUpdates: FirmwareUpdate[]
  loading: boolean
  error: string | null

  fetchDevices: (filter?: Partial<DeviceFilter>) => Promise<void>
  setFilters: (partial: Partial<DeviceFilter>) => void
  fetchDeviceStats: () => Promise<void>
  fetchDeviceById: (id: string) => Promise<void>
  createDevice: (data: CreateDeviceRequest) => Promise<Device | null>
  updateDevice: (id: string, data: UpdateDeviceRequest) => Promise<void>
  deleteDevice: (id: string) => Promise<boolean>
  updateStatus: (id: string, data: UpdateStatusRequest) => Promise<void>
  fetchShadow: (deviceId: string) => Promise<void>
  updateDesired: (deviceId: string, state: Record<string, unknown>) => Promise<void>
  fetchFirmwares: (status?: string) => Promise<void>
  fetchFirmwareUpdates: (deviceId?: string) => Promise<void>
  createFirmwareUpdate: (data: CreateUpdateRequest) => Promise<void>
  clearError: () => void
}

export const useDeviceStore = create<DeviceState>((set, get) => ({
  devices: [],
  totalCount: 0,
  filters: { page: 1, limit: 20 },
  deviceStats: null,
  selectedDevice: null,
  shadow: null,
  firmwares: [],
  firmwareUpdates: [],
  loading: false,
  error: null,

  fetchDevices: async (filter?) => {
    const merged = { ...get().filters, ...filter }
    set({ filters: merged, loading: true, error: null })
    try {
      const data = await fetchDevices(merged)
      set({ devices: data.devices, totalCount: data.count, loading: false })
    } catch (err) {
      set({ error: (err as Error).message, loading: false })
    }
  },

  setFilters: (partial) => {
    const merged = { ...get().filters, ...partial, page: partial.page ?? 1 }
    set({ filters: merged })
    get().fetchDevices(merged)
  },

  fetchDeviceStats: async () => {
    try { set({ deviceStats: await fetchDeviceStats() }) } catch { /* non-critical */ }
  },

  fetchDeviceById: async (id) => {
    set({ loading: true, error: null })
    try {
      const d = await fetchDevice(id)
      set({ selectedDevice: d, loading: false })
    } catch (err) {
      set({ error: (err as Error).message, loading: false })
    }
  },

  createDevice: async (data) => {
    set({ loading: true, error: null })
    try {
      const d = await createDevice(data)
      await get().fetchDevices()
      return d
    } catch (err) {
      set({ error: (err as Error).message, loading: false })
      return null
    }
  },

  updateDevice: async (id, data) => {
    set({ loading: true, error: null })
    try {
      await updateDevice(id, data)
      await get().fetchDevices()
      set({ selectedDevice: await fetchDevice(id), loading: false })
    } catch (err) {
      set({ error: (err as Error).message, loading: false })
    }
  },

  deleteDevice: async (id) => {
    try {
      await deleteDevice(id)
      await get().fetchDevices()
      return true
    } catch (err) {
      set({ error: (err as Error).message })
      return false
    }
  },

  updateStatus: async (id, data) => {
    try {
      await updateDeviceStatus(id, data)
      await get().fetchDevices()
    } catch (err) {
      set({ error: (err as Error).message })
    }
  },

  fetchShadow: async (deviceId) => {
    try { set({ shadow: await fetchShadow(deviceId) }) } catch { /* non-critical */ }
  },

  updateDesired: async (deviceId, state) => {
    try { set({ shadow: await updateShadowDesired(deviceId, state) }) } catch (err) {
      set({ error: (err as Error).message })
    }
  },

  fetchFirmwares: async (status?) => {
    try { set({ firmwares: (await fetchFirmwares(status)).firmwares }) } catch { /* */ }
  },

  fetchFirmwareUpdates: async (deviceId?) => {
    try { set({ firmwareUpdates: (await fetchFirmwareUpdates(deviceId)).updates }) } catch { /* */ }
  },

  createFirmwareUpdate: async (data) => {
    try {
      await createFirmwareUpdate(data)
      await get().fetchFirmwareUpdates(data.deviceId)
    } catch (err) { set({ error: (err as Error).message }) }
  },

  clearError: () => set({ error: null }),
}))
