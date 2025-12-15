<script lang="ts">
	import TeamVacationSection from './TeamVacationSection.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';
	import { ArrowRight } from 'lucide-svelte';
	import { formatDateISO } from '$lib/utils/date';
	import type { TeamVacation } from '$lib/types';

	interface Props {
		vacations: TeamVacation[];
		isLoading?: boolean;
		maxCurrentlyOut?: number;
		maxUpcoming?: number;
		calendarLink?: string;
	}

	let {
		vacations,
		isLoading = false,
		maxCurrentlyOut = 5,
		maxUpcoming = 4,
		calendarLink = '/calendar'
	}: Props = $props();

	// Categorize vacations into "Currently Out" and "Upcoming"
	const today = $derived(formatDateISO(new Date()));

	const currentlyOut = $derived(
		vacations
			.filter((v) => v.startDate <= today && v.endDate >= today)
			.sort((a, b) => a.endDate.localeCompare(b.endDate))
	);

	const upcoming = $derived(
		vacations
			.filter((v) => v.startDate > today)
			.sort((a, b) => a.startDate.localeCompare(b.startDate))
	);

	const totalVacations = $derived(currentlyOut.length + upcoming.length);
</script>

{#if isLoading}
	<!-- Loading skeleton -->
	<div class="space-y-6">
		<!-- Currently Out skeleton -->
		<div class="space-y-3">
			<div class="h-4 w-32 bg-sand-200 rounded skeleton-shimmer"></div>
			<div class="flex gap-4">
				{#each [1, 2, 3] as _}
					<div class="flex flex-col items-center gap-2">
						<Skeleton variant="avatar" />
						<div class="h-3 w-16 bg-sand-200 rounded skeleton-shimmer"></div>
					</div>
				{/each}
			</div>
		</div>
		<!-- Upcoming skeleton -->
		<div class="space-y-3">
			<div class="h-4 w-24 bg-sand-200 rounded skeleton-shimmer"></div>
			<div class="flex -space-x-2">
				{#each [1, 2, 3, 4] as _}
					<Skeleton variant="avatar" class="ring-2 ring-white" />
				{/each}
			</div>
		</div>
	</div>
{:else}
	<div class="space-y-6 content-fade-in">
		<!-- Currently Out Section -->
		<TeamVacationSection
			title="Currently Out"
			vacations={currentlyOut}
			variant="current"
			maxVisible={maxCurrentlyOut}
			emptyMessage="No one is currently on vacation"
		/>

		<!-- Upcoming Section -->
		<TeamVacationSection
			title="Upcoming"
			vacations={upcoming}
			variant="upcoming"
			maxVisible={maxUpcoming}
			emptyMessage="No upcoming vacations this month"
		/>

		<!-- View Calendar Link -->
		{#if totalVacations > 0}
			<a
				href={calendarLink}
				class="inline-flex items-center gap-1.5 text-sm text-ocean-600 hover:text-ocean-700 transition-colors group"
			>
				<span>View full calendar</span>
				<ArrowRight class="w-4 h-4 transition-transform group-hover:translate-x-0.5" />
			</a>
		{/if}
	</div>
{/if}
