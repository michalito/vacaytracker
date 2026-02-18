import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import type { VacationRequest } from '$lib/types';

vi.mock('$lib/api/vacation', () => ({
	vacationApi: {
		list: vi.fn(),
		create: vi.fn(),
		cancel: vi.fn(),
		get: vi.fn(),
		team: vi.fn()
	}
}));

function mockRequest(overrides: Partial<VacationRequest> = {}): VacationRequest {
	return {
		id: '1',
		userId: 'user-1',
		startDate: '2025-06-01',
		endDate: '2025-06-05',
		totalDays: 5,
		status: 'pending',
		createdAt: '2025-05-01T00:00:00Z',
		updatedAt: '2025-05-01T00:00:00Z',
		...overrides
	};
}

describe('vacation store', () => {
	let vacation: (typeof import('$lib/stores/vacation.svelte'))['vacation'];
	let vacationApi: (typeof import('$lib/api/vacation'))['vacationApi'];

	beforeEach(async () => {
		vi.useFakeTimers();
		vi.setSystemTime(new Date(2025, 5, 15)); // June 15, 2025
		vi.resetModules();

		const storeModule = await import('$lib/stores/vacation.svelte');
		const apiModule = await import('$lib/api/vacation');
		vacation = storeModule.vacation;
		vacationApi = apiModule.vacationApi;

		// Clear call history accumulated from previous tests
		vi.mocked(vacationApi.list).mockClear();
		vi.mocked(vacationApi.create).mockClear();
		vi.mocked(vacationApi.cancel).mockClear();
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	// ── fetchRequests ──────────────────────────────────────────────

	describe('fetchRequests', () => {
		it('fetches and stores requests', async () => {
			const reqs = [mockRequest(), mockRequest({ id: '2', status: 'approved' })];
			vi.mocked(vacationApi.list).mockResolvedValueOnce({ requests: reqs, total: 2 });

			await vacation.fetchRequests();

			expect(vacationApi.list).toHaveBeenCalledOnce();
			expect(vacation.requests).toEqual(reqs);
		});

		it('sets isLoading during fetch', async () => {
			let resolveList!: (value: { requests: VacationRequest[]; total: number }) => void;
			vi.mocked(vacationApi.list).mockReturnValueOnce(
				new Promise((resolve) => {
					resolveList = resolve;
				})
			);

			const fetchPromise = vacation.fetchRequests();

			expect(vacation.isLoading).toBe(true);

			resolveList({ requests: [], total: 0 });
			await fetchPromise;

			expect(vacation.isLoading).toBe(false);
		});

		it('handles error', async () => {
			vi.mocked(vacationApi.list).mockRejectedValueOnce(new Error('Network failure'));

			await vacation.fetchRequests();

			expect(vacation.error).toBe('Network failure');
			expect(vacation.isLoading).toBe(false);
		});

		it('uses cache on second call', async () => {
			vi.mocked(vacationApi.list).mockResolvedValueOnce({ requests: [], total: 0 });

			await vacation.fetchRequests();
			await vacation.fetchRequests();

			expect(vacationApi.list).toHaveBeenCalledTimes(1);
		});

		it('cache expires after TTL', async () => {
			vi.mocked(vacationApi.list).mockResolvedValue({ requests: [], total: 0 });

			await vacation.fetchRequests();
			expect(vacationApi.list).toHaveBeenCalledTimes(1);

			// Advance past the 5-minute TTL
			vi.advanceTimersByTime(5 * 60 * 1000 + 1);

			await vacation.fetchRequests();
			expect(vacationApi.list).toHaveBeenCalledTimes(2);
		});

		it('force bypasses cache', async () => {
			vi.mocked(vacationApi.list).mockResolvedValue({ requests: [], total: 0 });

			await vacation.fetchRequests();
			await vacation.fetchRequests(undefined, undefined, { force: true });

			expect(vacationApi.list).toHaveBeenCalledTimes(2);
		});

		it('filtered requests bypass cache', async () => {
			vi.mocked(vacationApi.list).mockResolvedValue({ requests: [], total: 0 });

			await vacation.fetchRequests();
			await vacation.fetchRequests('pending');
			await vacation.fetchRequests(undefined, 2025);

			expect(vacationApi.list).toHaveBeenCalledTimes(3);
		});

		it('filtered requests do not update cache', async () => {
			vi.mocked(vacationApi.list).mockResolvedValue({ requests: [], total: 0 });

			// Fetch with a status filter — should not set cachedAt
			await vacation.fetchRequests('pending');

			// Next unfiltered call should still hit the API since cache was not set
			await vacation.fetchRequests();

			expect(vacationApi.list).toHaveBeenCalledTimes(2);
		});
	});

	// ── Derived lists ──────────────────────────────────────────────

	describe('derived lists', () => {
		it('pendingRequests filters status === pending', async () => {
			const pending = mockRequest({ id: '1', status: 'pending' });
			const approved = mockRequest({ id: '2', status: 'approved' });
			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [pending, approved],
				total: 2
			});

			await vacation.fetchRequests();

			expect(vacation.pendingRequests).toEqual([pending]);
		});

		it('approvedRequests filters status === approved', async () => {
			const pending = mockRequest({ id: '1', status: 'pending' });
			const approved = mockRequest({ id: '2', status: 'approved' });
			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [pending, approved],
				total: 2
			});

			await vacation.fetchRequests();

			expect(vacation.approvedRequests).toEqual([approved]);
		});

		it('rejectedRequests filters status === rejected', async () => {
			const rejected = mockRequest({ id: '1', status: 'rejected' });
			const approved = mockRequest({ id: '2', status: 'approved' });
			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [rejected, approved],
				total: 2
			});

			await vacation.fetchRequests();

			expect(vacation.rejectedRequests).toEqual([rejected]);
		});
	});

	// ── Date-dependent derived (system time: June 15, 2025) ────────

	describe('date-dependent derived values', () => {
		// today = '2025-06-15', currentYear = 2025

		it('usedDays sums approved requests in 2025 where startDate <= today', async () => {
			const pastApproved = mockRequest({
				id: '1',
				status: 'approved',
				startDate: '2025-06-01',
				totalDays: 3
			});
			const currentApproved = mockRequest({
				id: '2',
				status: 'approved',
				startDate: '2025-06-15',
				totalDays: 2
			});
			const futureApproved = mockRequest({
				id: '3',
				status: 'approved',
				startDate: '2025-07-01',
				totalDays: 4
			});
			const lastYearApproved = mockRequest({
				id: '4',
				status: 'approved',
				startDate: '2024-06-01',
				totalDays: 10
			});

			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [pastApproved, currentApproved, futureApproved, lastYearApproved],
				total: 4
			});

			await vacation.fetchRequests();

			// pastApproved (3) + currentApproved (2) = 5  (startDate <= '2025-06-15')
			// futureApproved excluded (startDate > today)
			// lastYearApproved excluded (year 2024 !== 2025)
			expect(vacation.usedDays).toBe(5);
		});

		it('upcomingDays sums approved requests in 2025 where startDate > today', async () => {
			const pastApproved = mockRequest({
				id: '1',
				status: 'approved',
				startDate: '2025-06-01',
				totalDays: 3
			});
			const futureApproved = mockRequest({
				id: '2',
				status: 'approved',
				startDate: '2025-07-01',
				totalDays: 4
			});
			const anotherFuture = mockRequest({
				id: '3',
				status: 'approved',
				startDate: '2025-08-10',
				totalDays: 6
			});

			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [pastApproved, futureApproved, anotherFuture],
				total: 3
			});

			await vacation.fetchRequests();

			// futureApproved (4) + anotherFuture (6) = 10
			expect(vacation.upcomingDays).toBe(10);
		});

		it('totalDaysUsed equals usedDays + upcomingDays', async () => {
			const past = mockRequest({
				id: '1',
				status: 'approved',
				startDate: '2025-05-01',
				totalDays: 3
			});
			const future = mockRequest({
				id: '2',
				status: 'approved',
				startDate: '2025-09-01',
				totalDays: 7
			});

			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [past, future],
				total: 2
			});

			await vacation.fetchRequests();

			expect(vacation.totalDaysUsed).toBe(vacation.usedDays + vacation.upcomingDays);
			expect(vacation.totalDaysUsed).toBe(10);
		});

		it('upcomingRequests returns approved with startDate > today sorted ascending', async () => {
			const futureB = mockRequest({
				id: '1',
				status: 'approved',
				startDate: '2025-08-01'
			});
			const futureA = mockRequest({
				id: '2',
				status: 'approved',
				startDate: '2025-07-01'
			});
			const past = mockRequest({
				id: '3',
				status: 'approved',
				startDate: '2025-05-01',
				endDate: '2025-05-05'
			});

			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [futureB, futureA, past],
				total: 3
			});

			await vacation.fetchRequests();

			expect(vacation.upcomingRequests).toHaveLength(2);
			expect(vacation.upcomingRequests[0].id).toBe('2'); // July before August
			expect(vacation.upcomingRequests[1].id).toBe('1');
		});

		it('pastRequests returns approved with endDate < today plus rejected, sorted by createdAt desc', async () => {
			const pastApproved = mockRequest({
				id: '1',
				status: 'approved',
				startDate: '2025-03-01',
				endDate: '2025-03-05',
				createdAt: '2025-02-01T00:00:00Z'
			});
			const rejected = mockRequest({
				id: '2',
				status: 'rejected',
				startDate: '2025-07-01',
				endDate: '2025-07-10',
				createdAt: '2025-06-01T00:00:00Z'
			});
			const ongoingApproved = mockRequest({
				id: '3',
				status: 'approved',
				startDate: '2025-06-10',
				endDate: '2025-06-20'
			});

			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [pastApproved, rejected, ongoingApproved],
				total: 3
			});

			await vacation.fetchRequests();

			// pastApproved: endDate '2025-03-05' < '2025-06-15' -> included
			// rejected: always included in pastRequests
			// ongoingApproved: endDate '2025-06-20' >= '2025-06-15' -> excluded
			expect(vacation.pastRequests).toHaveLength(2);
			// Sorted by createdAt desc: rejected (2025-06-01) before pastApproved (2025-02-01)
			expect(vacation.pastRequests[0].id).toBe('2');
			expect(vacation.pastRequests[1].id).toBe('1');
		});
	});

	// ── Mutations ──────────────────────────────────────────────────

	describe('mutations', () => {
		it('createRequest calls API and prepends to requests', async () => {
			const existing = mockRequest({ id: '1' });
			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [existing],
				total: 1
			});
			await vacation.fetchRequests();

			const newReq = mockRequest({ id: '2', startDate: '2025-07-01', endDate: '2025-07-05' });
			vi.mocked(vacationApi.create).mockResolvedValueOnce(newReq);

			const result = await vacation.createRequest({
				startDate: '2025-07-01',
				endDate: '2025-07-05'
			});

			expect(vacationApi.create).toHaveBeenCalledWith({
				startDate: '2025-07-01',
				endDate: '2025-07-05'
			});
			expect(result).toEqual(newReq);
			expect(vacation.requests[0]).toEqual(newReq);
			expect(vacation.requests).toHaveLength(2);
		});

		it('cancelRequest calls API and removes from requests', async () => {
			const req1 = mockRequest({ id: '1' });
			const req2 = mockRequest({ id: '2' });
			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [req1, req2],
				total: 2
			});
			await vacation.fetchRequests();

			vi.mocked(vacationApi.cancel).mockResolvedValueOnce(undefined);

			await vacation.cancelRequest('1');

			expect(vacationApi.cancel).toHaveBeenCalledWith('1');
			expect(vacation.requests).toHaveLength(1);
			expect(vacation.requests[0].id).toBe('2');
		});

		it('updateRequest merges partial update into matching request', async () => {
			const req = mockRequest({ id: '1', status: 'pending' });
			vi.mocked(vacationApi.list).mockResolvedValueOnce({
				requests: [req],
				total: 1
			});
			await vacation.fetchRequests();

			vacation.updateRequest('1', { status: 'approved', reviewedBy: 'admin-1' });

			expect(vacation.requests[0].status).toBe('approved');
			expect(vacation.requests[0].reviewedBy).toBe('admin-1');
			// Original fields preserved
			expect(vacation.requests[0].startDate).toBe('2025-06-01');
		});

		it('invalidateCache causes next fetchRequests to hit API', async () => {
			vi.mocked(vacationApi.list).mockResolvedValue({ requests: [], total: 0 });

			await vacation.fetchRequests();
			expect(vacationApi.list).toHaveBeenCalledTimes(1);

			// Without invalidation, cache would prevent another call
			vacation.invalidateCache();
			await vacation.fetchRequests();

			expect(vacationApi.list).toHaveBeenCalledTimes(2);
		});
	});
});
