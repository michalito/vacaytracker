<script lang="ts">
	import { createTooltip, melt } from '@melt-ui/svelte';
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

	// Melt-UI tooltip for accessible, styled hover info
	const {
		elements: { trigger, content, arrow },
		states: { open }
	} = createTooltip({
		positioning: { placement: 'top' },
		openDelay: 300,
		closeDelay: 100,
		group: 'calendar-events'
	});
</script>

<div
	use:melt={$trigger}
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

{#if $open}
	<div
		use:melt={$content}
		class="z-50 rounded-lg bg-ocean-800 text-white px-3 py-2 text-sm shadow-lg max-w-xs"
	>
		<div use:melt={$arrow} class="z-50"></div>
		<div class="font-medium">{vacation.userName}</div>
		<div class="text-ocean-200 text-xs mt-1">
			{vacation.startDate} - {vacation.endDate}
		</div>
		<div class="text-ocean-300 text-xs">
			{vacation.totalDays} day{vacation.totalDays !== 1 ? 's' : ''}
		</div>
	</div>
{/if}
