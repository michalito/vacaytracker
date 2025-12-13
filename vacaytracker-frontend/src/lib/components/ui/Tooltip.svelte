<script lang="ts">
	import { createTooltip, melt } from '@melt-ui/svelte';
	import { clsx } from 'clsx';
	import type { Snippet } from 'svelte';

	type Placement = 'top' | 'bottom' | 'left' | 'right';

	interface Props {
		content: string | Snippet;
		placement?: Placement;
		openDelay?: number;
		closeDelay?: number;
		disabled?: boolean;
		class?: string;
		children: Snippet;
	}

	let {
		content,
		placement = 'top',
		openDelay = 300,
		closeDelay = 150,
		disabled = false,
		class: className = '',
		children
	}: Props = $props();

	const {
		elements: { trigger, content: tooltipContent, arrow },
		states: { open }
	} = createTooltip({
		forceVisible: true,
		positioning: { placement },
		openDelay,
		closeDelay
	});
</script>

<span use:melt={$trigger} class={clsx('inline-flex', className)}>
	{@render children()}
</span>

{#if $open && !disabled}
	<div
		use:melt={$tooltipContent}
		class="z-50 px-3 py-2 rounded-lg bg-ocean-800 text-white text-sm shadow-lg
			transition-opacity duration-150
			data-[state=open]:opacity-100
			data-[state=closed]:opacity-0"
	>
		<div
			use:melt={$arrow}
			class="absolute h-2 w-2 bg-ocean-800 rotate-45
				data-[side=top]:bottom-[-4px]
				data-[side=bottom]:top-[-4px]
				data-[side=left]:right-[-4px]
				data-[side=right]:left-[-4px]"
		></div>
		{#if typeof content === 'string'}
			{content}
		{:else}
			{@render content()}
		{/if}
	</div>
{/if}
