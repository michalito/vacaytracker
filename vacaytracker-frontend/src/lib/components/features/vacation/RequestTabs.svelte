<script lang="ts">
	import { Tabs, Select } from 'bits-ui';
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

	// Year straight from the ISO date string; `new Date('YYYY-MM-DD')` is UTC and can
	// land in the previous year for users west of UTC
	function yearOf(isoDate: string): number {
		return Number(isoDate.slice(0, 4));
	}

	// Extract unique years from past requests, sorted descending
	const availableYears = $derived(
		Array.from(new Set(pastRequests.map((r) => yearOf(r.endDate)))).sort((a, b) => b - a)
	);

	// Selected year; falls back to the most recent year with data
	let selectedYear = $state<number | null>(null);
	const activeYear = $derived(
		selectedYear !== null && availableYears.includes(selectedYear)
			? selectedYear
			: (availableYears[0] ?? new Date().getFullYear())
	);

	// Filtered past requests by year
	const filteredPastRequests = $derived(
		pastRequests.filter((r) => yearOf(r.endDate) === activeYear)
	);

	// Year options for select dropdown
	const yearOptions = $derived(
		availableYears.map((year) => ({ value: year.toString(), label: year.toString() }))
	);

	let yearSelectOpen = $state(false);

	function handleYearChange(next: string) {
		if (next) selectedYear = Number(next);
	}

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

	// Default tab follows the data (pending, then upcoming, then past) until the user picks one,
	// so it is still right when requests arrive after the first render
	let chosenTab = $state<TabId | null>(null);
	const defaultTab = $derived<TabId>(
		pendingRequests.length > 0 ? 'pending' : upcomingRequests.length > 0 ? 'upcoming' : 'past'
	);
	const activeTab = $derived<TabId>(chosenTab ?? defaultTab);

	function isTabId(value: string): value is TabId {
		return value === 'pending' || value === 'upcoming' || value === 'past';
	}
</script>

<Card padding="none">
	{#snippet header()}
		<div class="flex items-center justify-between">
			<h2 class="text-lg font-semibold text-ocean-700">Your Requests</h2>
			<Umbrella class="w-5 h-5 text-ocean-400" />
		</div>
	{/snippet}

	<Tabs.Root
		bind:value={
			() => activeTab,
			(next) => {
				if (isTabId(next)) chosenTab = next;
			}
		}
		class="w-full"
	>
		<!-- Tab List -->
		<div class="px-4 pt-4">
			<Tabs.List class="inline-flex p-1 bg-sand-100 rounded-lg" aria-label="Request categories">
				{#each tabs as tab (tab.id)}
					<Tabs.Trigger
						value={tab.id}
						class="flex items-center gap-2 px-3 py-1.5 text-sm font-medium rounded-md transition-colors cursor-pointer
							data-[state=active]:bg-white data-[state=active]:text-ocean-700 data-[state=active]:shadow-sm
							data-[state=inactive]:text-ocean-500 data-[state=inactive]:hover:text-ocean-700"
					>
						<tab.icon class="w-4 h-4" />
						<span>{tab.label}</span>
						{#if tab.count > 0}
							<span
								class="px-1.5 py-0.5 text-xs rounded-md transition-colors
									{activeTab === tab.id ? 'bg-ocean-100 text-ocean-700' : 'bg-sand-200 text-ocean-500'}"
							>
								{tab.count}
							</span>
						{/if}
					</Tabs.Trigger>
				{/each}
			</Tabs.List>
		</div>

		<!-- Tab panels (the fade-in comes from the [data-tabs-content] rule in app.css) -->
		{#each tabs as tab (tab.id)}
			<Tabs.Content value={tab.id} class="focus:outline-none">
				{#if tab.id === activeTab}
					<div class="p-4">
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
										<Select.Root
											type="single"
											items={yearOptions}
											value={activeYear.toString()}
											onValueChange={handleYearChange}
											bind:open={yearSelectOpen}
										>
											<Select.Trigger
												aria-label="Filter past requests by year"
												class="flex items-center p-1 text-sm font-medium rounded-lg bg-sand-100 cursor-pointer
														focus:outline-none focus-visible:ring-2 focus-visible:ring-ocean-500"
											>
												<span
													class="flex items-center gap-1.5 px-3 py-1.5 rounded-md bg-white text-ocean-700 shadow-sm"
												>
													{activeYear}
													<ChevronDown
														class="w-4 h-4 transition-transform duration-200 {yearSelectOpen
															? 'rotate-180'
															: ''}"
														aria-hidden="true"
													/>
												</span>
											</Select.Trigger>

											<Select.Portal>
												<Select.Content
													class="presence-pop z-50 w-max rounded-lg bg-sand-100 p-1 shadow-lg"
													sideOffset={4}
													align="end"
												>
													{#each yearOptions as opt (opt.value)}
														<Select.Item
															value={opt.value}
															label={opt.label}
															class="flex items-center justify-between gap-3 rounded-md px-3 py-1.5 text-sm cursor-pointer outline-none whitespace-nowrap
																text-ocean-500
																data-[highlighted]:bg-white data-[highlighted]:text-ocean-700 data-[highlighted]:shadow-sm
																data-[selected]:font-medium"
														>
															{#snippet children({ selected })}
																<span>{opt.label}</span>
																{#if selected}
																	<Check class="w-4 h-4 text-ocean-500" />
																{/if}
															{/snippet}
														</Select.Item>
													{/each}
												</Select.Content>
											</Select.Portal>
										</Select.Root>
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
			</Tabs.Content>
		{/each}
	</Tabs.Root>
</Card>
