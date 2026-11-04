<script lang="ts">
	import { DateRangePicker, Portal, type DateRange } from 'bits-ui';
	import { today, getLocalTimeZone } from '@internationalized/date';
	import { untrack } from 'svelte';
	import { MediaQuery } from 'svelte/reactivity';
	import { clsx } from 'clsx';
	import type { WeekendPolicy } from '#lib/types/index.js';
	import { vacation } from '#lib/stores/vacation.svelte.js';
	import { auth } from '#lib/stores/auth.svelte.js';
	import { toast } from '#lib/stores/toast.svelte.js';
	import {
		dateValueToApiFormat,
		calculateBusinessDays,
		isExcludedDay,
		DEFAULT_EXCLUDED_DAYS
	} from '#lib/utils/date.js';
	import Modal from '#lib/components/ui/Modal.svelte';
	import Button from '#lib/components/ui/Button.svelte';
	import { Palmtree, CalendarDays, ChevronLeft, ChevronRight } from 'lucide-svelte';

	interface Props {
		open: boolean;
		/** Company weekend policy, so the preview counts days the same way the API does. */
		weekendPolicy?: WeekendPolicy;
	}

	let { open = $bindable(false), weekendPolicy }: Props = $props();

	// Two months side by side from the `sm` breakpoint. Bits must be told how many months are
	// shown: hiding one with CSS would leave its days in keyboard navigation.
	const wideCalendar = new MediaQuery('min-width: 40rem', true);

	const excludeWeekends = $derived(weekendPolicy?.excludeWeekends ?? true);
	const excludedDays = $derived(weekendPolicy?.excludedDays ?? DEFAULT_EXCLUDED_DAYS);

	let dateRange = $state<DateRange>({ start: undefined, end: undefined });
	let reason = $state('');
	let isSubmitting = $state(false);
	let dateError = $state<string | null>(null);

	// Every opening starts from a clean form, no matter how the previous one was closed
	$effect.pre(() => {
		if (open) untrack(resetForm);
	});

	function resetForm() {
		dateRange = { start: undefined, end: undefined };
		reason = '';
		dateError = null;
	}

	const balance = $derived(auth.user?.vacationBalance ?? 0);

	const estimatedDays = $derived(
		dateRange.start && dateRange.end
			? calculateBusinessDays(dateRange.start, dateRange.end, excludeWeekends, excludedDays)
			: 0
	);

	const remainingBalance = $derived(balance - estimatedDays);

	function validate(): boolean {
		if (!dateRange.start || !dateRange.end) {
			dateError = 'Please select a date range';
		} else if (estimatedDays === 0) {
			dateError = 'The selected dates contain no working days';
		} else if (estimatedDays > balance) {
			dateError = `Request exceeds your available balance of ${balance} days`;
		} else {
			dateError = null;
		}
		return dateError === null;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!validate() || !dateRange.start || !dateRange.end) return;

		isSubmitting = true;
		try {
			await vacation.createRequest({
				startDate: dateValueToApiFormat(dateRange.start),
				endDate: dateValueToApiFormat(dateRange.end),
				reason: reason.trim() || undefined
			});
			toast.success('Vacation request submitted!');
			open = false;
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to submit request');
		} finally {
			isSubmitting = false;
		}
	}
</script>

