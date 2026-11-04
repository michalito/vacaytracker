<script lang="ts">
	import { Tooltip, Popover } from 'bits-ui';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { Calendar, Clock, X, Palmtree } from 'lucide-svelte';
	import {
		formatDateRangeShort,
		formatDateShort,
		formatDateISO,
		parseISODate,
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

	let popoverOpen = $state(false);

	// Derived values for popover content
	const isCurrentlyOut = $derived(isCurrentlyOngoing(vacation.startDate, vacation.endDate));
	const daysUntilEnd = $derived(getDaysUntilEnd(vacation.endDate));
	const daysUntilStart = $derived(getDaysUntilStart(vacation.startDate));
	// Local-date arithmetic throughout (no UTC round-trip through toISOString)
	const returnDate = $derived(
		formatDateShort(formatDateISO(addDays(parseISODate(vacation.endDate), 1)))
	);
</script>

<div class="flex flex-col items-center {className}">
	<Popover.Root bind:open={popoverOpen}>
		<!-- Hide the hover tooltip while the details popover is open -->
		<Tooltip.Root delayDuration={300} disabled={popoverOpen}>
			<!-- One button carries both triggers: tooltip on hover/focus, popover on click -->
			<Tooltip.Trigger>
				{#snippet child({ props: tooltipProps })}
					<Popover.Trigger
						{...tooltipProps}
						aria-label="{vacation.userName}'s vacation details"
						class="relative block group rounded-full focus:outline-none focus-visible:ring-2 focus-visible:ring-ocean-500 focus-visible:ring-offset-2"
					>
						<Avatar
							name={vacation.userName}
							{size}
							class="transition-transform duration-200 group-hover:scale-105 ring-[3px] ring-white shadow-md"
						/>

						<!-- Status indicator (palm tree) for currently out -->
						{#if showStatusIndicator}
							<span
								class="absolute -bottom-0.5 -right-0.5 w-4 h-4 bg-amber-400 border-2 border-white rounded-full flex items-center justify-center shadow-sm"
							>
								<Palmtree class="w-2.5 h-2.5 text-amber-800" aria-hidden="true" />
								<span class="sr-only">Currently on vacation</span>
							</span>
						{/if}
					</Popover.Trigger>
				{/snippet}
			</Tooltip.Trigger>

			<Tooltip.Portal>
				<Tooltip.Content
					side="top"
					sideOffset={6}
					class="presence-pop z-50 px-3 py-2 rounded-lg bg-ocean-800 text-white text-sm shadow-lg"
				>
					<Tooltip.Arrow class="text-ocean-800" />
					<div class="text-center">
						<p class="font-medium">{vacation.userName}</p>
						<p class="text-xs text-ocean-200 mt-0.5">
							{formatDateRangeShort(vacation.startDate, vacation.endDate)}
						</p>
						<p class="text-xs text-ocean-300">
							{vacation.totalDays} day{vacation.totalDays !== 1 ? 's' : ''}
						</p>
					</div>
				</Tooltip.Content>
			</Tooltip.Portal>
		</Tooltip.Root>

		<Popover.Portal>
			<Popover.Content
				side="bottom"
				sideOffset={4}
				class="presence-pop z-50 w-64 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/40 overflow-hidden"
			>
				<!-- Header with close button -->
				<div
					class="flex items-center justify-between px-4 py-3 bg-ocean-50/50 border-b border-ocean-100/50"
				>
					<span class="font-semibold text-ocean-800">Vacation Details</span>
					<Popover.Close
						aria-label="Close vacation details"
						class="p-1 rounded text-ocean-400 hover:text-ocean-600 hover:bg-ocean-100 transition-colors cursor-pointer"
					>
						<X class="w-4 h-4" aria-hidden="true" />
					</Popover.Close>
				</div>

				<!-- Content -->
				<div class="p-4 space-y-3">
					<!-- User info -->
					<div class="flex items-center gap-3">
						<Avatar name={vacation.userName} size="md" />
						<div>
							<p class="font-medium text-ocean-800">{vacation.userName}</p>
							{#if isCurrentlyOut}
								<Badge variant="success" size="sm">On Vacation</Badge>
							{:else}
								<Badge variant="ocean" size="sm">Upcoming</Badge>
							{/if}
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
			</Popover.Content>
		</Popover.Portal>
	</Popover.Root>

	<!-- End date label -->
	{#if showEndDate}
		<span class="text-xs text-ocean-500 mt-1.5 text-center whitespace-nowrap">
			ends {formatDateShort(vacation.endDate)}
		</span>
	{/if}
</div>
