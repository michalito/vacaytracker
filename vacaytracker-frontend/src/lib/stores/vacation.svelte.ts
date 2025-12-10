import { vacationApi } from '$lib/api/vacation';
import type { VacationRequest, VacationStatus } from '$lib/types';

function createVacationStore() {
	let requests = $state<VacationRequest[]>([]);
	let isLoading = $state(false);
	let error = $state<string | null>(null);

	const pendingRequests = $derived(requests.filter((r) => r.status === 'pending'));
	const approvedRequests = $derived(requests.filter((r) => r.status === 'approved'));
	const rejectedRequests = $derived(requests.filter((r) => r.status === 'rejected'));

	async function fetchRequests(status?: VacationStatus, year?: number): Promise<void> {
		isLoading = true;
		error = null;

		try {
			const response = await vacationApi.list({ status, year });
			requests = response.requests;
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
		fetchRequests,
		createRequest,
		cancelRequest,
		updateRequest
	};
}

export const vacation = createVacationStore();
