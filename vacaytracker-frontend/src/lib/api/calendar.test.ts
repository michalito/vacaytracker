import { describe, it, expect, vi, beforeEach } from 'vitest';
import type { TeamVacationResponse } from '$lib/types';

vi.mock('./vacation', () => ({
	vacationApi: {
		team: vi.fn()
	}
}));

import { vacationApi } from './vacation';
import { calendarApi } from './calendar';

const mockTeam = vi.mocked(vacationApi.team);

describe('calendarApi', () => {
	beforeEach(() => {
		mockTeam.mockReset();
	});

	describe('getMonth', () => {
		it('delegates to vacationApi.team with month and year', async () => {
			const mockResponse: TeamVacationResponse = { vacations: [] };
			mockTeam.mockResolvedValueOnce(mockResponse);

			const result = await calendarApi.getMonth(6, 2025);

			expect(mockTeam).toHaveBeenCalledWith(6, 2025);
			expect(result).toEqual(mockResponse);
		});
	});

	describe('getMonths', () => {
		it('fetches a single month and returns vacations', async () => {
			const vacation = {
				id: 'v-1',
				userName: 'Alice',
				startDate: '2025-06-01',
				endDate: '2025-06-05',
				status: 'approved'
			};
			mockTeam.mockResolvedValueOnce({ vacations: [vacation] });

			const result = await calendarApi.getMonths([{ month: 6, year: 2025 }]);

			expect(mockTeam).toHaveBeenCalledTimes(1);
			expect(mockTeam).toHaveBeenCalledWith(6, 2025);
			expect(result).toEqual([vacation]);
		});

		it('fetches multiple months and merges vacations', async () => {
			const vacation1 = {
				id: 'v-1',
				userName: 'Alice',
				startDate: '2025-06-01',
				endDate: '2025-06-05',
				status: 'approved'
			};
			const vacation2 = {
				id: 'v-2',
				userName: 'Bob',
				startDate: '2025-07-10',
				endDate: '2025-07-15',
				status: 'pending'
			};
			mockTeam.mockResolvedValueOnce({ vacations: [vacation1] });
			mockTeam.mockResolvedValueOnce({ vacations: [vacation2] });

			const result = await calendarApi.getMonths([
				{ month: 6, year: 2025 },
				{ month: 7, year: 2025 }
			]);

			expect(mockTeam).toHaveBeenCalledTimes(2);
			expect(mockTeam).toHaveBeenCalledWith(6, 2025);
			expect(mockTeam).toHaveBeenCalledWith(7, 2025);
			expect(result).toEqual([vacation1, vacation2]);
		});

		it('deduplicates vacations with the same id across months', async () => {
			const sharedVacation = {
				id: 'v-shared',
				userName: 'Alice',
				startDate: '2025-06-28',
				endDate: '2025-07-05',
				status: 'approved'
			};
			const uniqueVacation = {
				id: 'v-unique',
				userName: 'Bob',
				startDate: '2025-07-10',
				endDate: '2025-07-12',
				status: 'approved'
			};
			mockTeam.mockResolvedValueOnce({ vacations: [sharedVacation] });
			mockTeam.mockResolvedValueOnce({ vacations: [sharedVacation, uniqueVacation] });

			const result = await calendarApi.getMonths([
				{ month: 6, year: 2025 },
				{ month: 7, year: 2025 }
			]);

			expect(result).toHaveLength(2);
			expect(result).toEqual([sharedVacation, uniqueVacation]);
		});

		it('returns empty array when given empty months array', async () => {
			const result = await calendarApi.getMonths([]);

			expect(mockTeam).not.toHaveBeenCalled();
			expect(result).toEqual([]);
		});
	});
});
