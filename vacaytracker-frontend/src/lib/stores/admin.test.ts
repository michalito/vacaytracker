import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import type { User, VacationRequest, Settings, NewsletterPreview } from '$lib/types';

// Mock the admin API -- hoisted above resetModules
vi.mock('$lib/api/admin', () => ({
	adminApi: {
		listUsers: vi.fn(),
		createUser: vi.fn(),
		getUser: vi.fn(),
		updateUser: vi.fn(),
		deleteUser: vi.fn(),
		updateBalance: vi.fn(),
		resetAllBalances: vi.fn(),
		pendingRequests: vi.fn(),
		reviewRequest: vi.fn(),
		getSettings: vi.fn(),
		updateSettings: vi.fn(),
		sendNewsletter: vi.fn(),
		getNewsletterPreview: vi.fn(),
		sendTestEmail: vi.fn(),
		previewEmail: vi.fn()
	}
}));

const mockUser = (overrides: Partial<User> = {}): User => ({
	id: '1',
	email: 'user@test.com',
	name: 'Test User',
	role: 'employee',
	vacationBalance: 20,
	emailPreferences: {
		vacationUpdates: true,
		weeklyDigest: true,
		teamNotifications: true
	},
	...overrides
});

const mockVacationRequest = (overrides: Partial<VacationRequest> = {}): VacationRequest => ({
	id: 'req1',
	userId: '1',
	userName: 'Test User',
	startDate: '2025-06-01',
	endDate: '2025-06-05',
	totalDays: 5,
	status: 'pending',
	createdAt: '2025-05-20T10:00:00Z',
	updatedAt: '2025-05-20T10:00:00Z',
	...overrides
});

const mockSettings: Settings = {
	id: 'settings-1',
	weekendPolicy: { excludeWeekends: true, excludedDays: [0, 6] },
	newsletter: { enabled: true, frequency: 'weekly', dayOfMonth: 1 },
	defaultVacationDays: 25,
	vacationResetMonth: 1,
	updatedAt: '2025-01-01T00:00:00Z'
};

