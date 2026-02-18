import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { CalendarDate } from '@internationalized/date';
import {
	parseISODate,
	formatDateISO,
	formatMonthYear,
	formatWeekRange,
	addDays,
	addWeeks,
	addMonths,
	isSameDay,
	isDateInRange,
	getMonday,
	getWeekDays,
	getWeekBoundaries,
	getMonthBoundaries,
	getCalendarGridDays,
	getWeekNumber,
	getMonthsForWeek,
	formatDateShort,
	formatDateMedium,
	formatDateLong,
	formatDateRangeShort,
	toApiDateFormat,
	fromApiDateFormat,
	isCurrentlyOngoing,
	getDaysUntilEnd,
	getDaysUntilStart,
	dateValueToApiFormat,
	dateValueToISO,
	calculateBusinessDays
} from '$lib/utils/date';

describe('parseISODate', () => {
	it('parses a standard YYYY-MM-DD string correctly', () => {
		const date = parseISODate('2025-06-15');
		expect(date.getFullYear()).toBe(2025);
		expect(date.getMonth()).toBe(5); // 0-indexed: June = 5
		expect(date.getDate()).toBe(15);
	});

	it('parses January 1st correctly', () => {
		const date = parseISODate('2024-01-01');
		expect(date.getFullYear()).toBe(2024);
		expect(date.getMonth()).toBe(0);
		expect(date.getDate()).toBe(1);
	});

	it('parses December 31st correctly', () => {
		const date = parseISODate('2024-12-31');
		expect(date.getFullYear()).toBe(2024);
		expect(date.getMonth()).toBe(11);
		expect(date.getDate()).toBe(31);
	});

	it('parses a leap day correctly', () => {
		const date = parseISODate('2024-02-29');
		expect(date.getFullYear()).toBe(2024);
		expect(date.getMonth()).toBe(1);
		expect(date.getDate()).toBe(29);
	});
});

describe('formatDateISO', () => {
	it('formats a date to YYYY-MM-DD string', () => {
		const date = new Date(2025, 5, 15); // June 15, 2025
		expect(formatDateISO(date)).toBe('2025-06-15');
	});

	it('zero-pads single digit months and days', () => {
		const date = new Date(2025, 0, 5); // January 5, 2025
		expect(formatDateISO(date)).toBe('2025-01-05');
	});

	it('handles double-digit months and days', () => {
		const date = new Date(2025, 11, 25); // December 25, 2025
		expect(formatDateISO(date)).toBe('2025-12-25');
	});

	it('handles year correctly', () => {
		const date = new Date(2000, 0, 1);
		expect(formatDateISO(date)).toBe('2000-01-01');
	});
});

describe('formatMonthYear', () => {
	it('formats date as "Month Year" with default locale', () => {
		const date = new Date(2025, 11, 1); // December 2025
		expect(formatMonthYear(date)).toBe('December 2025');
	});

	it('formats January correctly', () => {
		const date = new Date(2025, 0, 15); // January 2025
		expect(formatMonthYear(date)).toBe('January 2025');
	});
});

describe('formatWeekRange', () => {
	it('formats a week range within the same month', () => {
		// Monday June 2 to Sunday June 8, 2025
		const monday = new Date(2025, 5, 2);
		const result = formatWeekRange(monday);
		expect(result).toBe('Jun 2 - 8, 2025');
	});

	it('formats a week range spanning two months', () => {
		// Monday June 30 to Sunday July 6, 2025
		const monday = new Date(2025, 5, 30);
		const result = formatWeekRange(monday);
		expect(result).toBe('Jun 30 - Jul 6, 2025');
	});
});

