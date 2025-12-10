<script lang="ts">
	import { page } from '$app/stores';
	import { Calendar, LayoutDashboard, Settings } from 'lucide-svelte';
	import UserMenu from './UserMenu.svelte';
	import { clsx } from 'clsx';

	const navItems = [
		{ href: '/employee', icon: LayoutDashboard, label: 'Dashboard', exact: true },
		{ href: '/employee/calendar', icon: Calendar, label: 'Calendar' },
		{ href: '/employee/settings', icon: Settings, label: 'Settings' }
	];

	function isActive(href: string, exact = false): boolean {
		if (exact) {
			return $page.url.pathname === href;
		}
		return $page.url.pathname.startsWith(href);
	}
</script>

<header class="bg-white shadow-sm border-b border-sand-200">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex items-center justify-between h-16">
			<!-- Logo -->
			<a href="/employee" class="flex items-center gap-2">
				<img src="/logo.png" alt="VacayTracker" class="w-8 h-8" />
				<span class="text-xl font-bold text-ocean-700">VacayTracker</span>
			</a>

			<!-- Navigation -->
			<nav class="hidden md:flex items-center gap-1">
				{#each navItems as item}
					{@const active = isActive(item.href, item.exact)}
					<a
						href={item.href}
						class={clsx(
							'flex items-center gap-2 px-3 py-2 rounded-lg font-medium transition-colors',
							active
								? 'bg-ocean-50 text-ocean-700'
								: 'text-ocean-600 hover:text-ocean-800 hover:bg-sand-50'
						)}
					>
						<item.icon class="w-4 h-4" />
						{item.label}
					</a>
				{/each}
			</nav>

			<!-- User Menu -->
			<UserMenu />
		</div>
	</div>

	<!-- Mobile Navigation -->
	<nav class="md:hidden border-t border-sand-200">
		<div class="flex justify-around py-2">
			{#each navItems as item}
				{@const active = isActive(item.href, item.exact)}
				<a
					href={item.href}
					class={clsx(
						'flex flex-col items-center gap-1 px-3 py-2 rounded-lg text-xs font-medium transition-colors',
						active ? 'text-ocean-700' : 'text-ocean-500 hover:text-ocean-700'
					)}
				>
					<item.icon class="w-5 h-5" />
					{item.label}
				</a>
			{/each}
		</div>
	</nav>
</header>
