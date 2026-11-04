import { calendarApi } from '$lib/api/calendar';
import type { TeamVacation } from '$lib/types';
import { addWeeks, addMonths, getMonthsForWeek } from '$lib/utils/date';

export type CalendarView = 'week' | 'month';

export interface CachedMonth {
	vacations: TeamVacation[];
	fetchedAt: number;
}

export interface UserInfo {
	id: string;
	name: string;
}

const CACHE_TTL = 5 * 60 * 1000; // 5 minutes

function createCalendarStore() {
	// State
	let currentDate = $state(new Date());
	let viewType = $state<CalendarView>('month');
	// Filter state: showAll=true means show everyone, showAll=false uses selectedUserIds
	let showAll = $state(true);
	let selectedUserIds = $state<Set<string>>(new Set());
	let monthCache = $state<Map<string, CachedMonth>>(new Map());
	let isLoading = $state(false);
	let error = $state<string | null>(null);

	// Derived values
	const currentMonth = $derived(currentDate.getMonth() + 1);
	const currentYear = $derived(currentDate.getFullYear());

	// Collections below are rebuilt (never mutated) so plain Map/Set stay correct under runes;
	// state updates reassign a new instance instead.

	// Get all unique users from cached data (for filter panel)
	const availableUsers = $derived.by(() => {
		const users = new Map<string, UserInfo>(
			[...monthCache.values()]
				.flatMap((cached) => cached.vacations)
				.map((v) => [v.userId, { id: v.userId, name: v.userName }])
		);
		return [...users.values()].sort((a, b) => a.name.localeCompare(b.name));
	});

	// Get vacations for current view
	const currentViewVacations = $derived.by(() => {
		if (viewType === 'month') {
			const key = getCacheKey(currentMonth, currentYear);
			return monthCache.get(key)?.vacations ?? [];
		}

		// Week view might span two months; a vacation crossing the boundary appears in both
		const byId = new Map<string, TeamVacation>(
			getMonthsForWeek(currentDate)
				.flatMap((m) => monthCache.get(getCacheKey(m.month, m.year))?.vacations ?? [])
				.map((v) => [v.id, v])
		);
		return [...byId.values()];
	});

	// Filtered vacations based on showAll flag and selectedUserIds
	const filteredVacations = $derived.by(() => {
		if (showAll) {
			return currentViewVacations;
		}
		// When not showing all, filter by selected users (empty = show nothing)
		return currentViewVacations.filter((v) => selectedUserIds.has(v.userId));
	});

	// Current month vacations (for dashboard)
	const currentMonthVacations = $derived.by(() => {
		const now = new Date();
		const key = getCacheKey(now.getMonth() + 1, now.getFullYear());
		return monthCache.get(key)?.vacations ?? [];
	});

	// Cache helpers
	function getCacheKey(month: number, year: number): string {
		return `${year}-${month}`;
	}

	function isCacheValid(key: string): boolean {
		const cached = monthCache.get(key);
		if (!cached) return false;
		return Date.now() - cached.fetchedAt < CACHE_TTL;
	}

	// Fetch a single month and cache it
	async function fetchMonth(month: number, year: number): Promise<void> {
		const key = getCacheKey(month, year);
		if (isCacheValid(key)) return;

		const response = await calendarApi.getMonth(month, year);
		// Reassign a new Map so dependents re-run (Maps in $state are not deeply reactive)
		monthCache = new Map([
			...monthCache,
			[key, { vacations: response.vacations, fetchedAt: Date.now() }]
		]);
	}

	// Ensure data is loaded for current view
	async function ensureDataForCurrentView(): Promise<void> {
		isLoading = true;
		error = null;

		try {
			if (viewType === 'month') {
				await fetchMonth(currentMonth, currentYear);
			} else {
				// Week view might span two months
				const months = getMonthsForWeek(currentDate);
				await Promise.all(months.map((m) => fetchMonth(m.month, m.year)));
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load calendar data';
		} finally {
			isLoading = false;
		}
	}

	// Navigation
	function goToToday(): void {
		currentDate = new Date();
	}

	function goToPrevious(): void {
		if (viewType === 'week') {
			currentDate = addWeeks(currentDate, -1);
		} else {
			currentDate = addMonths(currentDate, -1);
		}
	}

	function goToNext(): void {
		if (viewType === 'week') {
			currentDate = addWeeks(currentDate, 1);
		} else {
			currentDate = addMonths(currentDate, 1);
		}
	}

	function goToDate(date: Date): void {
		currentDate = new Date(date);
	}

	function setViewType(view: CalendarView): void {
		viewType = view;
	}

	// Filtering - simple and clear
	function setShowAll(value: boolean): void {
		showAll = value;
		if (value) {
			selectedUserIds = new Set();
		}
	}

	function toggleUserFilter(userId: string): void {
		// Always switch to custom filtering mode
		showAll = false;

		const ids = selectedUserIds.has(userId)
			? [...selectedUserIds].filter((id) => id !== userId)
			: [...selectedUserIds, userId];
		selectedUserIds = new Set(ids);
	}

	function clearFilters(): void {
		showAll = true;
		selectedUserIds = new Set();
	}

	function selectNone(): void {
		showAll = false;
		selectedUserIds = new Set();
	}

	function setSelectedUsers(userIds: string[]): void {
		showAll = false;
		selectedUserIds = new Set(userIds);
	}

	// Clear cache (useful for forcing refresh)
	function clearCache(): void {
		monthCache = new Map();
	}

	// Fetch current month (convenience method for dashboard)
	async function fetchCurrentMonth(): Promise<void> {
		const now = new Date();
		await fetchMonth(now.getMonth() + 1, now.getFullYear());
	}

	return {
		// Getters
		get currentDate() {
			return currentDate;
		},
		get viewType() {
			return viewType;
		},
		get showAll() {
			return showAll;
		},
		get currentMonth() {
			return currentMonth;
		},
		get currentYear() {
			return currentYear;
		},
		get isLoading() {
			return isLoading;
		},
		get error() {
			return error;
		},
		get availableUsers() {
			return availableUsers;
		},
		get filteredVacations() {
			return filteredVacations;
		},
		get selectedUserIds() {
			return selectedUserIds;
		},
		get currentMonthVacations() {
			return currentMonthVacations;
		},

		// Actions
		goToToday,
		goToPrevious,
		goToNext,
		goToDate,
		setViewType,
		setShowAll,
		toggleUserFilter,
		clearFilters,
		selectNone,
		setSelectedUsers,
		ensureDataForCurrentView,
		clearCache,
		fetchCurrentMonth
	};
}

export const calendar = createCalendarStore();