{#snippet dateInput(type: 'start' | 'end')}
	<DateRangePicker.Input {type} class="flex items-center tabular-nums">
		{#snippet children({ segments })}
			{#each segments as { part, value }, i (i)}
				<DateRangePicker.Segment
					{part}
					class={part === 'literal'
						? 'px-0.5 text-ocean-400'
						: 'rounded px-0.5 outline-none focus:bg-ocean-100 data-[placeholder]:text-ocean-400'}
				>
					{value}
				</DateRangePicker.Segment>
			{/each}
		{/snippet}
	</DateRangePicker.Input>
{/snippet}

<Modal bind:open title="Request Vacation" icon={Palmtree} closeLabel="Close request form">
	<form onsubmit={handleSubmit} class="space-y-4 p-4" novalidate>
		<DateRangePicker.Root
			bind:value={dateRange}
			onValueChange={() => (dateError = null)}
			locale="en-GB"
			weekStartsOn={1}
			numberOfMonths={wideCalendar.current ? 2 : 1}
			minValue={today(getLocalTimeZone())}
			weekdayFormat="short"
			fixedWeeks
		>
			<div class="flex flex-col gap-1.5">
				<DateRangePicker.Label class="text-sm font-semibold text-ocean-800">
					Vacation Period <span class="text-coral-500" aria-hidden="true">*</span>
				</DateRangePicker.Label>

				<div class="flex items-center gap-2">
					<div
						class={clsx(
							'flex flex-1 flex-wrap items-center gap-x-2 gap-y-1 rounded-lg border-2 bg-white px-4 py-2.5',
							'transition-all duration-200 hover:border-ocean-500',
							'focus-within:border-ocean-500 focus-within:ring-2 focus-within:ring-ocean-500/30',
							dateError ? 'border-error' : 'border-ocean-500/40'
						)}
					>
						{@render dateInput('start')}
						<span class="text-ocean-400" aria-hidden="true">to</span>
						{@render dateInput('end')}
					</div>

					<DateRangePicker.Trigger
						aria-label="Open calendar"
						class="cursor-pointer rounded-lg border-2 border-ocean-500/40 bg-white p-2.5 transition-all duration-200
							hover:border-ocean-500 hover:bg-ocean-50
							focus:outline-none focus:ring-2 focus:ring-ocean-500/30"
					>
						<CalendarDays class="h-5 w-5 text-ocean-500" aria-hidden="true" />
					</DateRangePicker.Trigger>
				</div>

				{#if dateError}
					<p class="text-sm text-error" role="alert">{dateError}</p>
				{/if}
			</div>

			<!--
				Portalled so the dialog's transform/backdrop-filter (which create a containing
				block for fixed elements) cannot clip or offset the calendar.
			-->
			<Portal>
				<DateRangePicker.Content
					sideOffset={6}
					align="end"
					collisionPadding={8}
					class="presence-pop z-[60] rounded-xl border border-ocean-200 bg-white/95 p-4 shadow-xl backdrop-blur-md"
				>
					<DateRangePicker.Calendar>
						{#snippet children({ months, weekdays })}
							<DateRangePicker.Header class="flex items-center justify-between">
								<DateRangePicker.PrevButton
									aria-label="Previous month"
									class="cursor-pointer rounded-lg p-1.5 transition-colors hover:bg-ocean-100
										disabled:cursor-not-allowed disabled:opacity-40"
								>
									<ChevronLeft class="h-5 w-5 text-ocean-600" aria-hidden="true" />
								</DateRangePicker.PrevButton>
								<DateRangePicker.Heading class="font-semibold text-ocean-800" />
								<DateRangePicker.NextButton
									aria-label="Next month"
									class="cursor-pointer rounded-lg p-1.5 transition-colors hover:bg-ocean-100
										disabled:cursor-not-allowed disabled:opacity-40"
								>
									<ChevronRight class="h-5 w-5 text-ocean-600" aria-hidden="true" />
								</DateRangePicker.NextButton>
							</DateRangePicker.Header>

							<div class="flex gap-6 pt-4">
								{#each months as month (month.value.toString())}
									<DateRangePicker.Grid class="border-collapse select-none">
										<DateRangePicker.GridHead>
											<DateRangePicker.GridRow class="flex">
												{#each weekdays as day (day)}
													<DateRangePicker.HeadCell
														class="w-9 pb-2 text-center text-xs font-medium text-ocean-500"
													>
														{day.slice(0, 2)}
													</DateRangePicker.HeadCell>
												{/each}
											</DateRangePicker.GridRow>
										</DateRangePicker.GridHead>
										<DateRangePicker.GridBody>
											{#each month.weeks as week, w (w)}
												<DateRangePicker.GridRow class="flex">
													{#each week as date (date.toString())}
														<DateRangePicker.Cell {date} month={month.value} class="size-9 p-0">
															<DateRangePicker.Day
																class="calendar-day"
																data-excluded={isExcludedDay(date, excludeWeekends, excludedDays)
																	? ''
																	: undefined}
															>
																{date.day}
															</DateRangePicker.Day>
														</DateRangePicker.Cell>
													{/each}
												</DateRangePicker.GridRow>
											{/each}
										</DateRangePicker.GridBody>
									</DateRangePicker.Grid>
								{/each}
							</div>
						{/snippet}
					</DateRangePicker.Calendar>

					<p class="mt-4 border-t border-ocean-100 pt-4 text-sm text-ocean-600">
						{#if estimatedDays > 0}
							<span class="font-medium">{estimatedDays}</span>
							business {estimatedDays === 1 ? 'day' : 'days'}
							{#if excludeWeekends}
								<span class="text-ocean-400">(non-working days excluded)</span>
							{/if}
						{:else}
							Select start and end dates
						{/if}
					</p>
				</DateRangePicker.Content>
			</Portal>
		</DateRangePicker.Root>

		<div>
			<label for="vacation-reason" class="mb-1.5 block text-sm font-semibold text-ocean-800">
				Reason (optional)
			</label>
			<textarea
				id="vacation-reason"
				bind:value={reason}
				rows="3"
				placeholder="Family vacation, personal time, etc."
				class="w-full resize-none rounded-lg border-2 border-ocean-500/40 bg-white px-4 py-2.5 text-ocean-900
					placeholder-ocean-500/50 transition-all duration-200 hover:border-ocean-500
					focus:border-ocean-500 focus:outline-none focus:ring-2 focus:ring-ocean-500/30"
			></textarea>
		</div>

		<div class="space-y-2 rounded-xl bg-ocean-500/10 p-4">
			<div class="flex justify-between text-sm">
				<span class="text-ocean-700">Current Balance:</span>
				<span class="font-medium text-ocean-900">{balance} days</span>
			</div>
			{#if estimatedDays > 0}
				<div class="flex justify-between text-sm">
					<span class="text-ocean-700">Days Requested:</span>
					<span class="font-medium text-ocean-900">{estimatedDays} days</span>
				</div>
				<div class="flex justify-between text-sm">
					<span class="text-ocean-700">Remaining After:</span>
					<span class="font-medium {remainingBalance < 0 ? 'text-error' : 'text-ocean-900'}">
						{remainingBalance} days
					</span>
				</div>
				{#if excludeWeekends}
					<p class="mt-1 text-xs text-ocean-500">
						* Non-working days are excluded per company policy
					</p>
				{/if}
			{/if}
		</div>

		<div class="flex gap-3 pt-2">
			<Button type="button" variant="outline" class="flex-1" onclick={() => (open = false)}>
				Cancel
			</Button>
			<Button type="submit" variant="primary" class="flex-1" loading={isSubmitting}>
				Submit Request
			</Button>
		</div>
	</form>
</Modal>
