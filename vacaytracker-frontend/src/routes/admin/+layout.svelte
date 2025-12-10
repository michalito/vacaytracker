<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import AdminSidebar from '$lib/components/layout/AdminSidebar.svelte';
	import AdminHeader from '$lib/components/layout/AdminHeader.svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();

	$effect(() => {
		if (!auth.isLoading && (!auth.isAuthenticated || !auth.isAdmin)) {
			goto('/');
		}
	});
</script>

{#if auth.isAuthenticated && auth.isAdmin}
	<div class="min-h-screen bg-sand-50 flex">
		<AdminSidebar />
		<div class="flex-1 flex flex-col">
			<AdminHeader />
			<main class="flex-1 p-6">
				{@render children()}
			</main>
		</div>
	</div>
{/if}
