<script lang="ts">
	import { clsx } from 'clsx';
	import type { TeamVacation } from '$lib/types';
	import { getWeekDays, isDateInRange, parseISODate, type CalendarDay } from '$lib/utils/date';
	import DayCell from './DayCell.svelte';

	interface Props {
		currentDate: Date;
		vacations: TeamVacation[];
		class?: string;
	}

	let { currentDate, vacations, class: className = '' }: Props = $props();

	const weekDays = $derived(getWeekDays(currentDate));
	const weekDayLabels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

	function getVacationsForDay(day: CalendarDay): TeamVacation[] {
		return vacations.filter((v) =>
			isDateInRange(day.date, parseISODate(v.startDate), parseISODate(v.endDate))
		);
	}
</script>

<div
	class={clsx('bg-white rounded-lg shadow-md border border-sand-200 overflow-hidden', className)}
>
	<!-- Header row with day labels and dates -->
	<div class="grid grid-cols-7 gap-px bg-sand-200">
		{#each weekDayLabels as label, i}
			{@const day = weekDays[i]}
			<div
				class={clsx(
					'bg-ocean-50 py-2 px-3 text-center',
					day?.isWeekend && 'bg-sand-100',
					day?.isToday && 'bg-ocean-100'
				)}
			>
				<div class="text-xs font-medium text-ocean-600">{label}</div>
				<div
					class={clsx(
						'text-lg font-semibold',
						day?.isToday ? 'text-ocean-600' : 'text-ocean-800'
					)}
				>
					{day?.dayOfMonth}
				</div>
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
