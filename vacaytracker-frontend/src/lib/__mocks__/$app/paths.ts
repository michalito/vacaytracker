// Test double for $app/paths (SvelteKit 3), with no base path.
// Paths have no leading slash ('dashboard'), route IDs do ('/(app)/dashboard');
// both resolve to an absolute pathname with route groups removed.
export function resolve(pathOrRouteId: string): string {
	const absolute = pathOrRouteId.startsWith('/') ? pathOrRouteId : `/${pathOrRouteId}`;
	return absolute.replace(/\/\([^/)]+\)/g, '') || '/';
}

export function asset(file: string): string {
	return `/${file}`;
}
