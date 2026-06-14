import { request } from './client'
import { snakeToCamel, camelToSnake } from '../utils/transformKeys'
import type { Alarm, AlarmRule, CreateAlarmRuleRequest, UpdateAlarmRuleRequest, AcknowledgeAlarmRequest } from '../types/alarm'

// ========== Alarm Rules ==========

export async function fetchAlarmRules(): Promise<AlarmRule[]> {
  const data = await request<unknown[]>('/api/v1/alarms/rules')
  if (!data) return []
  return data.map(snakeToCamel as (obj: unknown) => AlarmRule)
}

export async function fetchAlarmRule(id: string): Promise<AlarmRule> {
  return (snakeToCamel as (obj: unknown) => AlarmRule)(await request(`/api/v1/alarms/rules/${id}`))
}

export async function createAlarmRule(data: CreateAlarmRuleRequest): Promise<AlarmRule> {
  return (snakeToCamel as (obj: unknown) => AlarmRule)(
    await request('/api/v1/alarms/rules', {
      method: 'POST',
      body: JSON.stringify(camelToSnake(data)),
    }),
  )
}

export async function updateAlarmRule(id: string, data: UpdateAlarmRuleRequest): Promise<AlarmRule> {
  return (snakeToCamel as (obj: unknown) => AlarmRule)(
    await request(`/api/v1/alarms/rules/${id}`, {
      method: 'PUT',
      body: JSON.stringify(camelToSnake(data)),
    }),
  )
}

export async function deleteAlarmRule(id: string): Promise<void> {
  await request(`/api/v1/alarms/rules/${id}`, { method: 'DELETE' })
}

// ========== Active Alarms ==========

export async function fetchActiveAlarms(): Promise<Alarm[]> {
  const data = await request<unknown[]>('/api/v1/alarms/active')
  if (!data) return []
  return data.map(snakeToCamel as (obj: unknown) => Alarm)
}

export async function acknowledgeAlarm(id: string, data: AcknowledgeAlarmRequest): Promise<Alarm> {
  return (snakeToCamel as (obj: unknown) => Alarm)(
    await request(`/api/v1/alarms/${id}/acknowledge`, {
      method: 'POST',
      body: JSON.stringify(camelToSnake(data)),
    }),
  )
}

export async function resolveAlarm(id: string): Promise<Alarm> {
  return (snakeToCamel as (obj: unknown) => Alarm)(
    await request(`/api/v1/alarms/${id}/resolve`, {
      method: 'POST',
    }),
  )
}
