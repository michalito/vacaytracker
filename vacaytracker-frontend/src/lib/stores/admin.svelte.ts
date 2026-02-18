import { adminApi } from '$lib/api/admin';
import type { User, Settings, VacationRequest, PaginationInfo, NewsletterPreview } from '$lib/types';

const CACHE_TTL = 5 * 60 * 1000; // 5 minutes

interface CacheInfo {
	pendingRequestsAt: number | null;
	usersAt: number | null;
	settingsAt: number | null;
}

function createAdminStore() {
	let users = $state<User[]>([]);
	let pendingRequests = $state<VacationRequest[]>([]);
	let settings = $state<Settings | null>(null);
	let pagination = $state<PaginationInfo>({ page: 1, limit: 10, total: 0, totalPages: 0 });
	let isLoading = $state(false);
	let error = $state<string | null>(null);
	let newsletterPreview = $state<NewsletterPreview | null>(null);
	let isSendingNewsletter = $state(false);
	let cacheInfo = $state<CacheInfo>({
		pendingRequestsAt: null,
		usersAt: null,
		settingsAt: null
	});

	function isCacheValid(key: keyof CacheInfo): boolean {
		const timestamp = cacheInfo[key];
		if (!timestamp) return false;
		return Date.now() - timestamp < CACHE_TTL;
	}

	async function fetchUsers(
		params?: { page?: number; limit?: number; role?: string; search?: string },
		options?: { force?: boolean }
	) {
		// Only use cache for default list (page 1, no filters)
		const isDefaultList = !params?.page || params.page === 1;
		const hasFilters = params?.search || params?.role;

		if (!options?.force && isDefaultList && !hasFilters && isCacheValid('usersAt')) {
			return;
		}

		isLoading = true;
		error = null;
		try {
			const response = await adminApi.listUsers(params);
			users = response.users;
			pagination = response.pagination;
			if (isDefaultList && !hasFilters) {
				cacheInfo = { ...cacheInfo, usersAt: Date.now() };
			}
		} catch (e) {
			console.error('Failed to fetch users:', e);
			error = e instanceof Error ? e.message : 'Failed to fetch users';
			users = [];
		} finally {
			isLoading = false;
		}
	}

	async function fetchPendingRequests(options?: { force?: boolean }) {
		if (!options?.force && isCacheValid('pendingRequestsAt')) {
			return;
		}

		error = null;
		try {
			const response = await adminApi.pendingRequests();
			pendingRequests = response.requests;
			cacheInfo = { ...cacheInfo, pendingRequestsAt: Date.now() };
		} catch (e) {
			console.error('Failed to fetch pending requests:', e);
			error = e instanceof Error ? e.message : 'Failed to fetch pending requests';
			pendingRequests = [];
		}
	}

	async function fetchSettings(options?: { force?: boolean }) {
		if (!options?.force && isCacheValid('settingsAt')) {
			return;
		}

		error = null;
		try {
			settings = await adminApi.getSettings();
			cacheInfo = { ...cacheInfo, settingsAt: Date.now() };
		} catch (e) {
			console.error('Failed to fetch settings:', e);
			error = e instanceof Error ? e.message : 'Failed to fetch settings';
			settings = null;
		}
	}

	async function reviewRequest(
		id: string,
		status: 'approved' | 'rejected',
		reason?: string
	): Promise<VacationRequest> {
		const reviewed = await adminApi.reviewRequest(id, status, reason);
		// Optimistic update - remove from pending requests
		pendingRequests = pendingRequests.filter((r) => r.id !== id);
		return reviewed;
	}

	function invalidateCache(keys?: (keyof CacheInfo)[]) {
		if (keys) {
			const updates = keys.reduce((acc, key) => ({ ...acc, [key]: null }), {} as Partial<CacheInfo>);
			cacheInfo = { ...cacheInfo, ...updates };
		} else {
			cacheInfo = { pendingRequestsAt: null, usersAt: null, settingsAt: null };
		}
	}

	async function createUser(data: Parameters<typeof adminApi.createUser>[0]) {
		const user = await adminApi.createUser(data);
		users = [user, ...users];
		return user;
	}

	async function updateUser(id: string, data: Parameters<typeof adminApi.updateUser>[1]) {
		const user = await adminApi.updateUser(id, data);
		users = users.map((u) => (u.id === id ? user : u));
		return user;
	}

	async function deleteUser(id: string) {
		await adminApi.deleteUser(id);
		users = users.filter((u) => u.id !== id);
	}

	async function updateBalance(id: string, balance: number) {
		const user = await adminApi.updateBalance(id, balance);
		users = users.map((u) => (u.id === id ? user : u));
		return user;
	}

	async function updateSettings(data: Parameters<typeof adminApi.updateSettings>[0]) {
		settings = await adminApi.updateSettings(data);
		return settings;
	}

	async function fetchNewsletterPreview() {
		try {
			newsletterPreview = await adminApi.getNewsletterPreview();
			return newsletterPreview;
		} catch (e) {
			console.error('Failed to fetch newsletter preview:', e);
			newsletterPreview = null;
			throw e;
		}
	}

	async function sendNewsletter() {
		isSendingNewsletter = true;
		try {
			const response = await adminApi.sendNewsletter();
			// Refresh settings to get updated lastSentAt
			await fetchSettings();
			return response;
		} finally {
			isSendingNewsletter = false;
		}
	}

	return {
		get users() {
			return users;
		},
		get pendingRequests() {
			return pendingRequests;
		},
		get settings() {
			return settings;
		},
		get pagination() {
			return pagination;
		},
		get isLoading() {
			return isLoading;
		},
		get error() {
			return error;
		},
		get newsletterPreview() {
			return newsletterPreview;
		},
		get isSendingNewsletter() {
			return isSendingNewsletter;
		},
		fetchUsers,
		fetchPendingRequests,
		fetchSettings,
		createUser,
		updateUser,
		deleteUser,
		updateBalance,
		updateSettings,
		fetchNewsletterPreview,
		sendNewsletter,
		reviewRequest,
		invalidateCache
	};
}

export const admin = createAdminStore();
