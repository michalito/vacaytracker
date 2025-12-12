<script lang="ts">
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import { Users } from 'lucide-svelte';
	import { formatDateRangeShort } from '$lib/utils/date';
	import type { TeamVacation } from '$lib/types';

	interface Props {
		vacations: TeamVacation[];
		isLoading?: boolean;
		maxItems?: number;
		showMoreLink?: string;
		emptyMessage?: string;
	}

	let {
		vacations,
		isLoading = false,
		maxItems = 5,
		showMoreLink = '/calendar',
		emptyMessage = 'No team vacations this month'
	}: Props = $props();

	const displayedVacations = $derived(vacations.slice(0, maxItems));
	const hasMore = $derived(vacations.length > maxItems);
</script>

{#if isLoading}
	<div class="space-y-3 py-2">
		{#each [1, 2, 3] as _}
			<div class="flex items-center gap-3 p-2">
				<Skeleton variant="avatar" />
				<div class="flex-1 space-y-2">
					<Skeleton variant="text" width="30%" />
					<Skeleton variant="text" width="50%" />
				</div>
			</div>
		{/each}
	</div>
{:else if vacations.length === 0}
	<EmptyState icon={Users} message={emptyMessage} />
{:else}
	<div class="space-y-2 content-fade-in">
		{#each displayedVacations as vacation (vacation.id)}
			<div class="flex items-center gap-3 p-2 bg-sand-50 rounded-lg">
				<Avatar name={vacation.userName} size="sm" />
				<div class="flex-1 min-w-0">
					<p class="font-medium text-ocean-800 truncate">{vacation.userName}</p>
					<p class="text-sm text-ocean-500">
						{formatDateRangeShort(vacation.startDate, vacation.endDate)} ({vacation.totalDays} days)
					</p>
				</div>
			</div>
		{/each}
		{#if hasMore && showMoreLink}
			<a
				href={showMoreLink}
				class="block text-center text-sm text-ocean-600 hover:text-ocean-700 pt-2"
			>
				View all {vacations.length} vacations
			</a>
		{/if}
	</div>
{/if}
