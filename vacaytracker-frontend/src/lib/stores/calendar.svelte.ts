import { calendarApi, type MonthKey } from '$lib/api/calendar';
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

	// Get all unique users from cached data (for filter panel)
	const availableUsers = $derived.by(() => {
		const userMap = new Map<string, UserInfo>();
		for (const [, cached] of monthCache) {
			for (const v of cached.vacations) {
				if (!userMap.has(v.userId)) {
					userMap.set(v.userId, { id: v.userId, name: v.userName });
				}
			}
		}
		return Array.from(userMap.values()).sort((a, b) => a.name.localeCompare(b.name));
	});

	// Get vacations for current view
	const currentViewVacations = $derived.by(() => {
		if (viewType === 'month') {
			const key = getCacheKey(currentMonth, currentYear);
			return monthCache.get(key)?.vacations ?? [];
		} else {
			// Week view might span two months
			const months = getMonthsForWeek(currentDate);
			const vacations: TeamVacation[] = [];
			const seen = new Set<string>();

			for (const m of months) {
				const key = getCacheKey(m.month, m.year);
				const cached = monthCache.get(key);
				if (cached) {
					for (const v of cached.vacations) {
						if (!seen.has(v.id)) {
							seen.add(v.id);
							vacations.push(v);
						}
					}
				}
			}
			return vacations;
		}
	});

	// Filtered vacations based on showAll flag and selectedUserIds
	const filteredVacations = $derived.by(() => {
		if (showAll) {
			return currentViewVacations;
		}
		// When not showing all, filter by selected users (empty = show nothing)
		return currentViewVacations.filter((v) => selectedUserIds.has(v.userId));
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

		try {
			const response = await calendarApi.getMonth(month, year);
			monthCache.set(key, {
				vacations: response.vacations,
				fetchedAt: Date.now()
			});
			// Trigger reactivity by reassigning
			monthCache = new Map(monthCache);
		} catch (e) {
			throw e;
		}
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

		const newSet = new Set(selectedUserIds);
		if (newSet.has(userId)) {
			newSet.delete(userId);
		} else {
			newSet.add(userId);
		}
		selectedUserIds = newSet;
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
		clearCache
	};
}

export const calendar = createCalendarStore();
