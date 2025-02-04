import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(),
    VitePWA({
      registerType: 'autoUpdate', // Optional: set how service worker updates
      manifest: {
        name: 'Reppy - Workout and Diet Tracker',
        short_name: 'Reppy',
        description: "Reppy is a smart tracking app that helps you reach your fitness & health goals.",
        icons: [
          {
            src: '/public/reppy-app-logo-180.png',
            sizes: '180x180',
            type: 'image/png',
            purpose: 'any'
          },
          {
            src: '/public/reppy-app-logo-512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any'
          },
          {
            src: '/public/reppy-app-logo-192.png',
            sizes: '192x192',
            type: 'image/png',
            purpose: 'any maskable'
          }
        ],
        theme_color: '#ffffff',
        background_color: '#ffffff',
        display: 'standalone',
        scope: '/',
        start_url: '/',
        orientation: 'portrait'
      },
      includeAssets: ['favicon.ico'],
      workbox: {
        sourcemap: true
      }
    })
  ],
  server: {
    host: '0.0.0.0',
    port: 8080,
    strictPort: true,
    cors: true,
    allowedHosts: [
      'reppy.io',
      'localhost',
    ],
  }
})
