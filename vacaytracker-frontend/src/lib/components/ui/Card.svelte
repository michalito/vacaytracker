<script lang="ts">
	import { clsx } from 'clsx';
	import type { Snippet } from 'svelte';

	interface Props {
		padding?: 'none' | 'sm' | 'md' | 'lg';
		class?: string;
		header?: Snippet;
		footer?: Snippet;
		children: Snippet;
	}

	let { padding = 'md', class: className = '', header, footer, children }: Props = $props();

	const paddingStyles = {
		none: '',
		sm: 'p-3',
		md: 'p-4',
		lg: 'p-6'
	};

	const cardClasses = $derived(
		clsx('bg-white rounded-lg shadow-md border border-sand-200', className)
	);

	const contentClasses = $derived(paddingStyles[padding]);
</script>

<div class={cardClasses}>
	{#if header}
		<div class="px-4 py-3 border-b border-sand-200">
			{@render header()}
		</div>
	{/if}

	<div class={contentClasses}>
		{@render children()}
	</div>

	{#if footer}
		<div class="px-4 py-3 border-t border-sand-200 bg-sand-50 rounded-b-lg">
			{@render footer()}
		</div>
	{/if}
</div>
