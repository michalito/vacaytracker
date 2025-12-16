<script lang="ts">
	import { createTabs, createSelect, melt } from '@melt-ui/svelte';
	import type { VacationRequest } from '$lib/types';
	import EnhancedRequestCard from './EnhancedRequestCard.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import {
		Clock,
		CalendarCheck,
		History,
		Umbrella,
		Plane,
		Archive,
		ChevronDown,
		Check
	} from 'lucide-svelte';

	interface Props {
		pendingRequests: VacationRequest[];
		upcomingRequests: VacationRequest[];
		pastRequests: VacationRequest[];
		onRequestVacation?: () => void;
	}

	let { pendingRequests, upcomingRequests, pastRequests, onRequestVacation }: Props = $props();

	// Extract unique years from past requests, sorted descending (most recent first)
	const availableYears = $derived(() => {
		const years = new Set<number>();
		pastRequests.forEach((r) => {
			const year = new Date(r.endDate).getFullYear();
			years.add(year);
		});
		return Array.from(years).sort((a, b) => b - a);
	});

	// Default to most recent year with data
	const defaultYear = $derived(() => availableYears()[0] ?? new Date().getFullYear());

	// Year filter state for Past tab
	let selectedYear = $state<number | null>(null);

	// Initialize selectedYear to defaultYear on first load
	$effect(() => {
		if (selectedYear === null && availableYears().length > 0) {
			selectedYear = defaultYear();
		}
	});

	// The actual year to filter by (use defaultYear as fallback)
	const activeYear = $derived(selectedYear ?? defaultYear());

	// Year filter options for Melt UI Select
	const yearOptions = $derived(() => {
		return availableYears().map((year) => ({
			value: year.toString(),
			label: year.toString()
		}));
	});

	// Melt UI Select for year filter
	const {
		elements: {
			trigger: yearTrigger,
			menu: yearMenu,
			option: yearOption
		},
		states: { open: yearSelectOpen, selected: yearSelected },
		helpers: { isSelected: isYearSelected }
	} = createSelect<string>({
		forceVisible: true,
		portal: null,
		onSelectedChange: ({ next }) => {
			if (next) {
				selectedYear = parseInt(next.value);
			}
			return next;
		}
	});

	// Sync the select's displayed value with activeYear
	$effect(() => {
		const year = activeYear;
		if (year && (!$yearSelected || $yearSelected.value !== year.toString())) {
			yearSelected.set({ value: year.toString(), label: year.toString() });
		}
	});

	// Filter past requests by active year
	const filteredPastRequests = $derived(() => {
		return pastRequests.filter((r) => {
			const requestYear = new Date(r.endDate).getFullYear();
			return requestYear === activeYear;
		});
	});

	// Tab configuration
	const tabs = $derived([
		{
			id: 'pending',
			label: 'Pending',
			icon: Clock,
			count: pendingRequests.length,
			requests: pendingRequests,
			variant: 'pending' as const
		},
		{
			id: 'upcoming',
			label: 'Upcoming',
			icon: CalendarCheck,
			count: upcomingRequests.length,
			requests: upcomingRequests,
			variant: 'upcoming' as const
		},
		{
			id: 'past',
			label: 'Past',
			icon: History,
			count: pastRequests.length,
			requests: filteredPastRequests(),
			variant: 'past' as const
		}
	]);

	// Determine default tab: show pending if any, otherwise upcoming, otherwise past
	const defaultTab = $derived(
		pendingRequests.length > 0
			? 'pending'
			: upcomingRequests.length > 0
				? 'upcoming'
				: 'past'
	);

	const {
		elements: { root, list, trigger, content },
		states: { value }
	} = createTabs({
		defaultValue: defaultTab,
		loop: true,
		activateOnFocus: true
	});
</script>

