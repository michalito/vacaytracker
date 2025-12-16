<script lang="ts">
	import { admin } from '$lib/stores/admin.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import StatsCard from '$lib/components/ui/StatsCard.svelte';
	import StatsCardSkeleton from '$lib/components/ui/StatsCardSkeleton.svelte';
	import ListSkeleton from '$lib/components/ui/ListSkeleton.svelte';
	import PendingRequests from '$lib/components/features/admin/PendingRequests.svelte';
	import { Users, Umbrella, CalendarDays } from 'lucide-svelte';
	import type { TeamVacation } from '$lib/types';

	interface Props {
		isLoading?: boolean;
		onUpdate?: () => void;
		teamVacations?: TeamVacation[];
	}

	let { isLoading = false, onUpdate = () => {}, teamVacations = [] }: Props = $props();

	function getDateStr(date: Date): string {
		return date.toISOString().split('T')[0];
	}

	const onLeaveToday = $derived.by(() => {
		const today = getDateStr(new Date());
		return teamVacations.filter((v) => v.startDate <= today && v.endDate >= today).length;
	});

	const upcomingThisWeek = $derived.by(() => {
		const now = new Date();
		const today = getDateStr(now);
		const weekFromNow = getDateStr(new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000));
		return teamVacations.filter((v) => v.startDate > today && v.startDate <= weekFromNow).length;
	});
</script>

<div class="space-y-6">
	<!-- Stats Grid -->
	<div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
		{#if isLoading}
			<StatsCardSkeleton count={3} />
		{:else}
			<div class="contents content-fade-in">
				<StatsCard
					title="Total Crew"
					value={admin.pagination.total}
					icon={Users}
					color="ocean"
				/>
				<StatsCard
					title="On Leave Today"
					value={onLeaveToday}
					icon={Umbrella}
					color="coral"
				/>
				<StatsCard
					title="Upcoming This Week"
					value={upcomingThisWeek}
					icon={CalendarDays}
					color="ocean"
				/>
			</div>
		{/if}
	</div>

	<!-- Pending Requests (only shown when there are requests or loading) -->
	{#if isLoading || admin.pendingRequests.length > 0}
		<Card>
			{#snippet header()}
				<h2 class="text-lg font-semibold text-ocean-700">
					Pending Requests ({admin.pendingRequests.length})
				</h2>
			{/snippet}

			{#if isLoading}
				<ListSkeleton count={3} variant="withActions" />
			{:else}
				<div class="content-fade-in">
					<PendingRequests requests={admin.pendingRequests} {onUpdate} />
				</div>
			{/if}
		</Card>
	{/if}
</div>
