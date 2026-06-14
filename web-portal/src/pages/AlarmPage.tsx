import React, { useState, useEffect } from 'react'
import {
  Box,
  Button,
  Card,
  CardContent,
  Typography,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  Switch,
  FormControlLabel,
  CircularProgress,
  Alert,
  Tooltip,
} from '@mui/material'
import AddIcon from '@mui/icons-material/Add'
import EditIcon from '@mui/icons-material/Edit'
import DeleteIcon from '@mui/icons-material/Delete'
import CheckIcon from '@mui/icons-material/Check'
import RefreshIcon from '@mui/icons-material/Refresh'
import ThumbUpIcon from '@mui/icons-material/ThumbUp'
import NotificationsIcon from '@mui/icons-material/Notifications'
import NotificationsActiveIcon from '@mui/icons-material/NotificationsActive'

import {
  fetchAlarmRules,
  createAlarmRule,
  updateAlarmRule,
  deleteAlarmRule,
  fetchActiveAlarms,
  acknowledgeAlarm,
  resolveAlarm,
} from '../api/alarmApi'
import type { Alarm, AlarmRule } from '../types/alarm'

interface TabPanelProps {
  children?: React.ReactNode
  index: number
  value: number
}

function CustomTabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`alarm-tabpanel-${index}`}
      aria-labelledby={`alarm-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ py: 3 }}>{children}</Box>}
    </div>
  )
}

export default function AlarmPage() {
  const [tabValue, setTabValue] = useState(0)
  const [activeAlarms, setActiveAlarms] = useState<Alarm[]>([])
  const [alarmRules, setAlarmRules] = useState<AlarmRule[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Dialog State for Rule Form
  const [ruleFormOpen, setRuleFormOpen] = useState(false)
  const [editingRule, setEditingRule] = useState<Partial<AlarmRule> | null>(null)

  // Dialog State for Acknowledge Action
  const [ackDialogOpen, setAckDialogOpen] = useState(false)
  const [selectedAlarmId, setSelectedAlarmId] = useState<string | null>(null)
  const [ackBy, setAckBy] = useState('admin')

  const loadData = async () => {
    setLoading(true)
    setError(null)
    try {
      if (tabValue === 0) {
        const alarms = await fetchActiveAlarms()
        setActiveAlarms(alarms)
      } else {
        const rules = await fetchAlarmRules()
        setAlarmRules(rules)
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Đã có lỗi xảy ra khi tải dữ liệu.')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tabValue])

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue)
  }

  // --- Alarm Rules Handlers ---

  const handleOpenCreateRule = () => {
    setEditingRule({
      name: '',
      description: '',
      type: 'threshold',
      metricName: '',
      conditionOperator: '>',
      conditionValue: 0,
      durationSeconds: 0,
      severity: 'critical',
      factoryId: 'factory-01',
      isEnabled: true,
    })
    setRuleFormOpen(true)
  }

  const handleOpenEditRule = (rule: AlarmRule) => {
    setEditingRule({ ...rule })
    setRuleFormOpen(true)
  }

  const handleSaveRule = async () => {
    if (!editingRule || !editingRule.name || !editingRule.metricName) {
      setError('Vui lòng điền đầy đủ Tên và Tên Metric.')
      return
    }

    try {
      if (editingRule.id) {
        // Update
        const updated = await updateAlarmRule(editingRule.id, {
          name: editingRule.name,
          description: editingRule.description ?? '',
          type: editingRule.type ?? 'threshold',
          metricName: editingRule.metricName,
          conditionOperator: editingRule.conditionOperator ?? '>',
          conditionValue: editingRule.conditionValue ?? 0,
          durationSeconds: editingRule.durationSeconds ?? 0,
          severity: editingRule.severity ?? 'critical',
          factoryId: editingRule.factoryId ?? 'factory-01',
          isEnabled: editingRule.isEnabled ?? true,
        })
        setAlarmRules(alarmRules.map((r) => (r.id === updated.id ? updated : r)))
      } else {
        // Create
        const created = await createAlarmRule({
          name: editingRule.name,
          description: editingRule.description ?? '',
          type: editingRule.type ?? 'threshold',
          metricName: editingRule.metricName,
          conditionOperator: editingRule.conditionOperator ?? '>',
          conditionValue: editingRule.conditionValue ?? 0,
          durationSeconds: editingRule.durationSeconds ?? 0,
          severity: editingRule.severity ?? 'critical',
          factoryId: editingRule.factoryId ?? 'factory-01',
          isEnabled: editingRule.isEnabled ?? true,
        })
        setAlarmRules([...alarmRules, created])
      }
      setRuleFormOpen(false)
      setEditingRule(null)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Lỗi khi lưu quy tắc.')
    }
  }

  const handleDeleteRule = async (id: string) => {
    if (!window.confirm('Bạn có chắc chắn muốn xóa quy tắc này không?')) return
    try {
      await deleteAlarmRule(id)
      setAlarmRules(alarmRules.filter((r) => r.id !== id))
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Lỗi khi xóa quy tắc.')
    }
  }

  const handleToggleRuleStatus = async (rule: AlarmRule) => {
    try {
      const updated = await updateAlarmRule(rule.id, {
        name: rule.name,
        description: rule.description,
        type: rule.type,
        metricName: rule.metricName,
        conditionOperator: rule.conditionOperator,
        conditionValue: rule.conditionValue,
        durationSeconds: rule.durationSeconds,
        severity: rule.severity,
        factoryId: rule.factoryId,
        isEnabled: !rule.isEnabled,
      })
      setAlarmRules(alarmRules.map((r) => (r.id === updated.id ? updated : r)))
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Lỗi khi thay đổi trạng thái quy tắc.')
    }
  }

  // --- Active Alarms Handlers ---

  const handleOpenAckDialog = (alarmId: string) => {
    setSelectedAlarmId(alarmId)
    setAckBy('admin')
    setAckDialogOpen(true)
  }

  const handleConfirmAck = async () => {
    if (!selectedAlarmId || !ackBy.trim()) return
    try {
      const updated = await acknowledgeAlarm(selectedAlarmId, { ackBy })
      setActiveAlarms(activeAlarms.map((a) => (a.id === selectedAlarmId ? updated : a)))
      setAckDialogOpen(false)
      setSelectedAlarmId(null)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Lỗi khi xác nhận cảnh báo.')
    }
  }

  const handleResolveAlarm = async (alarmId: string) => {
    try {
      await resolveAlarm(alarmId)
      // Since it's resolved, it's no longer "active", remove it from the list
      setActiveAlarms(activeAlarms.filter((a) => a.id !== alarmId))
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Lỗi khi dập cảnh báo.')
    }
  }

  // Helper colors
  const getSeverityChip = (severity: string) => {
    let color: 'warning' | 'error' | 'default' = 'default'
    if (severity === 'warning') color = 'warning'
    if (severity === 'critical') color = 'error'
    if (severity === 'emergency') color = 'error' // emergency will use error (red) too

    return (
      <Chip
        label={severity.toUpperCase()}
        color={color}
        size="small"
        sx={{ fontWeight: 'bold' }}
      />
    )
  }

  const getStatusChip = (status: string) => {
    let color: 'error' | 'info' | 'success' | 'default' = 'default'
    if (status === 'RAISED') color = 'error'
    if (status === 'ACKNOWLEDGED') color = 'info'
    if (status === 'RESOLVED') color = 'success'

    return (
      <Chip
        label={status}
        color={color}
        variant="outlined"
        size="small"
        sx={{ fontWeight: 'bold' }}
      />
    )
  }

  return (
    <Box sx={{ p: 1 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1" color="primary.main" sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
          <NotificationsActiveIcon fontSize="large" color="secondary" />
          Trung Tâm Cảnh Báo (Alarm Center)
        </Typography>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            startIcon={<RefreshIcon />}
            onClick={loadData}
            disabled={loading}
          >
            Làm mới
          </Button>
          {tabValue === 1 && (
            <Button
              variant="contained"
              color="secondary"
              startIcon={<AddIcon />}
              onClick={handleOpenCreateRule}
            >
              Thêm Quy Tắc
            </Button>
          )}
        </Box>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }} onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={tabValue} onChange={handleTabChange} aria-label="alarm tabs">
          <Tab
            label={
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <NotificationsIcon fontSize="small" />
                Cảnh Báo Đang Hoạt Động ({activeAlarms.length})
              </Box>
            }
            sx={{ fontWeight: 'bold' }}
          />
          <Tab label="Quy Tắc Cảnh Báo (Rules)" sx={{ fontWeight: 'bold' }} />
        </Tabs>
      </Box>

      {loading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 5 }}>
          <CircularProgress color="secondary" />
        </Box>
      )}

      {!loading && (
        <>
          {/* TAB 0: ACTIVE ALARMS */}
          <CustomTabPanel value={tabValue} index={0}>
            {activeAlarms.length === 0 ? (
              <Card sx={{ textAlign: 'center', py: 5, bgcolor: '#fbfbfb' }}>
                <CardContent>
                  <Typography variant="h6" color="text.secondary" gutterBottom>
                    🎉 Tuyệt vời! Không có cảnh báo nào đang hoạt động.
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Hệ thống hoạt động ổn định và các chỉ số đều trong ngưỡng an toàn.
                  </Typography>
                </CardContent>
              </Card>
            ) : (
              <TableContainer component={Paper}>
                <Table sx={{ minWidth: 650 }} aria-label="active alarms table">
                  <TableHead sx={{ bgcolor: 'primary.main' }}>
                    <TableRow>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Tên Cảnh Báo (Rule)</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Mã Máy (Machine ID)</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Mức Độ</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Trạng Thái</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Giá Trị Lỗi</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Thời Gian Bắt Đầu</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Xác Nhận Bởi</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold', textAlign: 'center' }}>Thao Tác</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {activeAlarms.map((alarm) => (
                      <TableRow key={alarm.id} hover>
                        <TableCell sx={{ fontWeight: 'medium' }}>{alarm.ruleName}</TableCell>
                        <TableCell>{alarm.machineId}</TableCell>
                        <TableCell>{getSeverityChip(alarm.severity)}</TableCell>
                        <TableCell>{getStatusChip(alarm.status)}</TableCell>
                        <TableCell color="error.main" sx={{ fontWeight: 'bold' }}>
                          {alarm.value}
                        </TableCell>
                        <TableCell>{new Date(alarm.raisedAt).toLocaleString('vi-VN')}</TableCell>
                        <TableCell>
                          {alarm.ackBy ? (
                            <Tooltip title={alarm.ackAt ? new Date(alarm.ackAt).toLocaleString('vi-VN') : ''}>
                              <Chip label={alarm.ackBy} size="small" color="default" />
                            </Tooltip>
                          ) : (
                            <Typography variant="caption" color="text.secondary">Chưa xác nhận</Typography>
                          )}
                        </TableCell>
                        <TableCell sx={{ textAlign: 'center' }}>
                          {alarm.status === 'RAISED' && (
                            <Button
                              variant="outlined"
                              color="info"
                              size="small"
                              startIcon={<ThumbUpIcon />}
                              onClick={() => handleOpenAckDialog(alarm.id)}
                              sx={{ mr: 1 }}
                            >
                              Xác Nhận (Ack)
                            </Button>
                          )}
                          <Button
                            variant="outlined"
                            color="success"
                            size="small"
                            startIcon={<CheckIcon />}
                            onClick={() => handleResolveAlarm(alarm.id)}
                          >
                            Dập Cảnh Báo
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            )}
          </CustomTabPanel>

          {/* TAB 1: ALARM RULES */}
          <CustomTabPanel value={tabValue} index={1}>
            {alarmRules.length === 0 ? (
              <Card sx={{ textAlign: 'center', py: 5, bgcolor: '#fbfbfb' }}>
                <CardContent>
                  <Typography variant="h6" color="text.secondary" gutterBottom>
                    Chưa cấu hình quy tắc cảnh báo nào.
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                    Hãy thêm một quy tắc so sánh ngưỡng (threshold) để tự động phát hiện lỗi telemetry.
                  </Typography>
                  <Button
                    variant="contained"
                    color="secondary"
                    startIcon={<AddIcon />}
                    onClick={handleOpenCreateRule}
                  >
                    Tạo Quy Tắc Đầu Tiên
                  </Button>
                </CardContent>
              </Card>
            ) : (
              <TableContainer component={Paper}>
                <Table sx={{ minWidth: 650 }} aria-label="alarm rules table">
                  <TableHead sx={{ bgcolor: 'primary.main' }}>
                    <TableRow>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Tên Quy Tắc</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Metric</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Điều Kiện Ngưỡng</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Trì Hoãn</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Mức Độ</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Nhà Máy</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold' }}>Hoạt Động</TableCell>
                      <TableCell sx={{ color: 'white', fontWeight: 'bold', textAlign: 'center' }}>Thao Tác</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {alarmRules.map((rule) => (
                      <TableRow key={rule.id} hover>
                        <TableCell>
                          <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
                            {rule.name}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            {rule.description}
                          </Typography>
                        </TableCell>
                        <TableCell sx={{ fontFamily: 'monospace' }}>{rule.metricName}</TableCell>
                        <TableCell sx={{ fontWeight: 'bold' }}>
                          {rule.conditionOperator} {rule.conditionValue}
                        </TableCell>
                        <TableCell>{rule.durationSeconds} giây</TableCell>
                        <TableCell>{getSeverityChip(rule.severity)}</TableCell>
                        <TableCell>{rule.factoryId}</TableCell>
                        <TableCell>
                          <Switch
                            checked={rule.isEnabled}
                            onChange={() => handleToggleRuleStatus(rule)}
                            color="secondary"
                          />
                        </TableCell>
                        <TableCell sx={{ textAlign: 'center' }}>
                          <IconButton
                            color="primary"
                            onClick={() => handleOpenEditRule(rule)}
                            sx={{ mr: 1 }}
                          >
                            <EditIcon />
                          </IconButton>
                          <IconButton
                            color="error"
                            onClick={() => handleDeleteRule(rule.id)}
                          >
                            <DeleteIcon />
                          </IconButton>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            )}
          </CustomTabPanel>
        </>
      )}

      {/* --- DIALOGS --- */}

      {/* Rule Form Dialog */}
      <Dialog open={ruleFormOpen} onClose={() => setRuleFormOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ bgcolor: 'primary.main', color: 'white', fontWeight: 'bold' }}>
          {editingRule?.id ? 'Cập Nhật Quy Tắc Cảnh Báo' : 'Tạo Mới Quy Tắc Cảnh Báo'}
        </DialogTitle>
        <DialogContent sx={{ pt: 3 }}>
          {editingRule && (
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2.5, mt: 1 }}>
              <TextField
                label="Tên Quy Tắc"
                required
                value={editingRule.name || ''}
                onChange={(e) => setEditingRule({ ...editingRule, name: e.target.value })}
                fullWidth
              />
              <TextField
                label="Mô Tả"
                multiline
                rows={2}
                value={editingRule.description || ''}
                onChange={(e) => setEditingRule({ ...editingRule, description: e.target.value })}
                fullWidth
              />
              <Box sx={{ display: 'flex', gap: 2 }}>
                <TextField
                  label="Tên Metric"
                  required
                  placeholder="ví dụ: temperature"
                  value={editingRule.metricName || ''}
                  onChange={(e) => setEditingRule({ ...editingRule, metricName: e.target.value })}
                  sx={{ flex: 2 }}
                />
                <TextField
                  select
                  label="Phân Loại Ngưỡng"
                  value={editingRule.type || 'threshold'}
                  onChange={(e) => setEditingRule({ ...editingRule, type: e.target.value })}
                  sx={{ flex: 1 }}
                >
                  <MenuItem value="threshold">Threshold</MenuItem>
                </TextField>
              </Box>

              <Box sx={{ display: 'flex', gap: 2 }}>
                <TextField
                  select
                  label="Toán Tử"
                  value={editingRule.conditionOperator || '>'}
                  onChange={(e) => setEditingRule({ ...editingRule, conditionOperator: e.target.value })}
                  sx={{ flex: 1 }}
                >
                  <MenuItem value=">">&gt; (Lớn hơn)</MenuItem>
                  <MenuItem value="<">&lt; (Nhỏ hơn)</MenuItem>
                  <MenuItem value="==">== (Bằng)</MenuItem>
                  <MenuItem value=">=">&gt;= (Lớn hơn hoặc bằng)</MenuItem>
                  <MenuItem value="<=">&lt;= (Nhỏ hơn hoặc bằng)</MenuItem>
                </TextField>
                <TextField
                  label="Giá Trị Ngưỡng"
                  type="number"
                  value={editingRule.conditionValue ?? 0}
                  onChange={(e) => setEditingRule({ ...editingRule, conditionValue: parseFloat(e.target.value) || 0 })}
                  sx={{ flex: 1.5 }}
                />
                <TextField
                  label="Trì Hoãn (giây)"
                  type="number"
                  helperText="0 tức cảnh báo ngay"
                  value={editingRule.durationSeconds ?? 0}
                  onChange={(e) => setEditingRule({ ...editingRule, durationSeconds: parseInt(e.target.value) || 0 })}
                  sx={{ flex: 1.5 }}
                />
              </Box>

              <Box sx={{ display: 'flex', gap: 2 }}>
                <TextField
                  select
                  label="Mức Độ Nghiêm Trọng"
                  value={editingRule.severity || 'critical'}
                  onChange={(e) => setEditingRule({ ...editingRule, severity: e.target.value })}
                  sx={{ flex: 1 }}
                >
                  <MenuItem value="warning">Warning</MenuItem>
                  <MenuItem value="critical">Critical</MenuItem>
                  <MenuItem value="emergency">Emergency</MenuItem>
                </TextField>
                <TextField
                  label="Nhà Máy (Factory ID)"
                  value={editingRule.factoryId || ''}
                  onChange={(e) => setEditingRule({ ...editingRule, factoryId: e.target.value })}
                  sx={{ flex: 1 }}
                />
              </Box>

              <FormControlLabel
                control={
                  <Switch
                    checked={editingRule.isEnabled ?? true}
                    onChange={(e) => setEditingRule({ ...editingRule, isEnabled: e.target.checked })}
                    color="secondary"
                  />
                }
                label="Cho phép quy tắc hoạt động ngay"
              />
            </Box>
          )}
        </DialogContent>
        <DialogActions sx={{ p: 3 }}>
          <Button onClick={() => setRuleFormOpen(false)} variant="outlined">
            Hủy
          </Button>
          <Button onClick={handleSaveRule} color="secondary" variant="contained">
            Lưu Quy Tắc
          </Button>
        </DialogActions>
      </Dialog>

      {/* Acknowledge Dialog */}
      <Dialog open={ackDialogOpen} onClose={() => setAckDialogOpen(false)}>
        <DialogTitle sx={{ fontWeight: 'bold' }}>Xác Nhận Xử Lý Cảnh Báo (Acknowledge)</DialogTitle>
        <DialogContent>
          <Typography variant="body2" sx={{ mb: 2 }}>
            Nhập tên của bạn hoặc mã định danh của nhân viên vận hành xử lý cảnh báo này:
          </Typography>
          <TextField
            autoFocus
            label="Tên Người Xác Nhận"
            fullWidth
            value={ackBy}
            onChange={(e) => setAckBy(e.target.value)}
          />
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button onClick={() => setAckDialogOpen(false)} variant="outlined">
            Hủy
          </Button>
          <Button onClick={handleConfirmAck} color="secondary" variant="contained" disabled={!ackBy.trim()}>
            Xác Nhận
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}
