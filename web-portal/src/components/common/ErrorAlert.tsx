import { Alert, Box, Button } from '@mui/material'

interface ErrorAlertProps {
  message: string
  onRetry?: () => void
}

export default function ErrorAlert({ message, onRetry }: ErrorAlertProps) {
  return (
    <Box sx={{ py: 4, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
      <Alert severity="error" sx={{ maxWidth: 500 }}>
        {message}
      </Alert>
      {onRetry && (
        <Button variant="outlined" onClick={onRetry}>
          Retry
        </Button>
      )}
    </Box>
  )
}
