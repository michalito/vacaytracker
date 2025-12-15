<script lang="ts">
	import AvatarRow from './AvatarRow.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import { Users, Calendar } from 'lucide-svelte';
	import type { TeamVacation } from '$lib/types';

	interface Props {
		title: string;
		vacations: TeamVacation[];
		variant: 'current' | 'upcoming';
		maxVisible?: number;
		emptyMessage?: string;
	}

	let {
		title,
		vacations,
		variant,
		maxVisible = 5,
		emptyMessage = 'No vacations'
	}: Props = $props();

	const rowVariant = $derived(variant === 'current' ? 'spread' : 'stacked');
	const showEndDates = $derived(variant === 'current');
	const emptyIcon = $derived(variant === 'current' ? Users : Calendar);
</script>

<div class="space-y-3">
	<!-- Section title -->
	<h3 class="text-sm font-semibold text-ocean-600 uppercase tracking-wide">
		{title}
		{#if vacations.length > 0}
			<span class="text-ocean-400 font-normal">({vacations.length})</span>
		{/if}
	</h3>

	<!-- Content -->
	{#if vacations.length === 0}
		<EmptyState icon={emptyIcon} message={emptyMessage} iconSize="sm" />
	{:else}
		<AvatarRow
			{vacations}
			{maxVisible}
			variant={rowVariant}
			{showEndDates}
		/>
	{/if}
</div>
