<script lang="ts">
	import { page } from '$app/stores';
	import { clsx } from 'clsx';
	import { LayoutDashboard, Users, Settings, Wallet, Palmtree, Calendar } from 'lucide-svelte';

	const adminNavItems = [
		{ href: '/admin', icon: LayoutDashboard, label: 'Dashboard' },
		{ href: '/admin/users', icon: Users, label: 'Users' },
		{ href: '/admin/balances', icon: Wallet, label: 'Balances' },
		{ href: '/admin/settings', icon: Settings, label: 'Settings' }
	];

	const myNavItems = [
		{ href: '/employee', icon: Palmtree, label: 'My Vacations' },
		{ href: '/employee/calendar', icon: Calendar, label: 'Team Calendar' }
	];
</script>

<aside class="w-64 bg-ocean-900 text-white min-h-screen flex flex-col">
	<!-- Logo -->
	<div class="p-4 border-b border-ocean-700">
		<a href="/admin" class="flex items-center gap-2">
			<img src="/logo.png" alt="VacayTracker" class="w-8 h-8" />
			<div>
				<span class="text-lg font-bold">VacayTracker</span>
				<span class="block text-xs text-ocean-400">Captain's Deck</span>
			</div>
		</a>
	</div>

	<!-- Navigation -->
	<nav class="flex-1 p-4">
		<!-- Admin Section -->
		<div class="mb-6">
			<h3 class="text-xs font-semibold text-ocean-500 uppercase tracking-wider mb-2 px-3">
				Admin
			</h3>
			<ul class="space-y-1">
				{#each adminNavItems as item}
					{@const isActive = $page.url.pathname === item.href}
					<li>
						<a
							href={item.href}
							class={clsx(
								'flex items-center gap-3 px-3 py-2 rounded-md transition-colors',
								isActive
									? 'bg-ocean-700 text-white'
									: 'text-ocean-300 hover:bg-ocean-800 hover:text-white'
							)}
						>
							<item.icon class="w-5 h-5" />
							<span>{item.label}</span>
						</a>
					</li>
				{/each}
			</ul>
		</div>

		<!-- My Vacations Section -->
		<div>
			<h3 class="text-xs font-semibold text-ocean-500 uppercase tracking-wider mb-2 px-3">
				My Vacations
			</h3>
			<ul class="space-y-1">
				{#each myNavItems as item}
					{@const isActive = $page.url.pathname === item.href || $page.url.pathname.startsWith(item.href + '/')}
					<li>
						<a
							href={item.href}
							class={clsx(
								'flex items-center gap-3 px-3 py-2 rounded-md transition-colors',
								isActive
									? 'bg-ocean-700 text-white'
									: 'text-ocean-300 hover:bg-ocean-800 hover:text-white'
							)}
						>
							<item.icon class="w-5 h-5" />
							<span>{item.label}</span>
						</a>
					</li>
				{/each}
			</ul>
		</div>
	</nav>

	<!-- Footer -->
	<div class="p-4 border-t border-ocean-700 text-ocean-400 text-sm">
		<p>VacayTracker v1.0</p>
	</div>
</aside>
