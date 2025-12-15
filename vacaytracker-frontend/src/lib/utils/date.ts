// Date utility types
export interface DateRange {
	start: Date;
	end: Date;
}

export interface CalendarDay {
	date: Date;
	dateString: string; // YYYY-MM-DD format
	dayOfMonth: number;
	dayOfWeek: number; // 0=Monday, 6=Sunday (EU week)
	isToday: boolean;
	isCurrentMonth: boolean;
	isWeekend: boolean;
}

// Parse YYYY-MM-DD string to Date
export function parseISODate(dateString: string): Date {
	const [year, month, day] = dateString.split('-').map(Number);
	return new Date(year, month - 1, day);
}

// Format Date to YYYY-MM-DD string
export function formatDateISO(date: Date): string {
	const year = date.getFullYear();
	const month = String(date.getMonth() + 1).padStart(2, '0');
	const day = String(date.getDate()).padStart(2, '0');
	return `${year}-${month}-${day}`;
}

// Format "December 2025"
export function formatMonthYear(date: Date, locale = 'en-US'): string {
	return date.toLocaleDateString(locale, { month: 'long', year: 'numeric' });
}

// Format "Dec 1 - Dec 7, 2025"
export function formatWeekRange(startDate: Date, locale = 'en-US'): string {
	const endDate = addDays(startDate, 6);
	const startMonth = startDate.toLocaleDateString(locale, { month: 'short' });
	const endMonth = endDate.toLocaleDateString(locale, { month: 'short' });
	const startDay = startDate.getDate();
	const endDay = endDate.getDate();
	const year = endDate.getFullYear();

	if (startMonth === endMonth) {
		return `${startMonth} ${startDay} - ${endDay}, ${year}`;
	}
	return `${startMonth} ${startDay} - ${endMonth} ${endDay}, ${year}`;
}

// Add days to a date
export function addDays(date: Date, days: number): Date {
	const result = new Date(date);
	result.setDate(result.getDate() + days);
	return result;
}

// Add weeks to a date
export function addWeeks(date: Date, weeks: number): Date {
	return addDays(date, weeks * 7);
}

// Add months to a date
export function addMonths(date: Date, months: number): Date {
	const result = new Date(date);
	result.setMonth(result.getMonth() + months);
	return result;
}

// Check if two dates are the same day
export function isSameDay(a: Date, b: Date): boolean {
	return (
		a.getFullYear() === b.getFullYear() &&
		a.getMonth() === b.getMonth() &&
		a.getDate() === b.getDate()
	);
}

// Check if a date is within a range (inclusive)
export function isDateInRange(date: Date, start: Date, end: Date): boolean {
	const d = new Date(date.getFullYear(), date.getMonth(), date.getDate()).getTime();
	const s = new Date(start.getFullYear(), start.getMonth(), start.getDate()).getTime();
	const e = new Date(end.getFullYear(), end.getMonth(), end.getDate()).getTime();
	return d >= s && d <= e;
}

// Get Monday of the week containing the given date (EU standard: week starts Monday)
export function getMonday(date: Date): Date {
	const d = new Date(date);
	const day = d.getDay();
	// Sunday is 0, so we need to go back 6 days; otherwise go back (day - 1) days
	const diff = day === 0 ? -6 : 1 - day;
	d.setDate(d.getDate() + diff);
	return d;
}

// Create a CalendarDay object
function createCalendarDay(date: Date, currentMonth: number): CalendarDay {
	const today = new Date();
	const dayOfWeek = date.getDay();
	// Convert Sunday=0..Saturday=6 to Monday=0..Sunday=6
	const euDayOfWeek = dayOfWeek === 0 ? 6 : dayOfWeek - 1;

	return {
		date: new Date(date),
		dateString: formatDateISO(date),
		dayOfMonth: date.getDate(),
		dayOfWeek: euDayOfWeek,
		isToday: isSameDay(date, today),
		isCurrentMonth: date.getMonth() === currentMonth,
		isWeekend: euDayOfWeek >= 5 // Saturday=5, Sunday=6
	};
}

