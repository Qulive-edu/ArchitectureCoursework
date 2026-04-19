import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  return {
    plugins: [react()],
    define: {
      // Инъекция переменных в код при сборке
      'import.meta.env.VITE_API_BASE': JSON.stringify(env.VITE_API_BASE || '/api'),
    },
  }
})