import { AppBar, Toolbar, Typography, IconButton, LinearProgress, useMediaQuery, useTheme } from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu'

interface TopBarProps {
  loading?: boolean
  onMenuClick: () => void
}

const DRAWER_WIDTH = 240

export default function TopBar({ loading, onMenuClick }: TopBarProps) {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('md'))

  return (
    <AppBar
      position="fixed"
      sx={{
        zIndex: theme.zIndex.drawer + 1,
        width: { md: `calc(100% - ${DRAWER_WIDTH}px)` },
        ml: { md: `${DRAWER_WIDTH}px` },
      }}
    >
      <Toolbar>
        {isMobile && (
          <IconButton color="inherit" edge="start" onClick={onMenuClick} sx={{ mr: 2 }}>
            <MenuIcon />
          </IconButton>
        )}
        <Typography variant="h6" sx={{ flexGrow: 1 }}>
          Industrial IoT Platform
        </Typography>
      </Toolbar>
      {loading && <LinearProgress color="secondary" />}
    </AppBar>
  )
}

