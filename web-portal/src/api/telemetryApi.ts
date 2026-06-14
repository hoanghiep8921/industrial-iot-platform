import { request } from './client'
import { pascalToCamel } from '../utils/transformKeys'
import type { TelemetryRecord, TelemetryResponse, LatestValuesResponse, TelemetryStats } from '../types/telemetry'

export async function fetchRecentTelemetry(limit = 100): Promise<TelemetryResponse> {
  const data = await request<{ count: number; records: unknown[] }>(
    `/api/v1/telemetry/recent?limit=${limit}`,
  )
  return { count: data.count, records: data.records.map(pascalToCamel as (obj: unknown) => TelemetryRecord) }
}

export async function fetchDeviceTelemetry(
  deviceId: string,
  from?: string,
  to?: string,
): Promise<{ deviceId: string; from: string; to: string; count: number; records: TelemetryRecord[] }> {
  const params = new URLSearchParams()
  if (from) params.set('from', from)
  if (to) params.set('to', to)
  const qs = params.toString()
  const data = await request<{ device_id: string; from: string; to: string; count: number; records: unknown[] }>(
    `/api/v1/telemetry/device/${deviceId}${qs ? '?' + qs : ''}`,
  )
  return {
    deviceId: data.device_id,
    from: data.from,
    to: data.to,
    count: data.count,
    records: data.records.map(pascalToCamel as (obj: unknown) => TelemetryRecord),
  }
}

export async function fetchLatestValues(factoryId?: string): Promise<LatestValuesResponse> {
  const url = factoryId
    ? `/api/v1/telemetry/latest?factory_id=${factoryId}`
    : '/api/v1/telemetry/latest'
  const data = await request<{ factory_id: string | null; count: number; records: unknown[] }>(url)
  return {
    factoryId: data.factory_id,
    count: data.count,
    records: data.records.map(pascalToCamel as (obj: unknown) => TelemetryRecord),
  }
}

export async function fetchTelemetryStats(): Promise<TelemetryStats> {
  return request<TelemetryStats>('/api/v1/stats?service=telemetry')
}
