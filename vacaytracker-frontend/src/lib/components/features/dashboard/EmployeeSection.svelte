<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { vacation } from '$lib/stores/vacation.svelte';
	import BalanceDisplay from '$lib/components/features/vacation/BalanceDisplay.svelte';
	import RequestList from '$lib/components/features/vacation/RequestList.svelte';
	import TeamVacationList from '$lib/components/features/vacation/TeamVacationList.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import StatsCard from '$lib/components/ui/StatsCard.svelte';
	import ListSkeleton from '$lib/components/ui/ListSkeleton.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import { Umbrella, Users, Clock, CheckCircle, Calculator } from 'lucide-svelte';
	import type { TeamVacation } from '$lib/types';

	interface Props {
		teamVacations?: TeamVacation[];
		isLoadingTeam?: boolean;
		onRequestVacation?: () => void;
	}

	let { teamVacations = [], isLoadingTeam = false, onRequestVacation }: Props = $props();
</script>

<div class="space-y-6">
	<!-- Balance and Quick Stats -->
	<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
		<!-- Balance Card -->
		<Card class="md:col-span-1">
			<div class="flex flex-col items-center py-4">
				<h2 class="text-lg font-semibold text-ocean-700 mb-4">Your Balance</h2>
				<BalanceDisplay current={auth.user?.vacationBalance ?? 0} total={25} />
				<p class="mt-4 text-sm text-ocean-500">Days remaining this year</p>
			</div>
		</Card>

		<!-- Quick Stats using StatsCard -->
		<StatsCard
			title="Pending"
			value={vacation.pendingRequests.length}
			icon={Clock}
			color="yellow"
		/>
		<StatsCard
			title="Approved"
			value={vacation.approvedRequests.length}
			icon={CheckCircle}
			color="green"
		/>
		<StatsCard title="Days Used" value={vacation.totalDaysUsed} icon={Calculator} color="ocean" />
	</div>

	<!-- Team Vacations This Month -->
	<Card>
		{#snippet header()}
			<div class="flex items-center justify-between">
				<h2 class="text-lg font-semibold text-ocean-700">Team Vacations This Month</h2>
				<Users class="w-5 h-5 text-ocean-400" />
			</div>
		{/snippet}

		<TeamVacationList vacations={teamVacations} isLoading={isLoadingTeam} />
	</Card>

	<!-- Your Requests -->
	<Card>
		{#snippet header()}
			<div class="flex items-center justify-between">
				<h2 class="text-lg font-semibold text-ocean-700">Your Requests</h2>
				<Umbrella class="w-5 h-5 text-ocean-400" />
			</div>
		{/snippet}

		{#if vacation.isLoading}
			<ListSkeleton count={3} variant="simple" />
		{:else if vacation.requests.length === 0}
			<div class="content-fade-in">
				<EmptyState icon={Umbrella} message="No vacation requests yet" iconSize="lg">
					{#if onRequestVacation}
						<Button variant="outline" onclick={onRequestVacation}>
							Request Your First Vacation
						</Button>
					{/if}
				</EmptyState>
			</div>
		{:else}
			<div class="content-fade-in">
				<RequestList requests={vacation.requests} />
			</div>
		{/if}
	</Card>
</div>
