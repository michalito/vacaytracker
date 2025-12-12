<script lang="ts">
	import { admin } from '$lib/stores/admin.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import Pagination from '$lib/components/ui/Pagination.svelte';
	import UserTable from '$lib/components/features/admin/UserTable.svelte';
	import UserModal from '$lib/components/features/admin/UserModal.svelte';
	import DeleteConfirm from '$lib/components/features/admin/DeleteConfirm.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';
	import { UserPlus, Search, Users } from 'lucide-svelte';
	import type { User, Role } from '$lib/types';

	const roleFilterOptions = [
		{ value: '', label: 'All Roles' },
		{ value: 'admin', label: 'Admin' },
		{ value: 'employee', label: 'Employee' }
	];

	let searchQuery = $state('');
	let roleFilter = $state('');
	let isUserModalOpen = $state(false);
	let editingUser = $state<User | null>(null);
	let deletingUser = $state<User | null>(null);
	let isDeleting = $state(false);

	$effect(() => {
		loadUsers();
	});

	function loadUsers(page = 1) {
		admin.fetchUsers({
			page,
			limit: 10,
			search: searchQuery || undefined,
			role: (roleFilter as Role) || undefined
		});
	}

	function handleSearch() {
		loadUsers(1);
	}

	function openCreateModal() {
		editingUser = null;
		isUserModalOpen = true;
	}

	function openEditModal(user: User) {
		editingUser = user;
		isUserModalOpen = true;
	}

	function openDeleteModal(user: User) {
		deletingUser = user;
	}

	async function handleDelete() {
		if (!deletingUser) return;

		isDeleting = true;
		try {
			await admin.deleteUser(deletingUser.id);
			toast.success('User deleted');
			deletingUser = null;
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to delete user');
		} finally {
			isDeleting = false;
		}
	}
</script>

<svelte:head>
	<title>User Management - VacayTracker</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-ocean-800">User Management</h1>
			<p class="text-ocean-600">Manage your team members</p>
		</div>
		<Button onclick={openCreateModal}>
			<UserPlus class="w-4 h-4 mr-2" />
			Add User
		</Button>
	</div>

	<Card>
		{#snippet header()}
			<div class="flex flex-col sm:flex-row gap-4">
				<div class="flex-1 relative">
					<Search class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-ocean-400" />
					<input
						type="text"
						placeholder="Search users..."
						bind:value={searchQuery}
						onkeydown={(e) => e.key === 'Enter' && handleSearch()}
						class="w-full pl-11 pr-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 placeholder-ocean-500/50 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500 hover:border-ocean-500"
					/>
				</div>
				<div class="w-40">
					<Select
						bind:value={roleFilter}
						options={roleFilterOptions}
						onchange={() => loadUsers(1)}
					/>
				</div>
				<Button variant="outline" onclick={handleSearch}>
					<Search class="w-4 h-4 mr-2" />
					Search
				</Button>
			</div>
		{/snippet}

		{#if admin.isLoading}
			<div class="space-y-3 py-4">
				{#each [1, 2, 3, 4, 5] as _}
					<div class="flex items-center gap-4 p-4 border border-sand-200 rounded-lg">
						<Skeleton variant="avatar" />
						<div class="flex-1 space-y-2">
							<Skeleton variant="text" width="30%" />
							<Skeleton variant="text" width="50%" />
						</div>
						<Skeleton variant="text" width="60px" />
						<Skeleton variant="text" width="80px" />
					</div>
				{/each}
			</div>
		{:else if admin.users.length === 0}
			<div class="py-8 text-center">
				<Users class="w-12 h-12 mx-auto text-ocean-300 mb-2" />
				<p class="text-ocean-500">No users found</p>
			</div>
		{:else}
			<UserTable users={admin.users} onEdit={openEditModal} onDelete={openDeleteModal} />

			{#if admin.pagination.totalPages > 1}
				<div class="mt-4 pt-4 border-t border-sand-200">
					<Pagination
						page={admin.pagination.page}
						totalPages={admin.pagination.totalPages}
						onPageChange={loadUsers}
					/>
				</div>
			{/if}
		{/if}

		{#snippet footer()}
			<p class="text-sm text-ocean-500">
				Showing {admin.users.length} of {admin.pagination.total} users
			</p>
		{/snippet}
	</Card>
</div>

<UserModal
	bind:open={isUserModalOpen}
	user={editingUser}
	onClose={() => {
		isUserModalOpen = false;
		editingUser = null;
	}}
/>

<DeleteConfirm
	open={!!deletingUser}
	title="Delete User"
	message="Are you sure you want to delete {deletingUser?.name}? This will remove all their data and cannot be undone."
	isLoading={isDeleting}
	onConfirm={handleDelete}
	onCancel={() => (deletingUser = null)}
/>
