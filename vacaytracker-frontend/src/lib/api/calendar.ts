import { vacationApi } from './vacation';
import type { TeamVacation, TeamVacationResponse } from '$lib/types';

export interface MonthKey {
	month: number;
	year: number;
}

export const calendarApi = {
	// Fetch team vacations for a single month
	getMonth: (month: number, year: number): Promise<TeamVacationResponse> => {
		return vacationApi.team(month, year);
	},

	// Fetch team vacations for multiple months and merge results
	getMonths: async (months: MonthKey[]): Promise<TeamVacation[]> => {
		const responses = await Promise.all(
			months.map((m) => vacationApi.team(m.month, m.year))
		);

		// Merge and deduplicate by vacation id
		const vacationMap = new Map<string, TeamVacation>();
		for (const response of responses) {
			for (const vacation of response.vacations) {
				vacationMap.set(vacation.id, vacation);
			}
		}

		return Array.from(vacationMap.values());
	}
};
