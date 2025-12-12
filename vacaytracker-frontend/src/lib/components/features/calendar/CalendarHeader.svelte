<script lang="ts">
	import { createToggleGroup, melt } from '@melt-ui/svelte';
	import { clsx } from 'clsx';
	import { calendar } from '$lib/stores/calendar.svelte';
	import { formatMonthYear, formatWeekRange, getMonday } from '$lib/utils/date';
	import Button from '$lib/components/ui/Button.svelte';
	import { ChevronLeft, ChevronRight, Calendar, CalendarDays } from 'lucide-svelte';

	interface Props {
		class?: string;
	}

	let { class: className = '' }: Props = $props();

	const {
		elements: { root, item },
		states: { value }
	} = createToggleGroup({
		type: 'single',
		defaultValue: calendar.viewType,
		onValueChange: ({ next }) => {
			if (next) {
				calendar.setViewType(next as 'week' | 'month');
			}
			return next;
		}
	});

	// Sync with calendar store changes
	$effect(() => {
		if (calendar.viewType !== $value) {
			value.set(calendar.viewType);
		}
	});

	const title = $derived(
		calendar.viewType === 'month'
			? formatMonthYear(calendar.currentDate)
			: formatWeekRange(getMonday(calendar.currentDate))
	);
</script>

<div class={clsx('flex items-center justify-between flex-wrap gap-3', className)}>
	<!-- Title and navigation -->
	<div class="flex items-center gap-2">
		<h2 class="text-lg font-semibold text-ocean-700">{title}</h2>
		<div class="flex items-center gap-1">
			<Button variant="ghost" size="sm" onclick={() => calendar.goToPrevious()}>
				<ChevronLeft class="w-4 h-4" />
			</Button>
			<Button variant="outline" size="sm" onclick={() => calendar.goToToday()}>Today</Button>
			<Button variant="ghost" size="sm" onclick={() => calendar.goToNext()}>
				<ChevronRight class="w-4 h-4" />
			</Button>
		</div>
	</div>

	<!-- View toggle -->
	<div
		use:melt={$root}
		class="flex items-center gap-1 bg-sand-100 rounded-lg p-1"
		aria-label="Calendar view"
	>
		<button
			use:melt={$item('week')}
			class="flex items-center gap-1 px-3 py-1.5 rounded-md text-sm font-medium transition-colors cursor-pointer
				data-[state=on]:bg-white data-[state=on]:text-ocean-700 data-[state=on]:shadow-sm
				data-[state=off]:text-ocean-500 data-[state=off]:hover:text-ocean-700"
		>
			<CalendarDays class="w-4 h-4" />
			Week
		</button>
		<button
			use:melt={$item('month')}
			class="flex items-center gap-1 px-3 py-1.5 rounded-md text-sm font-medium transition-colors cursor-pointer
				data-[state=on]:bg-white data-[state=on]:text-ocean-700 data-[state=on]:shadow-sm
				data-[state=off]:text-ocean-500 data-[state=off]:hover:text-ocean-700"
		>
			<Calendar class="w-4 h-4" />
			Month
		</button>
	</div>
</div>
