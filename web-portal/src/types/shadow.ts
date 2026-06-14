export interface DeviceShadow {
  deviceId: string
  desiredState: Record<string, unknown>
  reportedState: Record<string, unknown>
  delta: Record<string, unknown>
  version: number
  updatedAt: string
}
