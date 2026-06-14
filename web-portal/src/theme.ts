import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    primary: {
      main: '#37474f', // Blue-grey 800
      light: '#62727b',
      dark: '#102027',
    },
    secondary: {
      main: '#ef6c00', // Orange 800
      light: '#ff9d3f',
      dark: '#b53d00',
    },
    success: {
      main: '#2e7d32',
    },
    error: {
      main: '#c62828',
    },
    warning: {
      main: '#ed6c02',
    },
    background: {
      default: '#f5f5f5',
      paper: '#ffffff',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h4: {
      fontWeight: 600,
    },
    h6: {
      fontWeight: 600,
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiCard: {
      styleOverrides: {
        root: {
          boxShadow: '0 2px 8px rgba(0,0,0,0.08)',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
        },
      },
    },
  },
});

export default theme;