// Get 7 days for week view (Monday-Sunday)
export function getWeekDays(date: Date): CalendarDay[] {
	const monday = getMonday(date);
	const currentMonth = date.getMonth();
	const days: CalendarDay[] = [];

	for (let i = 0; i < 7; i++) {
		const d = addDays(monday, i);
		days.push(createCalendarDay(d, currentMonth));
	}

	return days;
}

// Get week boundaries (Monday to Sunday)
export function getWeekBoundaries(date: Date): DateRange {
	const monday = getMonday(date);
	const sunday = addDays(monday, 6);
	return { start: monday, end: sunday };
}

// Get month boundaries
export function getMonthBoundaries(year: number, month: number): DateRange {
	const start = new Date(year, month - 1, 1);
	const end = new Date(year, month, 0); // Last day of month
	return { start, end };
}

// Get 42 days (6 weeks) for month view, starting from Monday
export function getCalendarGridDays(year: number, month: number): CalendarDay[] {
	// First day of the month
	const firstOfMonth = new Date(year, month - 1, 1);
	// Find the Monday of the week containing the first day
	const monday = getMonday(firstOfMonth);

	const days: CalendarDay[] = [];
	const currentMonth = month - 1; // 0-indexed for comparison

	for (let i = 0; i < 42; i++) {
		const d = addDays(monday, i);
		days.push(createCalendarDay(d, currentMonth));
	}

	return days;
}

// Get ISO week number (EU standard)
export function getWeekNumber(date: Date): number {
	const d = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
	const dayNum = d.getUTCDay() || 7;
	d.setUTCDate(d.getUTCDate() + 4 - dayNum);
	const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
	return Math.ceil(((d.getTime() - yearStart.getTime()) / 86400000 + 1) / 7);
}

// Get the months needed for a week view (handles month boundaries)
export function getMonthsForWeek(date: Date): { month: number; year: number }[] {
	const { start, end } = getWeekBoundaries(date);
	const months: { month: number; year: number }[] = [];

	const startMonth = start.getMonth() + 1;
	const startYear = start.getFullYear();
	months.push({ month: startMonth, year: startYear });

	const endMonth = end.getMonth() + 1;
	const endYear = end.getFullYear();

	if (startMonth !== endMonth || startYear !== endYear) {
		months.push({ month: endMonth, year: endYear });
	}

	return months;
}

// ============================================
// UI Display Formatting Functions
// ============================================

/**
 * Format date string to short display format: "15 Jan"
 * @param dateStr - ISO date string (YYYY-MM-DD) or any valid date string
 * @param locale - Locale for formatting (default: 'en-GB')
 */
export function formatDateShort(dateStr: string, locale = 'en-GB'): string {
	return new Date(dateStr).toLocaleDateString(locale, {
		day: 'numeric',
		month: 'short'
	});
}

/**
 * Format date string to medium display format: "15 Jan 2025"
 * @param dateStr - ISO date string (YYYY-MM-DD) or any valid date string
 * @param locale - Locale for formatting (default: 'en-GB')
 */
export function formatDateMedium(dateStr: string, locale = 'en-GB'): string {
	return new Date(dateStr).toLocaleDateString(locale, {
		day: 'numeric',
		month: 'short',
		year: 'numeric'
	});
}

/**
 * Format date string to long display format: "15 January 2025"
 * @param dateStr - ISO date string (YYYY-MM-DD) or any valid date string
 * @param locale - Locale for formatting (default: 'en-GB')
 */
export function formatDateLong(dateStr: string | undefined, locale = 'en-GB'): string {
	if (!dateStr) return 'Not set';
	return new Date(dateStr).toLocaleDateString(locale, {
		day: 'numeric',
		month: 'long',
		year: 'numeric'
	});
}

/**
 * Format date range to display format: "15 Jan - 20 Jan"
 * @param startDate - Start date string
 * @param endDate - End date string
 * @param locale - Locale for formatting (default: 'en-GB')
 */
