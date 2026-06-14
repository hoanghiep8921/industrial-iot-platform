export type FirmwareStatus = 'draft' | 'released' | 'deprecated'
export type UpdateStatus = 'pending' | 'downloading' | 'installing' | 'success' | 'failed' | 'rolled_back'

export interface Firmware {
  id: string
  version: string
  model?: string | null
  description?: string | null
  fileName?: string | null
  fileSize?: number | null
  checksumSha256?: string | null
  status: FirmwareStatus
  createdBy?: string | null
  createdAt: string
  releasedAt?: string | null
}

export interface FirmwareUpdate {
  id: string
  firmwareId: string
  deviceId: string
  campaignName?: string | null
  status: UpdateStatus
  fromVersion?: string | null
  toVersion: string
  priority: number
  retryCount: number
  maxRetries: number
  startedAt?: string | null
  completedAt?: string | null
  errorMessage?: string | null
  createdAt: string
  updatedAt: string
}

export interface CreateUpdateRequest {
  firmwareId: string
  deviceId?: string
  deviceIds?: string[]
  campaignName?: string
  priority?: number
  maxRetries?: number
}
