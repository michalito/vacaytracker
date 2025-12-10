<script lang="ts">
	import { admin } from '$lib/stores/admin.svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { vacationApi } from '$lib/api/vacation';
	import Card from '$lib/components/ui/Card.svelte';
	import StatsCard from '$lib/components/features/admin/StatsCard.svelte';
	import PendingRequests from '$lib/components/features/admin/PendingRequests.svelte';
	import BalanceDisplay from '$lib/components/features/vacation/BalanceDisplay.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';
	import { Users, Clock, CheckCircle, XCircle, Palmtree, ArrowRight } from 'lucide-svelte';
	import type { TeamVacation } from '$lib/types';

	let approvedToday = $state(0);
	let rejectedToday = $state(0);
	let isLoading = $state(true);
	let teamVacations = $state<TeamVacation[]>([]);

	const currentMonth = new Date().getMonth() + 1;
	const currentYear = new Date().getFullYear();

	$effect(() => {
		loadData();
	});

	async function loadData() {
		isLoading = true;
		await Promise.all([
			admin.fetchPendingRequests(),
			admin.fetchUsers({ limit: 1 }),
			loadTeamVacations()
		]);
		isLoading = false;
	}

	async function loadTeamVacations() {
		try {
			const response = await vacationApi.team(currentMonth, currentYear);
			teamVacations = response.vacations || [];
		} catch {
			teamVacations = [];
		}
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('en-GB', {
			day: 'numeric',
			month: 'short'
		});
	}
</script>

<svelte:head>
	<title>Admin Dashboard - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
	<h1 class="text-2xl font-bold text-ocean-800">Captain's Dashboard</h1>

	<!-- Stats Grid -->
	<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
		{#if isLoading}
			{#each [1, 2, 3, 4] as _}
				<div class="bg-white rounded-lg shadow-sm p-6 border border-sand-200">
					<div class="flex items-center justify-between">
						<div class="space-y-2 flex-1">
							<Skeleton variant="text" width="60%" />
							<Skeleton variant="text" width="40%" height="2rem" />
						</div>
						<Skeleton variant="circle" width="48px" height="48px" />
					</div>
				</div>
			{/each}
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
			<div class="space-y-4 py-4">
				{#each [1, 2, 3] as _}
					<div class="flex items-center gap-4 p-4 border border-sand-200 rounded-lg">
						<Skeleton variant="avatar" />
						<div class="flex-1 space-y-2">
							<Skeleton variant="text" width="35%" />
							<Skeleton variant="text" width="55%" />
						</div>
						<div class="flex gap-2">
							<Skeleton variant="text" width="70px" height="32px" />
							<Skeleton variant="text" width="70px" height="32px" />
						</div>
					</div>
				{/each}
			</div>
		{:else}
			<PendingRequests requests={admin.pendingRequests} onUpdate={loadData} />
		{/if}
	</Card>

	<!-- My Vacations Section -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<!-- Balance Card -->
		<Card class="md:col-span-1">
			{#snippet header()}
				<div class="flex items-center gap-2">
					<Palmtree class="w-5 h-5 text-ocean-600" />
					<h2 class="text-lg font-semibold text-ocean-700">My Balance</h2>
				</div>
			{/snippet}

			<div class="flex flex-col items-center py-4">
				{#if isLoading}
					<Skeleton variant="circle" width="120px" height="120px" />
				{:else}
					<BalanceDisplay current={auth.user?.vacationBalance ?? 0} total={25} size={120} />
					<p class="mt-3 text-sm text-ocean-500">Days remaining</p>
					<a
						href="/employee"
						class="mt-4 inline-flex items-center gap-1 text-sm font-medium text-ocean-600 hover:text-ocean-700"
					>
						Request vacation
						<ArrowRight class="w-4 h-4" />
					</a>
				{/if}
			</div>
		</Card>

		<!-- Team Vacations This Month -->
		<Card class="md:col-span-2">
			{#snippet header()}
				<div class="flex items-center justify-between">
					<h2 class="text-lg font-semibold text-ocean-700">Team Vacations This Month</h2>
					<a
						href="/employee/calendar"
						class="text-sm font-medium text-ocean-600 hover:text-ocean-700 flex items-center gap-1"
					>
						View calendar
						<ArrowRight class="w-4 h-4" />
					</a>
				</div>
			{/snippet}

			{#if isLoading}
				<div class="space-y-3 py-2">
					{#each [1, 2, 3] as _}
						<div class="flex items-center gap-3">
							<Skeleton variant="circle" width="32px" height="32px" />
							<div class="flex-1 space-y-1">
								<Skeleton variant="text" width="40%" />
								<Skeleton variant="text" width="60%" />
							</div>
						</div>
					{/each}
				</div>
			{:else if teamVacations.length === 0}
				<div class="py-8 text-center text-ocean-500">
					<p>No team vacations this month</p>
				</div>
			{:else}
				<ul class="divide-y divide-sand-200">
					{#each teamVacations.slice(0, 5) as vacation}
						<li class="py-3 flex items-center gap-3">
							<div
								class="w-8 h-8 rounded-full bg-ocean-100 text-ocean-700 flex items-center justify-center text-sm font-medium"
							>
								{vacation.userName.charAt(0).toUpperCase()}
							</div>
							<div class="flex-1 min-w-0">
								<p class="font-medium text-ocean-800 truncate">{vacation.userName}</p>
								<p class="text-sm text-ocean-500">
									{formatDate(vacation.startDate)} - {formatDate(vacation.endDate)}
									<span class="text-ocean-400">({vacation.totalDays} days)</span>
								</p>
							</div>
						</li>
					{/each}
				</ul>
				{#if teamVacations.length > 5}
					<p class="pt-3 text-sm text-ocean-500 text-center">
						+{teamVacations.length - 5} more
					</p>
				{/if}
			{/if}
		</Card>
	</div>
</div>
