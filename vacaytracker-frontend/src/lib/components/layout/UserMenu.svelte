<script lang="ts">
	import { createDropdownMenu, melt } from '@melt-ui/svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import { Settings, LogOut, ChevronDown } from 'lucide-svelte';

	// Create Melt-UI dropdown menu
	const {
		elements: { trigger, menu, item, separator },
		states: { open }
	} = createDropdownMenu({
		forceVisible: true,
		positioning: { placement: 'bottom-end' }
	});

	function handleLogout() {
		auth.logout();
		goto('/');
	}
</script>

<div class="relative">
	<button
		use:melt={$trigger}
		class="flex items-center gap-2 p-1.5 rounded-xl hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
	>
		<Avatar name={auth.user?.name} size="sm" />
		<span class="hidden sm:block text-sm font-medium text-ocean-700">
			{auth.user?.name}
		</span>
		<ChevronDown class="w-4 h-4 text-ocean-500 transition-transform {$open ? 'rotate-180' : ''}" />
	</button>

	{#if $open}
		<div
			use:melt={$menu}
			class="absolute right-0 mt-2 w-56 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/30 py-1 z-50
				transition-all duration-200 data-[state=open]:animate-scale-in"
		>
			<!-- User Info -->
			<div class="px-4 py-3 border-b border-ocean-100/50">
				<p class="font-medium text-ocean-800">{auth.user?.name}</p>
				<p class="text-sm text-ocean-500 truncate">{auth.user?.email}</p>
			</div>

			<!-- Menu Items -->
			<div class="py-1">
				<a
					use:melt={$item}
					href="/settings"
					class="flex items-center gap-3 px-4 py-2 text-ocean-700 transition-all duration-200
						data-[highlighted]:bg-ocean-500/10"
				>
					<Settings class="w-4 h-4" />
					<span>Settings</span>
				</a>
			</div>

			<!-- Separator -->
			<div use:melt={$separator} class="h-px bg-ocean-100/50"></div>

			<!-- Logout -->
			<div class="py-1">
				<button
					use:melt={$item}
					onclick={handleLogout}
					class="flex items-center gap-3 w-full px-4 py-2 text-error transition-all duration-200 cursor-pointer
						data-[highlighted]:bg-red-50/80"
				>
					<LogOut class="w-4 h-4" />
					<span>Sign out</span>
				</button>
			</div>
		</div>
	{/if}
</div>
