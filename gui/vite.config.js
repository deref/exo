import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// Vite server runs on 3001 in development, which is proxied from the
// `exo` server running on 4001. In production, exo serves the GUI and
// the API on port 43643.
const port = 3001;

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    port,
  },
});
