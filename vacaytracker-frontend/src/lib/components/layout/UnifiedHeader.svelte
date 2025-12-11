<script lang="ts">
	import { page } from '$app/stores';
	import { auth } from '$lib/stores/auth.svelte';
	import { Calendar, LayoutDashboard, Settings } from 'lucide-svelte';
	import UserMenu from './UserMenu.svelte';
	import AdminDropdown from './AdminDropdown.svelte';
	import { clsx } from 'clsx';

	const navItems = [
		{ href: '/dashboard', icon: LayoutDashboard, label: 'Dashboard', exact: true },
		{ href: '/calendar', icon: Calendar, label: 'Calendar' },
		{ href: '/settings', icon: Settings, label: 'Settings' }
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
			<a href="/dashboard" class="flex items-center gap-2">
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

				<!-- Admin Dropdown (only for admins) -->
				{#if auth.isAdmin}
					<AdminDropdown />
				{/if}
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
			{#if auth.isAdmin}
				<a
					href="/admin/users"
					class={clsx(
						'flex flex-col items-center gap-1 px-3 py-2 rounded-lg text-xs font-medium transition-colors',
						$page.url.pathname.startsWith('/admin')
							? 'text-ocean-700'
							: 'text-ocean-500 hover:text-ocean-700'
					)}
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="w-5 h-5"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path
							d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"
						/>
					</svg>
					Admin
				</a>
			{/if}
		</div>
	</nav>
</header>
