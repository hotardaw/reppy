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
        description: "Reppy is a smart tracking app that helps you reach your fitness & health goals - track your workouts, whether you're into cardio or lifting, and create reproducible recipes to conveniently log your meals.",
        icons: [
          // Add your app icons in different sizes
          {
            src: '/src/assets/reppy-app-logo.png',
            sizes: '1154x1154',
            type: 'image/png',
            purpose: 'any maskable'
          }
          // {
          //   src: '/src/assets/reppy-logo-192.jpg',
          //   sizes: '192x192',
          //   type: 'image/jpg',
          // },
          // {
          //   src: '/src/assets/reppy-logo-180.jpg',
          //   sizes: '180x180',
          //   type: 'image/jpg',
          // }
        ],
        // theme_color: '#ffffff', // Add your theme color
        // background_color: '#ffffff', // Add your background color
        // display: 'standalone',
        // scope: '/',
        // start_url: '/',
        // orientation: 'portrait'
      },
      // includeAssets: ['favicon.ico'], // Add favicon and other assets
      // workbox: {
      //   sourcemap: true
      // }
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
