/**
 * Compatibility shims for third-party libraries.
 *
 * Some dependencies still import `svelte/elements.js`, but Svelte exports these
 * types from `svelte/elements` (no `.js` extension).
 */
declare module 'svelte/elements.js' {
	export * from 'svelte/elements';
}



