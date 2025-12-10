<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { vacation } from '$lib/stores/vacation.svelte';
	import { vacationApi } from '$lib/api/vacation';
	import BalanceDisplay from '$lib/components/features/vacation/BalanceDisplay.svelte';
	import RequestList from '$lib/components/features/vacation/RequestList.svelte';
	import RequestModal from '$lib/components/features/vacation/RequestModal.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import { Plus, Umbrella, Users } from 'lucide-svelte';
	import type { TeamVacation } from '$lib/types';

	let isRequestModalOpen = $state(false);
	let teamVacations = $state<TeamVacation[]>([]);
	let isLoadingTeam = $state(true);

	$effect(() => {
		vacation.fetchRequests();
		loadTeamVacations();
	});

	async function loadTeamVacations() {
		isLoadingTeam = true;
		try {
			const now = new Date();
			const response = await vacationApi.team(now.getMonth() + 1, now.getFullYear());
			teamVacations = response.vacations;
		} catch (error) {
			console.error('Failed to load team vacations:', error);
			teamVacations = [];
		} finally {
			isLoadingTeam = false;
		}
	}

	function formatDateRange(startDate: string, endDate: string): string {
		const start = new Date(startDate);
		const end = new Date(endDate);
		const startStr = start.toLocaleDateString('en-GB', { day: 'numeric', month: 'short' });
		const endStr = end.toLocaleDateString('en-GB', { day: 'numeric', month: 'short' });
		return `${startStr} - ${endStr}`;
	}
</script>

<svelte:head>
	<title>Dashboard - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
	<!-- Welcome Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-ocean-800">
				Welcome back, {auth.user?.name?.split(' ')[0]}!
			</h1>
			<p class="text-ocean-600">Here's your vacation overview</p>
		</div>
		<Button onclick={() => (isRequestModalOpen = true)}>
			<Plus class="w-4 h-4 mr-2" />
			Request Vacation
		</Button>
	</div>

	<!-- Balance Section -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<Card class="md:col-span-1">
			<div class="flex flex-col items-center py-4">
				<h2 class="text-lg font-semibold text-ocean-700 mb-4">Your Balance</h2>
				<BalanceDisplay current={auth.user?.vacationBalance ?? 0} total={25} />
				<p class="mt-4 text-sm text-ocean-500">Days remaining this year</p>
			</div>
		</Card>

		<!-- Quick Stats -->
		<Card class="md:col-span-2">
			<h2 class="text-lg font-semibold text-ocean-700 mb-4">Quick Stats</h2>
			<div class="grid grid-cols-3 gap-4">
				<div class="text-center p-4 bg-yellow-50 rounded-lg">
					<p class="text-2xl font-bold text-yellow-600">
						{vacation.pendingRequests.length}
					</p>
					<p class="text-sm text-yellow-700">Pending</p>
				</div>
				<div class="text-center p-4 bg-green-50 rounded-lg">
					<p class="text-2xl font-bold text-green-600">
						{vacation.approvedRequests.length}
					</p>
					<p class="text-sm text-green-700">Approved</p>
				</div>
				<div class="text-center p-4 bg-ocean-50 rounded-lg">
					<p class="text-2xl font-bold text-ocean-600">
						{vacation.approvedRequests.reduce((sum, r) => sum + r.totalDays, 0)}
					</p>
					<p class="text-sm text-ocean-700">Days Used</p>
				</div>
			</div>
		</Card>
	</div>

	<!-- Team Vacations This Month -->
	<Card>
		{#snippet header()}
			<div class="flex items-center justify-between">
				<h2 class="text-lg font-semibold text-ocean-700">Team Vacations This Month</h2>
				<Users class="w-5 h-5 text-ocean-400" />
			</div>
		{/snippet}

		{#if isLoadingTeam}
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
		{:else if teamVacations.length === 0}
			<div class="py-6 text-center">
				<Users class="w-10 h-10 mx-auto text-ocean-300 mb-2" />
				<p class="text-ocean-500">No team vacations this month</p>
			</div>
		{:else}
			<div class="space-y-2">
				{#each teamVacations.slice(0, 5) as v (v.id)}
					<div class="flex items-center gap-3 p-2 bg-sand-50 rounded-lg">
						<Avatar name={v.userName} size="sm" />
						<div class="flex-1 min-w-0">
							<p class="font-medium text-ocean-800 truncate">{v.userName}</p>
							<p class="text-sm text-ocean-500">
								{formatDateRange(v.startDate, v.endDate)} ({v.totalDays} days)
							</p>
						</div>
					</div>
				{/each}
				{#if teamVacations.length > 5}
					<a href="/employee/calendar" class="block text-center text-sm text-ocean-600 hover:text-ocean-700 pt-2">
						View all {teamVacations.length} vacations
					</a>
				{/if}
			</div>
		{/if}
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
			<div class="space-y-4 py-4">
				{#each [1, 2, 3] as _}
					<div class="flex items-center gap-4 p-4 border border-sand-200 rounded-lg">
						<Skeleton variant="avatar" />
						<div class="flex-1 space-y-2">
							<Skeleton variant="text" width="40%" />
							<Skeleton variant="text" width="60%" />
						</div>
						<Skeleton variant="text" width="80px" />
					</div>
				{/each}
			</div>
		{:else if vacation.requests.length === 0}
			<div class="py-8 text-center">
				<Umbrella class="w-12 h-12 mx-auto text-ocean-300 mb-2" />
				<p class="text-ocean-500">No vacation requests yet</p>
				<Button variant="outline" class="mt-4" onclick={() => (isRequestModalOpen = true)}>
					Request Your First Vacation
				</Button>
			</div>
		{:else}
			<RequestList requests={vacation.requests} />
		{/if}
	</Card>
</div>

<RequestModal bind:open={isRequestModalOpen} />
