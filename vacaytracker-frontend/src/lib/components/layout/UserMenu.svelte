<script lang="ts">
	import { DropdownMenu } from 'bits-ui';
	import { auth } from '#lib/stores/auth.svelte.js';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import Avatar from '#lib/components/ui/Avatar.svelte';
	import { Settings, LogOut, ChevronDown } from 'lucide-svelte';

	let menuOpen = $state(false);

	function handleLogout() {
		auth.logout();
		goto(resolve('/'));
	}
</script>

<div class="relative">
	<DropdownMenu.Root bind:open={menuOpen}>
		<DropdownMenu.Trigger
			aria-label="Account menu"
			class="flex items-center gap-2 p-1.5 rounded-xl hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
		>
			<Avatar name={auth.user?.name} size="sm" />
			<span class="hidden sm:block text-sm font-medium text-ocean-700">
				{auth.user?.name}
			</span>
			<ChevronDown
				class="w-4 h-4 text-ocean-500 transition-transform {menuOpen ? 'rotate-180' : ''}"
				aria-hidden="true"
			/>
		</DropdownMenu.Trigger>

		<DropdownMenu.Portal>
			<DropdownMenu.Content
				side="bottom"
				align="end"
				sideOffset={8}
				class="presence-pop w-56 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/40 py-1 z-50"
			>
				<!-- User Info -->
				<div class="px-4 py-3 border-b border-ocean-100/50">
					<p class="font-medium text-ocean-800">{auth.user?.name}</p>
					<p class="text-sm text-ocean-500 truncate">{auth.user?.email}</p>
				</div>

				<!-- Menu Items -->
				<div class="py-1">
					<DropdownMenu.Item>
						{#snippet child({ props })}
							<a
								{...props}
								href={resolve('settings')}
								class="flex items-center gap-3 px-4 py-2 text-ocean-700 transition-all duration-200 cursor-pointer
									data-[highlighted]:bg-ocean-500/10"
							>
								<Settings class="w-4 h-4" />
								<span>Settings</span>
							</a>
						{/snippet}
					</DropdownMenu.Item>
				</div>

				<!-- Separator -->
				<DropdownMenu.Separator class="h-px bg-ocean-100/50" />

				<!-- Logout -->
				<div class="py-1">
					<DropdownMenu.Item
						onSelect={handleLogout}
						class="flex items-center gap-3 w-full px-4 py-2 text-error transition-all duration-200 cursor-pointer
							data-[highlighted]:bg-red-50/80"
					>
						<LogOut class="w-4 h-4" />
						<span>Sign out</span>
					</DropdownMenu.Item>
				</div>
			</DropdownMenu.Content>
		</DropdownMenu.Portal>
	</DropdownMenu.Root>
</div>
