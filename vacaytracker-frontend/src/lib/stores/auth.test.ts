import { describe, it, expect, beforeEach, vi } from 'vitest';
import type { User, LoginResponse } from '$lib/types';

// Mock the API module -- this persists across resetModules because vi.mock is hoisted
vi.mock('$lib/api/auth', () => ({
	authApi: {
		login: vi.fn(),
		logout: vi.fn(),
		me: vi.fn(),
		changePassword: vi.fn(),
		updateEmailPreferences: vi.fn()
	}
}));

const mockUser: User = {
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

const mockAdmin: User = {
	...mockUser,
	id: '2',
	email: 'admin@example.com',
	name: 'Admin User',
	role: 'admin'
};

describe('auth store', () => {
	let auth: (typeof import('$lib/stores/auth.svelte'))['auth'];
	let authApi: (typeof import('$lib/api/auth'))['authApi'];

	beforeEach(async () => {
		vi.resetModules();
		localStorage.clear();

		const storeModule = await import('$lib/stores/auth.svelte');
		const apiModule = await import('$lib/api/auth');
		auth = storeModule.auth;
		authApi = apiModule.authApi;
	});

	describe('initial state', () => {
		it('has correct defaults before any actions', () => {
			expect(auth.user).toBeNull();
			expect(auth.isLoading).toBe(true);
			expect(auth.error).toBeNull();
			expect(auth.isAuthenticated).toBe(false);
			expect(auth.isAdmin).toBe(false);
			expect(auth.isEmployee).toBe(false);
		});
	});

	describe('initialize', () => {
		it('sets isLoading to false when no token is in localStorage', async () => {
			await auth.initialize();

			expect(auth.isLoading).toBe(false);
			expect(auth.user).toBeNull();
			expect(auth.isAuthenticated).toBe(false);
		});

		it('fetches user when a valid token exists in localStorage', async () => {
			localStorage.setItem('auth_token', 'valid-jwt-token');
			vi.mocked(authApi.me).mockResolvedValue(mockUser);

			await auth.initialize();

			expect(authApi.me).toHaveBeenCalled();
			expect(auth.user).toEqual(mockUser);
			expect(auth.isAuthenticated).toBe(true);
			expect(auth.isLoading).toBe(false);
		});

		it('clears token and sets user to null when token is expired/invalid', async () => {
			localStorage.setItem('auth_token', 'expired-token');
			vi.mocked(authApi.me).mockRejectedValue(new Error('Unauthorized'));

			await auth.initialize();

			expect(localStorage.getItem('auth_token')).toBeNull();
			expect(auth.user).toBeNull();
			expect(auth.isAuthenticated).toBe(false);
			expect(auth.isLoading).toBe(false);
		});
	});

	describe('login', () => {
		it('sets user and derived state on successful login', async () => {
			const loginResponse: LoginResponse = { token: 'new-token', user: mockUser };
			vi.mocked(authApi.login).mockResolvedValue(loginResponse);

			await auth.login('test@example.com', 'password123');

			expect(authApi.login).toHaveBeenCalledWith('test@example.com', 'password123');
			expect(auth.user).toEqual(mockUser);
			expect(auth.isAuthenticated).toBe(true);
			expect(auth.isLoading).toBe(false);
			expect(auth.error).toBeNull();
		});

		it('sets isAdmin true when logging in as admin', async () => {
			const loginResponse: LoginResponse = { token: 'admin-token', user: mockAdmin };
			vi.mocked(authApi.login).mockResolvedValue(loginResponse);

			await auth.login('admin@example.com', 'adminpass');

			expect(auth.isAdmin).toBe(true);
			expect(auth.isEmployee).toBe(false);
		});

		it('sets error, keeps user null, and re-throws on login failure', async () => {
			const loginError = new Error('Invalid credentials');
			vi.mocked(authApi.login).mockRejectedValue(loginError);

			await expect(auth.login('bad@example.com', 'wrong')).rejects.toThrow(
				'Invalid credentials'
			);

			expect(auth.error).toBe('Invalid credentials');
			expect(auth.user).toBeNull();
			expect(auth.isAuthenticated).toBe(false);
			expect(auth.isLoading).toBe(false);
		});

		it('sets isLoading to true during the login call', async () => {
			// Create a deferred promise so we can check state mid-call
			let resolveLogin!: (value: LoginResponse) => void;
			const loginPromise = new Promise<LoginResponse>((resolve) => {
				resolveLogin = resolve;
			});
			vi.mocked(authApi.login).mockReturnValue(loginPromise);

			const loginCall = auth.login('test@example.com', 'password123');

			// While the promise is pending, isLoading should be true
			expect(auth.isLoading).toBe(true);

			// Resolve and let it complete
			resolveLogin({ token: 'tok', user: mockUser });
			await loginCall;

			expect(auth.isLoading).toBe(false);
		});
	});

	describe('logout', () => {
		it('calls authApi.logout and resets user to null', async () => {
			// First login
			const loginResponse: LoginResponse = { token: 'token', user: mockUser };
			vi.mocked(authApi.login).mockResolvedValue(loginResponse);
			await auth.login('test@example.com', 'password123');
			expect(auth.isAuthenticated).toBe(true);

			auth.logout();

			expect(authApi.logout).toHaveBeenCalled();
			expect(auth.user).toBeNull();
			expect(auth.isAuthenticated).toBe(false);
		});
	});

	describe('updateUser', () => {
		it('merges partial updates into the current user', async () => {
			const loginResponse: LoginResponse = { token: 'token', user: mockUser };
			vi.mocked(authApi.login).mockResolvedValue(loginResponse);
			await auth.login('test@example.com', 'password123');

			auth.updateUser({ name: 'Updated Name', vacationBalance: 15 });

			expect(auth.user).toMatchObject({
				id: '1',
				email: 'test@example.com',
				name: 'Updated Name',
				vacationBalance: 15
			});
		});

		it('does nothing when no user is set', () => {
			expect(auth.user).toBeNull();

			// Should not throw
			auth.updateUser({ name: 'Ghost' });

			expect(auth.user).toBeNull();
		});
	});

	describe('derived properties', () => {
		it('isAuthenticated is true when user is set, false when null', async () => {
			expect(auth.isAuthenticated).toBe(false);

			const loginResponse: LoginResponse = { token: 'token', user: mockUser };
			vi.mocked(authApi.login).mockResolvedValue(loginResponse);
			await auth.login('test@example.com', 'password123');

			expect(auth.isAuthenticated).toBe(true);

			auth.logout();

			expect(auth.isAuthenticated).toBe(false);
		});

		it('isAdmin is true only for admin role', async () => {
			const loginResponse: LoginResponse = { token: 'token', user: mockAdmin };
			vi.mocked(authApi.login).mockResolvedValue(loginResponse);
			await auth.login('admin@example.com', 'adminpass');

			expect(auth.isAdmin).toBe(true);
			expect(auth.isEmployee).toBe(false);
		});

		it('isEmployee is true only for employee role', async () => {
			const loginResponse: LoginResponse = { token: 'token', user: mockUser };
			vi.mocked(authApi.login).mockResolvedValue(loginResponse);
			await auth.login('test@example.com', 'password123');

			expect(auth.isEmployee).toBe(true);
			expect(auth.isAdmin).toBe(false);
		});
	});
});
