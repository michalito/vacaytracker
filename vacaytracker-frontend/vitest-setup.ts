import '@testing-library/jest-dom/vitest';

// Mock localStorage
const localStorageMock = (() => {
	let store: Record<string, string> = {};
	return {
		getItem: (key: string) => store[key] ?? null,
		setItem: (key: string, value: string) => {
			store[key] = value;
		},
		removeItem: (key: string) => {
			delete store[key];
		},
		clear: () => {
			store = {};
		},
		get length() {
			return Object.keys(store).length;
		},
		key: (index: number) => Object.keys(store)[index] ?? null
	};
})();

Object.defineProperty(globalThis, 'localStorage', { value: localStorageMock });

// jsdom lacks pointer capture and scrollIntoView, which Bits UI's Select uses while opening
// and highlighting items. No-op versions are enough for interaction tests.
Element.prototype.hasPointerCapture ??= () => false;
Element.prototype.setPointerCapture ??= () => {};
Element.prototype.releasePointerCapture ??= () => {};
Element.prototype.scrollIntoView ??= () => {};

// jsdom has no matchMedia, which svelte/reactivity's MediaQuery needs.
// Defaults to "no match" (a narrow screen); tests can override it with vi.spyOn(window, 'matchMedia').
if (!window.matchMedia) {
	Object.defineProperty(window, 'matchMedia', {
		writable: true,
		configurable: true,
		value: (query: string): MediaQueryList =>
			({
				matches: false,
				media: query,
				onchange: null,
				addEventListener: () => {},
				removeEventListener: () => {},
				addListener: () => {},
				removeListener: () => {},
				dispatchEvent: () => false
			}) as MediaQueryList
	});
}

// Mock crypto.randomUUID
if (!globalThis.crypto) {
	Object.defineProperty(globalThis, 'crypto', {
		value: {
			randomUUID: () =>
				'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
					const r = (Math.random() * 16) | 0;
					const v = c === 'x' ? r : (r & 0x3) | 0x8;
					return v.toString(16);
				})
		}
	});
}
