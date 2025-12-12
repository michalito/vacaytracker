<script lang="ts">
	import { admin } from '$lib/stores/admin.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import BalanceEditor from '$lib/components/features/admin/BalanceEditor.svelte';
	import Pagination from '$lib/components/ui/Pagination.svelte';
	import { Wallet, Users } from 'lucide-svelte';

	$effect(() => {
		admin.fetchUsers({ limit: 20 });
	});

	function handlePageChange(page: number) {
		admin.fetchUsers({ page, limit: 20 });
	}

	const totalBalance = $derived(admin.users.reduce((sum, u) => sum + u.vacationBalance, 0));
	const avgBalance = $derived(
		admin.users.length > 0 ? Math.round(totalBalance / admin.users.length) : 0
	);
</script>

<svelte:head>
	<title>Vacation Balances - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold text-ocean-800">Vacation Balances</h1>
		<p class="text-ocean-600">Manage vacation days for your team</p>
	</div>

	<!-- Summary Stats -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
		<div class="bg-white/80 backdrop-blur-sm rounded-xl shadow-md border border-white/30 p-5 transition-all duration-200 hover:shadow-lg">
			<div class="flex items-center gap-4">
				<div class="p-3 rounded-xl bg-ocean-500/10 text-ocean-600">
					<Users class="w-6 h-6" />
				</div>
				<div>
					<p class="text-2xl font-bold text-ocean-900">{admin.pagination.total}</p>
					<p class="text-sm text-ocean-500">Total Users</p>
				</div>
			</div>
		</div>

		<div class="bg-white/80 backdrop-blur-sm rounded-xl shadow-md border border-white/30 p-5 transition-all duration-200 hover:shadow-lg">
			<div class="flex items-center gap-4">
				<div class="p-3 rounded-xl bg-ocean-400/10 text-ocean-500">
					<Wallet class="w-6 h-6" />
				</div>
				<div>
					<p class="text-2xl font-bold text-ocean-900">{totalBalance}</p>
					<p class="text-sm text-ocean-500">Total Days Available</p>
				</div>
			</div>
		</div>

		<div class="bg-white/80 backdrop-blur-sm rounded-xl shadow-md border border-white/30 p-5 transition-all duration-200 hover:shadow-lg">
			<div class="flex items-center gap-4">
				<div class="p-3 rounded-xl bg-coral-400/10 text-coral-500">
					<Wallet class="w-6 h-6" />
				</div>
				<div>
					<p class="text-2xl font-bold text-ocean-900">{avgBalance}</p>
					<p class="text-sm text-ocean-500">Average Balance</p>
				</div>
			</div>
		</div>
	</div>

	<Card>
		{#snippet header()}
			<div class="flex items-center gap-2">
				<Wallet class="w-5 h-5 text-ocean-500" />
				<h2 class="text-lg font-semibold text-ocean-700">User Balances</h2>
			</div>
		{/snippet}

		{#if admin.isLoading}
			<div class="py-8 text-center text-ocean-500">Loading balances...</div>
		{:else if admin.users.length === 0}
			<div class="py-8 text-center">
				<Users class="w-12 h-12 mx-auto text-ocean-300 mb-2" />
				<p class="text-ocean-500">No users found</p>
			</div>
		{:else}
			<div class="divide-y divide-sand-100">
				{#each admin.users as user (user.id)}
					<BalanceEditor {user} />
				{/each}
			</div>

			{#if admin.pagination.totalPages > 1}
				<div class="mt-4 pt-4 border-t border-sand-200">
					<Pagination
						page={admin.pagination.page}
						totalPages={admin.pagination.totalPages}
						onPageChange={handlePageChange}
					/>
				</div>
			{/if}
		{/if}
	</Card>
</div>
