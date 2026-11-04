<script lang="ts">
	import { Tooltip } from 'bits-ui';
	import type { Snippet } from 'svelte';

	type Placement = 'top' | 'bottom' | 'left' | 'right';

	interface Props {
		content: string | Snippet;
		placement?: Placement;
		openDelay?: number;
		disabled?: boolean;
		/**
		 * Renders the trigger. Spread `props` onto exactly one element, e.g.
		 * `{#snippet child({ props })}<button type="button" {...props}>…</button>{/snippet}`,
		 * so the tooltip attaches to that element instead of wrapping it in another button.
		 */
		child: Snippet<[{ props: Record<string, unknown> }]>;
	}

	let {
		content,
		placement = 'top',
		openDelay = 300,
		disabled = false,
		child: childSnippet
	}: Props = $props();
</script>

<Tooltip.Root delayDuration={openDelay} {disabled}>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<!-- Bits adds type="button" for its default <button>; the consumer chooses the element -->
			{@const { type: _type, ...triggerProps } = props}
			{@render childSnippet({ props: triggerProps })}
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Portal>
		<Tooltip.Content
			side={placement}
			sideOffset={8}
			class="presence-pop z-50 rounded-lg bg-ocean-800 px-3 py-2 text-sm text-white shadow-lg"
		>
			<Tooltip.Arrow class="text-ocean-800" />
			{#if typeof content === 'string'}
				{content}
			{:else}
				{@render content()}
			{/if}
		</Tooltip.Content>
	</Tooltip.Portal>
</Tooltip.Root>
