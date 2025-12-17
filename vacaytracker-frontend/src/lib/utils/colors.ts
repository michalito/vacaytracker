export interface UserColor {
	background: string;
	text: string;
	combined: string;
}

// Theme-aligned color palette for user assignments in team calendar
export const USER_COLOR_PALETTE: UserColor[] = [
	{ background: 'bg-ocean-200', text: 'text-ocean-800', combined: 'bg-ocean-200 text-ocean-800' },
	{ background: 'bg-coral-300', text: 'text-coral-600', combined: 'bg-coral-300 text-coral-600' },
	{ background: 'bg-teal-200', text: 'text-teal-800', combined: 'bg-teal-200 text-teal-800' },
	{ background: 'bg-amber-200', text: 'text-amber-800', combined: 'bg-amber-200 text-amber-800' },
	{ background: 'bg-sand-300', text: 'text-sand-500', combined: 'bg-sand-300 text-sand-500' },
	{ background: 'bg-slate-200', text: 'text-slate-700', combined: 'bg-slate-200 text-slate-700' },
	{
		background: 'bg-purple-200',
		text: 'text-purple-800',
		combined: 'bg-purple-200 text-purple-800'
	},
	{
		background: 'bg-indigo-200',
		text: 'text-indigo-800',
		combined: 'bg-indigo-200 text-indigo-800'
	}
];

// Generate a hash from a string
function hashString(str: string): number {
	let hash = 0;
	for (let i = 0; i < str.length; i++) {
		hash = str.charCodeAt(i) + ((hash << 5) - hash);
	}
	return hash;
}

// Get the color index for a user
export function getUserColorIndex(userId: string): number {
	const hash = hashString(userId);
	return Math.abs(hash) % USER_COLOR_PALETTE.length;
}

// Get the color for a user (consistent based on userId)
export function getUserColor(userId: string): UserColor {
	return USER_COLOR_PALETTE[getUserColorIndex(userId)];
}
