<script lang="ts">
	import type { User } from '$lib/types';
	import { admin } from '$lib/stores/admin.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Check, X, Edit2 } from 'lucide-svelte';

	interface Props {
		user: User;
	}

	let { user }: Props = $props();

	let isEditing = $state(false);
	let editValue = $state(0);
	let isSaving = $state(false);

	// Sync editValue with user.vacationBalance when not editing
	$effect(() => {
		if (!isEditing) {
			editValue = user.vacationBalance;
		}
	});

	function startEdit() {
		editValue = user.vacationBalance;
		isEditing = true;
	}

	function cancelEdit() {
		isEditing = false;
		editValue = user.vacationBalance;
	}

	async function saveEdit() {
		if (editValue === user.vacationBalance) {
			isEditing = false;
			return;
		}

		isSaving = true;
		try {
			await admin.updateBalance(user.id, editValue);
			toast.success(`Updated balance for ${user.name}`);
			isEditing = false;
		} catch (error) {
			toast.error('Failed to update balance');
		} finally {
			isSaving = false;
		}
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			saveEdit();
		} else if (e.key === 'Escape') {
			cancelEdit();
		}
	}
</script>

<div class="flex items-center justify-between py-3 px-4 hover:bg-sand-50 rounded-lg transition-colors">
	<div class="flex items-center gap-3">
		<Avatar name={user.name} size="sm" />
		<div>
			<p class="font-medium text-ocean-800">{user.name}</p>
			<p class="text-sm text-ocean-500">{user.email}</p>
		</div>
	</div>

	<div class="flex items-center gap-2">
		{#if isEditing}
			<input
				type="number"
				bind:value={editValue}
				onkeydown={handleKeyDown}
				min="0"
				max="365"
				class="w-20 px-2 py-1 text-center rounded border border-ocean-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent"
			/>
			<Button variant="ghost" size="sm" onclick={saveEdit} loading={isSaving}>
				<Check class="w-4 h-4 text-green-600" />
			</Button>
			<Button variant="ghost" size="sm" onclick={cancelEdit} disabled={isSaving}>
				<X class="w-4 h-4 text-error" />
			</Button>
		{:else}
			<span class="font-bold text-ocean-800 w-12 text-center">{user.vacationBalance}</span>
			<span class="text-ocean-500 text-sm">days</span>
			<Button variant="ghost" size="sm" onclick={startEdit}>
				<Edit2 class="w-4 h-4" />
			</Button>
		{/if}
	</div>
</div>