describe('addDays', () => {
	it('adds positive days', () => {
		const date = new Date(2025, 5, 10);
		const result = addDays(date, 5);
		expect(result.getDate()).toBe(15);
		expect(result.getMonth()).toBe(5);
	});

	it('adds negative days', () => {
		const date = new Date(2025, 5, 10);
		const result = addDays(date, -3);
		expect(result.getDate()).toBe(7);
		expect(result.getMonth()).toBe(5);
	});

	it('crosses month boundary forwards', () => {
		const date = new Date(2025, 5, 28); // June 28
		const result = addDays(date, 5);
		expect(result.getMonth()).toBe(6); // July
		expect(result.getDate()).toBe(3);
	});

	it('crosses month boundary backwards', () => {
		const date = new Date(2025, 6, 2); // July 2
		const result = addDays(date, -5);
		expect(result.getMonth()).toBe(5); // June
		expect(result.getDate()).toBe(27);
	});

	it('does not mutate the original date', () => {
		const date = new Date(2025, 5, 10);
		addDays(date, 5);
		expect(date.getDate()).toBe(10);
	});
});

describe('addWeeks', () => {
	it('adds positive weeks', () => {
		const date = new Date(2025, 5, 2); // June 2
		const result = addWeeks(date, 2);
		expect(result.getDate()).toBe(16);
		expect(result.getMonth()).toBe(5);
	});

	it('adds one week correctly', () => {
		const date = new Date(2025, 5, 10);
		const result = addWeeks(date, 1);
		expect(formatDateISO(result)).toBe('2025-06-17');
	});
});

describe('addMonths', () => {
	it('adds positive months', () => {
		const date = new Date(2025, 0, 15); // January 15
		const result = addMonths(date, 3);
		expect(result.getMonth()).toBe(3); // April
		expect(result.getDate()).toBe(15);
	});

	it('subtracts months with negative value', () => {
		const date = new Date(2025, 5, 15); // June 15
		const result = addMonths(date, -2);
		expect(result.getMonth()).toBe(3); // April
		expect(result.getDate()).toBe(15);
	});

	it('crosses year boundary forwards', () => {
		const date = new Date(2025, 10, 15); // November 15
		const result = addMonths(date, 3);
		expect(result.getFullYear()).toBe(2026);
		expect(result.getMonth()).toBe(1); // February
	});

	it('crosses year boundary backwards', () => {
		const date = new Date(2025, 1, 15); // February 15
		const result = addMonths(date, -3);
		expect(result.getFullYear()).toBe(2024);
		expect(result.getMonth()).toBe(10); // November
	});

	it('does not mutate the original date', () => {
		const date = new Date(2025, 5, 10);
		addMonths(date, 3);
		expect(date.getMonth()).toBe(5);
	});
});

describe('isSameDay', () => {
	it('returns true for the same day', () => {
		const a = new Date(2025, 5, 15, 10, 30);
		const b = new Date(2025, 5, 15, 22, 0);
		expect(isSameDay(a, b)).toBe(true);
	});

	it('returns false for different days in the same month', () => {
		const a = new Date(2025, 5, 15);
		const b = new Date(2025, 5, 16);
		expect(isSameDay(a, b)).toBe(false);
	});

	it('returns false for same day number in different months', () => {
		const a = new Date(2025, 5, 15);
		const b = new Date(2025, 6, 15);
		expect(isSameDay(a, b)).toBe(false);
	});

	it('returns false for same day and month in different years', () => {
		const a = new Date(2025, 5, 15);
		const b = new Date(2026, 5, 15);
		expect(isSameDay(a, b)).toBe(false);
	});
});

describe('isDateInRange', () => {
	const start = new Date(2025, 5, 10); // June 10
	const end = new Date(2025, 5, 20); // June 20

	it('returns true for a date inside the range', () => {
		const date = new Date(2025, 5, 15);
		expect(isDateInRange(date, start, end)).toBe(true);
	});

	it('returns true for the start boundary', () => {
		expect(isDateInRange(start, start, end)).toBe(true);
	});

	it('returns true for the end boundary', () => {
		expect(isDateInRange(end, start, end)).toBe(true);
	});

	it('returns false for a date before the range', () => {
		const date = new Date(2025, 5, 9);
		expect(isDateInRange(date, start, end)).toBe(false);
	});

	it('returns false for a date after the range', () => {
		const date = new Date(2025, 5, 21);
		expect(isDateInRange(date, start, end)).toBe(false);
	});

	it('ignores time components', () => {
		const date = new Date(2025, 5, 15, 23, 59, 59);
		expect(isDateInRange(date, start, end)).toBe(true);
	});
});

