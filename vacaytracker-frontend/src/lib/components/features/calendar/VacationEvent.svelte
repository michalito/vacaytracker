<script lang="ts">
	import { clsx } from 'clsx';
	import type { TeamVacation } from '$lib/types';
	import { getUserColor } from '$lib/utils/colors';
	import { formatDateRangeShort } from '$lib/utils/date';
	import Tooltip from '$lib/components/ui/Tooltip.svelte';

	type EventPosition = 'start' | 'middle' | 'end' | 'single';

	interface Props {
		vacation: TeamVacation;
		showName?: boolean;
		position?: EventPosition;
		compact?: boolean;
		class?: string;
	}

	let {
		vacation,
		showName = true,
		position = 'single',
		compact = false,
		class: className = ''
	}: Props = $props();

	const color = $derived(getUserColor(vacation.userId));

	const borderRadius = $derived.by(() => {
		switch (position) {
			case 'start':
				return 'rounded-l-md rounded-r-none';
			case 'middle':
				return 'rounded-none';
			case 'end':
				return 'rounded-r-md rounded-l-none';
			default:
				return 'rounded-md';
		}
	});

	const daysLabel = $derived(vacation.totalDays === 1 ? 'day' : 'days');
</script>

<Tooltip placement="top">
	{#snippet content()}
		<div class="text-center">
			<div class="font-medium">{vacation.userName}</div>
			<div class="text-ocean-200 text-xs mt-0.5">
				{formatDateRangeShort(vacation.startDate, vacation.endDate)}
			</div>
			<div class="text-ocean-300 text-xs">
				{vacation.totalDays} {daysLabel}
			</div>
		</div>
	{/snippet}
	<div
		class={clsx(
			color.combined,
			borderRadius,
			compact ? 'text-xs px-1 py-0.5' : 'text-sm px-2 py-1',
			'truncate cursor-default',
			className
		)}
	>
		{#if showName}
			{vacation.userName}
		{/if}
	</div>
</Tooltip>
