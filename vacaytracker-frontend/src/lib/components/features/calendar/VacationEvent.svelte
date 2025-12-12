<script lang="ts">
	import { clsx } from 'clsx';
	import type { TeamVacation } from '$lib/types';
	import { getUserColor } from '$lib/utils/colors';

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

	const title = $derived(
		`${vacation.userName}: ${vacation.startDate} - ${vacation.endDate} (${vacation.totalDays} day${vacation.totalDays !== 1 ? 's' : ''})`
	);
</script>

<div
	class={clsx(
		color.combined,
		borderRadius,
		compact ? 'text-xs px-1 py-0.5' : 'text-sm px-2 py-1',
		'truncate cursor-default',
		className
	)}
	{title}
>
	{#if showName}
		{vacation.userName}
	{/if}
</div>
