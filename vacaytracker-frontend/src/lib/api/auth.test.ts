import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('./client', () => ({
	request: vi.fn(),
	setAuthToken: vi.fn(),
	clearAuthToken: vi.fn()
}));

import { request, setAuthToken, clearAuthToken } from './client';
import { authApi } from './auth';
import type { LoginResponse, User, EmailPreferences } from '$lib/types';

describe('authApi', () => {
	beforeEach(() => {
		vi.mocked(request).mockReset();
		vi.mocked(setAuthToken).mockReset();
		vi.mocked(clearAuthToken).mockReset();
	});

	describe('login', () => {
		it('calls request with POST /auth/login, stores token, and returns response', async () => {
			const loginResponse: LoginResponse = {
				token: 'jwt-token-abc',
				user: {
					id: '1',
					email: 'test@example.com',
					name: 'Test User',
					role: 'employee',
					vacationBalance: 20,
					emailPreferences: {
						vacationUpdates: true,
						weeklyDigest: true,
						teamNotifications: true
					}
				}
			};
			vi.mocked(request).mockResolvedValue(loginResponse);

			const result = await authApi.login('test@example.com', 'password123');

			expect(request).toHaveBeenCalledWith('/auth/login', {
				method: 'POST',
				body: JSON.stringify({ email: 'test@example.com', password: 'password123' })
			});
			expect(setAuthToken).toHaveBeenCalledWith('jwt-token-abc');
			expect(result).toEqual(loginResponse);
		});
	});

	describe('logout', () => {
		it('calls clearAuthToken', () => {
			authApi.logout();

			expect(clearAuthToken).toHaveBeenCalled();
		});
	});

	describe('me', () => {
		it('calls request with /auth/me', async () => {
			const user: User = {
				id: '1',
				email: 'test@example.com',
				name: 'Test User',
				role: 'employee',
				vacationBalance: 20,
				emailPreferences: {
					vacationUpdates: true,
					weeklyDigest: true,
					teamNotifications: true
				}
			};
			vi.mocked(request).mockResolvedValue(user);

			const result = await authApi.me();

			expect(request).toHaveBeenCalledWith('/auth/me');
			expect(result).toEqual(user);
		});
	});

	describe('changePassword', () => {
		it('calls request with PUT /auth/password and correct body', async () => {
			vi.mocked(request).mockResolvedValue(undefined);

			await authApi.changePassword('oldPass', 'newPass');

			expect(request).toHaveBeenCalledWith('/auth/password', {
				method: 'PUT',
				body: JSON.stringify({ currentPassword: 'oldPass', newPassword: 'newPass' })
			});
		});
	});

	describe('updateEmailPreferences', () => {
		it('calls request with PUT /auth/email-preferences and returns updated preferences', async () => {
			const prefs: Partial<EmailPreferences> = { weeklyDigest: false };
			const responseData = {
				emailPreferences: {
					vacationUpdates: true,
					weeklyDigest: false,
					teamNotifications: true
				}
			};
			vi.mocked(request).mockResolvedValue(responseData);

			const result = await authApi.updateEmailPreferences(prefs);

			expect(request).toHaveBeenCalledWith('/auth/email-preferences', {
				method: 'PUT',
				body: JSON.stringify({ weeklyDigest: false })
			});
			expect(result).toEqual(responseData);
		});
	});
});