describe('getMonday', () => {
	it('returns the same date when given a Monday', () => {
		// June 2, 2025 is a Monday
		const monday = new Date(2025, 5, 2);
		const result = getMonday(monday);
		expect(result.getDate()).toBe(2);
		expect(result.getMonth()).toBe(5);
	});

	it('returns the previous Monday when given a Wednesday', () => {
		// June 4, 2025 is a Wednesday
		const wednesday = new Date(2025, 5, 4);
		const result = getMonday(wednesday);
		expect(result.getDate()).toBe(2);
		expect(result.getMonth()).toBe(5);
	});

	it('returns the previous Monday when given a Sunday (goes back 6 days)', () => {
		// June 8, 2025 is a Sunday
		const sunday = new Date(2025, 5, 8);
		const result = getMonday(sunday);
		expect(result.getDate()).toBe(2);
		expect(result.getMonth()).toBe(5);
	});

	it('crosses month boundary when needed', () => {
		// July 1, 2025 is a Tuesday; Monday is June 30
		const tuesday = new Date(2025, 6, 1);
		const result = getMonday(tuesday);
		expect(result.getDate()).toBe(30);
		expect(result.getMonth()).toBe(5); // June
	});
});

describe('getWeekDays', () => {
	it('returns exactly 7 days', () => {
		const date = new Date(2025, 5, 4); // Wednesday June 4, 2025
		const days = getWeekDays(date);
		expect(days).toHaveLength(7);
	});

	it('starts on Monday and ends on Sunday', () => {
		const date = new Date(2025, 5, 4); // Wednesday June 4, 2025
		const days = getWeekDays(date);
		// First day should be Monday (dayOfWeek = 0)
		expect(days[0].dayOfWeek).toBe(0);
		// Last day should be Sunday (dayOfWeek = 6)
		expect(days[6].dayOfWeek).toBe(6);
	});

	it('has correct dayOfWeek values 0 through 6', () => {
		const date = new Date(2025, 5, 4);
		const days = getWeekDays(date);
		for (let i = 0; i < 7; i++) {
			expect(days[i].dayOfWeek).toBe(i);
		}
	});

	it('marks Saturday and Sunday as weekend', () => {
		const date = new Date(2025, 5, 4);
		const days = getWeekDays(date);
		expect(days[0].isWeekend).toBe(false); // Monday
		expect(days[4].isWeekend).toBe(false); // Friday
		expect(days[5].isWeekend).toBe(true); // Saturday
		expect(days[6].isWeekend).toBe(true); // Sunday
	});

	it('has dateString in YYYY-MM-DD format', () => {
		const date = new Date(2025, 5, 4);
		const days = getWeekDays(date);
		expect(days[0].dateString).toMatch(/^\d{4}-\d{2}-\d{2}$/);
	});

	it('flags days in the current month correctly', () => {
		// Week spanning June/July: Mon Jun 30 - Sun Jul 6, 2025
		// Pass a date in June so currentMonth = June
		const date = new Date(2025, 5, 30); // June 30
		const days = getWeekDays(date);
		// Monday Jun 30 is in current month
		expect(days[0].isCurrentMonth).toBe(true);
		// Tuesday Jul 1 is NOT in current month (June)
		expect(days[1].isCurrentMonth).toBe(false);
	});
});

describe('getWeekBoundaries', () => {
	it('returns Monday to Sunday range', () => {
		// June 4, 2025 is a Wednesday
		const date = new Date(2025, 5, 4);
		const { start, end } = getWeekBoundaries(date);
		expect(start.getDay()).toBe(1); // Monday
		expect(end.getDay()).toBe(0); // Sunday
	});

	it('start is 6 days before end', () => {
		const date = new Date(2025, 5, 4);
		const { start, end } = getWeekBoundaries(date);
		const diff = (end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24);
		expect(diff).toBe(6);
	});

	it('returns correct dates for a known week', () => {
		// June 4, 2025 (Wednesday) => Mon Jun 2 - Sun Jun 8
		const date = new Date(2025, 5, 4);
		const { start, end } = getWeekBoundaries(date);
		expect(formatDateISO(start)).toBe('2025-06-02');
		expect(formatDateISO(end)).toBe('2025-06-08');
	});
});

