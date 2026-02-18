import { describe, it, expect, vi, beforeEach } from 'vitest';
import { request, setAuthToken, clearAuthToken, ApiException } from './client';

function mockFetchResponse(
	body: unknown,
	init: { ok?: boolean; status?: number } = {}
): Response {
	const { ok = true, status = 200 } = init;
	return {
		ok,
		status,
		json: vi.fn().mockResolvedValue(body),
		headers: new Headers(),
		redirected: false,
		statusText: 'OK',
		type: 'basic',
		url: '',
		clone: vi.fn(),
		body: null,
		bodyUsed: false,
		arrayBuffer: vi.fn(),
		blob: vi.fn(),
		formData: vi.fn(),
		text: vi.fn(),
		bytes: vi.fn()
	} as unknown as Response;
}

describe('client', () => {
	beforeEach(() => {
		globalThis.fetch = vi.fn();
		localStorage.clear();
	});

	describe('request', () => {
		it('sends a GET request with correct URL and Content-Type header', async () => {
			const data = { id: 1, name: 'Test' };
			vi.mocked(globalThis.fetch).mockResolvedValue(mockFetchResponse(data));

			const result = await request('/users');

			expect(globalThis.fetch).toHaveBeenCalledWith('/api/users', {
				headers: { 'Content-Type': 'application/json' }
			});
			expect(result).toEqual(data);
		});

		it('sends a POST request with method and body', async () => {
			const payload = { name: 'New User' };
			const responseData = { id: 1, name: 'New User' };
			vi.mocked(globalThis.fetch).mockResolvedValue(mockFetchResponse(responseData));

			const result = await request('/users', {
				method: 'POST',
				body: JSON.stringify(payload)
			});

			expect(globalThis.fetch).toHaveBeenCalledWith('/api/users', {
				method: 'POST',
				body: JSON.stringify(payload),
				headers: { 'Content-Type': 'application/json' }
			});
			expect(result).toEqual(responseData);
		});

		it('includes Authorization header when auth token is set', async () => {
			localStorage.setItem('auth_token', 'my-jwt-token');
			vi.mocked(globalThis.fetch).mockResolvedValue(mockFetchResponse({ ok: true }));

			await request('/auth/me');

			expect(globalThis.fetch).toHaveBeenCalledWith('/api/auth/me', {
				headers: {
					'Content-Type': 'application/json',
					Authorization: 'Bearer my-jwt-token'
				}
			});
		});

		it('does not include Authorization header when no auth token exists', async () => {
			vi.mocked(globalThis.fetch).mockResolvedValue(mockFetchResponse({ ok: true }));

			await request('/auth/me');

			const callArgs = vi.mocked(globalThis.fetch).mock.calls[0];
			const headers = (callArgs[1] as RequestInit).headers as Record<string, string>;
			expect(headers).not.toHaveProperty('Authorization');
		});

		it('throws ApiException with correct fields on error response with JSON body', async () => {
			const errorBody = {
				code: 'VALIDATION_ERROR',
				message: 'Invalid input',
				details: { field: 'email' }
			};
			vi.mocked(globalThis.fetch).mockResolvedValue(
				mockFetchResponse(errorBody, { ok: false, status: 422 })
			);

			await expect(request('/users')).rejects.toThrow(ApiException);

			try {
				await request('/users');
			} catch (err) {
				const error = err as ApiException;
				expect(error.code).toBe('VALIDATION_ERROR');
				expect(error.message).toBe('Invalid input');
				expect(error.status).toBe(422);
				expect(error.details).toEqual({ field: 'email' });
				expect(error.name).toBe('ApiException');
			}
		});

		it('throws ApiException with UNKNOWN_ERROR when error response has non-JSON body', async () => {
			const response = mockFetchResponse(null, { ok: false, status: 500 });
			(response.json as ReturnType<typeof vi.fn>).mockRejectedValue(
				new Error('Invalid JSON')
			);
			vi.mocked(globalThis.fetch).mockResolvedValue(response);

			try {
				await request('/users');
			} catch (err) {
				const error = err as ApiException;
				expect(error).toBeInstanceOf(ApiException);
				expect(error.code).toBe('UNKNOWN_ERROR');
				expect(error.message).toBe('An unknown error occurred');
				expect(error.status).toBe(500);
			}
		});

		it('returns undefined for 204 No Content responses', async () => {
			vi.mocked(globalThis.fetch).mockResolvedValue(
				mockFetchResponse(null, { ok: true, status: 204 })
			);

			const result = await request('/users/1');

			expect(result).toBeUndefined();
		});

		it('propagates network errors from fetch', async () => {
			const networkError = new TypeError('Failed to fetch');
			vi.mocked(globalThis.fetch).mockRejectedValue(networkError);

			await expect(request('/users')).rejects.toThrow('Failed to fetch');
		});
	});

	describe('setAuthToken', () => {
		it('stores the token in localStorage', () => {
			setAuthToken('test-token-123');

			expect(localStorage.getItem('auth_token')).toBe('test-token-123');
		});
	});

	describe('clearAuthToken', () => {
		it('removes the token from localStorage', () => {
			localStorage.setItem('auth_token', 'existing-token');

			clearAuthToken();

			expect(localStorage.getItem('auth_token')).toBeNull();
		});
	});
});