describe('admin store', () => {
	let admin: (typeof import('$lib/stores/admin.svelte'))['admin'];
	let adminApi: (typeof import('$lib/api/admin'))['adminApi'];

	beforeEach(async () => {
		vi.useFakeTimers();
		vi.setSystemTime(new Date(2025, 5, 15)); // June 15, 2025
		vi.clearAllMocks();
		vi.resetModules();

		const storeModule = await import('$lib/stores/admin.svelte');
		const apiModule = await import('$lib/api/admin');
		admin = storeModule.admin;
		adminApi = apiModule.adminApi;
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	describe('fetchUsers', () => {
		it('fetches and stores users with pagination', async () => {
			const users = [mockUser(), mockUser({ id: '2', name: 'User Two' })];
			const pagination = { page: 1, limit: 10, total: 2, totalPages: 1 };

			vi.mocked(adminApi.listUsers).mockResolvedValue({ users, pagination });

			await admin.fetchUsers();

			expect(adminApi.listUsers).toHaveBeenCalledTimes(1);
			expect(admin.users).toHaveLength(2);
			expect(admin.users[0].name).toBe('Test User');
			expect(admin.pagination).toEqual(pagination);
			expect(admin.isLoading).toBe(false);
			expect(admin.error).toBeNull();
		});

		it('uses cache on second call without force', async () => {
			const users = [mockUser()];
			const pagination = { page: 1, limit: 10, total: 1, totalPages: 1 };

			vi.mocked(adminApi.listUsers).mockResolvedValue({ users, pagination });

			await admin.fetchUsers();
			expect(adminApi.listUsers).toHaveBeenCalledTimes(1);

			// Second call should be served from cache
			await admin.fetchUsers();
			expect(adminApi.listUsers).toHaveBeenCalledTimes(1);
		});

		it('force option bypasses cache', async () => {
			const users = [mockUser()];
			const pagination = { page: 1, limit: 10, total: 1, totalPages: 1 };

			vi.mocked(adminApi.listUsers).mockResolvedValue({ users, pagination });

			await admin.fetchUsers();
			expect(adminApi.listUsers).toHaveBeenCalledTimes(1);

			await admin.fetchUsers(undefined, { force: true });
			expect(adminApi.listUsers).toHaveBeenCalledTimes(2);
		});

		it('filtered requests with search/role bypass cache', async () => {
			const users = [mockUser()];
			const pagination = { page: 1, limit: 10, total: 1, totalPages: 1 };

			vi.mocked(adminApi.listUsers).mockResolvedValue({ users, pagination });

			// First default fetch
			await admin.fetchUsers();
			expect(adminApi.listUsers).toHaveBeenCalledTimes(1);

			// Filtered request always calls API even though default cache is valid
			await admin.fetchUsers({ search: 'alice' });
			expect(adminApi.listUsers).toHaveBeenCalledTimes(2);

			// Role filter also bypasses cache
			await admin.fetchUsers({ role: 'admin' });
			expect(adminApi.listUsers).toHaveBeenCalledTimes(3);
		});

		it('handles error and clears users', async () => {
			vi.mocked(adminApi.listUsers).mockRejectedValue(new Error('Network failure'));

			await admin.fetchUsers();

			expect(admin.error).toBe('Network failure');
			expect(admin.users).toEqual([]);
			expect(admin.isLoading).toBe(false);
		});
	});

	describe('fetchPendingRequests', () => {
		it('fetches and stores pending requests', async () => {
			const requests = [
				mockVacationRequest({ id: 'req1' }),
				mockVacationRequest({ id: 'req2', userId: '2' })
			];

			vi.mocked(adminApi.pendingRequests).mockResolvedValue({
				requests,
				total: 2
			});

			await admin.fetchPendingRequests();

			expect(adminApi.pendingRequests).toHaveBeenCalledTimes(1);
			expect(admin.pendingRequests).toHaveLength(2);
			expect(admin.error).toBeNull();
		});

		it('uses cache on second call', async () => {
			vi.mocked(adminApi.pendingRequests).mockResolvedValue({
				requests: [mockVacationRequest()],
				total: 1
			});

			await admin.fetchPendingRequests();
			expect(adminApi.pendingRequests).toHaveBeenCalledTimes(1);

			await admin.fetchPendingRequests();
			expect(adminApi.pendingRequests).toHaveBeenCalledTimes(1);
		});

		it('handles error and clears pendingRequests', async () => {
			vi.mocked(adminApi.pendingRequests).mockRejectedValue(
				new Error('Failed to fetch pending requests')
			);

			await admin.fetchPendingRequests();

			expect(admin.error).toBe('Failed to fetch pending requests');
			expect(admin.pendingRequests).toEqual([]);
		});
	});

	describe('fetchSettings', () => {
		it('fetches and stores settings', async () => {
			vi.mocked(adminApi.getSettings).mockResolvedValue(mockSettings);

			await admin.fetchSettings();

			expect(adminApi.getSettings).toHaveBeenCalledTimes(1);
			expect(admin.settings).toEqual(mockSettings);
			expect(admin.error).toBeNull();
		});

		it('handles error and sets settings to null', async () => {
			vi.mocked(adminApi.getSettings).mockRejectedValue(
				new Error('Failed to fetch settings')
			);

			await admin.fetchSettings();

			expect(admin.error).toBe('Failed to fetch settings');
			expect(admin.settings).toBeNull();
		});
	});

	describe('CRUD operations', () => {
		it('createUser calls API and prepends user to list', async () => {
			const existingUser = mockUser({ id: '1', name: 'Existing' });
			const newUser = mockUser({ id: '2', name: 'New User', email: 'new@test.com' });

			// Set up existing users first
			vi.mocked(adminApi.listUsers).mockResolvedValue({
				users: [existingUser],
				pagination: { page: 1, limit: 10, total: 1, totalPages: 1 }
			});
			await admin.fetchUsers();

			vi.mocked(adminApi.createUser).mockResolvedValue(newUser);

			const result = await admin.createUser({
				email: 'new@test.com',
				password: 'pass123',
				name: 'New User',
				role: 'employee'
			});

			expect(adminApi.createUser).toHaveBeenCalledWith({
				email: 'new@test.com',
				password: 'pass123',
				name: 'New User',
				role: 'employee'
			});
			expect(result).toEqual(newUser);
			// New user should be prepended
			expect(admin.users).toHaveLength(2);
			expect(admin.users[0].id).toBe('2');
			expect(admin.users[1].id).toBe('1');
		});

		it('updateUser calls API and updates user in list', async () => {
			const user = mockUser({ id: '1', name: 'Original' });
			const updatedUser = mockUser({ id: '1', name: 'Updated' });

			vi.mocked(adminApi.listUsers).mockResolvedValue({
				users: [user],
				pagination: { page: 1, limit: 10, total: 1, totalPages: 1 }
			});
			await admin.fetchUsers();

			vi.mocked(adminApi.updateUser).mockResolvedValue(updatedUser);

			const result = await admin.updateUser('1', { name: 'Updated' });

			expect(adminApi.updateUser).toHaveBeenCalledWith('1', { name: 'Updated' });
			expect(result).toEqual(updatedUser);
			expect(admin.users[0].name).toBe('Updated');
		});

		it('deleteUser calls API and removes user from list', async () => {
			const user1 = mockUser({ id: '1', name: 'User One' });
			const user2 = mockUser({ id: '2', name: 'User Two' });

			vi.mocked(adminApi.listUsers).mockResolvedValue({
				users: [user1, user2],
				pagination: { page: 1, limit: 10, total: 2, totalPages: 1 }
			});
			await admin.fetchUsers();
			expect(admin.users).toHaveLength(2);

			vi.mocked(adminApi.deleteUser).mockResolvedValue(undefined);

			await admin.deleteUser('1');

			expect(adminApi.deleteUser).toHaveBeenCalledWith('1');
			expect(admin.users).toHaveLength(1);
			expect(admin.users[0].id).toBe('2');
		});

		it('updateBalance calls API and updates user in list', async () => {
			const user = mockUser({ id: '1', vacationBalance: 20 });
			const updatedUser = mockUser({ id: '1', vacationBalance: 25 });

			vi.mocked(adminApi.listUsers).mockResolvedValue({
				users: [user],
				pagination: { page: 1, limit: 10, total: 1, totalPages: 1 }
			});
			await admin.fetchUsers();

			vi.mocked(adminApi.updateBalance).mockResolvedValue(updatedUser);

			const result = await admin.updateBalance('1', 25);

			expect(adminApi.updateBalance).toHaveBeenCalledWith('1', 25);
			expect(result).toEqual(updatedUser);
			expect(admin.users[0].vacationBalance).toBe(25);
		});
	});

	describe('reviewRequest', () => {
		it('calls API and removes request from pendingRequests', async () => {
			const requests = [
				mockVacationRequest({ id: 'req1' }),
				mockVacationRequest({ id: 'req2' })
			];

			vi.mocked(adminApi.pendingRequests).mockResolvedValue({
				requests,
				total: 2
			});
			await admin.fetchPendingRequests();
			expect(admin.pendingRequests).toHaveLength(2);

			const reviewed = mockVacationRequest({ id: 'req1', status: 'approved' });
			vi.mocked(adminApi.reviewRequest).mockResolvedValue(reviewed);

			const result = await admin.reviewRequest('req1', 'approved');

			expect(adminApi.reviewRequest).toHaveBeenCalledWith('req1', 'approved', undefined);
			expect(result).toEqual(reviewed);
			expect(admin.pendingRequests).toHaveLength(1);
			expect(admin.pendingRequests[0].id).toBe('req2');
		});

		it('passes reason when rejecting a request', async () => {
			const requests = [mockVacationRequest({ id: 'req1' })];

			vi.mocked(adminApi.pendingRequests).mockResolvedValue({
				requests,
				total: 1
			});
			await admin.fetchPendingRequests();

			const reviewed = mockVacationRequest({
				id: 'req1',
				status: 'rejected',
				rejectionReason: 'Team capacity'
			});
			vi.mocked(adminApi.reviewRequest).mockResolvedValue(reviewed);

			const result = await admin.reviewRequest('req1', 'rejected', 'Team capacity');

			expect(adminApi.reviewRequest).toHaveBeenCalledWith('req1', 'rejected', 'Team capacity');
			expect(result.status).toBe('rejected');
		});
	});

	describe('updateSettings', () => {
		it('calls API and updates settings in store', async () => {
			const updatedSettings: Settings = {
				...mockSettings,
				defaultVacationDays: 30
			};

			vi.mocked(adminApi.updateSettings).mockResolvedValue(updatedSettings);

			const result = await admin.updateSettings({ defaultVacationDays: 30 });

			expect(adminApi.updateSettings).toHaveBeenCalledWith({ defaultVacationDays: 30 });
			expect(result).toEqual(updatedSettings);
			expect(admin.settings).toEqual(updatedSettings);
		});
	});

	describe('newsletter', () => {
		it('sendNewsletter calls API, refreshes settings, and manages isSendingNewsletter flag', async () => {
			const sendResponse = {
				success: true,
				recipientCount: 5,
				message: 'Sent to 5 recipients'
			};

			vi.mocked(adminApi.sendNewsletter).mockResolvedValue(sendResponse);
			vi.mocked(adminApi.getSettings).mockResolvedValue(mockSettings);

			const result = await admin.sendNewsletter();

			expect(adminApi.sendNewsletter).toHaveBeenCalledTimes(1);
			// sendNewsletter should also refresh settings afterwards
			expect(adminApi.getSettings).toHaveBeenCalledTimes(1);
			expect(result).toEqual(sendResponse);
			expect(admin.isSendingNewsletter).toBe(false);
		});

		it('fetchNewsletterPreview populates newsletterPreview', async () => {
			const preview: NewsletterPreview = {
				subject: 'Weekly Update',
				htmlBody: '<h1>Hello</h1>',
				textBody: 'Hello',
				recipients: ['a@test.com', 'b@test.com'],
				recipientCount: 2
			};

			vi.mocked(adminApi.getNewsletterPreview).mockResolvedValue(preview);

			const result = await admin.fetchNewsletterPreview();

			expect(adminApi.getNewsletterPreview).toHaveBeenCalledTimes(1);
			expect(result).toEqual(preview);
			expect(admin.newsletterPreview).toEqual(preview);
		});
	});

	describe('cache invalidation', () => {
		it('invalidateCache without arguments invalidates all caches', async () => {
			// Populate all caches
			vi.mocked(adminApi.listUsers).mockResolvedValue({
				users: [mockUser()],
				pagination: { page: 1, limit: 10, total: 1, totalPages: 1 }
			});
			vi.mocked(adminApi.pendingRequests).mockResolvedValue({
				requests: [mockVacationRequest()],
				total: 1
			});
			vi.mocked(adminApi.getSettings).mockResolvedValue(mockSettings);

			await admin.fetchUsers();
			await admin.fetchPendingRequests();
			await admin.fetchSettings();

			expect(adminApi.listUsers).toHaveBeenCalledTimes(1);
			expect(adminApi.pendingRequests).toHaveBeenCalledTimes(1);
			expect(adminApi.getSettings).toHaveBeenCalledTimes(1);

			// Invalidate all
			admin.invalidateCache();

			// All subsequent fetches should hit the API
			await admin.fetchUsers();
			await admin.fetchPendingRequests();
			await admin.fetchSettings();

			expect(adminApi.listUsers).toHaveBeenCalledTimes(2);
			expect(adminApi.pendingRequests).toHaveBeenCalledTimes(2);
			expect(adminApi.getSettings).toHaveBeenCalledTimes(2);
		});

		it('invalidateCache with specific keys only invalidates those caches', async () => {
			// Populate all caches
			vi.mocked(adminApi.listUsers).mockResolvedValue({
				users: [mockUser()],
				pagination: { page: 1, limit: 10, total: 1, totalPages: 1 }
			});
			vi.mocked(adminApi.pendingRequests).mockResolvedValue({
				requests: [mockVacationRequest()],
				total: 1
			});
			vi.mocked(adminApi.getSettings).mockResolvedValue(mockSettings);

			await admin.fetchUsers();
			await admin.fetchPendingRequests();
			await admin.fetchSettings();

			// Invalidate only usersAt
			admin.invalidateCache(['usersAt']);

			await admin.fetchUsers();
			await admin.fetchPendingRequests();
			await admin.fetchSettings();

			// Users should have been re-fetched (2 calls total)
			expect(adminApi.listUsers).toHaveBeenCalledTimes(2);
			// Pending requests and settings should still be cached (1 call total)
			expect(adminApi.pendingRequests).toHaveBeenCalledTimes(1);
			expect(adminApi.getSettings).toHaveBeenCalledTimes(1);
		});
	});
});
