import { describe, it, expect, vi, beforeEach } from 'vitest';
import type { TeamVacation, TeamVacationResponse } from '#lib/types/index.js';

vi.mock('./vacation', () => ({
	vacationApi: {
		team: vi.fn()
	}
}));

import { vacationApi } from './vacation';
import { calendarApi } from './calendar';

const mockTeam = vi.mocked(vacationApi.team);

function makeVacation(overrides: Partial<TeamVacation> & Pick<TeamVacation, 'id'>): TeamVacation {
	return {
		userId: `user-${overrides.id}`,
		userName: 'Alice',
		startDate: '2025-06-01',
		endDate: '2025-06-05',
		totalDays: 5,
		...overrides
	};
}

function teamResponse(month: number, vacations: TeamVacation[]): TeamVacationResponse {
	return { vacations, month, year: 2025 };
}

describe('calendarApi', () => {
	beforeEach(() => {
		mockTeam.mockReset();
	});

	describe('getMonth', () => {
		it('delegates to vacationApi.team with month and year', async () => {
			const mockResponse = teamResponse(6, []);
			mockTeam.mockResolvedValueOnce(mockResponse);

			const result = await calendarApi.getMonth(6, 2025);

			expect(mockTeam).toHaveBeenCalledWith(6, 2025);
			expect(result).toEqual(mockResponse);
		});
	});

	describe('getMonths', () => {
		it('fetches a single month and returns vacations', async () => {
			const vacation = makeVacation({ id: 'v-1' });
			mockTeam.mockResolvedValueOnce(teamResponse(6, [vacation]));

			const result = await calendarApi.getMonths([{ month: 6, year: 2025 }]);

			expect(mockTeam).toHaveBeenCalledTimes(1);
			expect(mockTeam).toHaveBeenCalledWith(6, 2025);
			expect(result).toEqual([vacation]);
		});

		it('fetches multiple months and merges vacations', async () => {
			const vacation1 = makeVacation({ id: 'v-1' });
			const vacation2 = makeVacation({
				id: 'v-2',
				userName: 'Bob',
				startDate: '2025-07-10',
				endDate: '2025-07-15',
				totalDays: 4
			});
			mockTeam.mockResolvedValueOnce(teamResponse(6, [vacation1]));
			mockTeam.mockResolvedValueOnce(teamResponse(7, [vacation2]));

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
			const sharedVacation = makeVacation({
				id: 'v-shared',
				startDate: '2025-06-28',
				endDate: '2025-07-05',
				totalDays: 6
			});
			const uniqueVacation = makeVacation({
				id: 'v-unique',
				userName: 'Bob',
				startDate: '2025-07-10',
				endDate: '2025-07-12',
				totalDays: 3
			});
			mockTeam.mockResolvedValueOnce(teamResponse(6, [sharedVacation]));
			mockTeam.mockResolvedValueOnce(teamResponse(7, [sharedVacation, uniqueVacation]));

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
