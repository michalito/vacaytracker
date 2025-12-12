<script lang="ts">
	import type { Snippet } from 'svelte';
	import Button from './Button.svelte';
	import Card from './Card.svelte';
	import { AlertTriangle, RefreshCw } from 'lucide-svelte';

	interface Props {
		/** Content to render when no error */
		children: Snippet;
		/** Optional custom fallback UI */
		fallback?: Snippet<[{ error: Error; reset: () => void }]>;
		/** Called when an error is caught */
		onError?: (error: Error) => void;
	}

	let { children, fallback, onError }: Props = $props();

	let error = $state<Error | null>(null);

	function handleError(e: unknown, reset: () => void) {
		// Normalize unknown error to Error object
		const normalizedError = e instanceof Error ? e : new Error(String(e));
		error = normalizedError;
		onError?.(normalizedError);

		// Log to console in development
		if (import.meta.env.DEV) {
			console.error('[ErrorBoundary]', normalizedError);
		}
	}

	function reset() {
		error = null;
	}
</script>

<svelte:boundary onerror={handleError}>
	{#if error}
		{#if fallback}
			{@render fallback({ error, reset })}
		{:else}
			<Card class="border-coral-400/30">
				<div class="flex flex-col items-center py-8 px-4 text-center">
					<div class="p-3 bg-coral-400/10 rounded-full mb-4">
						<AlertTriangle class="w-8 h-8 text-coral-500" />
					</div>

					<h3 class="text-lg font-semibold text-ocean-800 mb-2">
						Something went wrong
					</h3>

					<p class="text-sm text-ocean-600 mb-6 max-w-sm">
						{error.message || 'An unexpected error occurred. Please try again.'}
					</p>

					<Button variant="outline" onclick={reset}>
						<RefreshCw class="w-4 h-4 mr-2" />
						Try Again
					</Button>
				</div>
			</Card>
		{/if}
	{:else}
		{@render children()}
	{/if}
</svelte:boundary>
