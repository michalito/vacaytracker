import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';

// Mock the calendar API -- hoisted above resetModules
vi.mock('$lib/api/calendar', () => ({
	calendarApi: {
		getMonth: vi.fn(),
		getMonths: vi.fn()
	}
}));

// Mock date utilities -- real implementations are pure functions, but we need
// getMonthsForWeek to return predictable values when testing week view
vi.mock('$lib/utils/date', async (importOriginal) => {
	const actual = (await importOriginal()) as Record<string, unknown>;
	return {
		...actual
	};
});

const mockVacation = (overrides = {}) => ({
	id: 'v1',
	userId: 'u1',
	userName: 'Alice',
	startDate: '2025-06-01',
	endDate: '2025-06-05',
	totalDays: 5,
	...overrides
});

describe('calendar store', () => {
	let calendar: (typeof import('$lib/stores/calendar.svelte'))['calendar'];
	let calendarApi: (typeof import('$lib/api/calendar'))['calendarApi'];

	beforeEach(async () => {
		vi.useFakeTimers();
		vi.setSystemTime(new Date(2025, 5, 15)); // June 15, 2025
		vi.clearAllMocks();
		vi.resetModules();

		const storeModule = await import('$lib/stores/calendar.svelte');
		const apiModule = await import('$lib/api/calendar');
		calendar = storeModule.calendar;
		calendarApi = apiModule.calendarApi;
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	describe('navigation', () => {
		it('has correct initial state', () => {
			expect(calendar.viewType).toBe('month');
			expect(calendar.currentMonth).toBe(6); // June
			expect(calendar.currentYear).toBe(2025);
			expect(calendar.isLoading).toBe(false);
			expect(calendar.error).toBeNull();
			expect(calendar.showAll).toBe(true);
		});

		it('goToNext in month view advances by 1 month', () => {
			expect(calendar.currentMonth).toBe(6);
			expect(calendar.currentYear).toBe(2025);

			calendar.goToNext();

			expect(calendar.currentMonth).toBe(7);
			expect(calendar.currentYear).toBe(2025);
		});

		it('goToPrevious in month view goes back by 1 month', () => {
			expect(calendar.currentMonth).toBe(6);

			calendar.goToPrevious();

			expect(calendar.currentMonth).toBe(5);
			expect(calendar.currentYear).toBe(2025);
		});

		it('goToNext in week view advances by 7 days', () => {
			calendar.setViewType('week');
			const dateBefore = new Date(calendar.currentDate);

			calendar.goToNext();

			const dateAfter = new Date(calendar.currentDate);
			const diffMs = dateAfter.getTime() - dateBefore.getTime();
			const diffDays = diffMs / (1000 * 60 * 60 * 24);
			expect(diffDays).toBe(7);
		});

		it('goToPrevious in week view goes back by 7 days', () => {
			calendar.setViewType('week');
			const dateBefore = new Date(calendar.currentDate);

			calendar.goToPrevious();

			const dateAfter = new Date(calendar.currentDate);
			const diffMs = dateBefore.getTime() - dateAfter.getTime();
			const diffDays = diffMs / (1000 * 60 * 60 * 24);
			expect(diffDays).toBe(7);
		});

		it('goToToday resets to current date', () => {
			// Navigate away first
			calendar.goToNext();
			calendar.goToNext();
			expect(calendar.currentMonth).toBe(8); // August

			calendar.goToToday();

			// Should be back to June 2025 (the faked system time)
			expect(calendar.currentMonth).toBe(6);
			expect(calendar.currentYear).toBe(2025);
		});

		it('goToDate sets a specific date', () => {
			const targetDate = new Date(2025, 11, 25); // Dec 25, 2025

			calendar.goToDate(targetDate);

			expect(calendar.currentMonth).toBe(12);
			expect(calendar.currentYear).toBe(2025);
			expect(calendar.currentDate.getDate()).toBe(25);
		});
	});

	describe('view switching', () => {
		it('setViewType changes the viewType', () => {
			expect(calendar.viewType).toBe('month');

			calendar.setViewType('week');

			expect(calendar.viewType).toBe('week');

			calendar.setViewType('month');

			expect(calendar.viewType).toBe('month');
		});
	});

	describe('filtering', () => {
		it('showAll defaults to true and filteredVacations returns all vacations', async () => {
			const vacations = [
				mockVacation({ id: 'v1', userId: 'u1', userName: 'Alice' }),
				mockVacation({ id: 'v2', userId: 'u2', userName: 'Bob' })
			];

			vi.mocked(calendarApi.getMonth).mockResolvedValue({
				vacations,
				month: 6,
				year: 2025
			});

			await calendar.ensureDataForCurrentView();

			expect(calendar.showAll).toBe(true);
			expect(calendar.filteredVacations).toHaveLength(2);
		});

		it('toggleUserFilter sets showAll false and adds userId to selectedUserIds', async () => {
			const vacations = [
				mockVacation({ id: 'v1', userId: 'u1', userName: 'Alice' }),
				mockVacation({ id: 'v2', userId: 'u2', userName: 'Bob' })
			];

			vi.mocked(calendarApi.getMonth).mockResolvedValue({
				vacations,
				month: 6,
				year: 2025
			});

			await calendar.ensureDataForCurrentView();

			calendar.toggleUserFilter('u1');

			expect(calendar.showAll).toBe(false);
			expect(calendar.selectedUserIds.has('u1')).toBe(true);
			expect(calendar.filteredVacations).toHaveLength(1);
			expect(calendar.filteredVacations[0].userId).toBe('u1');
		});

		it('toggleUserFilter twice removes the userId', () => {
			calendar.toggleUserFilter('u1');
			expect(calendar.selectedUserIds.has('u1')).toBe(true);

			calendar.toggleUserFilter('u1');
			expect(calendar.selectedUserIds.has('u1')).toBe(false);
		});

		it('clearFilters sets showAll true and clears selectedUserIds', () => {
			calendar.toggleUserFilter('u1');
			calendar.toggleUserFilter('u2');
			expect(calendar.showAll).toBe(false);
			expect(calendar.selectedUserIds.size).toBe(2);

			calendar.clearFilters();

			expect(calendar.showAll).toBe(true);
			expect(calendar.selectedUserIds.size).toBe(0);
		});

		it('selectNone sets showAll false with empty selectedUserIds resulting in empty filteredVacations', async () => {
			const vacations = [mockVacation({ id: 'v1', userId: 'u1' })];

			vi.mocked(calendarApi.getMonth).mockResolvedValue({
				vacations,
				month: 6,
				year: 2025
			});

			await calendar.ensureDataForCurrentView();

			calendar.selectNone();

			expect(calendar.showAll).toBe(false);
			expect(calendar.selectedUserIds.size).toBe(0);
			expect(calendar.filteredVacations).toHaveLength(0);
		});

		it('setSelectedUsers sets specific user IDs for filtering', async () => {
			const vacations = [
				mockVacation({ id: 'v1', userId: 'u1', userName: 'Alice' }),
				mockVacation({ id: 'v2', userId: 'u2', userName: 'Bob' }),
				mockVacation({ id: 'v3', userId: 'u3', userName: 'Charlie' })
			];

			vi.mocked(calendarApi.getMonth).mockResolvedValue({
				vacations,
				month: 6,
				year: 2025
			});

			await calendar.ensureDataForCurrentView();

			calendar.setSelectedUsers(['u1', 'u3']);

			expect(calendar.showAll).toBe(false);
			expect(calendar.selectedUserIds.has('u1')).toBe(true);
			expect(calendar.selectedUserIds.has('u3')).toBe(true);
			expect(calendar.selectedUserIds.has('u2')).toBe(false);
			expect(calendar.filteredVacations).toHaveLength(2);
		});
	});

	describe('data fetching', () => {
		it('ensureDataForCurrentView in month view calls calendarApi.getMonth with current month/year', async () => {
			vi.mocked(calendarApi.getMonth).mockResolvedValue({
				vacations: [mockVacation()],
				month: 6,
				year: 2025
			});

			await calendar.ensureDataForCurrentView();

			expect(calendarApi.getMonth).toHaveBeenCalledWith(6, 2025);
			expect(calendar.filteredVacations).toHaveLength(1);
			expect(calendar.isLoading).toBe(false);
			expect(calendar.error).toBeNull();
		});

		it('ensureDataForCurrentView uses cache on second call', async () => {
			vi.mocked(calendarApi.getMonth).mockResolvedValue({
				vacations: [mockVacation()],
				month: 6,
				year: 2025
			});

			await calendar.ensureDataForCurrentView();
			expect(calendarApi.getMonth).toHaveBeenCalledTimes(1);

			// Second call within cache TTL should not hit API
			await calendar.ensureDataForCurrentView();
			expect(calendarApi.getMonth).toHaveBeenCalledTimes(1);
		});

		it('clearCache causes next ensureDataForCurrentView to hit API again', async () => {
			vi.mocked(calendarApi.getMonth).mockResolvedValue({
				vacations: [mockVacation()],
				month: 6,
				year: 2025
			});

			await calendar.ensureDataForCurrentView();
			expect(calendarApi.getMonth).toHaveBeenCalledTimes(1);

			calendar.clearCache();

			await calendar.ensureDataForCurrentView();
			expect(calendarApi.getMonth).toHaveBeenCalledTimes(2);
		});

		it('ensureDataForCurrentView sets error on failure', async () => {
			vi.mocked(calendarApi.getMonth).mockRejectedValue(new Error('Network error'));

			await calendar.ensureDataForCurrentView();

			expect(calendar.error).toBe('Network error');
			expect(calendar.isLoading).toBe(false);
		});
	});
});
