<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import { Settings, LogOut, ChevronDown, User } from 'lucide-svelte';

	let isOpen = $state(false);

	function handleLogout() {
		auth.logout();
		goto('/');
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as Element;
		if (!target.closest('.user-menu')) {
			isOpen = false;
		}
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			isOpen = false;
		}
	}

	$effect(() => {
		if (isOpen) {
			document.addEventListener('click', handleClickOutside);
			document.addEventListener('keydown', handleKeyDown);
		}
		return () => {
			document.removeEventListener('click', handleClickOutside);
			document.removeEventListener('keydown', handleKeyDown);
		};
	});
</script>

<div class="relative user-menu">
	<button
		type="button"
		onclick={() => (isOpen = !isOpen)}
		class="flex items-center gap-2 p-1.5 rounded-lg hover:bg-sand-100 transition-colors"
	>
		<Avatar name={auth.user?.name} size="sm" />
		<span class="hidden sm:block text-sm font-medium text-ocean-700">
			{auth.user?.name}
		</span>
		<ChevronDown class="w-4 h-4 text-ocean-500 transition-transform {isOpen ? 'rotate-180' : ''}" />
	</button>

	{#if isOpen}
		<div
			class="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-lg border border-sand-200 py-1 z-50"
		>
			<!-- User Info -->
			<div class="px-4 py-3 border-b border-sand-100">
				<p class="font-medium text-ocean-800">{auth.user?.name}</p>
				<p class="text-sm text-ocean-500 truncate">{auth.user?.email}</p>
			</div>

			<!-- Menu Items -->
			<div class="py-1">
				<a
					href="/settings"
					onclick={() => (isOpen = false)}
					class="flex items-center gap-3 px-4 py-2 text-ocean-700 hover:bg-sand-50 transition-colors"
				>
					<User class="w-4 h-4" />
					<span>Profile</span>
				</a>
				<a
					href="/settings"
					onclick={() => (isOpen = false)}
					class="flex items-center gap-3 px-4 py-2 text-ocean-700 hover:bg-sand-50 transition-colors"
				>
					<Settings class="w-4 h-4" />
					<span>Settings</span>
				</a>
			</div>

			<!-- Logout -->
			<div class="border-t border-sand-100 py-1">
				<button
					type="button"
					onclick={handleLogout}
					class="flex items-center gap-3 w-full px-4 py-2 text-error hover:bg-red-50 transition-colors"
				>
					<LogOut class="w-4 h-4" />
					<span>Sign out</span>
				</button>
			</div>
		</div>
	{/if}
</div>
