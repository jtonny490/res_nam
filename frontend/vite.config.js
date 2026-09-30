import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const backend = process.env.BACKEND_URL || 'http://localhost:8080'

export default defineConfig({
  plugins: [react()],
  base: './',
  server: {
    proxy: {
      '/api': backend,
      '/uploads': backend,
      '/health': backend,
    },
  },
})