<Card padding="none">
	{#snippet header()}
		<div class="flex items-center justify-between">
			<h2 class="text-lg font-semibold text-ocean-700">Your Requests</h2>
			<Umbrella class="w-5 h-5 text-ocean-400" />
		</div>
	{/snippet}

	<div use:melt={$root} class="w-full">
		<!-- Tab List -->
		<div class="px-4 pt-4">
			<div
				use:melt={$list}
				class="inline-flex p-1 bg-sand-100 rounded-lg"
				aria-label="Request categories"
			>
				{#each tabs as tab (tab.id)}
					<button
						use:melt={$trigger(tab.id)}
						class="flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-md transition-colors cursor-pointer
							data-[state=active]:bg-white data-[state=active]:text-ocean-700 data-[state=active]:shadow-sm
							data-[state=inactive]:text-ocean-500 data-[state=inactive]:hover:text-ocean-700"
					>
						<tab.icon class="w-4 h-4" />
						<span>{tab.label}</span>
						{#if tab.count > 0}
							<span
								class="px-1.5 py-0.5 text-xs rounded-md transition-colors
									data-[state=active]:bg-ocean-100 data-[state=active]:text-ocean-700
									data-[state=inactive]:bg-sand-200 data-[state=inactive]:text-ocean-500"
								data-state={$value === tab.id ? 'active' : 'inactive'}
							>
								{tab.count}
							</span>
						{/if}
					</button>
				{/each}
			</div>
		</div>

		<!-- Tab Content -->
		{#each tabs as tab (tab.id)}
			<div use:melt={$content(tab.id)} class="p-4 focus:outline-none">
				{#if tab.requests.length === 0}
					<!-- Empty States -->
					<div class="content-fade-in">
						{#if tab.id === 'pending'}
							<EmptyState icon={Plane} message="No requests in the queue" iconSize="lg">
								<p class="text-sm text-ocean-400 -mt-2">Time to dream up your next getaway?</p>
							</EmptyState>
						{:else if tab.id === 'upcoming'}
							<EmptyState icon={CalendarCheck} message="No upcoming vacations" iconSize="lg">
								{#if onRequestVacation}
									<Button variant="outline" onclick={onRequestVacation}>
										Plan your next getaway
									</Button>
								{/if}
							</EmptyState>
						{:else}
							<EmptyState icon={Archive} message="No vacation history yet" iconSize="lg">
								<p class="text-sm text-ocean-400 -mt-2">Your completed requests will appear here</p>
							</EmptyState>
						{/if}
					</div>
				{:else}
					<!-- Past Tab: Year Filter + Cards -->
					{#if tab.id === 'past' && availableYears().length > 0}
						<div class="flex items-center justify-between mb-4 content-fade-in">
							<p class="text-sm text-ocean-500">
								{filteredPastRequests().length} request{filteredPastRequests().length !== 1 ? 's' : ''}
							</p>
							<div class="relative inline-block">
								<button
									use:melt={$yearTrigger}
									class="flex items-center p-1 text-sm font-medium rounded-lg
										bg-sand-100
										focus:outline-none
										cursor-pointer"
								>
									<span class="flex items-center gap-1.5 px-3 py-1.5 rounded-md bg-white text-ocean-700 shadow-sm">
										{$yearSelected?.label ?? activeYear}
										<ChevronDown
											class="w-4 h-4 transition-transform duration-200 {$yearSelectOpen
												? 'rotate-180'
												: ''}"
										/>
									</span>
								</button>

								{#if $yearSelectOpen}
									<div
										use:melt={$yearMenu}
										class="absolute right-0 top-full z-50 mt-1 w-max rounded-lg bg-sand-100 p-1 shadow-lg"
									>
										{#each yearOptions() as opt}
											<div
												use:melt={$yearOption({ value: opt.value, label: opt.label })}
												class="flex items-center justify-between gap-3 rounded-md px-3 py-1.5 text-sm cursor-pointer outline-none whitespace-nowrap
													text-ocean-500
													data-[highlighted]:bg-white data-[highlighted]:text-ocean-700 data-[highlighted]:shadow-sm
													data-[selected]:font-medium"
											>
												<span>{opt.label}</span>
												{#if $isYearSelected(opt.value)}
													<Check class="w-4 h-4 text-ocean-500" />
												{/if}
											</div>
										{/each}
									</div>
								{/if}
							</div>
						</div>
					{/if}

					<!-- Request Cards -->
					<div class="space-y-3 content-fade-in">
						{#each tab.requests as request (request.id)}
							<EnhancedRequestCard {request} variant={tab.variant} />
						{/each}
					</div>
				{/if}
			</div>
		{/each}
	</div>
</Card>