describe('getMonthBoundaries', () => {
	it('returns correct start and end for June 2025', () => {
		const { start, end } = getMonthBoundaries(2025, 6);
		expect(formatDateISO(start)).toBe('2025-06-01');
		expect(formatDateISO(end)).toBe('2025-06-30');
	});

	it('returns correct end for February in a leap year', () => {
		const { start, end } = getMonthBoundaries(2024, 2);
		expect(formatDateISO(start)).toBe('2024-02-01');
		expect(formatDateISO(end)).toBe('2024-02-29');
	});

	it('returns correct end for February in a non-leap year', () => {
		const { end } = getMonthBoundaries(2025, 2);
		expect(formatDateISO(end)).toBe('2025-02-28');
	});

	it('returns correct boundaries for December', () => {
		const { start, end } = getMonthBoundaries(2025, 12);
		expect(formatDateISO(start)).toBe('2025-12-01');
		expect(formatDateISO(end)).toBe('2025-12-31');
	});
});

describe('getCalendarGridDays', () => {
	it('returns exactly 42 days (6 weeks)', () => {
		const days = getCalendarGridDays(2025, 6);
		expect(days).toHaveLength(42);
	});

	it('first day is a Monday', () => {
		const days = getCalendarGridDays(2025, 6);
		expect(days[0].dayOfWeek).toBe(0); // EU Monday = 0
		expect(days[0].date.getDay()).toBe(1); // JS Monday = 1
	});

	it('correctly flags days of the target month as isCurrentMonth', () => {
		const days = getCalendarGridDays(2025, 6); // June 2025
		const juneDays = days.filter((d) => d.isCurrentMonth);
		// June has 30 days
		expect(juneDays).toHaveLength(30);
	});

	it('includes padding days from previous and next months', () => {
		const days = getCalendarGridDays(2025, 6);
		const nonJuneDays = days.filter((d) => !d.isCurrentMonth);
		// 42 total - 30 June days = 12 padding days
		expect(nonJuneDays).toHaveLength(12);
	});

	it('has consecutive dates', () => {
		const days = getCalendarGridDays(2025, 6);
		for (let i = 1; i < days.length; i++) {
			const prevDate = days[i - 1].date;
			const currDate = days[i].date;
			const diff = (currDate.getTime() - prevDate.getTime()) / (1000 * 60 * 60 * 24);
			expect(diff).toBe(1);
		}
	});
});

describe('getWeekNumber', () => {
	it('returns week 1 for January 1, 2024 (Monday)', () => {
		// Jan 1, 2024 is a Monday => ISO week 1
		const date = new Date(2024, 0, 1);
		expect(getWeekNumber(date)).toBe(1);
	});

	it('returns week 1 for December 31, 2024 (Tuesday of week 1 of 2025)', () => {
		// Dec 31, 2024 is a Tuesday => ISO week 1 of 2025
		const date = new Date(2024, 11, 31);
		expect(getWeekNumber(date)).toBe(1);
	});

	it('returns week 52 for December 28, 2025 (Sunday)', () => {
		// Dec 28, 2025 is a Sunday => last day of ISO week 52
		const date = new Date(2025, 11, 28);
		expect(getWeekNumber(date)).toBe(52);
	});

	it('returns week 1 for December 29, 2025 (Monday of week 1 of 2026)', () => {
		// Dec 29, 2025 is a Monday => ISO week 1 of 2026
		const date = new Date(2025, 11, 29);
		expect(getWeekNumber(date)).toBe(1);
	});
});

