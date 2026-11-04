<script lang="ts">
	import '../app.css';
	import { auth } from '$lib/stores/auth.svelte';
	import Toaster from '$lib/components/ui/Toaster.svelte';
	import { Tooltip } from 'bits-ui';
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();

	$effect(() => {
		auth.initialize();
	});
</script>

<Tooltip.Provider delayDuration={300} skipDelayDuration={150}>
	<Toaster />

	{#if auth.isLoading}
		<div class="min-h-screen flex items-center justify-center">
			<div
				class="animate-spin rounded-full h-12 w-12 border-4 border-ocean-500 border-t-transparent"
			></div>
		</div>
	{:else}
		{@render children()}
	{/if}
</Tooltip.Provider>
