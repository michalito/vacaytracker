import { request } from './client';
import type {
	User,
	UserListResponse,
	VacationRequest,
	Settings,
	CreateUserForm,
	UpdateUserForm,
	NewsletterPreview,
	NewsletterSendResponse
} from '$lib/types';

export const adminApi = {
	// Users
	listUsers: (params?: {
		page?: number;
		limit?: number;
		role?: string;
		search?: string;
	}): Promise<UserListResponse> => {
		const searchParams = new URLSearchParams();
		if (params?.page) searchParams.set('page', params.page.toString());
		if (params?.limit) searchParams.set('limit', params.limit.toString());
		if (params?.role) searchParams.set('role', params.role);
		if (params?.search) searchParams.set('search', params.search);

		const query = searchParams.toString();
		return request(`/admin/users${query ? `?${query}` : ''}`);
	},

	createUser: (data: CreateUserForm): Promise<User> => {
		return request('/admin/users', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	getUser: (id: string): Promise<User> => {
		return request(`/admin/users/${id}`);
	},

	updateUser: (id: string, data: UpdateUserForm): Promise<User> => {
		return request(`/admin/users/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},

	deleteUser: (id: string): Promise<void> => {
		return request(`/admin/users/${id}`, { method: 'DELETE' });
	},

	updateBalance: (id: string, balance: number): Promise<User> => {
		return request(`/admin/users/${id}/balance`, {
			method: 'PUT',
			body: JSON.stringify({ vacationBalance: balance })
		});
	},

	resetAllBalances: (): Promise<{
		success: boolean;
		usersUpdated: number;
		newBalance: number;
		message: string;
	}> => {
		return request('/admin/users/reset-balances', { method: 'POST' });
	},

	// Vacation Management
	pendingRequests: (): Promise<{ requests: VacationRequest[]; total: number }> => {
		return request('/admin/vacation/pending');
	},

	reviewRequest: (
		id: string,
		status: 'approved' | 'rejected',
		reason?: string
	): Promise<VacationRequest> => {
		return request(`/admin/vacation/${id}/review`, {
			method: 'PUT',
			body: JSON.stringify({ status, reason })
		});
	},

	// Settings
	getSettings: (): Promise<Settings> => {
		return request('/admin/settings');
	},

	updateSettings: (data: Partial<Settings>): Promise<Settings> => {
		return request('/admin/settings', {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},

	// Newsletter
	sendNewsletter: (): Promise<NewsletterSendResponse> => {
		return request('/admin/newsletter/send', { method: 'POST' });
	},

	getNewsletterPreview: (): Promise<NewsletterPreview> => {
		return request('/admin/newsletter/preview');
	}
};
