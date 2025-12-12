<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import UnifiedHeader from '$lib/components/layout/UnifiedHeader.svelte';
	import Footer from '$lib/components/layout/Footer.svelte';
	import DecorativeBackground from '$lib/components/layout/DecorativeBackground.svelte';

	let { children } = $props();

	$effect(() => {
		if (!auth.isLoading && !auth.isAuthenticated) {
			goto('/');
		}
	});
</script>

{#if auth.isAuthenticated}
	<div class="min-h-screen flex flex-col bg-sand-50 relative">
		<!-- Decorative Background Elements -->
		<DecorativeBackground variant="standard" />

		<!-- Header -->
		<UnifiedHeader />

		<!-- Main Content -->
		<main class="relative z-10 flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-8 pb-56 w-full">
			{@render children()}
		</main>

		<!-- Footer with Waves -->
		<Footer />
	</div>
{/if}
