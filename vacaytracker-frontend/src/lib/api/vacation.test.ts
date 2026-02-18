import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('./client', () => ({
	request: vi.fn()
}));

import { request } from './client';
import { vacationApi } from './vacation';

const mockRequest = vi.mocked(request);

describe('vacationApi', () => {
	beforeEach(() => {
		mockRequest.mockReset();
	});

	describe('create', () => {
		it('sends POST to /vacation/request with body', async () => {
			const formData = {
				startDate: '15/06/2025',
				endDate: '20/06/2025',
				reason: 'Summer holiday'
			};
			const mockResponse = { id: 'req-1', status: 'pending', ...formData };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await vacationApi.create(formData);

			expect(mockRequest).toHaveBeenCalledWith('/vacation/request', {
				method: 'POST',
				body: JSON.stringify(formData)
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('list', () => {
		it('sends GET to /vacation/requests with no query string when no params', async () => {
			const mockResponse = { requests: [], total: 0 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await vacationApi.list();

			expect(mockRequest).toHaveBeenCalledWith('/vacation/requests');
			expect(result).toEqual(mockResponse);
		});

		it('sends GET with status query param', async () => {
			const mockResponse = { requests: [], total: 0 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			await vacationApi.list({ status: 'pending' });

			expect(mockRequest).toHaveBeenCalledWith('/vacation/requests?status=pending');
		});

		it('sends GET with year query param', async () => {
			const mockResponse = { requests: [], total: 0 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			await vacationApi.list({ year: 2025 });

			expect(mockRequest).toHaveBeenCalledWith('/vacation/requests?year=2025');
		});

		it('sends GET with both status and year query params', async () => {
			const mockResponse = { requests: [], total: 0 };
			mockRequest.mockResolvedValueOnce(mockResponse);

			await vacationApi.list({ status: 'approved', year: 2025 });

			expect(mockRequest).toHaveBeenCalledWith(
				'/vacation/requests?status=approved&year=2025'
			);
		});
	});

	describe('get', () => {
		it('sends GET to /vacation/requests/{id}', async () => {
			const mockResponse = { id: 'req-42', status: 'approved' };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await vacationApi.get('req-42');

			expect(mockRequest).toHaveBeenCalledWith('/vacation/requests/req-42');
			expect(result).toEqual(mockResponse);
		});
	});

	describe('cancel', () => {
		it('sends DELETE to /vacation/requests/{id}', async () => {
			mockRequest.mockResolvedValueOnce(undefined);

			await vacationApi.cancel('req-42');

			expect(mockRequest).toHaveBeenCalledWith('/vacation/requests/req-42', {
				method: 'DELETE'
			});
		});
	});

	describe('team', () => {
		it('sends GET to /vacation/team with month and year', async () => {
			const mockResponse = { vacations: [] };
			mockRequest.mockResolvedValueOnce(mockResponse);

			const result = await vacationApi.team(6, 2025);

			expect(mockRequest).toHaveBeenCalledWith('/vacation/team?month=6&year=2025');
			expect(result).toEqual(mockResponse);
		});
	});
});
