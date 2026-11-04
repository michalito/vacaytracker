// Test double for $app/state
export const page = {
	url: new URL('http://localhost:5173'),
	params: {},
	route: { id: '' },
	status: 200,
	error: null,
	data: {},
	form: null,
	state: {}
};

export const navigating = {
	from: null,
	to: null,
	type: null,
	willUnload: null,
	delta: null,
	complete: null
};

export const updated = {
	current: false,
	check: async () => false
};
