export interface AlarmRule {
  id: string
  name: string
  description: string
  type: string // "threshold", etc.
  metricName: string
  conditionOperator: string // ">", "<", "==", etc.
  conditionValue: number
  durationSeconds: number
  severity: string // "warning", "critical", "emergency"
  factoryId: string
  isEnabled: boolean
  createdAt: string
  updatedAt: string
}

export interface Alarm {
  id: string
  ruleId: string
  ruleName: string
  machineId: string
  metricName: string
  triggerValue: number
  triggerTime: string
  ackTime?: string | null
  ackBy?: string
  resolveTime?: string | null
  status: string // "RAISED", "ACKNOWLEDGED", "RESOLVED"
  severity: string // "warning", "critical", "emergency"
  createdAt: string
  updatedAt: string
}

export interface CreateAlarmRuleRequest {
  name: string
  description: string
  type: string
  metricName: string
  conditionOperator: string
  conditionValue: number
  durationSeconds: number
  severity: string
  factoryId: string
  isEnabled: boolean
}

export interface UpdateAlarmRuleRequest {
  name: string
  description: string
  type: string
  metricName: string
  conditionOperator: string
  conditionValue: number
  durationSeconds: number
  severity: string
  factoryId: string
  isEnabled: boolean
}

export interface AcknowledgeAlarmRequest {
  ackBy: string
}
