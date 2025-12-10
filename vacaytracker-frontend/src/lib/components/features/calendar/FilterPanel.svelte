<script lang="ts">
	import { clsx } from 'clsx';
	import { calendar, type UserInfo } from '$lib/stores/calendar.svelte';
	import { getUserColor } from '$lib/utils/colors';
	import Button from '$lib/components/ui/Button.svelte';
	import { Filter, X } from 'lucide-svelte';

	interface Props {
		users: UserInfo[];
		class?: string;
	}

	let { users, class: className = '' }: Props = $props();

	let isExpanded = $state(false);

	const selectedCount = $derived(calendar.selectedUserIds.size);
	const hasFilters = $derived(selectedCount > 0);

	function isUserSelected(userId: string): boolean {
		// If no filters are set, all users are effectively "selected"
		if (calendar.selectedUserIds.size === 0) return true;
		return calendar.selectedUserIds.has(userId);
	}

	function handleToggle(userId: string) {
		calendar.toggleUserFilter(userId);
	}

	function handleSelectAll() {
		calendar.clearFilters();
	}

	function handleSelectNone() {
		// Select all users first, then they can toggle individually
		calendar.selectAllUsers();
	}
</script>

<div class={clsx('relative', className)}>
	<Button
		variant={hasFilters ? 'primary' : 'outline'}
		size="sm"
		onclick={() => (isExpanded = !isExpanded)}
	>
		<Filter class="w-4 h-4 mr-1" />
		Filter
		{#if hasFilters}
			<span class="ml-1">({selectedCount})</span>
		{/if}
	</Button>

	{#if isExpanded}
		<!-- Backdrop -->
		<button
			class="fixed inset-0 z-10"
			onclick={() => (isExpanded = false)}
			aria-label="Close filter panel"
		></button>

		<!-- Panel -->
		<div
			class="absolute top-full right-0 mt-2 w-64 bg-white rounded-lg shadow-lg border border-sand-200 p-3 z-20"
		>
			<div class="flex items-center justify-between mb-2">
				<span class="text-sm font-medium text-ocean-700">Filter by User</span>
				<button
					onclick={() => (isExpanded = false)}
					class="text-ocean-400 hover:text-ocean-600 p-1"
				>
					<X class="w-4 h-4" />
				</button>
			</div>

			<div class="flex gap-2 mb-3">
				<Button variant="ghost" size="sm" onclick={handleSelectAll}>All</Button>
				<Button variant="ghost" size="sm" onclick={handleSelectNone}>None</Button>
				{#if hasFilters}
					<Button variant="ghost" size="sm" onclick={() => calendar.clearFilters()}>
						Clear
					</Button>
				{/if}
			</div>

			{#if users.length === 0}
				<p class="text-sm text-ocean-500 text-center py-2">No users with vacations</p>
			{:else}
				<div class="max-h-48 overflow-y-auto space-y-1">
					{#each users as user (user.id)}
						{@const color = getUserColor(user.id)}
						{@const selected = isUserSelected(user.id)}
						<label
							class={clsx(
								'flex items-center gap-2 p-2 rounded cursor-pointer transition-colors',
								selected ? 'bg-ocean-50' : 'hover:bg-sand-50 opacity-60'
							)}
						>
							<input
								type="checkbox"
								checked={selected}
								onchange={() => handleToggle(user.id)}
								class="rounded border-ocean-300 text-ocean-500 focus:ring-ocean-500"
							/>
							<span class={clsx('w-3 h-3 rounded-full flex-shrink-0', color.background)}></span>
							<span class="text-sm text-ocean-700 truncate">{user.name}</span>
						</label>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>