describe('getMonthsForWeek', () => {
	it('returns a single month for a week fully within one month', () => {
		// June 4, 2025 => week Mon Jun 2 - Sun Jun 8 (all in June)
		const date = new Date(2025, 5, 4);
		const months = getMonthsForWeek(date);
		expect(months).toHaveLength(1);
		expect(months[0]).toEqual({ month: 6, year: 2025 });
	});

	it('returns two months for a week spanning a month boundary', () => {
		// June 30, 2025 => week Mon Jun 30 - Sun Jul 6
		const date = new Date(2025, 5, 30);
		const months = getMonthsForWeek(date);
		expect(months).toHaveLength(2);
		expect(months[0]).toEqual({ month: 6, year: 2025 });
		expect(months[1]).toEqual({ month: 7, year: 2025 });
	});

	it('returns two months for a week spanning a year boundary', () => {
		// Dec 31, 2025 is a Wednesday => week Mon Dec 29 - Sun Jan 4, 2026
		const date = new Date(2025, 11, 31);
		const months = getMonthsForWeek(date);
		expect(months).toHaveLength(2);
		expect(months[0]).toEqual({ month: 12, year: 2025 });
		expect(months[1]).toEqual({ month: 1, year: 2026 });
	});
});

describe('formatDateShort', () => {
	it('formats a date string to short format', () => {
		const result = formatDateShort('2025-06-15');
		// en-GB: "15 Jun"
		expect(result).toBe('15 Jun');
	});

	it('formats January date correctly', () => {
		const result = formatDateShort('2025-01-05');
		expect(result).toBe('5 Jan');
	});
});

describe('formatDateMedium', () => {
	it('formats a date string to medium format with year', () => {
		const result = formatDateMedium('2025-06-15');
		// en-GB: "15 Jun 2025"
		expect(result).toBe('15 Jun 2025');
	});

	it('formats December date correctly', () => {
		const result = formatDateMedium('2025-12-25');
		expect(result).toBe('25 Dec 2025');
	});
});

describe('formatDateLong', () => {
	it('formats a date string to long format', () => {
		const result = formatDateLong('2025-06-15');
		// en-GB: "15 June 2025"
		expect(result).toBe('15 June 2025');
	});

	it('returns "Not set" when dateStr is undefined', () => {
		expect(formatDateLong(undefined)).toBe('Not set');
	});

	it('returns "Not set" when dateStr is empty string', () => {
		expect(formatDateLong('')).toBe('Not set');
	});
});

describe('formatDateRangeShort', () => {
	it('formats a date range correctly', () => {
		const result = formatDateRangeShort('2025-06-10', '2025-06-20');
		expect(result).toBe('10 Jun - 20 Jun');
	});

	it('formats a cross-month range correctly', () => {
		const result = formatDateRangeShort('2025-06-28', '2025-07-05');
		expect(result).toBe('28 Jun - 5 Jul');
	});
});

describe('toApiDateFormat', () => {
	it('converts YYYY-MM-DD to DD/MM/YYYY', () => {
		expect(toApiDateFormat('2025-06-15')).toBe('15/06/2025');
	});

	it('preserves zero-padding', () => {
		expect(toApiDateFormat('2025-01-05')).toBe('05/01/2025');
	});
});

describe('fromApiDateFormat', () => {
	it('converts DD/MM/YYYY to YYYY-MM-DD', () => {
		expect(fromApiDateFormat('15/06/2025')).toBe('2025-06-15');
	});

	it('pads single digit day and month', () => {
		expect(fromApiDateFormat('5/1/2025')).toBe('2025-01-05');
	});
});

