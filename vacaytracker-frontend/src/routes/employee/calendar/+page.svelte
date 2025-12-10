<script lang="ts">
	import { calendar } from '$lib/stores/calendar.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import CalendarHeader from '$lib/components/features/calendar/CalendarHeader.svelte';
	import WeekView from '$lib/components/features/calendar/WeekView.svelte';
	import MonthView from '$lib/components/features/calendar/MonthView.svelte';
	import FilterPanel from '$lib/components/features/calendar/FilterPanel.svelte';
	import { Calendar as CalendarIcon, Users } from 'lucide-svelte';

	// Fetch data when component mounts and when view/date changes
	$effect(() => {
		// Dependencies: viewType, currentMonth, currentYear
		const _ = [calendar.viewType, calendar.currentMonth, calendar.currentYear];
		calendar.ensureDataForCurrentView();
	});
</script>

<svelte:head>
	<title>Calendar - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
	<!-- Page header -->
	<div class="flex items-center justify-between flex-wrap gap-4">
		<div>
			<h1 class="text-2xl font-bold text-ocean-800">Team Calendar</h1>
			<p class="text-ocean-600">View and filter team vacation schedules</p>
		</div>
		<FilterPanel users={calendar.availableUsers} />
	</div>

	<Card padding="none">
		{#snippet header()}
			<CalendarHeader />
		{/snippet}

		<div class="p-4">
			{#if calendar.isLoading}
				<div class="py-12 text-center text-ocean-500">
					<CalendarIcon class="w-8 h-8 mx-auto mb-2 animate-pulse" />
					<p>Loading calendar...</p>
				</div>
			{:else if calendar.error}
				<div class="py-12 text-center">
					<CalendarIcon class="w-8 h-8 mx-auto mb-2 text-error" />
					<p class="text-error font-medium">Failed to load calendar data</p>
					<p class="text-sm text-ocean-500 mt-1">{calendar.error}</p>
					<button
						onclick={() => calendar.ensureDataForCurrentView()}
						class="mt-3 text-sm text-ocean-600 hover:text-ocean-700 underline"
					>
						Try again
					</button>
				</div>
			{:else if calendar.filteredVacations.length === 0 && calendar.availableUsers.length > 0}
				<div class="py-12 text-center text-ocean-500">
					<Users class="w-8 h-8 mx-auto mb-2" />
					<p>No vacations match your filter</p>
					<button
						onclick={() => calendar.clearFilters()}
						class="mt-2 text-sm text-ocean-600 hover:text-ocean-700 underline"
					>
						Clear filters
					</button>
				</div>
			{:else}
				{#if calendar.viewType === 'week'}
					<WeekView currentDate={calendar.currentDate} vacations={calendar.filteredVacations} />
				{:else}
					<MonthView
						year={calendar.currentYear}
						month={calendar.currentMonth}
						vacations={calendar.filteredVacations}
					/>
				{/if}

				{#if calendar.filteredVacations.length === 0}
					<div class="mt-4 py-8 text-center text-ocean-500">
						<Users class="w-8 h-8 mx-auto mb-2" />
						<p>No team vacations scheduled for this period</p>
					</div>
				{/if}
			{/if}
		</div>
	</Card>
</div>
