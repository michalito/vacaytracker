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
		clsx(
			'bg-white/90 backdrop-blur-sm rounded-xl shadow-md border border-white/40',
			className
		)
	);

	const contentClasses = $derived(paddingStyles[padding]);
</script>

<div class={cardClasses}>
	{#if header}
		<div class="px-4 py-3 border-b border-ocean-100/50">
			{@render header()}
		</div>
	{/if}

	<div class={contentClasses}>
		{@render children()}
	</div>

	{#if footer}
		<div class="px-4 py-3 border-t border-ocean-100/50 bg-ocean-50/30 rounded-b-xl">
			{@render footer()}
		</div>
	{/if}
</div>
