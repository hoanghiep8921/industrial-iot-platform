import { useEffect, useState, useMemo } from 'react'
import { Card, CardContent, CardHeader, ToggleButtonGroup, ToggleButton } from '@mui/material'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import type { TelemetryRecord } from '../../types/telemetry'
import { formatDateTime, formatNumber } from '../../utils/formatters'
import EmptyState from '../common/EmptyState'

interface TelemetryHistoryChartProps {
  deviceId: string
  records: TelemetryRecord[]
  onFetch: (from: string, to: string) => void
}

const timeRanges: { label: string; hours: number }[] = [
  { label: '1h', hours: 1 },
  { label: '6h', hours: 6 },
  { label: '24h', hours: 24 },
  { label: '7d', hours: 168 },
]

export default function TelemetryHistoryChart({ deviceId, records, onFetch }: TelemetryHistoryChartProps) {
  const [hours, setHours] = useState(24)

  useEffect(() => {
    const to = new Date().toISOString()
    const from = new Date(Date.now() - hours * 3600_000).toISOString()
    onFetch(from, to)
  }, [deviceId, hours, onFetch])

  const metrics = useMemo(() => [...new Set(records.map((r) => r.metricName))], [records])

  const chartData = useMemo(() => {
    const byTime = new Map<string, Record<string, number>>()
    for (const r of records) {
      const t = formatDateTime(r.time)
      const entry = byTime.get(t) ?? {}
      entry[r.metricName] = r.value
      byTime.set(t, entry)
    }
    return [...byTime.entries()].map(([time, vals]) => ({ time, ...vals }))
  }, [records])

  if (chartData.length === 0) {
    return (
      <Card>
        <CardHeader title="Telemetry" />
        <CardContent><EmptyState message="No telemetry data for this period" /></CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader
        title="Telemetry"
        action={
          <ToggleButtonGroup size="small" value={hours} exclusive onChange={(_, v) => v && setHours(v)}>
            {timeRanges.map((tr) => (
              <ToggleButton key={tr.hours} value={tr.hours}>{tr.label}</ToggleButton>
            ))}
          </ToggleButtonGroup>
        }
      />
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="time" fontSize={12} />
            <YAxis fontSize={12} tickFormatter={(v) => formatNumber(v, 0)} />
            <Tooltip />
            <Legend />
            {metrics.map((m, i) => (
              <Line key={m} type="monotone" dataKey={m} stroke={['#1976d2', '#ef6c00', '#2e7d32'][i % 3]} dot={false} />
            ))}
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}
