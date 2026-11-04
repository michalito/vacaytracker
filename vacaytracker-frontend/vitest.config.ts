import { svelte } from '@sveltejs/vite-plugin-svelte';
import { svelteTesting } from '@testing-library/svelte/vite';
import path from 'node:path';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [svelte(), svelteTesting()],
	test: {
		environment: 'jsdom',
		include: ['src/**/*.{test,spec}.{js,ts}'],
		setupFiles: ['./vitest-setup.ts'],
		css: false,
		// `#lib/*` resolves through package.json "imports"; only SvelteKit's $app modules need doubles
		alias: {
			$app: path.resolve('./src/lib/__mocks__/$app')
		},
		restoreMocks: true
	},
	resolve: {
		conditions: ['browser']
	}
});
