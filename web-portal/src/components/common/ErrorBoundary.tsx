import { Component, type ReactNode } from 'react'
import { Typography, Button, Paper } from '@mui/material'

interface Props {
  children: ReactNode
  fallback?: ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
}

export default class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, info: { componentStack: string }) {
    console.error('[ErrorBoundary]', error, info.componentStack)
  }

  handleReset = () => {
    this.setState({ hasError: false, error: null })
  }

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) return this.props.fallback
      return (
        <Paper sx={{ p: 3, textAlign: 'center', my: 2 }}>
          <Typography variant="h3" sx={{ mb: 1 }}>⚠️</Typography>
          <Typography variant="h6" color="error" gutterBottom>
            Something went wrong
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            {this.state.error?.message ?? 'An unexpected error occurred'}
          </Typography>
          <Button variant="outlined" onClick={this.handleReset}>
            Retry
          </Button>
        </Paper>
      )
    }
    return this.props.children
  }
}
