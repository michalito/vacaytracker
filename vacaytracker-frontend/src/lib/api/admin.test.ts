import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('./client', () => ({
	request: vi.fn()
}));

import { request } from './client';
import { adminApi } from './admin';

const mockRequest = vi.mocked(request);

describe('adminApi', () => {
	beforeEach(() => {
		mockRequest.mockReset();
	});

	describe('listUsers', () => {
		it('sends GET to /admin/users with no query string when no params', async () => {
			const mockResponse = { users: [], total: 0 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.listUsers();

			expect(mockRequest).toHaveBeenCalledWith('/admin/users');
			expect(result).toEqual(mockResponse);
		});

		it('sends GET with all query params', async () => {
			const mockResponse = { users: [], total: 0 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			await adminApi.listUsers({ page: 2, limit: 10, role: 'admin', search: 'john' });

			expect(mockRequest).toHaveBeenCalledWith(
				'/admin/users?page=2&limit=10&role=admin&search=john'
			);
		});
	});

	describe('createUser', () => {
		it('sends POST to /admin/users with body', async () => {
			const userData = { email: 'test@example.com', name: 'Test User', role: 'employee' };
			const mockResponse = { id: 'user-1', ...userData };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.createUser(userData);

			expect(mockRequest).toHaveBeenCalledWith('/admin/users', {
				method: 'POST',
				body: JSON.stringify(userData)
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('getUser', () => {
		it('sends GET to /admin/users/{id}', async () => {
			const mockResponse = { id: 'user-1', name: 'Test User' };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.getUser('user-1');

			expect(mockRequest).toHaveBeenCalledWith('/admin/users/user-1');
			expect(result).toEqual(mockResponse);
		});
	});

	describe('updateUser', () => {
		it('sends PUT to /admin/users/{id} with body', async () => {
			const updateData = { name: 'Updated Name' };
			const mockResponse = { id: 'user-1', name: 'Updated Name' };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.updateUser('user-1', updateData);

			expect(mockRequest).toHaveBeenCalledWith('/admin/users/user-1', {
				method: 'PUT',
				body: JSON.stringify(updateData)
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('deleteUser', () => {
		it('sends DELETE to /admin/users/{id}', async () => {
			mockRequest.mockResolvedValueOnce(undefined);

			await adminApi.deleteUser('user-1');

			expect(mockRequest).toHaveBeenCalledWith('/admin/users/user-1', {
				method: 'DELETE'
			});
		});
	});

	describe('updateBalance', () => {
		it('sends PUT to /admin/users/{id}/balance with vacationBalance body', async () => {
			const mockResponse = { id: 'user-1', vacationBalance: 25 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.updateBalance('user-1', 25);

			expect(mockRequest).toHaveBeenCalledWith('/admin/users/user-1/balance', {
				method: 'PUT',
				body: JSON.stringify({ vacationBalance: 25 })
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('resetAllBalances', () => {
		it('sends POST to /admin/users/reset-balances', async () => {
			const mockResponse = {
				success: true,
				usersUpdated: 10,
				newBalance: 25,
				message: 'Balances reset'
			};
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.resetAllBalances();

			expect(mockRequest).toHaveBeenCalledWith('/admin/users/reset-balances', {
				method: 'POST'
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('pendingRequests', () => {
		it('sends GET to /admin/vacation/pending', async () => {
			const mockResponse = { requests: [], total: 0 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.pendingRequests();

			expect(mockRequest).toHaveBeenCalledWith('/admin/vacation/pending');
			expect(result).toEqual(mockResponse);
		});
	});

	describe('reviewRequest', () => {
		it('sends PUT to /admin/vacation/{id}/review for approve', async () => {
			const mockResponse = { id: 'req-1', status: 'approved' };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.reviewRequest('req-1', 'approved');

			expect(mockRequest).toHaveBeenCalledWith('/admin/vacation/req-1/review', {
				method: 'PUT',
				body: JSON.stringify({ status: 'approved', reason: undefined })
			});
			expect(result).toEqual(mockResponse);
		});

		it('sends PUT to /admin/vacation/{id}/review for reject with reason', async () => {
			const mockResponse = { id: 'req-1', status: 'rejected' };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.reviewRequest('req-1', 'rejected', 'Team conflict');

			expect(mockRequest).toHaveBeenCalledWith('/admin/vacation/req-1/review', {
				method: 'PUT',
				body: JSON.stringify({ status: 'rejected', reason: 'Team conflict' })
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('getSettings', () => {
		it('sends GET to /admin/settings', async () => {
			const mockResponse = { excludeWeekends: true, defaultBalance: 25 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.getSettings();

			expect(mockRequest).toHaveBeenCalledWith('/admin/settings');
			expect(result).toEqual(mockResponse);
		});
	});

	describe('updateSettings', () => {
		it('sends PUT to /admin/settings with body', async () => {
			const settingsData = { excludeWeekends: false };
			const mockResponse = { excludeWeekends: false, defaultBalance: 25 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.updateSettings(settingsData);

			expect(mockRequest).toHaveBeenCalledWith('/admin/settings', {
				method: 'PUT',
				body: JSON.stringify(settingsData)
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('sendNewsletter', () => {
		it('sends POST to /admin/newsletter/send', async () => {
			const mockResponse = { success: true, sentCount: 5 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.sendNewsletter();

			expect(mockRequest).toHaveBeenCalledWith('/admin/newsletter/send', {
				method: 'POST'
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('getNewsletterPreview', () => {
		it('sends GET to /admin/newsletter/preview', async () => {
			const mockResponse = { html: '<p>Preview</p>', subject: 'Newsletter' };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.getNewsletterPreview();

			expect(mockRequest).toHaveBeenCalledWith('/admin/newsletter/preview');
			expect(result).toEqual(mockResponse);
		});
	});

	describe('sendTestEmail', () => {
		it('sends POST to /admin/email/test with template body', async () => {
			const mockResponse = { success: true, message: 'Test email sent' };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.sendTestEmail('welcome');

			expect(mockRequest).toHaveBeenCalledWith('/admin/email/test', {
				method: 'POST',
				body: JSON.stringify({ template: 'welcome' })
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('previewEmail', () => {
		it('sends POST to /admin/email/preview with template body', async () => {
			const mockResponse = { html: '<p>Email</p>', subject: 'Welcome' };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await adminApi.previewEmail('welcome');

			expect(mockRequest).toHaveBeenCalledWith('/admin/email/preview', {
				method: 'POST',
				body: JSON.stringify({ template: 'welcome' })
			});
			expect(result).toEqual(mockResponse);
		});
	});
});
