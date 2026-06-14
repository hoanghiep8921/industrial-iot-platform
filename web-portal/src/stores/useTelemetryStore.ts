import { create } from 'zustand'
import { fetchRecentTelemetry, fetchLatestValues, fetchTelemetryStats } from '../api/telemetryApi'
import type { TelemetryRecord, TelemetryStats } from '../types/telemetry'

interface TelemetryState {
  records: TelemetryRecord[]
  latestValues: TelemetryRecord[]
  telemetryStats: TelemetryStats | null
  loading: boolean
  error: string | null

  fetchRecent: (limit?: number) => Promise<void>
  fetchLatest: (factoryId?: string) => Promise<void>
  fetchStats: () => Promise<void>
  clearError: () => void
}

export const useTelemetryStore = create<TelemetryState>((set) => ({
  records: [],
  latestValues: [],
  telemetryStats: null,
  loading: false,
  error: null,

  fetchRecent: async (limit = 100) => {
    set({ loading: true, error: null })
    try {
      const data = await fetchRecentTelemetry(limit)
      set({ records: data.records, loading: false })
    } catch (err) {
      set({ error: (err as Error).message, loading: false })
    }
  },

  fetchLatest: async (factoryId?) => {
    set({ loading: true, error: null })
    try {
      const data = await fetchLatestValues(factoryId)
      set({ latestValues: data.records, loading: false })
    } catch (err) {
      set({ error: (err as Error).message, loading: false })
    }
  },

  fetchStats: async () => {
    try {
      set({ telemetryStats: await fetchTelemetryStats() })
    } catch {
      // Non-critical
    }
  },

  clearError: () => set({ error: null }),
}))
