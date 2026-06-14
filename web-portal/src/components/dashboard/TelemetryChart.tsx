import { useMemo, useState } from 'react'
import { Card, CardContent, CardHeader, FormControl, InputLabel, Select, MenuItem } from '@mui/material'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import type { TelemetryRecord } from '../../types/telemetry'
import { formatDateTime, formatNumber } from '../../utils/formatters'
import EmptyState from '../common/EmptyState'

interface TelemetryChartProps {
  records: TelemetryRecord[]
}

export default function TelemetryChart({ records }: TelemetryChartProps) {
  const metrics = useMemo(() => [...new Set(records.map((r) => r.metricName))], [records])
  const [selectedMetric, setSelectedMetric] = useState<string>(metrics[0] ?? '')

  const chartData = useMemo(() => {
    if (!selectedMetric) return []
    return records
      .filter((r) => r.metricName === selectedMetric)
      .slice(-100)
      .map((r) => ({
        time: formatDateTime(r.time ?? ''),
        value: typeof r.value === 'number' ? r.value : null,
        device: (r.deviceId ?? '').slice(0, 8),
      }))
  }, [records, selectedMetric])

  if (chartData.length === 0) {
    return (
      <Card>
        <CardHeader title="Telemetry History" />
        <CardContent>
          <EmptyState message="No telemetry data available" />
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader
        title="Telemetry History"
        action={
          <FormControl size="small" sx={{ minWidth: 180 }}>
            <InputLabel>Metric</InputLabel>
            <Select
              value={selectedMetric}
              label="Metric"
              onChange={(e) => setSelectedMetric(e.target.value)}
            >
              {metrics.map((m) => (
                <MenuItem key={m} value={m}>{m}</MenuItem>
              ))}
            </Select>
          </FormControl>
        }
      />
      <CardContent>
        <ResponsiveContainer width="100%" height={350}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="time" fontSize={12} />
            <YAxis fontSize={12} tickFormatter={(v) => formatNumber(v, 0)} />
            <Tooltip labelFormatter={(l) => 'Time: ' + l} />
            <Legend />
            <Line type="monotone" dataKey="value" stroke="#1976d2" dot={false} name={selectedMetric} />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}
