<script lang="ts">
	import { createTabs, melt } from '@melt-ui/svelte';
	import type { VacationRequest } from '$lib/types';
	import EnhancedRequestCard from './EnhancedRequestCard.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Clock, CalendarCheck, History, Umbrella, PartyPopper, Archive } from 'lucide-svelte';

	interface Props {
		pendingRequests: VacationRequest[];
		upcomingRequests: VacationRequest[];
		pastRequests: VacationRequest[];
		onRequestVacation?: () => void;
	}

	let { pendingRequests, upcomingRequests, pastRequests, onRequestVacation }: Props = $props();

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
			requests: pastRequests,
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
		<!-- Tab List (Pill Style) -->
		<div class="px-4 pt-4">
			<div
				use:melt={$list}
				class="inline-flex p-1 bg-ocean-100/60 backdrop-blur-sm rounded-xl border border-white/30"
				aria-label="Request categories"
			>
				{#each tabs as tab (tab.id)}
					<button
						use:melt={$trigger(tab.id)}
						class="flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-lg transition-all cursor-pointer
							data-[state=active]:bg-white data-[state=active]:text-ocean-800 data-[state=active]:shadow-sm
							data-[state=inactive]:text-ocean-600 data-[state=inactive]:hover:text-ocean-800"
					>
						<tab.icon class="w-4 h-4" />
						<span>{tab.label}</span>
						{#if tab.count > 0}
							<span
								class="px-1.5 py-0.5 text-xs rounded-md transition-colors
									data-[state=active]:bg-ocean-100 data-[state=active]:text-ocean-700
									data-[state=inactive]:bg-ocean-200/80 data-[state=inactive]:text-ocean-600"
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
							<EmptyState icon={PartyPopper} message="No pending requests" iconSize="lg">
								<p class="text-sm text-ocean-400 -mt-2">All caught up!</p>
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
