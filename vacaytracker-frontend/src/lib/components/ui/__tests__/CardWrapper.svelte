<script lang="ts">
	import Card from '../Card.svelte';

	interface Props {
		bodyText?: string;
		headerText?: string;
		footerText?: string;
		padding?: 'none' | 'sm' | 'md' | 'lg';
		class?: string;
		showHeader?: boolean;
		showFooter?: boolean;
	}

	let {
		bodyText = 'Card content',
		headerText = 'Header',
		footerText = 'Footer',
		showHeader = false,
		showFooter = false,
		...rest
	}: Props = $props();
</script>

{#if showHeader && showFooter}
	<Card {...rest}>
		{#snippet header()}<span>{headerText}</span>{/snippet}
		{#snippet footer()}<span>{footerText}</span>{/snippet}
		{bodyText}
	</Card>
{:else if showHeader}
	<Card {...rest}>
		{#snippet header()}<span>{headerText}</span>{/snippet}
		{bodyText}
	</Card>
{:else if showFooter}
	<Card {...rest}>
		{#snippet footer()}<span>{footerText}</span>{/snippet}
		{bodyText}
	</Card>
{:else}
	<Card {...rest}>
		{bodyText}
	</Card>
{/if}
