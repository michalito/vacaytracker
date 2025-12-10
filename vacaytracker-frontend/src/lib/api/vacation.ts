import { request } from './client';
import type {
	VacationRequest,
	VacationListResponse,
	TeamVacationResponse,
	CreateVacationForm
} from '$lib/types';

export const vacationApi = {
	create: (data: CreateVacationForm): Promise<VacationRequest> => {
		return request('/vacation/request', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},

	list: (params?: { status?: string; year?: number }): Promise<VacationListResponse> => {
		const searchParams = new URLSearchParams();
		if (params?.status) searchParams.set('status', params.status);
		if (params?.year) searchParams.set('year', params.year.toString());

		const query = searchParams.toString();
		return request(`/vacation/requests${query ? `?${query}` : ''}`);
	},

	get: (id: string): Promise<VacationRequest> => {
		return request(`/vacation/requests/${id}`);
	},

	cancel: (id: string): Promise<void> => {
		return request(`/vacation/requests/${id}`, { method: 'DELETE' });
	},

	team: (month: number, year: number): Promise<TeamVacationResponse> => {
		return request(`/vacation/team?month=${month}&year=${year}`);
	}
};
