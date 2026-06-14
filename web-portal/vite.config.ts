import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api/v1/telemetry': { target: 'http://localhost:8082', changeOrigin: true },
      '/api/v1/alarms': { target: 'http://localhost:8086', changeOrigin: true },
      '/api/v1/devices': { target: 'http://localhost:8085', changeOrigin: true },
      '/api/v1/firmwares': { target: 'http://localhost:8085', changeOrigin: true },
      '/api/v1/firmware-updates': { target: 'http://localhost:8085', changeOrigin: true },
      '/api/v1': { target: 'http://localhost:8085', changeOrigin: true },
    },
  },
  build: { outDir: 'dist', sourcemap: true },
})
