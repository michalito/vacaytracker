<script lang="ts">
	import { createPopover, createCheckbox, melt } from '@melt-ui/svelte';
	import { clsx } from 'clsx';
	import { calendar, type UserInfo } from '$lib/stores/calendar.svelte';
	import { getUserColor } from '$lib/utils/colors';
	import { Filter, X, Check } from 'lucide-svelte';

	interface Props {
		users: UserInfo[];
		class?: string;
	}

	let { users, class: className = '' }: Props = $props();

	// Popover for the filter dropdown
	const {
		elements: { trigger, content, close },
		states: { open }
	} = createPopover({
		forceVisible: true,
		positioning: { placement: 'bottom-end' }
	});

	// Derive filter display state from store
	const showAll = $derived(calendar.showAll);
	const selectedUserIds = $derived(calendar.selectedUserIds);
	const selectedCount = $derived(selectedUserIds.size);

	// For trigger button display
	const hasActiveFilter = $derived(!showAll);
	const filterLabel = $derived.by(() => {
		if (showAll) return 'Filter';
		if (selectedCount === 0) return 'Filter (None)';
		return `Filter (${selectedCount})`;
	});

	// Check if a user is visually selected
	function isUserSelected(userId: string): boolean {
		if (showAll) return true;
		return selectedUserIds.has(userId);
	}

	// Handlers
	function handleSelectAll() {
		calendar.clearFilters();
	}

	function handleSelectNone() {
		calendar.selectNone();
	}

	function handleToggleUser(userId: string) {
		calendar.toggleUserFilter(userId);
	}
</script>

<div class={clsx('relative', className)}>
	<button
		use:melt={$trigger}
		class={clsx(
			'inline-flex items-center justify-center gap-1 px-3 py-1.5 text-sm font-medium rounded-lg transition-all duration-200 cursor-pointer',
			hasActiveFilter
				? 'bg-ocean-500 text-white hover:bg-ocean-600'
				: 'border-2 border-ocean-500/40 text-ocean-700 hover:border-ocean-500 hover:bg-ocean-50'
		)}
	>
		<Filter class="w-4 h-4" />
		{filterLabel}
	</button>

	{#if $open}
		<div
			use:melt={$content}
			class="absolute right-0 mt-2 z-50 w-64 bg-white rounded-lg shadow-lg border border-sand-200 p-3"
		>
			<div class="flex items-center justify-between mb-3">
				<span class="text-sm font-medium text-ocean-700">Filter by User</span>
				<button
					use:melt={$close}
					class="text-ocean-400 hover:text-ocean-600 p-1 cursor-pointer"
				>
					<X class="w-4 h-4" />
				</button>
			</div>

			<!-- Mode toggle buttons -->
			<div class="flex gap-1 mb-3 p-1 bg-sand-100 rounded-lg">
				<button
					onclick={handleSelectAll}
					class={clsx(
						'flex-1 px-3 py-1.5 text-sm font-medium rounded-md transition-all duration-200 cursor-pointer',
						showAll
							? 'bg-white text-ocean-700 shadow-sm'
							: 'text-ocean-500 hover:text-ocean-700'
					)}
				>
					All
				</button>
				<button
					onclick={handleSelectNone}
					class={clsx(
						'flex-1 px-3 py-1.5 text-sm font-medium rounded-md transition-all duration-200 cursor-pointer',
						!showAll && selectedCount === 0
							? 'bg-white text-ocean-700 shadow-sm'
							: 'text-ocean-500 hover:text-ocean-700'
					)}
				>
					None
				</button>
			</div>

			{#if users.length === 0}
				<p class="text-sm text-ocean-500 text-center py-2">No users with vacations</p>
			{:else}
				<div class="max-h-48 overflow-y-auto space-y-1">
					{#each users as user (user.id)}
						{@const color = getUserColor(user.id)}
						{@const selected = isUserSelected(user.id)}
						{@const checkbox = createCheckbox({ defaultChecked: selected })}
						<button
							use:melt={checkbox.elements.root}
							onclick={() => handleToggleUser(user.id)}
							class={clsx(
								'w-full flex items-center gap-2 p-2 rounded cursor-pointer transition-colors text-left',
								selected ? 'bg-ocean-50' : 'hover:bg-sand-50 opacity-60'
							)}
						>
							<span
								class={clsx(
									'h-5 w-5 shrink-0 rounded border-2 flex items-center justify-center transition-all duration-200',
									selected
										? 'border-ocean-500 bg-ocean-500'
										: 'border-ocean-400 bg-white'
								)}
							>
								{#if selected}
									<Check class="h-3.5 w-3.5 text-white" />
								{/if}
							</span>
							<span class={clsx('w-3 h-3 rounded-full flex-shrink-0', color.background)}></span>
							<span class="text-sm text-ocean-700 truncate">{user.name}</span>
						</button>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>
