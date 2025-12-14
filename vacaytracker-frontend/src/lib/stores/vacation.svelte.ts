import { vacationApi } from '$lib/api/vacation';
import type { VacationRequest, VacationStatus } from '$lib/types';

const CACHE_TTL = 5 * 60 * 1000; // 5 minutes

function createVacationStore() {
	let requests = $state<VacationRequest[]>([]);
	let isLoading = $state(false);
	let error = $state<string | null>(null);
	let cachedAt = $state<number | null>(null);

	const pendingRequests = $derived(requests.filter((r) => r.status === 'pending'));
	const approvedRequests = $derived(requests.filter((r) => r.status === 'approved'));
	const rejectedRequests = $derived(requests.filter((r) => r.status === 'rejected'));
	const totalDaysUsed = $derived(
		approvedRequests.reduce((sum, r) => sum + r.totalDays, 0)
	);

	function isCacheValid(): boolean {
		if (!cachedAt) return false;
		return Date.now() - cachedAt < CACHE_TTL;
	}

	async function fetchRequests(
		status?: VacationStatus,
		year?: number,
		options?: { force?: boolean }
	): Promise<void> {
		// Skip if cache is valid and not forcing refresh (only for unfiltered requests)
		if (!options?.force && isCacheValid() && !status && !year) {
			return;
		}

		isLoading = true;
		error = null;

		try {
			const response = await vacationApi.list({ status, year });
			requests = response.requests;
			// Only update cache timestamp for unfiltered requests
			if (!status && !year) {
				cachedAt = Date.now();
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to fetch requests';
		} finally {
			isLoading = false;
		}
	}

	async function createRequest(data: {
		startDate: string;
		endDate: string;
		reason?: string;
	}): Promise<VacationRequest> {
		const newRequest = await vacationApi.create(data);
		requests = [newRequest, ...requests];
		return newRequest;
	}

	async function cancelRequest(id: string): Promise<void> {
		await vacationApi.cancel(id);
		requests = requests.filter((r) => r.id !== id);
	}

	function updateRequest(id: string, updates: Partial<VacationRequest>): void {
		requests = requests.map((r) => (r.id === id ? { ...r, ...updates } : r));
	}

	function invalidateCache(): void {
		cachedAt = null;
	}

	return {
		get requests() {
			return requests;
		},
		get pendingRequests() {
			return pendingRequests;
		},
		get approvedRequests() {
			return approvedRequests;
		},
		get rejectedRequests() {
			return rejectedRequests;
		},
		get isLoading() {
			return isLoading;
		},
		get error() {
			return error;
		},
		get totalDaysUsed() {
			return totalDaysUsed;
		},
		fetchRequests,
		createRequest,
		cancelRequest,
		updateRequest,
		invalidateCache
	};
}

export const vacation = createVacationStore();
