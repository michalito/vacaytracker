<script lang="ts">
	import { ToggleGroup } from 'bits-ui';
	import { clsx } from 'clsx';
	import { calendar } from '$lib/stores/calendar.svelte';
	import { formatMonthYear, formatWeekRange, getMonday, getWeekDays } from '$lib/utils/date';
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

	const prevLabel = $derived(calendar.viewType === 'month' ? 'Previous month' : 'Previous week');
	const nextLabel = $derived(calendar.viewType === 'month' ? 'Next month' : 'Next week');

	// Check if current view contains today's date
	const isViewingToday = $derived.by(() => {
		const today = new Date();
		if (calendar.viewType === 'month') {
			return (
				today.getMonth() === calendar.currentDate.getMonth() &&
				today.getFullYear() === calendar.currentDate.getFullYear()
			);
		} else {
			// Week view: check if today is in the current week
			const weekDays = getWeekDays(calendar.currentDate);
			const todayStr = today.toISOString().split('T')[0];
			return weekDays.some((day) => day.dateString === todayStr);
		}
	});
</script>

<div class={clsx('flex items-center justify-between flex-wrap gap-3', className)}>
	<!-- Title and navigation -->
	<div class="flex items-center gap-3">
		<!-- Fixed width title to prevent layout shift -->
		<h2 class="text-lg font-semibold text-ocean-700 min-w-[180px]">{title}</h2>
		<!-- Navigation buttons -->
		<div class="flex items-center gap-1 bg-sand-100 rounded-lg p-1">
			<button
				type="button"
				onclick={() => calendar.goToPrevious()}
				aria-label={prevLabel}
				title={prevLabel}
				class="p-1.5 rounded-md text-ocean-600 hover:bg-white hover:shadow-sm transition-all duration-200 cursor-pointer"
			>
				<ChevronLeft class="w-4 h-4" />
			</button>
			<button
				type="button"
				onclick={() => calendar.goToToday()}
				class={clsx(
					'px-3 py-1.5 rounded-md text-sm font-medium transition-all duration-200 cursor-pointer',
					isViewingToday
						? 'bg-white text-ocean-700 shadow-sm hover:bg-ocean-50'
						: 'text-ocean-500 hover:bg-white hover:text-ocean-700 hover:shadow-sm'
				)}
			>
				Today
			</button>
			<button
				type="button"
				onclick={() => calendar.goToNext()}
				aria-label={nextLabel}
				title={nextLabel}
				class="p-1.5 rounded-md text-ocean-600 hover:bg-white hover:shadow-sm transition-all duration-200 cursor-pointer"
			>
				<ChevronRight class="w-4 h-4" />
			</button>
		</div>
	</div>

	<!-- View toggle -->
	<!--
		Single toggle groups emit "" when the active item is clicked again; ignoring it keeps
		one view selected, and the getter always reflects the store.
	-->
	<ToggleGroup.Root
		type="single"
		bind:value={
			() => calendar.viewType,
			(next) => {
				if (next === 'week' || next === 'month') calendar.setViewType(next);
			}
		}
		class="flex items-center gap-1 bg-sand-100 rounded-lg p-1"
		aria-label="Calendar view"
	>
		<ToggleGroup.Item
			value="week"
			class="flex items-center gap-1 px-3 py-1.5 rounded-md text-sm font-medium transition-colors cursor-pointer
				data-[state=on]:bg-white data-[state=on]:text-ocean-700 data-[state=on]:shadow-sm
				data-[state=off]:text-ocean-500 data-[state=off]:hover:text-ocean-700"
		>
			<CalendarDays class="w-4 h-4" />
			Week
		</ToggleGroup.Item>
		<ToggleGroup.Item
			value="month"
			class="flex items-center gap-1 px-3 py-1.5 rounded-md text-sm font-medium transition-colors cursor-pointer
				data-[state=on]:bg-white data-[state=on]:text-ocean-700 data-[state=on]:shadow-sm
				data-[state=off]:text-ocean-500 data-[state=off]:hover:text-ocean-700"
		>
			<Calendar class="w-4 h-4" />
			Month
		</ToggleGroup.Item>
	</ToggleGroup.Root>
</div>
