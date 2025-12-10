export interface UserColor {
	background: string;
	text: string;
	combined: string;
}

export const USER_COLOR_PALETTE: UserColor[] = [
	{ background: 'bg-ocean-200', text: 'text-ocean-800', combined: 'bg-ocean-200 text-ocean-800' },
	{ background: 'bg-green-200', text: 'text-green-800', combined: 'bg-green-200 text-green-800' },
	{
		background: 'bg-purple-200',
		text: 'text-purple-800',
		combined: 'bg-purple-200 text-purple-800'
	},
	{ background: 'bg-pink-200', text: 'text-pink-800', combined: 'bg-pink-200 text-pink-800' },
	{
		background: 'bg-yellow-200',
		text: 'text-yellow-800',
		combined: 'bg-yellow-200 text-yellow-800'
	},
	{ background: 'bg-red-200', text: 'text-red-800', combined: 'bg-red-200 text-red-800' },
	{ background: 'bg-teal-200', text: 'text-teal-800', combined: 'bg-teal-200 text-teal-800' },
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
