<script lang="ts">
	import { page } from '$app/state';
	import { resolve, asset } from '$app/paths';
	import { auth } from '$lib/stores/auth.svelte';
	import { Calendar, LayoutDashboard, Shield } from 'lucide-svelte';
	import UserMenu from './UserMenu.svelte';
	import AdminDropdown from './AdminDropdown.svelte';
	import { clsx } from 'clsx';

	const navItems = [
		{ path: '/dashboard', icon: LayoutDashboard, label: 'Dashboard', exact: true },
		{ path: '/calendar', icon: Calendar, label: 'Calendar', exact: false }
	] as const;

	// Compares against resolved hrefs so active states still work under a configured base path
	function isActive(href: string, exact = false): boolean {
		const current = page.url.pathname;
		return exact ? current === href : current === href || current.startsWith(`${href}/`);
	}
</script>

<header class="sticky top-0 z-50 bg-white/85 backdrop-blur-md border-b border-white/40 shadow-sm">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex items-center justify-between h-16">
			<!-- Logo -->
			<a href={resolve('/dashboard')} class="flex items-center gap-2 group">
				<img
					src={asset('/logo.png')}
					alt=""
					class="w-8 h-8 transition-transform group-hover:scale-105"
				/>
				<span class="text-xl font-bold text-ocean-700">VacayTracker</span>
			</a>

			<!-- Navigation -->
			<nav aria-label="Main" class="hidden md:flex items-center gap-1">
				{#each navItems as item (item.path)}
					{@const active = isActive(resolve(item.path), item.exact)}
					<a
						href={resolve(item.path)}
						aria-current={active ? 'page' : undefined}
						class={clsx(
							'flex items-center gap-2 px-4 py-2 rounded-xl font-medium transition-all duration-200',
							active
								? 'bg-ocean-500/15 text-ocean-700 shadow-sm'
								: 'text-ocean-600 hover:text-ocean-800 hover:bg-ocean-500/10'
						)}
					>
						<item.icon class="w-4 h-4" aria-hidden="true" />
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
	<nav aria-label="Main" class="md:hidden border-t border-white/40 bg-white/50 backdrop-blur-sm">
		<div class="flex justify-around py-2">
			{#each navItems as item (item.path)}
				{@const active = isActive(resolve(item.path), item.exact)}
				<a
					href={resolve(item.path)}
					aria-current={active ? 'page' : undefined}
					class={clsx(
						'flex flex-col items-center gap-1 px-3 py-2 rounded-xl text-xs font-medium transition-all duration-200',
						active
							? 'text-ocean-700 bg-ocean-500/10'
							: 'text-ocean-500 hover:text-ocean-700 hover:bg-ocean-500/5'
					)}
				>
					<item.icon class="w-5 h-5" aria-hidden="true" />
					{item.label}
				</a>
			{/each}
			{#if auth.isAdmin}
				{@const adminActive = isActive(resolve('/admin'))}
				<a
					href={resolve('/admin/users')}
					aria-current={adminActive ? 'page' : undefined}
					class={clsx(
						'flex flex-col items-center gap-1 px-3 py-2 rounded-xl text-xs font-medium transition-all duration-200',
						adminActive
							? 'text-ocean-700 bg-ocean-500/10'
							: 'text-ocean-500 hover:text-ocean-700 hover:bg-ocean-500/5'
					)}
				>
					<Shield class="w-5 h-5" aria-hidden="true" />
					Admin
				</a>
			{/if}
		</div>
	</nav>
</header>
