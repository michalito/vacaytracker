<script lang="ts">
	import { createTabs, createSelect, melt } from '@melt-ui/svelte';
	import { fade, fly } from 'svelte/transition';
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

	// ===================
	// Props
	// ===================
	interface Props {
		pendingRequests: VacationRequest[];
		upcomingRequests: VacationRequest[];
		pastRequests: VacationRequest[];
		onRequestVacation?: () => void;
	}

	const { pendingRequests, upcomingRequests, pastRequests, onRequestVacation }: Props = $props();

	// ===================
	// Constants
	// ===================
	const MAX_REQUESTS_PER_TAB = 50;

	// ===================
	// Tab Types
	// ===================
	type TabId = 'pending' | 'upcoming' | 'past';
	type TabVariant = 'pending' | 'upcoming' | 'past';

	interface TabConfig {
		id: TabId;
		label: string;
		icon: typeof Clock;
		count: number;
		requests: VacationRequest[];
		variant: TabVariant;
	}

	// ===================
	// Year Filter Logic (for Past tab)
	// ===================

	// Extract unique years from past requests, sorted descending
	const availableYears = $derived(
		Array.from(new Set(pastRequests.map((r) => new Date(r.endDate).getFullYear()))).sort(
			(a, b) => b - a
		)
	);

	// Default to most recent year with data, fallback to current year
	const mostRecentYear = $derived(availableYears[0] ?? new Date().getFullYear());

	// Selected year state (initialized lazily)
	let selectedYear = $state<number | null>(null);

	// Active year: user selection or default
	const activeYear = $derived(selectedYear ?? mostRecentYear);

	// Initialize on first render when data is available
	$effect(() => {
		if (selectedYear === null && availableYears.length > 0) {
			selectedYear = mostRecentYear;
		}
	});

	// Filtered past requests by year
	const filteredPastRequests = $derived(
		pastRequests.filter((r) => new Date(r.endDate).getFullYear() === activeYear)
	);

	// Year options for select dropdown
	const yearOptions = $derived(
		availableYears.map((year) => ({ value: year.toString(), label: year.toString() }))
	);

	// ===================
	// Melt UI: Year Select
	// ===================
	const {
		elements: { trigger: yearTrigger, menu: yearMenu, option: yearOption },
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

	// Sync select display with activeYear
	$effect(() => {
		const yearStr = activeYear.toString();
		if (!$yearSelected || $yearSelected.value !== yearStr) {
			yearSelected.set({ value: yearStr, label: yearStr });
		}
	});

	// ===================
	// Tab Configuration
	// ===================
	const tabs = $derived<TabConfig[]>([
		{
			id: 'pending',
			label: 'Pending',
			icon: Clock,
			count: pendingRequests.length,
			requests: pendingRequests.slice(0, MAX_REQUESTS_PER_TAB),
			variant: 'pending'
		},
		{
			id: 'upcoming',
			label: 'Upcoming',
			icon: CalendarCheck,
			count: upcomingRequests.length,
			requests: upcomingRequests.slice(0, MAX_REQUESTS_PER_TAB),
			variant: 'upcoming'
		},
		{
			id: 'past',
			label: 'Past',
			icon: History,
			count: pastRequests.length,
			requests: filteredPastRequests.slice(0, MAX_REQUESTS_PER_TAB),
			variant: 'past'
		}
	]);

	// Default tab: pending if any, then upcoming, then past
	const defaultTab = $derived<TabId>(
		pendingRequests.length > 0 ? 'pending' : upcomingRequests.length > 0 ? 'upcoming' : 'past'
	);

	// ===================
	// Melt UI: Tabs
	// ===================
	const {
		elements: { root, list, trigger, content },
		states: { value: activeTab }
	} = createTabs({
		defaultValue: defaultTab,
		loop: true,
		activateOnFocus: true
	});

	// Get current tab config
	const currentTab = $derived(tabs.find((t) => t.id === $activeTab) ?? tabs[0]);
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
								data-state={$activeTab === tab.id ? 'active' : 'inactive'}
							>
								{tab.count}
							</span>
						{/if}
					</button>
				{/each}
			</div>
		</div>

		<!-- Tab Content with Transitions -->
		{#each tabs as tab (tab.id)}
			<div use:melt={$content(tab.id)} class="focus:outline-none">
				{#key tab.id === $activeTab ? $activeTab : null}
					{#if tab.id === $activeTab}
						<div class="p-4" in:fly={{ y: 8, duration: 200, delay: 50 }} out:fade={{ duration: 100 }}>
							{#if tab.requests.length === 0}
								<!-- Empty States -->
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
										<p class="text-sm text-ocean-400 -mt-2">
											Your completed requests will appear here
										</p>
									</EmptyState>
								{/if}
							{:else}
								<!-- Year Filter (Past tab only) -->
								{#if tab.id === 'past' && availableYears.length > 0}
									<div class="flex items-center justify-between mb-4">
										<p class="text-sm text-ocean-500">
											{filteredPastRequests.length} request{filteredPastRequests.length !== 1
												? 's'
												: ''}
										</p>
										<div class="relative inline-block">
											<button
												use:melt={$yearTrigger}
												class="flex items-center p-1 text-sm font-medium rounded-lg bg-sand-100 focus:outline-none cursor-pointer"
											>
												<span
													class="flex items-center gap-1.5 px-3 py-1.5 rounded-md bg-white text-ocean-700 shadow-sm"
												>
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
													transition:fly={{ y: -4, duration: 150 }}
												>
													{#each yearOptions as opt (opt.value)}
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
								<div class="space-y-3">
									{#each tab.requests as request (request.id)}
										<EnhancedRequestCard {request} variant={tab.variant} />
									{/each}
								</div>
							{/if}
						</div>
					{/if}
				{/key}
			</div>
		{/each}
	</div>
</Card>
