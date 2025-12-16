<script lang="ts">
	import ProgressRing from '$lib/components/ui/ProgressRing.svelte';
	import Tooltip from '$lib/components/ui/Tooltip.svelte';
	import { Palmtree, Calendar, Info, Umbrella } from 'lucide-svelte';
	import type { VacationRequest } from '$lib/types';

	interface Props {
		available: number;
		used: number;
		nextVacation?: VacationRequest | null;
	}

	let { available, used, nextVacation = null }: Props = $props();

	// Derive total from available + used
	const total = $derived(available + used);
	const percentage = $derived(total > 0 ? Math.round((available / total) * 100) : 0);

	// Calculate days until next vacation
	const daysUntilVacation = $derived.by(() => {
		if (!nextVacation) return null;
		const today = new Date();
		today.setHours(0, 0, 0, 0);
		const startDate = new Date(nextVacation.startDate);
		const diffTime = startDate.getTime() - today.getTime();
		const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
		return diffDays;
	});

	// Format date range for display (e.g., "Dec 23-27")
	const formattedDateRange = $derived.by(() => {
		if (!nextVacation) return '';
		const start = new Date(nextVacation.startDate);
		const end = new Date(nextVacation.endDate);
		const startMonth = start.toLocaleDateString('en-US', { month: 'short' });
		const endMonth = end.toLocaleDateString('en-US', { month: 'short' });
		const startDay = start.getDate();
		const endDay = end.getDate();

		if (startMonth === endMonth) {
			return `${startMonth} ${startDay}-${endDay}`;
		}
		return `${startMonth} ${startDay} - ${endMonth} ${endDay}`;
	});

	// Dynamic message for next vacation
	const nextVacationMessage = $derived.by(() => {
		if (daysUntilVacation === null) return '';
		if (daysUntilVacation <= 0) return 'Enjoy your time off!';
		if (daysUntilVacation <= 7) return 'Time to start packing!';
		if (daysUntilVacation <= 30) return 'The countdown is on!';
		return 'Something to look forward to!';
	});

	// Dynamic message based on balance (when no upcoming vacation)
	const balanceMessage = $derived.by(() => {
		if (available === 0) return 'All adventured out this year!';
		if (percentage > 60) return 'Plenty of beach days ahead!';
		if (percentage > 30) return 'Still time for an adventure';
		return "Don't let these days slip away!";
	});
</script>

<div
	class="bg-white/90 backdrop-blur-sm rounded-xl shadow-md border border-white/40 p-6
		hover:shadow-lg transition-shadow duration-200"
>
	<!-- Header -->
	<div class="flex items-center justify-between mb-6">
		<h2 class="text-lg font-semibold text-ocean-700">Your Vacation Balance</h2>
		<Tooltip
			content="Your vacation balance shows how many days you have available this year. Days reset annually."
			placement="left"
		>
			<button
				type="button"
				class="p-1.5 rounded-full text-ocean-400 hover:text-ocean-600 hover:bg-ocean-100 transition-colors"
			>
				<Info class="w-4 h-4" />
			</button>
		</Tooltip>
	</div>

	<!-- Content -->
	<div class="flex flex-col md:flex-row gap-6 items-center md:items-start">
		<!-- Ring Section -->
		<div class="flex-shrink-0">
			<ProgressRing value={available} max={total} size={140} />
		</div>

		<!-- Stats Section -->
		<div class="flex-1 w-full space-y-4">
			<!-- Available Days -->
			<div class="flex items-center gap-3">
				<div class="p-2 rounded-lg bg-ocean-500/10">
					<Palmtree class="w-5 h-5 text-ocean-600" />
				</div>
				<div>
					<p class="text-xl font-bold text-ocean-800">{available} days</p>
					<p class="text-sm text-ocean-500">available</p>
				</div>
			</div>

			<!-- Used Days -->
			<div class="flex items-center gap-3">
				<div class="p-2 rounded-lg bg-sand-200">
					<Calendar class="w-5 h-5 text-ocean-500" />
				</div>
				<div>
					<p class="text-lg font-semibold text-ocean-700">{used} days</p>
					<p class="text-sm text-ocean-400">used</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Dynamic Status Message -->
	<div class="mt-5 pt-4 border-t border-ocean-100/50">
		{#if nextVacation}
			<div class="flex items-center gap-2">
				<div class="p-1.5 rounded-lg bg-coral-400/15">
					<Umbrella class="w-4 h-4 text-coral-500" />
				</div>
				<span class="font-medium text-ocean-700">Next escape: {formattedDateRange}</span>
				{#if daysUntilVacation !== null && daysUntilVacation > 0}
					<span class="text-sm text-ocean-400">({daysUntilVacation} days)</span>
				{/if}
			</div>
			<p class="text-sm text-ocean-500 mt-2 italic pl-9">{nextVacationMessage}</p>
		{:else}
			<p class="text-sm text-ocean-500 italic text-center md:text-left">{balanceMessage}</p>
		{/if}
	</div>
</div>
