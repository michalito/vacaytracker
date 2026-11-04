<script lang="ts">
	import { clsx } from 'clsx';
	import type { TeamVacation } from '#lib/types/index.js';
	import { getCalendarGridDays, groupByDay, type CalendarDay } from '#lib/utils/date.js';
	import DayCell from './DayCell.svelte';

	interface Props {
		year: number;
		month: number;
		vacations: TeamVacation[];
		class?: string;
	}

	let { year, month, vacations, class: className = '' }: Props = $props();

	const calendarDays = $derived(getCalendarGridDays(year, month));
	const weekDayLabels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

	// Pre-computed lookup (dateString -> vacations), rebuilt only when vacations or days change
	const vacationsByDay = $derived(groupByDay(calendarDays, vacations));

	function getVacationsForDay(day: CalendarDay): TeamVacation[] {
		return vacationsByDay.get(day.dateString) ?? [];
	}
</script>

<div
	class={clsx('bg-white rounded-lg shadow-md border border-sand-200 overflow-hidden', className)}
>
	<!-- Header row with day labels -->
	<div class="grid grid-cols-7 gap-px bg-sand-200">
		{#each weekDayLabels as label, i (label)}
			<div
				class={clsx(
					'bg-ocean-50 py-2 text-center text-sm font-medium text-ocean-700',
					i >= 5 && 'bg-sand-100 text-ocean-600'
				)}
			>
				{label}
			</div>
		{/each}
	</div>

	<!-- Day grid (6 weeks x 7 days = 42 cells) -->
	<div class="grid grid-cols-7 gap-px bg-sand-200">
		{#each calendarDays as day (day.dateString)}
			<DayCell {day} vacations={getVacationsForDay(day)} variant="month" maxEvents={3} />
		{/each}
	</div>
</div>