export function formatDateRangeShort(startDate: string, endDate: string, locale = 'en-GB'): string {
	const startStr = formatDateShort(startDate, locale);
	const endStr = formatDateShort(endDate, locale);
	return `${startStr} - ${endStr}`;
}

/**
 * Convert HTML date input (YYYY-MM-DD) to EU API format (DD/MM/YYYY)
 * @param dateStr - Date in YYYY-MM-DD format (from HTML date input)
 */
export function toApiDateFormat(dateStr: string): string {
	const [year, month, day] = dateStr.split('-');
	return `${day}/${month}/${year}`;
}

/**
 * Convert EU format (DD/MM/YYYY) to ISO format (YYYY-MM-DD)
 * @param dateStr - Date in DD/MM/YYYY format
 */
export function fromApiDateFormat(dateStr: string): string {
	const [day, month, year] = dateStr.split('/');
	return `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}`;
}

// ============================================
// @internationalized/date Integration
// ============================================

import type { DateValue } from '@internationalized/date';
import { getLocalTimeZone } from '@internationalized/date';

// ============================================
// Vacation Status Utilities
// ============================================

/**
 * Check if a vacation is currently ongoing (today is within the date range)
 * @param startDate - ISO date string (YYYY-MM-DD)
 * @param endDate - ISO date string (YYYY-MM-DD)
 */
export function isCurrentlyOngoing(startDate: string, endDate: string): boolean {
	const today = formatDateISO(new Date());
	return startDate <= today && endDate >= today;
}

/**
 * Get the number of days until a vacation ends
 * @param endDate - ISO date string (YYYY-MM-DD)
 * @returns Number of days (0 if ending today, negative if already ended)
 */
export function getDaysUntilEnd(endDate: string): number {
	const today = new Date();
	today.setHours(0, 0, 0, 0);
	const end = parseISODate(endDate);
	const diff = end.getTime() - today.getTime();
	return Math.ceil(diff / (1000 * 60 * 60 * 24));
}

/**
 * Get the number of days until a vacation starts
 * @param startDate - ISO date string (YYYY-MM-DD)
 * @returns Number of days (0 if starting today, negative if already started)
 */
export function getDaysUntilStart(startDate: string): number {
	const today = new Date();
	today.setHours(0, 0, 0, 0);
	const start = parseISODate(startDate);
	const diff = start.getTime() - today.getTime();
	return Math.ceil(diff / (1000 * 60 * 60 * 24));
}

/**
 * Convert @internationalized/date DateValue to API format (DD/MM/YYYY)
 * @param dateValue - DateValue from Melt UI date picker
 */
export function dateValueToApiFormat(dateValue: DateValue): string {
	const day = String(dateValue.day).padStart(2, '0');
	const month = String(dateValue.month).padStart(2, '0');
	return `${day}/${month}/${dateValue.year}`;
}

/**
 * Convert @internationalized/date DateValue to ISO format (YYYY-MM-DD)
 * @param dateValue - DateValue from Melt UI date picker
 */
export function dateValueToISO(dateValue: DateValue): string {
	return dateValue.toString(); // DateValue.toString() returns YYYY-MM-DD
}

/**
 * Calculate business days between two DateValue dates
 * @param start - Start date
 * @param end - End date
 * @param excludeWeekends - Whether to exclude Saturday and Sunday (default: true)
 */
export function calculateBusinessDays(
	start: DateValue,
	end: DateValue,
	excludeWeekends: boolean = true
): number {
	let days = 0;
	let current = start;

	while (current.compare(end) <= 0) {
		if (!excludeWeekends) {
			days++;
		} else {
			const dayOfWeek = current.toDate(getLocalTimeZone()).getDay();
			// Sunday = 0, Saturday = 6
			if (dayOfWeek !== 0 && dayOfWeek !== 6) {
				days++;
			}
		}
		current = current.add({ days: 1 });
	}

	return days;
}