describe('time-dependent functions', () => {
	beforeEach(() => {
		vi.useFakeTimers();
		vi.setSystemTime(new Date(2025, 5, 15)); // June 15, 2025
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	describe('isCurrentlyOngoing', () => {
		it('returns true when today is within the range', () => {
			expect(isCurrentlyOngoing('2025-06-10', '2025-06-20')).toBe(true);
		});

		it('returns true when today is the start date', () => {
			expect(isCurrentlyOngoing('2025-06-15', '2025-06-20')).toBe(true);
		});

		it('returns true when today is the end date', () => {
			expect(isCurrentlyOngoing('2025-06-10', '2025-06-15')).toBe(true);
		});

		it('returns false when today is before the range', () => {
			expect(isCurrentlyOngoing('2025-06-16', '2025-06-20')).toBe(false);
		});

		it('returns false when today is after the range', () => {
			expect(isCurrentlyOngoing('2025-06-01', '2025-06-14')).toBe(false);
		});
	});

	describe('getDaysUntilEnd', () => {
		it('returns positive days for a future end date', () => {
			expect(getDaysUntilEnd('2025-06-20')).toBe(5);
		});

		it('returns 0 when end date is today', () => {
			expect(getDaysUntilEnd('2025-06-15')).toBe(0);
		});

		it('returns negative days when end date is in the past', () => {
			expect(getDaysUntilEnd('2025-06-10')).toBe(-5);
		});
	});

	describe('getDaysUntilStart', () => {
		it('returns positive days for a future start date', () => {
			expect(getDaysUntilStart('2025-06-20')).toBe(5);
		});

		it('returns 0 when start date is today', () => {
			expect(getDaysUntilStart('2025-06-15')).toBe(0);
		});

		it('returns negative days when start date is in the past', () => {
			expect(getDaysUntilStart('2025-06-10')).toBe(-5);
		});
	});
});

describe('dateValueToApiFormat', () => {
	it('converts a CalendarDate to DD/MM/YYYY format', () => {
		const dateValue = new CalendarDate(2025, 6, 15);
		expect(dateValueToApiFormat(dateValue)).toBe('15/06/2025');
	});

	it('zero-pads single digit day and month', () => {
		const dateValue = new CalendarDate(2025, 1, 5);
		expect(dateValueToApiFormat(dateValue)).toBe('05/01/2025');
	});

	it('handles December 31st', () => {
		const dateValue = new CalendarDate(2025, 12, 31);
		expect(dateValueToApiFormat(dateValue)).toBe('31/12/2025');
	});
});

describe('dateValueToISO', () => {
	it('converts a CalendarDate to YYYY-MM-DD string', () => {
		const dateValue = new CalendarDate(2025, 6, 15);
		expect(dateValueToISO(dateValue)).toBe('2025-06-15');
	});

	it('handles single digit month and day', () => {
		const dateValue = new CalendarDate(2025, 1, 5);
		expect(dateValueToISO(dateValue)).toBe('2025-01-05');
	});
});

describe('calculateBusinessDays', () => {
	it('counts only weekdays when excludeWeekends is true', () => {
		// Mon Jun 2 to Fri Jun 6, 2025 => 5 business days
		const start = new CalendarDate(2025, 6, 2);
		const end = new CalendarDate(2025, 6, 6);
		expect(calculateBusinessDays(start, end, true)).toBe(5);
	});

	it('excludes Saturday and Sunday', () => {
		// Mon Jun 2 to Sun Jun 8, 2025 => 5 business days (Mon-Fri)
		const start = new CalendarDate(2025, 6, 2);
		const end = new CalendarDate(2025, 6, 8);
		expect(calculateBusinessDays(start, end, true)).toBe(5);
	});

	it('counts all days when excludeWeekends is false', () => {
		// Mon Jun 2 to Sun Jun 8, 2025 => 7 calendar days
		const start = new CalendarDate(2025, 6, 2);
		const end = new CalendarDate(2025, 6, 8);
		expect(calculateBusinessDays(start, end, false)).toBe(7);
	});

	it('returns 1 for a single weekday', () => {
		const date = new CalendarDate(2025, 6, 4); // Wednesday
		expect(calculateBusinessDays(date, date, true)).toBe(1);
	});

	it('returns 0 for a Saturday-only range with excludeWeekends', () => {
		const saturday = new CalendarDate(2025, 6, 7);
		expect(calculateBusinessDays(saturday, saturday, true)).toBe(0);
	});

	it('returns 1 for a Saturday-only range without excludeWeekends', () => {
		const saturday = new CalendarDate(2025, 6, 7);
		expect(calculateBusinessDays(saturday, saturday, false)).toBe(1);
	});

	it('counts business days across two weeks', () => {
		// Mon Jun 2 to Fri Jun 13, 2025 => 10 business days
		const start = new CalendarDate(2025, 6, 2);
		const end = new CalendarDate(2025, 6, 13);
		expect(calculateBusinessDays(start, end, true)).toBe(10);
	});

	it('defaults excludeWeekends to true', () => {
		// Mon Jun 2 to Sun Jun 8, 2025 => 5 business days
		const start = new CalendarDate(2025, 6, 2);
		const end = new CalendarDate(2025, 6, 8);
		expect(calculateBusinessDays(start, end)).toBe(5);
	});
});
