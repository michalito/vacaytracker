import { adminApi } from '$lib/api/admin';
import type { User, Settings, VacationRequest, PaginationInfo, NewsletterPreview } from '$lib/types';

function createAdminStore() {
	let users = $state<User[]>([]);
	let pendingRequests = $state<VacationRequest[]>([]);
	let settings = $state<Settings | null>(null);
	let pagination = $state<PaginationInfo>({ page: 1, limit: 10, total: 0, totalPages: 0 });
	let isLoading = $state(false);
	let newsletterPreview = $state<NewsletterPreview | null>(null);
	let isSendingNewsletter = $state(false);

	async function fetchUsers(params?: { page?: number; limit?: number; role?: string; search?: string }) {
		isLoading = true;
		try {
			const response = await adminApi.listUsers(params);
			users = response.users;
			pagination = response.pagination;
		} catch (error) {
			console.error('Failed to fetch users:', error);
			users = [];
		} finally {
			isLoading = false;
		}
	}

	async function fetchPendingRequests() {
		try {
			const response = await adminApi.pendingRequests();
			pendingRequests = response.requests;
		} catch (error) {
			console.error('Failed to fetch pending requests:', error);
			pendingRequests = [];
		}
	}

	async function fetchSettings() {
		try {
			settings = await adminApi.getSettings();
		} catch (error) {
			console.error('Failed to fetch settings:', error);
			settings = null;
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
		} catch (error) {
			console.error('Failed to fetch newsletter preview:', error);
			newsletterPreview = null;
			throw error;
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
		sendNewsletter
	};
}

export const admin = createAdminStore();
