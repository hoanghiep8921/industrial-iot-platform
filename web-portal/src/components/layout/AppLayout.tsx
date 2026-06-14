import { useState } from 'react'
import { Outlet } from 'react-router-dom'
import { Box } from '@mui/material'
import TopBar from './TopBar'
import Sidebar from './Sidebar'
import { useDeviceStore } from '../../stores/useDeviceStore'
import { useTelemetryStore } from '../../stores/useTelemetryStore'

const DRAWER_WIDTH = 240

export default function AppLayout() {
  const [mobileOpen, setMobileOpen] = useState(false)
  const deviceLoading = useDeviceStore((s) => s.loading)
  const telemetryLoading = useTelemetryStore((s) => s.loading)

  return (
    <Box sx={{ display: 'flex' }}>
      <TopBar loading={deviceLoading || telemetryLoading} onMenuClick={() => setMobileOpen(true)} />
      <Sidebar
        drawerWidth={DRAWER_WIDTH}
        mobileOpen={mobileOpen}
        onClose={() => setMobileOpen(false)}
      />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { md: `calc(100% - ${DRAWER_WIDTH}px)` },
          mt: '64px', // AppBar height
          minHeight: '100vh',
        }}
      >
        <Outlet />
      </Box>
    </Box>
  )
}
