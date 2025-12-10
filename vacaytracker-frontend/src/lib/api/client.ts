import type { ApiError } from '$lib/types';

const API_BASE = '/api';

export class ApiException extends Error {
	constructor(
		public code: string,
		message: string,
		public status: number,
		public details?: Record<string, unknown>
	) {
		super(message);
		this.name = 'ApiException';
	}
}

function getAuthToken(): string | null {
	if (typeof window === 'undefined') return null;
	return localStorage.getItem('auth_token');
}

export function setAuthToken(token: string): void {
	localStorage.setItem('auth_token', token);
}

export function clearAuthToken(): void {
	localStorage.removeItem('auth_token');
}

export async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
	const token = getAuthToken();

	const headers: HeadersInit = {
		'Content-Type': 'application/json',
		...options.headers
	};

	if (token) {
		(headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
	}

	const response = await fetch(`${API_BASE}${endpoint}`, {
		...options,
		headers
	});

	if (!response.ok) {
		let error: ApiError;
		try {
			error = await response.json();
		} catch {
			error = {
				code: 'UNKNOWN_ERROR',
				message: 'An unknown error occurred'
			};
		}
		throw new ApiException(error.code, error.message, response.status, error.details);
	}

	// Handle 204 No Content
	if (response.status === 204) {
		return undefined as T;
	}

	return response.json();
}
