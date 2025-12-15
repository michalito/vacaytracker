<script lang="ts">
	import { createPopover, melt } from '@melt-ui/svelte';
	import VacationAvatarItem from './VacationAvatarItem.svelte';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import { formatDateRangeShort } from '$lib/utils/date';
	import type { TeamVacation } from '$lib/types';

	interface Props {
		vacations: TeamVacation[];
		maxVisible?: number;
		variant?: 'spread' | 'stacked';
		showEndDates?: boolean;
		class?: string;
	}

	let {
		vacations,
		maxVisible = 5,
		variant = 'spread',
		showEndDates = false,
		class: className = ''
	}: Props = $props();

	const visibleVacations = $derived(vacations.slice(0, maxVisible));
	const hiddenVacations = $derived(vacations.slice(maxVisible));
	const overflowCount = $derived(hiddenVacations.length);

	// Popover for overflow indicator
	const {
		elements: { trigger: overflowTrigger, content: overflowContent, close: overflowClose },
		states: { open: overflowOpen }
	} = createPopover({
		forceVisible: true,
		positioning: { placement: 'bottom-end' },
		closeOnOutsideClick: true
	});

	const containerClasses = $derived(
		variant === 'stacked' ? 'flex items-center pl-1' : 'flex items-start gap-4'
	);
</script>

<div class="{containerClasses} {className}">
	{#each visibleVacations as vacation, index (vacation.id)}
		{@const zIndex = visibleVacations.length - index + 1}
		{#if variant === 'stacked'}
			<div
				class="stacked-avatar-wrapper rounded-full"
				style="z-index: {zIndex}; margin-left: {index === 0 ? '0' : '-12px'};"
			>
				<VacationAvatarItem
					{vacation}
					showStatusIndicator={false}
					showEndDate={false}
					size="md"
				/>
			</div>
		{:else}
			<div>
				<VacationAvatarItem
					{vacation}
					showStatusIndicator={true}
					showEndDate={showEndDates}
					size="md"
				/>
			</div>
		{/if}
	{/each}

	<!-- Overflow indicator -->
	{#if overflowCount > 0}
		<div style={variant === 'stacked' ? 'z-index: 0; margin-left: -12px;' : ''}>
			<button
				use:melt={$overflowTrigger}
				class="w-10 h-10 rounded-full bg-gradient-to-br from-ocean-100 to-ocean-200
					text-ocean-700 text-sm font-bold
					flex items-center justify-center
					shadow-md hover:shadow-lg hover:scale-105
					transition-all duration-200
					ring-[3px] ring-white
					focus:outline-none focus-visible:ring-2 focus-visible:ring-ocean-500 focus-visible:ring-offset-2"
			>
				+{overflowCount}
			</button>
		</div>

		<!-- Overflow popover -->
		{#if $overflowOpen}
			<div
				use:melt={$overflowContent}
				class="z-50 w-56 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/40 overflow-hidden
					transition-all duration-200
					data-[state=open]:opacity-100 data-[state=open]:scale-100
					data-[state=closed]:opacity-0 data-[state=closed]:scale-95"
			>
				<div class="px-3 py-2 bg-ocean-50/50 border-b border-ocean-100/50">
					<span class="text-sm font-semibold text-ocean-700">
						+{overflowCount} more
					</span>
				</div>
				<div class="max-h-48 overflow-y-auto">
					{#each hiddenVacations as vacation (vacation.id)}
						<div class="px-3 py-2 flex items-center gap-2 hover:bg-ocean-50/50 transition-colors">
							<Avatar name={vacation.userName} size="sm" />
							<div class="min-w-0 flex-1">
								<p class="text-sm font-medium text-ocean-800 truncate">{vacation.userName}</p>
								<p class="text-xs text-ocean-500">
									{formatDateRangeShort(vacation.startDate, vacation.endDate)}
								</p>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>
