<script lang="ts">
	import { clsx } from 'clsx';
	import { calendar } from '$lib/stores/calendar.svelte';
	import { formatMonthYear, formatWeekRange, getMonday } from '$lib/utils/date';
	import Button from '$lib/components/ui/Button.svelte';
	import { ChevronLeft, ChevronRight, Calendar, CalendarDays } from 'lucide-svelte';

	interface Props {
		class?: string;
	}

	let { class: className = '' }: Props = $props();

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
	<div class="flex items-center gap-1 bg-sand-100 rounded-lg p-1" role="group" aria-label="Calendar view">
		<button
			onclick={() => calendar.setViewType('week')}
			aria-pressed={calendar.viewType === 'week'}
			class={clsx(
				'flex items-center gap-1 px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
				calendar.viewType === 'week'
					? 'bg-white text-ocean-700 shadow-sm'
					: 'text-ocean-500 hover:text-ocean-700'
			)}
		>
			<CalendarDays class="w-4 h-4" />
			Week
		</button>
		<button
			onclick={() => calendar.setViewType('month')}
			aria-pressed={calendar.viewType === 'month'}
			class={clsx(
				'flex items-center gap-1 px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
				calendar.viewType === 'month'
					? 'bg-white text-ocean-700 shadow-sm'
					: 'text-ocean-500 hover:text-ocean-700'
			)}
		>
			<Calendar class="w-4 h-4" />
			Month
		</button>
	</div>
</div>
