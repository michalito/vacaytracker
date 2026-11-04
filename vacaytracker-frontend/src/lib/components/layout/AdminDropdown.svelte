<script lang="ts">
	import { DropdownMenu } from 'bits-ui';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { Shield, ChevronDown, Users, Wallet, Settings } from 'lucide-svelte';
	import { clsx } from 'clsx';

	let menuOpen = $state(false);

	// SvelteKit 3 paths have no leading slash (only route IDs do); resolve() adds it plus any base
	const adminNavItems = [
		{ path: 'admin/users', icon: Users, label: 'Users' },
		{ path: 'admin/balances', icon: Wallet, label: 'Balances' },
		{ path: 'admin/settings', icon: Settings, label: 'Settings' }
	] as const;

	// Compares against resolved hrefs so active states still work under a configured base path
	function isActive(href: string): boolean {
		const current = page.url.pathname;
		return current === href || current.startsWith(`${href}/`);
	}

	// /admin is a layout without its own page, so match on the route ID (base-path independent)
	const inAdmin = $derived(page.route.id?.startsWith('/admin/') ?? false);
</script>

<div class="relative">
	<DropdownMenu.Root bind:open={menuOpen}>
		<DropdownMenu.Trigger
			class={clsx(
				'flex items-center gap-2 px-4 py-2 rounded-xl font-medium transition-all duration-200 cursor-pointer',
				inAdmin
					? 'bg-ocean-500/15 text-ocean-700 shadow-sm'
					: 'text-ocean-600 hover:text-ocean-800 hover:bg-ocean-500/10'
			)}
		>
			<Shield class="w-4 h-4" aria-hidden="true" />
			<span>Admin</span>
			<ChevronDown
				class="w-4 h-4 transition-transform {menuOpen ? 'rotate-180' : ''}"
				aria-hidden="true"
			/>
		</DropdownMenu.Trigger>

		<DropdownMenu.Portal>
			<DropdownMenu.Content
				side="bottom"
				align="end"
				sideOffset={8}
				class="presence-pop w-48 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/40 py-1 z-50"
			>
				{#each adminNavItems as navItem (navItem.path)}
					{@const active = isActive(resolve(navItem.path))}
					<DropdownMenu.Item>
						{#snippet child({ props })}
							<a
								{...props}
								href={resolve(navItem.path)}
								aria-current={active ? 'page' : undefined}
								class={clsx(
									'flex items-center gap-3 px-4 py-2 transition-all duration-200 cursor-pointer',
									active
										? 'bg-ocean-500/15 text-ocean-700'
										: 'text-ocean-600 data-[highlighted]:bg-ocean-500/10 data-[highlighted]:text-ocean-800'
								)}
							>
								<navItem.icon class="w-4 h-4" aria-hidden="true" />
								<span>{navItem.label}</span>
							</a>
						{/snippet}
					</DropdownMenu.Item>
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Portal>
	</DropdownMenu.Root>
</div>
