export interface TelemetryRecord {
  time: string        // from "Time"
  deviceId: string    // from "DeviceID"
  factoryId: string   // from "FactoryID"
  area: string        // from "Area"
  metricName: string  // from "MetricName"
  value: number       // from "Value"
  unit: string        // from "Unit"
  quality: number     // from "Quality"
  tags: Record<string, string>
}

export interface TelemetryResponse {
  count: number
  records: TelemetryRecord[]
}

export interface LatestValuesResponse {
  factoryId: string | null
  count: number
  records: TelemetryRecord[]
}

export interface TelemetryStats {
  totalSizeBytes: number
  totalSizeMb: number
  totalSizeGb: number
  service: string
}
