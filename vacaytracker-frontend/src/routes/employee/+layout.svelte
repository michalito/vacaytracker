<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import EmployeeHeader from '$lib/components/layout/EmployeeHeader.svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
	}

	let { children }: Props = $props();

	$effect(() => {
		if (!auth.isLoading && !auth.isAuthenticated) {
			goto('/');
		}
	});
</script>

{#if auth.isAuthenticated}
	<div class="min-h-screen bg-sand-50">
		<EmployeeHeader />
		<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			{@render children()}
		</main>
	</div>
{/if}
