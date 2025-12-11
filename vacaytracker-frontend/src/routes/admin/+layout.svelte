<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import UnifiedHeader from '$lib/components/layout/UnifiedHeader.svelte';

	let { children } = $props();

	$effect(() => {
		if (!auth.isLoading && (!auth.isAuthenticated || !auth.isAdmin)) {
			goto('/');
		}
	});
</script>

{#if auth.isAuthenticated && auth.isAdmin}
	<div class="min-h-screen bg-sand-50">
		<UnifiedHeader />
		<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			{@render children()}
		</main>
	</div>
{/if}
