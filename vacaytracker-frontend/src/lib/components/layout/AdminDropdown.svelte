<script lang="ts">
	import { createDropdownMenu, melt } from '@melt-ui/svelte';
	import { page } from '$app/stores';
	import { Shield, ChevronDown, Users, Wallet, Settings } from 'lucide-svelte';
	import { clsx } from 'clsx';

	// Create Melt-UI dropdown menu
	const {
		elements: { trigger, menu, item },
		states: { open }
	} = createDropdownMenu({
		forceVisible: true,
		positioning: { placement: 'bottom-end' }
	});

	const adminNavItems = [
		{ href: '/admin/users', icon: Users, label: 'Users' },
		{ href: '/admin/balances', icon: Wallet, label: 'Balances' },
		{ href: '/admin/settings', icon: Settings, label: 'Settings' }
	];

	function isActive(href: string): boolean {
		return $page.url.pathname.startsWith(href);
	}
</script>

<div class="relative">
	<button
		use:melt={$trigger}
		class={clsx(
			'flex items-center gap-2 px-4 py-2 rounded-xl font-medium transition-all duration-200 cursor-pointer',
			$page.url.pathname.startsWith('/admin')
				? 'bg-ocean-500/15 text-ocean-700 shadow-sm'
				: 'text-ocean-600 hover:text-ocean-800 hover:bg-ocean-500/10'
		)}
	>
		<Shield class="w-4 h-4" />
		<span>Admin</span>
		<ChevronDown class="w-4 h-4 transition-transform {$open ? 'rotate-180' : ''}" />
	</button>

	{#if $open}
		<div
			use:melt={$menu}
			class="absolute right-0 mt-2 w-48 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/30 py-1 z-50
				transition-all duration-200 data-[state=open]:animate-scale-in"
		>
			{#each adminNavItems as navItem}
				{@const active = isActive(navItem.href)}
				<a
					use:melt={$item}
					href={navItem.href}
					class={clsx(
						'flex items-center gap-3 px-4 py-2 transition-all duration-200',
						active
							? 'bg-ocean-500/15 text-ocean-700'
							: 'text-ocean-600 data-[highlighted]:bg-ocean-500/10 data-[highlighted]:text-ocean-800'
					)}
				>
					<navItem.icon class="w-4 h-4" />
					<span>{navItem.label}</span>
				</a>
			{/each}
		</div>
	{/if}
</div>
