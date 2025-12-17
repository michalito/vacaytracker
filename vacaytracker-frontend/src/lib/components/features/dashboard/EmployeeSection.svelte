<script lang="ts">
	import { vacation } from '$lib/stores/vacation.svelte';
	import VacationBalanceCard from '$lib/components/features/dashboard/VacationBalanceCard.svelte';
	import RequestTabs from '$lib/components/features/vacation/RequestTabs.svelte';
	import TeamVacationList from '$lib/components/features/vacation/TeamVacationList.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import ListSkeleton from '$lib/components/ui/ListSkeleton.svelte';
	import { Users } from 'lucide-svelte';
	import type { TeamVacation } from '$lib/types';

	interface Props {
		teamVacations?: TeamVacation[];
		isLoadingTeam?: boolean;
		defaultVacationDays?: number;
		onRequestVacation?: () => void;
	}

	let {
		teamVacations = [],
		isLoadingTeam = false,
		defaultVacationDays = 25,
		onRequestVacation
	}: Props = $props();
</script>

<div class="space-y-6">
	<!-- Vacation Balance Card -->
	<VacationBalanceCard
		total={defaultVacationDays}
		used={vacation.usedDays}
		upcoming={vacation.upcomingDays}
		nextVacation={vacation.upcomingRequests[0] ?? null}
	/>

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

	<!-- Your Requests (Tabbed) -->
	{#if vacation.isLoading}
		<Card>
			{#snippet header()}
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-semibold text-ocean-700">Your Requests</h2>
				</div>
			{/snippet}
			<ListSkeleton count={3} variant="simple" />
		</Card>
	{:else}
		<RequestTabs
			pendingRequests={vacation.pendingRequests}
			upcomingRequests={vacation.upcomingRequests}
			pastRequests={vacation.pastRequests}
			{onRequestVacation}
		/>
	{/if}
</div>
