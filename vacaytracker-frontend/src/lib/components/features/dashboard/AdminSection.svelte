<script lang="ts">
	import { admin } from '$lib/stores/admin.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import StatsCard from '$lib/components/ui/StatsCard.svelte';
	import StatsCardSkeleton from '$lib/components/ui/StatsCardSkeleton.svelte';
	import ListSkeleton from '$lib/components/ui/ListSkeleton.svelte';
	import PendingRequests from '$lib/components/features/admin/PendingRequests.svelte';
	import { Users, Clock, CheckCircle, XCircle } from 'lucide-svelte';

	interface Props {
		isLoading?: boolean;
		onUpdate?: () => void;
	}

	let { isLoading = false, onUpdate = () => {} }: Props = $props();

	// TODO: Track approved/rejected today from API response if available
	let approvedToday = $state(0);
	let rejectedToday = $state(0);
</script>

<div class="space-y-6">
	<!-- Stats Grid -->
	<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
		{#if isLoading}
			<StatsCardSkeleton count={4} />
		{:else}
			<StatsCard title="Total Crew" value={admin.pagination.total} icon={Users} color="ocean" />
			<StatsCard
				title="Pending Requests"
				value={admin.pendingRequests.length}
				icon={Clock}
				color="yellow"
			/>
			<StatsCard title="Approved Today" value={approvedToday} icon={CheckCircle} color="green" />
			<StatsCard title="Rejected Today" value={rejectedToday} icon={XCircle} color="red" />
		{/if}
	</div>

	<!-- Pending Requests -->
	<Card>
		{#snippet header()}
			<h2 class="text-lg font-semibold text-ocean-700">
				Pending Requests ({admin.pendingRequests.length})
			</h2>
		{/snippet}

		{#if isLoading}
			<ListSkeleton count={3} variant="withActions" />
		{:else}
			<PendingRequests requests={admin.pendingRequests} {onUpdate} />
		{/if}
	</Card>
</div>
