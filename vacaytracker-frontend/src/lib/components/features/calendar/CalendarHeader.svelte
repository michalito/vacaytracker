<script lang="ts">
	import { createToggleGroup, melt } from '@melt-ui/svelte';
	import { clsx } from 'clsx';
	import { calendar } from '$lib/stores/calendar.svelte';
	import { formatMonthYear, formatWeekRange, getMonday, getWeekDays } from '$lib/utils/date';
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
		<!-- Navigation buttons styled to match Melt-UI toggle group -->
		<div class="flex items-center gap-1 bg-sand-100 rounded-lg p-1">
			<Button variant="ghost" size="sm" onclick={() => calendar.goToPrevious()} class="!p-1.5 !rounded-md hover:bg-white hover:shadow-sm focus:!ring-0 focus:!ring-offset-0">
				<ChevronLeft class="w-4 h-4" />
			</Button>
			<Button
				variant="ghost"
				size="sm"
				onclick={() => calendar.goToToday()}
				class={clsx(
					'!rounded-md focus:!ring-0 focus:!ring-offset-0',
					isViewingToday
						? '!bg-white !text-ocean-700 !shadow-sm hover:!bg-ocean-50'
						: 'hover:bg-white hover:shadow-sm'
				)}
			>
				Today
			</Button>
			<Button variant="ghost" size="sm" onclick={() => calendar.goToNext()} class="!p-1.5 !rounded-md hover:bg-white hover:shadow-sm focus:!ring-0 focus:!ring-offset-0">
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
