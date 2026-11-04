<script lang="ts">
	import { clsx } from 'clsx';
	import type { TeamVacation } from '$lib/types';
	import { getWeekDays, groupByDay, type CalendarDay } from '$lib/utils/date';
	import DayCell from './DayCell.svelte';

	interface Props {
		currentDate: Date;
		vacations: TeamVacation[];
		class?: string;
	}

	let { currentDate, vacations, class: className = '' }: Props = $props();

	const weekDays = $derived(getWeekDays(currentDate));
	const weekDayLabels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

	// Pre-computed lookup (dateString -> vacations), rebuilt only when vacations or days change
	const vacationsByDay = $derived(groupByDay(weekDays, vacations));

	function getVacationsForDay(day: CalendarDay): TeamVacation[] {
		return vacationsByDay.get(day.dateString) ?? [];
	}
</script>

<div
	class={clsx('bg-white rounded-lg shadow-md border border-sand-200 overflow-hidden', className)}
>
	<!-- Header row with day labels only -->
	<div class="grid grid-cols-7 gap-px bg-sand-200">
		{#each weekDayLabels as label, i (label)}
			{@const day = weekDays[i]}
			<div
				class={clsx(
					'bg-ocean-50 py-2 text-center text-sm font-medium text-ocean-700',
					day?.isWeekend && 'bg-sand-100 text-ocean-600'
				)}
			>
				{label}
			</div>
		{/each}
	</div>

	<!-- Day cells -->
	<div class="grid grid-cols-7 gap-px bg-sand-200">
		{#each weekDays as day (day.dateString)}
			<DayCell {day} vacations={getVacationsForDay(day)} variant="week" maxEvents={5} />
		{/each}
	</div>
</div>
