<script lang="ts">
	import MultiSegmentRing from '$lib/components/ui/MultiSegmentRing.svelte';
	import Tooltip from '$lib/components/ui/Tooltip.svelte';
	import { Palmtree, CheckCircle, Plane, Info, Umbrella } from 'lucide-svelte';
	import type { VacationRequest } from '$lib/types';

	type Theme = 'healthy' | 'warning' | 'critical';

	interface Props {
		total: number; // Total yearly allowance from admin settings
		used: number;
		upcoming: number;
		nextVacation?: VacationRequest | null;
	}

	let { total, used, upcoming, nextVacation = null }: Props = $props();

	// Calculate available days: total - used - upcoming
	const available = $derived(Math.max(0, total - used - upcoming));
	const percentage = $derived(total > 0 ? Math.round((available / total) * 100) : 0);

	// Determine theme based on available percentage (spa/resort traffic light)
	const theme = $derived<Theme>(
		percentage > 50 ? 'healthy' : percentage >= 20 ? 'warning' : 'critical'
	);

	// Ring segments: [available, upcoming, used] - clockwise from 12 o'clock, matches legend order
	const segments = $derived([
		{ value: available, label: 'Available' },
		{ value: upcoming, label: 'Upcoming' },
		{ value: used, label: 'Used' }
	]);

	// Spa/Resort legend colors (available=dynamic, upcoming=caribbean, used=stone)
	const legendColors = $derived({
		available:
			theme === 'healthy'
				? 'bg-mint-400'
				: theme === 'warning'
					? 'bg-sunshine-400'
					: 'bg-salmon-400',
		upcoming: 'bg-caribbean-400',
		used: 'bg-stone-300'
	});

	// Spa/Resort icon colors (darker shades for contrast)
	const iconColors = $derived({
		available:
			theme === 'healthy'
				? 'text-mint-600'
				: theme === 'warning'
					? 'text-sunshine-600'
					: 'text-salmon-600',
		upcoming: 'text-caribbean-500',
		used: 'text-stone-400'
	});

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
			<MultiSegmentRing
				{segments}
				{theme}
				size={140}
				strokeWidth={12}
				centerLabel={String(available)}
				centerSubLabel="of {total}"
			/>
		</div>

		<!-- Stats Section -->
		<div class="flex-1 w-full space-y-3">
			<!-- Available Days -->
			<div class="flex items-center gap-3">
				<div class="w-3 h-3 rounded-full {legendColors.available}"></div>
				<div class="flex items-center gap-2">
					<Palmtree class="w-4 h-4 {iconColors.available}" />
					<span class="text-ocean-800 font-semibold tabular-nums min-w-[2ch] text-right">{available}</span>
					<span class="text-ocean-500 text-sm">days available</span>
				</div>
			</div>

			<!-- Upcoming Days -->
			{#if upcoming > 0}
				<div class="flex items-center gap-3">
					<div class="w-3 h-3 rounded-full {legendColors.upcoming}"></div>
					<div class="flex items-center gap-2">
						<Plane class="w-4 h-4 {iconColors.upcoming}" />
						<span class="text-ocean-800 font-semibold tabular-nums min-w-[2ch] text-right">{upcoming}</span>
						<span class="text-ocean-500 text-sm">days upcoming</span>
					</div>
				</div>
			{/if}

			<!-- Used Days -->
			{#if used > 0}
				<div class="flex items-center gap-3">
					<div class="w-3 h-3 rounded-full {legendColors.used}"></div>
					<div class="flex items-center gap-2">
						<CheckCircle class="w-4 h-4 {iconColors.used}" />
						<span class="text-ocean-800 font-semibold tabular-nums min-w-[2ch] text-right">{used}</span>
						<span class="text-ocean-500 text-sm">days used</span>
					</div>
				</div>
			{/if}
		</div>
	</div>

	<!-- Dynamic Status Message -->
	<div class="mt-5 pt-4 border-t border-ocean-100/50">
		{#if nextVacation}
			<div class="flex items-center gap-2">
				<div class="p-1.5 rounded-lg bg-caribbean-400/15">
					<Umbrella class="w-4 h-4 text-caribbean-500" />
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
