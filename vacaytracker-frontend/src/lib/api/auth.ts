import { request, setAuthToken, clearAuthToken } from './client';
import type { User, LoginResponse, EmailPreferences } from '$lib/types';

export const authApi = {
	login: async (email: string, password: string): Promise<LoginResponse> => {
		const response = await request<LoginResponse>('/auth/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
		setAuthToken(response.token);
		return response;
	},

	logout: (): void => {
		clearAuthToken();
	},

	me: (): Promise<User> => {
		return request<User>('/auth/me');
	},

	changePassword: (currentPassword: string, newPassword: string): Promise<void> => {
		return request('/auth/password', {
			method: 'PUT',
			body: JSON.stringify({ currentPassword, newPassword })
		});
	},

	updateEmailPreferences: (
		prefs: Partial<EmailPreferences>
	): Promise<{ emailPreferences: EmailPreferences }> => {
		return request('/auth/email-preferences', {
			method: 'PUT',
			body: JSON.stringify(prefs)
		});
	}
};
