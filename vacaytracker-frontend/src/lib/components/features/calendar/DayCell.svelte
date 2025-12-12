<script lang="ts">
	import { clsx } from 'clsx';
	import type { TeamVacation } from '$lib/types';
	import type { CalendarDay } from '$lib/utils/date';
	import VacationEvent from './VacationEvent.svelte';

	type EventPosition = 'start' | 'middle' | 'end' | 'single';

	interface Props {
		day: CalendarDay;
		vacations: TeamVacation[];
		variant?: 'week' | 'month';
		maxEvents?: number;
		class?: string;
	}

	let {
		day,
		vacations,
		variant = 'month',
		maxEvents = 3,
		class: className = ''
	}: Props = $props();

	const visibleVacations = $derived(vacations.slice(0, maxEvents));
	const hiddenCount = $derived(Math.max(0, vacations.length - maxEvents));
	const hiddenVacations = $derived(vacations.slice(maxEvents));

	// Build title for overflow indicator
	const overflowTitle = $derived(
		hiddenVacations.map((v) => v.userName).join(', ')
	);

	function getEventPosition(vacation: TeamVacation): EventPosition {
		const isStart = day.dateString === vacation.startDate;
		const isEnd = day.dateString === vacation.endDate;
		if (isStart && isEnd) return 'single';
		if (isStart) return 'start';
		if (isEnd) return 'end';
		return 'middle';
	}
</script>

<div
	class={clsx(
		'bg-white p-2',
		variant === 'week' ? 'min-h-[120px]' : 'min-h-[100px]',
		!day.isCurrentMonth && 'bg-sand-50',
		day.isToday && 'ring-2 ring-inset ring-ocean-500',
		day.isWeekend && day.isCurrentMonth && 'bg-sand-50/50',
		className
	)}
>
	<!-- Day number -->
	<span
		class={clsx(
			'inline-flex items-center justify-center w-7 h-7 text-sm rounded-full mb-1',
			day.isToday && 'bg-ocean-500 text-white font-bold',
			!day.isToday && day.isCurrentMonth && 'text-ocean-700',
			!day.isCurrentMonth && 'text-ocean-400'
		)}
	>
		{day.dayOfMonth}
	</span>

	<!-- Events -->
	{#if vacations.length > 0}
		<div class="space-y-1">
			{#each visibleVacations as vacation (vacation.id)}
				<VacationEvent
					{vacation}
					position={getEventPosition(vacation)}
					compact={variant === 'month'}
				/>
			{/each}
			{#if hiddenCount > 0}
				<div class="text-xs text-ocean-500 pl-1 cursor-default" title={overflowTitle}>
					+{hiddenCount} more
				</div>
			{/if}
		</div>
	{/if}
</div>
