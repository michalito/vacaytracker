import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit({
			preprocess: [vitePreprocess()],
			// Production runs `node build` behind a reverse proxy (see Dockerfile, aidocs/09)
			adapter: adapter()
		})
	],
	server: {
		proxy: {
			'/api': {
				target: process.env.VITE_API_URL || 'http://localhost:3000',
				changeOrigin: true
			}
		}
	}
});
