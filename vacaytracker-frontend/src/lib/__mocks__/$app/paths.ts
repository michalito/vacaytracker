// Test double for $app/paths: no base path, so resolved routes and assets are unchanged
export const base = '';
export const assets = '';

export function resolve(path: string): string {
	return path;
}

export function asset(path: string): string {
	return path;
}
