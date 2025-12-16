<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { admin } from '$lib/stores/admin.svelte';
	import { vacation } from '$lib/stores/vacation.svelte';
	import { calendar } from '$lib/stores/calendar.svelte';
	import AdminSection from '$lib/components/features/dashboard/AdminSection.svelte';
	import EmployeeSection from '$lib/components/features/dashboard/EmployeeSection.svelte';
	import RequestModal from '$lib/components/features/vacation/RequestModal.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Plus } from 'lucide-svelte';

	let isRequestModalOpen = $state(false);
	let isLoadingTeam = $state(true);
	let isLoadingAdmin = $state(true);
	let hasMounted = $state(false);

	$effect(() => {
		if (hasMounted) return;
		hasMounted = true;
		loadData();
	});

	async function loadData() {
		// Load employee data (with caching)
		vacation.fetchRequests();

		// Use shared calendar store for team vacations (with caching)
		isLoadingTeam = true;
		await calendar.fetchCurrentMonth();
		isLoadingTeam = false;

		// Load admin data if user is admin (with caching)
		if (auth.isAdmin) {
			isLoadingAdmin = true;
			await Promise.all([admin.fetchPendingRequests(), admin.fetchUsers({ limit: 1 })]);
			isLoadingAdmin = false;
		}
	}

	async function handleUpdate() {
		// Only refresh pending requests after approve/reject (targeted update)
		await admin.fetchPendingRequests({ force: true });
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

	<!-- Admin Section (only for admins) -->
	{#if auth.isAdmin}
		<AdminSection
			isLoading={isLoadingAdmin}
			onUpdate={handleUpdate}
			teamVacations={calendar.currentMonthVacations}
		/>
	{/if}

	<!-- Employee Section (for everyone) -->
	<EmployeeSection
		teamVacations={calendar.currentMonthVacations}
		{isLoadingTeam}
		onRequestVacation={() => (isRequestModalOpen = true)}
	/>
</div>

<RequestModal bind:open={isRequestModalOpen} />
