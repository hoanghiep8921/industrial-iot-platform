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
  factoryId: string
  severity: string // "warning", "critical", "emergency"
  value: number
  status: string // "RAISED", "ACKNOWLEDGED", "RESOLVED"
  raisedAt: string
  ackAt?: string | null
  ackBy?: string
  resolvedAt?: string | null
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
