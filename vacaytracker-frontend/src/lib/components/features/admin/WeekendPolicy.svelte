<script lang="ts">
	import type { WeekendPolicy } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import { Calendar } from 'lucide-svelte';

	interface Props {
		policy: WeekendPolicy;
		onChange: (policy: WeekendPolicy) => void;
	}

	let { policy, onChange }: Props = $props();

	const weekDays = [
		{ value: 0, label: 'Sunday' },
		{ value: 1, label: 'Monday' },
		{ value: 2, label: 'Tuesday' },
		{ value: 3, label: 'Wednesday' },
		{ value: 4, label: 'Thursday' },
		{ value: 5, label: 'Friday' },
		{ value: 6, label: 'Saturday' }
	];

	function toggleExcludeWeekends() {
		onChange({
			...policy,
			excludeWeekends: !policy.excludeWeekends
		});
	}

	function toggleDay(day: number) {
		const newDays = policy.excludedDays.includes(day)
			? policy.excludedDays.filter((d) => d !== day)
			: [...policy.excludedDays, day];

		onChange({
			...policy,
			excludedDays: newDays
		});
	}
</script>

<Card>
	{#snippet header()}
		<div class="flex items-center gap-2">
			<Calendar class="w-5 h-5 text-ocean-500" />
			<h2 class="text-lg font-semibold text-ocean-700">Weekend Policy</h2>
		</div>
	{/snippet}

	<div class="space-y-4">
		<p class="text-sm text-ocean-600">
			Configure which days are excluded when calculating vacation days.
		</p>

		<label class="flex items-center gap-3 cursor-pointer">
			<input
				type="checkbox"
				checked={policy.excludeWeekends}
				onchange={toggleExcludeWeekends}
				class="w-5 h-5 rounded border-sand-300 text-ocean-500 focus:ring-ocean-500"
			/>
			<div>
				<p class="font-medium text-ocean-800">Exclude weekends from vacation days</p>
				<p class="text-sm text-ocean-500">
					When enabled, only business days are counted for vacations
				</p>
			</div>
		</label>

		{#if policy.excludeWeekends}
			<div class="pt-4 border-t border-sand-200">
				<p class="text-sm font-medium text-ocean-700 mb-3">Excluded Days</p>
				<div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
					{#each weekDays as day}
						<label
							class="flex items-center gap-2 p-2 rounded-lg border border-sand-200 cursor-pointer hover:bg-sand-50 transition-colors {policy.excludedDays.includes(
								day.value
							)
								? 'bg-ocean-50 border-ocean-200'
								: ''}"
						>
							<input
								type="checkbox"
								checked={policy.excludedDays.includes(day.value)}
								onchange={() => toggleDay(day.value)}
								class="w-4 h-4 rounded border-sand-300 text-ocean-500 focus:ring-ocean-500"
							/>
							<span class="text-sm text-ocean-700">{day.label}</span>
						</label>
					{/each}
				</div>
				<p class="text-xs text-ocean-500 mt-2">
					Selected days (0=Sunday, 6=Saturday) will not count as vacation days.
				</p>
			</div>
		{/if}
	</div>
</Card>
