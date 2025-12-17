<script lang="ts">
	import { createTooltip, createPopover, melt } from '@melt-ui/svelte';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { Calendar, Clock, X } from 'lucide-svelte';
	import {
		formatDateRangeShort,
		formatDateShort,
		isCurrentlyOngoing,
		getDaysUntilStart,
		getDaysUntilEnd,
		addDays
	} from '$lib/utils/date';
	import type { TeamVacation } from '$lib/types';

	interface Props {
		vacation: TeamVacation;
		showStatusIndicator?: boolean;
		showEndDate?: boolean;
		size?: 'sm' | 'md';
		class?: string;
	}

	let {
		vacation,
		showStatusIndicator = false,
		showEndDate = false,
		size = 'md',
		class: className = ''
	}: Props = $props();

	const {
		elements: { trigger: tooltipTrigger, content: tooltipContent, arrow: tooltipArrow },
		states: { open: tooltipOpen }
	} = createTooltip({
		openDelay: 300,
		closeDelay: 150,
		positioning: { placement: 'top' },
		group: 'vacation-avatars',
		forceVisible: true
	});

	const {
		elements: { trigger: popoverTrigger, content: popoverContent, close: popoverClose },
		states: { open: popoverOpen }
	} = createPopover({
		forceVisible: true,
		positioning: { placement: 'bottom' },
		closeOnOutsideClick: true
	});

	// Close tooltip when popover opens
	$effect(() => {
		if ($popoverOpen) {
			tooltipOpen.set(false);
		}
	});

	// Derived values for popover content
	const isCurrentlyOut = $derived(isCurrentlyOngoing(vacation.startDate, vacation.endDate));
	const daysUntilEnd = $derived(getDaysUntilEnd(vacation.endDate));
	const daysUntilStart = $derived(getDaysUntilStart(vacation.startDate));
	const returnDate = $derived(
		formatDateShort(addDays(new Date(vacation.endDate), 1).toISOString().split('T')[0])
	);
</script>

<div class="flex flex-col items-center {className}">
	<!-- Popover trigger wraps the avatar and tooltip trigger -->
	<button
		use:melt={$popoverTrigger}
		class="relative group focus:outline-none focus-visible:ring-2 focus-visible:ring-ocean-500 focus-visible:ring-offset-2 rounded-full"
	>
		<!-- Tooltip trigger wraps the avatar -->
		<span use:melt={$tooltipTrigger} class="block">
			<div class="relative">
				<Avatar
					name={vacation.userName}
					{size}
					class="transition-transform duration-200 group-hover:scale-105
						ring-[3px] ring-white shadow-md"
				/>

				<!-- Status indicator (green dot) for currently out -->
				{#if showStatusIndicator}
					<span
						class="absolute bottom-0 right-0 w-3 h-3 bg-green-500 border-2 border-white rounded-full"
						aria-label="Currently on vacation"
					></span>
				{/if}
			</div>
		</span>
	</button>

	<!-- End date label -->
	{#if showEndDate}
		<span class="text-xs text-ocean-500 mt-1.5 text-center whitespace-nowrap">
			ends {formatDateShort(vacation.endDate)}
		</span>
	{/if}
</div>

<!-- Tooltip Content (only show when popover is closed) -->
{#if $tooltipOpen && !$popoverOpen}
	<div
		use:melt={$tooltipContent}
		class="z-50 px-3 py-2 rounded-lg bg-ocean-800 text-white text-sm shadow-lg
			transition-opacity duration-150
			data-[state=open]:opacity-100 data-[state=closed]:opacity-0"
	>
		<div use:melt={$tooltipArrow} class="absolute h-2 w-2 bg-ocean-800 rotate-45"></div>
		<div class="text-center">
			<p class="font-medium">{vacation.userName}</p>
			<p class="text-xs text-ocean-200 mt-0.5">
				{formatDateRangeShort(vacation.startDate, vacation.endDate)}
			</p>
			<p class="text-xs text-ocean-300">
				{vacation.totalDays} day{vacation.totalDays !== 1 ? 's' : ''}
			</p>
		</div>
	</div>
{/if}

<!-- Popover Content -->
{#if $popoverOpen}
	<div
		use:melt={$popoverContent}
		class="z-50 w-64 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/40 overflow-hidden
			transition-all duration-200
			data-[state=open]:opacity-100 data-[state=open]:scale-100
			data-[state=closed]:opacity-0 data-[state=closed]:scale-95"
	>
		<!-- Header with close button -->
		<div
			class="flex items-center justify-between px-4 py-3 bg-ocean-50/50 border-b border-ocean-100/50"
		>
			<span class="font-semibold text-ocean-800">Vacation Details</span>
			<button
				use:melt={$popoverClose}
				class="p-1 rounded text-ocean-400 hover:text-ocean-600 hover:bg-ocean-100 transition-colors"
			>
				<X class="w-4 h-4" />
			</button>
		</div>

		<!-- Content -->
		<div class="p-4 space-y-3">
			<!-- User info -->
			<div class="flex items-center gap-3">
				<Avatar name={vacation.userName} size="md" />
				<div>
					<p class="font-medium text-ocean-800">{vacation.userName}</p>
					<Badge variant="success" size="sm">On Vacation</Badge>
				</div>
			</div>

			<!-- Date details -->
			<div class="space-y-2 text-sm">
				<div class="flex items-center gap-2 text-ocean-600">
					<Calendar class="w-4 h-4 text-ocean-400 shrink-0" />
					<span>{formatDateRangeShort(vacation.startDate, vacation.endDate)}</span>
				</div>
				<div class="flex items-center gap-2 text-ocean-600">
					<Clock class="w-4 h-4 text-ocean-400 shrink-0" />
					<span>{vacation.totalDays} day{vacation.totalDays !== 1 ? 's' : ''}</span>
				</div>
			</div>

			<!-- Status message -->
			{#if isCurrentlyOut}
				<p class="text-xs text-green-700 bg-green-50 px-3 py-2 rounded-lg">
					{#if daysUntilEnd === 0}
						Returns tomorrow ({returnDate})
					{:else if daysUntilEnd === 1}
						Returns in 2 days ({returnDate})
					{:else}
						Returns {returnDate} ({daysUntilEnd + 1} days)
					{/if}
				</p>
			{:else}
				<p class="text-xs text-ocean-600 bg-ocean-50 px-3 py-2 rounded-lg">
					{#if daysUntilStart === 0}
						Starts today
					{:else if daysUntilStart === 1}
						Starts tomorrow
					{:else}
						Starts in {daysUntilStart} days
					{/if}
				</p>
			{/if}
		</div>
	</div>
{/if}
