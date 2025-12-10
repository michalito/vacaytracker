import { authApi } from '$lib/api/auth';
import type { User } from '$lib/types';

function createAuthStore() {
	let user = $state<User | null>(null);
	let isLoading = $state(true);
	let error = $state<string | null>(null);

	const isAuthenticated = $derived(user !== null);
	const isAdmin = $derived(user?.role === 'admin');
	const isEmployee = $derived(user?.role === 'employee');

	async function initialize(): Promise<void> {
		if (typeof window === 'undefined') return;

		const token = localStorage.getItem('auth_token');
		if (!token) {
			isLoading = false;
			return;
		}

		try {
			user = await authApi.me();
		} catch {
			localStorage.removeItem('auth_token');
			user = null;
		} finally {
			isLoading = false;
		}
	}

	async function login(email: string, password: string): Promise<void> {
		error = null;
		isLoading = true;

		try {
			const response = await authApi.login(email, password);
			user = response.user;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Login failed';
			throw e;
		} finally {
			isLoading = false;
		}
	}

	function logout(): void {
		authApi.logout();
		user = null;
	}

	function updateUser(updates: Partial<User>): void {
		if (user) {
			user = { ...user, ...updates };
		}
	}

	return {
		get user() {
			return user;
		},
		get isAuthenticated() {
			return isAuthenticated;
		},
		get isAdmin() {
			return isAdmin;
		},
		get isEmployee() {
			return isEmployee;
		},
		get isLoading() {
			return isLoading;
		},
		get error() {
			return error;
		},
		initialize,
		login,
		logout,
		updateUser
	};
}

export const auth = createAuthStore();
