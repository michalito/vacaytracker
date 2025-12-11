<script lang="ts">
	import { page } from '$app/stores';
	import { Shield, ChevronDown, Users, Wallet, Settings } from 'lucide-svelte';
	import { clsx } from 'clsx';

	let isOpen = $state(false);

	const adminNavItems = [
		{ href: '/admin/users', icon: Users, label: 'Users' },
		{ href: '/admin/balances', icon: Wallet, label: 'Balances' },
		{ href: '/admin/settings', icon: Settings, label: 'Settings' }
	];

	function isActive(href: string): boolean {
		return $page.url.pathname.startsWith(href);
	}

	function handleClickOutside(event: MouseEvent) {
		const target = event.target as Element;
		if (!target.closest('.admin-dropdown')) {
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

<div class="relative admin-dropdown">
	<button
		type="button"
		onclick={() => (isOpen = !isOpen)}
		class={clsx(
			'flex items-center gap-2 px-3 py-2 rounded-lg font-medium transition-colors',
			$page.url.pathname.startsWith('/admin')
				? 'bg-ocean-50 text-ocean-700'
				: 'text-ocean-600 hover:text-ocean-800 hover:bg-sand-50'
		)}
	>
		<Shield class="w-4 h-4" />
		<span>Admin</span>
		<ChevronDown class="w-4 h-4 transition-transform {isOpen ? 'rotate-180' : ''}" />
	</button>

	{#if isOpen}
		<div
			class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-sand-200 py-1 z-50"
		>
			{#each adminNavItems as item}
				{@const active = isActive(item.href)}
				<a
					href={item.href}
					onclick={() => (isOpen = false)}
					class={clsx(
						'flex items-center gap-3 px-4 py-2 transition-colors',
						active
							? 'bg-ocean-50 text-ocean-700'
							: 'text-ocean-600 hover:bg-sand-50 hover:text-ocean-800'
					)}
				>
					<item.icon class="w-4 h-4" />
					<span>{item.label}</span>
				</a>
			{/each}
		</div>
	{/if}
</div>
