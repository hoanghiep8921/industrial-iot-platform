import { useLocation, useNavigate } from 'react-router-dom'
import {
  Box,
  Drawer,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  useMediaQuery,
  useTheme,
} from '@mui/material'
import DashboardIcon from '@mui/icons-material/Dashboard'
import DevicesIcon from '@mui/icons-material/PrecisionManufacturing'
import NotificationsActiveIcon from '@mui/icons-material/NotificationsActive'
import { useEffect } from 'react'

const navItems = [
  { path: '/', label: 'Dashboard', icon: <DashboardIcon /> },
  { path: '/devices', label: 'Devices', icon: <DevicesIcon /> },
  { path: '/alarms', label: 'Alarms', icon: <NotificationsActiveIcon /> },
]

interface SidebarProps {
  drawerWidth: number
  mobileOpen: boolean
  onClose: () => void
}

export default function Sidebar({ drawerWidth, mobileOpen, onClose }: SidebarProps) {
  const navigate = useNavigate()
  const location = useLocation()
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('md'))

  useEffect(() => {
    if (mobileOpen && isMobile) onClose()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname])

  const drawerContent = (
    <>
      <Toolbar>
        <Typography variant="h6" noWrap>
          🏭 IIoT Platform
        </Typography>
      </Toolbar>
      <List>
        {navItems.map((item) => (
          <ListItemButton
            key={item.path}
            selected={
              location.pathname === item.path ||
              (item.path !== '/' && location.pathname.startsWith(item.path))
            }
            onClick={() => navigate(item.path)}
          >
            <ListItemIcon>{item.icon}</ListItemIcon>
            <ListItemText primary={item.label} />
          </ListItemButton>
        ))}
      </List>
    </>
  )

  return (
    <Box
      component="nav"
      sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}
    >
      {/* Mobile drawer */}
      <Drawer
        variant="temporary"
        open={mobileOpen}
        onClose={onClose}
        ModalProps={{ keepMounted: true }}
        sx={{
          display: { xs: 'block', md: 'none' },
          '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
        }}
      >
        {drawerContent}
      </Drawer>

      {/* Desktop permanent drawer */}
      <Drawer
        variant="permanent"
        sx={{
          display: { xs: 'none', md: 'block' },
          '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
        }}
        open
      >
        {drawerContent}
      </Drawer>
    </Box>
  )
}
