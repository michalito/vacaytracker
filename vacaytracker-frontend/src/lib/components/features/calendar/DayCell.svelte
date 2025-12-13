<script lang="ts">
	import { createPopover, melt } from '@melt-ui/svelte';
	import { clsx } from 'clsx';
	import type { TeamVacation } from '$lib/types';
	import type { CalendarDay } from '$lib/utils/date';
	import { formatDateRangeShort } from '$lib/utils/date';
	import { getUserColor } from '$lib/utils/colors';
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

	// Popover for hidden vacations overflow
	const {
		elements: { trigger, content, close },
		states: { open }
	} = createPopover({
		forceVisible: true,
		positioning: { placement: 'bottom-start' }
	});

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
		day.isToday && 'ring-2 ring-inset ring-ocean-400 bg-ocean-50/30 rounded-lg',
		day.isWeekend && day.isCurrentMonth && !day.isToday && 'bg-sand-50/50',
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
				<button
					use:melt={$trigger}
					class="text-xs text-ocean-500 pl-1 hover:text-ocean-700 hover:underline cursor-pointer transition-colors"
				>
					+{hiddenCount} more
				</button>

				{#if $open}
					<div
						use:melt={$content}
						class="z-50 w-56 bg-white rounded-lg shadow-xl border border-ocean-200 overflow-hidden
							transition-all duration-150
							data-[state=open]:opacity-100 data-[state=open]:scale-100
							data-[state=closed]:opacity-0 data-[state=closed]:scale-95"
					>
						<div class="px-3 py-2 bg-ocean-50 border-b border-ocean-100">
							<span class="text-sm font-medium text-ocean-700">
								+{hiddenCount} more vacation{hiddenCount !== 1 ? 's' : ''}
							</span>
						</div>
						<div class="max-h-48 overflow-y-auto">
							{#each hiddenVacations as vacation (vacation.id)}
								{@const color = getUserColor(vacation.userId)}
								<div class="px-3 py-2 flex items-start gap-2 hover:bg-ocean-50 transition-colors">
									<span class={clsx('w-2 h-2 rounded-full mt-1.5 flex-shrink-0', color.background)}></span>
									<div class="min-w-0 flex-1">
										<div class="text-sm font-medium text-ocean-800 truncate">
											{vacation.userName}
										</div>
										<div class="text-xs text-ocean-500">
											{formatDateRangeShort(vacation.startDate, vacation.endDate)}
											<span class="text-ocean-400">
												({vacation.totalDays} day{vacation.totalDays !== 1 ? 's' : ''})
											</span>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			{/if}
		</div>
	{/if}
</div>
