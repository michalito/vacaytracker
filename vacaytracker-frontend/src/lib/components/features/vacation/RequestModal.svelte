<script lang="ts">
	import { createDialog, createDateRangePicker, melt } from '@melt-ui/svelte';
	import { today, getLocalTimeZone } from '@internationalized/date';
	import { vacation } from '$lib/stores/vacation.svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { admin } from '$lib/stores/admin.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { dateValueToApiFormat, calculateBusinessDays } from '$lib/utils/date';
	import Button from '$lib/components/ui/Button.svelte';
	import { X, Palmtree, CalendarDays, ChevronLeft, ChevronRight } from 'lucide-svelte';

	interface Props {
		open: boolean;
	}

	let { open = $bindable(false) }: Props = $props();

	// Fetch settings when modal opens (for weekend policy)
	$effect(() => {
		if (open) {
			admin.fetchSettings();
		}
	});

	// Create Melt-UI dialog with controlled state
	const {
		elements: { overlay, content, title, close, portalled },
		states: { open: dialogOpen }
	} = createDialog({
		forceVisible: true,
		onOpenChange: ({ next }) => {
			open = next;
			if (!next) resetForm();
			return next;
		}
	});

	// Sync external open prop with dialog state
	$effect(() => {
		dialogOpen.set(open);
	});

	// Create Date Range Picker
	const excludeWeekends = $derived(admin.settings?.weekendPolicy?.excludeWeekends ?? true);

	// Helper to check if a date is a weekend
	function isWeekend(date: import('@internationalized/date').DateValue): boolean {
		const day = date.toDate(getLocalTimeZone()).getDay();
		return day === 0 || day === 6; // Sunday = 0, Saturday = 6
	}

	const {
		elements: {
			trigger: calTrigger,
			content: calContent,
			calendar,
			cell,
			heading,
			grid,
			prevButton,
			nextButton,
			field,
			startSegment,
			endSegment,
			label: fieldLabel
		},
		states: { open: calendarOpen, value: dateRangeValue, months, weekdays, segmentContents }
	} = createDateRangePicker({
		locale: 'en-GB',
		weekStartsOn: 1,
		numberOfMonths: 2,
		minValue: today(getLocalTimeZone()),
		forceVisible: true
		// Note: We don't use isDateDisabled for weekends because users need to select
		// ranges that span across weekends. Weekends are styled differently and
		// excluded from the business days calculation instead.
	});

	let reason = $state('');
	let isSubmitting = $state(false);
	let errors = $state<{ dateRange?: string }>({});

	function resetForm() {
		dateRangeValue.set({ start: undefined, end: undefined });
		reason = '';
		errors = {};
	}

	function validate(): boolean {
		errors = {};
		const range = $dateRangeValue;

		if (!range?.start || !range?.end) {
			errors.dateRange = 'Please select a date range';
			return false;
		}

		// Check balance
		if (estimatedDays > (auth.user?.vacationBalance ?? 0)) {
			errors.dateRange = `Request exceeds your available balance of ${auth.user?.vacationBalance} days`;
			return false;
		}

		return true;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();

		if (!validate()) return;

		isSubmitting = true;
		const range = $dateRangeValue;

		try {
			await vacation.createRequest({
				startDate: dateValueToApiFormat(range.start!),
				endDate: dateValueToApiFormat(range.end!),
				reason: reason || undefined
			});

			toast.success('Vacation request submitted!');
			dialogOpen.set(false);
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to submit request');
		} finally {
			isSubmitting = false;
		}
	}

	// Calculate estimated business days
	const estimatedDays = $derived.by(() => {
		const range = $dateRangeValue;
		if (!range?.start || !range?.end) return 0;
		return calculateBusinessDays(range.start, range.end, excludeWeekends);
	});

	// Remaining balance after request
	const remainingBalance = $derived((auth.user?.vacationBalance ?? 0) - estimatedDays);
</script>

{#if open}
	<div use:melt={$portalled}>
		<!-- Overlay -->
		<div
			use:melt={$overlay}
			class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm transition-opacity duration-200
				data-[state=open]:opacity-100 data-[state=closed]:opacity-0"
		></div>

		<!-- Content Container -->
		<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
			<div
				use:melt={$content}
				class="bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30 w-full max-w-md
					transition-all duration-200 data-[state=open]:opacity-100 data-[state=open]:scale-100
					data-[state=closed]:opacity-0 data-[state=closed]:scale-95"
			>
				<!-- Header -->
				<div class="flex items-center justify-between p-4 border-b border-ocean-100/50">
					<div class="flex items-center gap-2">
						<Palmtree class="w-5 h-5 text-ocean-500" />
						<h2 use:melt={$title} class="text-lg font-semibold text-ocean-800">Request Vacation</h2>
					</div>
					<button
						use:melt={$close}
						type="button"
						aria-label="Close dialog"
						class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600 hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
					>
						<X class="w-5 h-5" />
					</button>
				</div>

				<!-- Body -->
				<form onsubmit={handleSubmit} class="p-4 space-y-4">
					<!-- Date Range Picker -->
					<div class="flex flex-col gap-1.5">
						<span use:melt={$fieldLabel} class="text-sm font-semibold text-ocean-800">
							Vacation Period <span class="text-coral-500">*</span>
						</span>

						<div class="flex items-center gap-2">
							<!-- Segmented Date Field -->
							<div
								use:melt={$field}
								class="flex-1 flex items-center gap-2 px-4 py-2.5 rounded-lg border-2
									bg-white transition-all duration-200
									{errors.dateRange ? 'border-error' : 'border-ocean-500/40'}
									focus-within:ring-2 focus-within:ring-ocean-500/30
									focus-within:border-ocean-500 hover:border-ocean-500"
							>
								<!-- Start Date Segments -->
								<div class="flex items-center gap-0.5 tabular-nums">
									{#each $segmentContents.start as seg}
										{#if seg.part === 'literal'}
											<span class="text-ocean-400">{seg.value}</span>
										{:else}
											<span
												use:melt={$startSegment(seg.part)}
												class="rounded px-0.5 outline-none focus:bg-ocean-100
													data-[placeholder]:text-ocean-400"
											>
												{seg.value}
											</span>
										{/if}
									{/each}
								</div>

								<span class="text-ocean-400 px-1">to</span>

								<!-- End Date Segments -->
								<div class="flex items-center gap-0.5 tabular-nums">
									{#each $segmentContents.end as seg}
										{#if seg.part === 'literal'}
											<span class="text-ocean-400">{seg.value}</span>
										{:else}
											<span
												use:melt={$endSegment(seg.part)}
												class="rounded px-0.5 outline-none focus:bg-ocean-100
													data-[placeholder]:text-ocean-400"
											>
												{seg.value}
											</span>
										{/if}
									{/each}
								</div>
							</div>

							<!-- Calendar Trigger Button -->
							<button
								use:melt={$calTrigger}
								type="button"
								class="p-2.5 rounded-lg border-2 border-ocean-500/40 bg-white
									hover:bg-ocean-50 hover:border-ocean-500
									transition-all duration-200 cursor-pointer
									focus:outline-none focus:ring-2 focus:ring-ocean-500/30"
								aria-label="Open calendar"
							>
								<CalendarDays class="h-5 w-5 text-ocean-500" />
							</button>
						</div>

						{#if errors.dateRange}
							<p class="text-sm text-error">{errors.dateRange}</p>
						{/if}
					</div>

					<!-- Calendar Popup -->
					{#if $calendarOpen}
						<div
							use:melt={$calContent}
							class="z-[60] bg-white/95 backdrop-blur-md rounded-xl shadow-xl
								border border-ocean-200 p-4 mt-2"
						>
							<div use:melt={$calendar} class="flex flex-col sm:flex-row gap-4">
								{#each $months as month, i}
									<div class="flex-1 min-w-[260px] calendar-month" class:hidden-on-mobile={i > 0}>
										<!-- Month Header -->
										<div class="flex items-center justify-between mb-4">
											{#if i === 0}
												<button
													use:melt={$prevButton}
													type="button"
													class="p-1.5 rounded-lg hover:bg-ocean-100
														transition-colors cursor-pointer
														disabled:opacity-40 disabled:cursor-not-allowed"
												>
													<ChevronLeft class="h-5 w-5 text-ocean-600" />
												</button>
											{:else}
												<div class="w-8"></div>
											{/if}

											<div use:melt={$heading} class="font-semibold text-ocean-800">
												{month.value
													.toDate(getLocalTimeZone())
													.toLocaleDateString('en-GB', {
														month: 'long',
														year: 'numeric'
													})}
											</div>

											{#if i === $months.length - 1}
												<button
													use:melt={$nextButton}
													type="button"
													class="p-1.5 rounded-lg hover:bg-ocean-100
														transition-colors cursor-pointer
														disabled:opacity-40 disabled:cursor-not-allowed"
												>
													<ChevronRight class="h-5 w-5 text-ocean-600" />
												</button>
											{:else}
												<div class="w-8"></div>
											{/if}
										</div>

										<!-- Calendar Grid -->
										<table use:melt={$grid} class="w-full border-collapse">
											<thead>
												<tr>
													{#each $weekdays as day}
														<th class="text-xs font-medium text-ocean-500 pb-2 w-9 text-center">
															{day.slice(0, 2)}
														</th>
													{/each}
												</tr>
											</thead>
											<tbody>
												{#each month.weeks as week}
													<tr>
														{#each week as day}
															<td class="p-0.5 text-center">
																<div
																	use:melt={$cell(day, month.value)}
																	class="h-9 w-9 mx-auto flex items-center justify-center
																		text-sm cursor-pointer transition-all rounded-lg
																		data-[outside-month]:text-ocean-300
																		data-[outside-month]:opacity-50
																		data-[selection-start]:bg-ocean-500
																		data-[selection-start]:text-white
																		data-[selection-start]:rounded-l-lg
																		data-[selection-start]:rounded-r-none
																		data-[selection-end]:bg-ocean-500
																		data-[selection-end]:text-white
																		data-[selection-end]:rounded-r-lg
																		data-[selection-end]:rounded-l-none
																		data-[highlighted]:bg-ocean-100
																		data-[today]:ring-2 data-[today]:ring-ocean-400
																		data-[disabled]:opacity-30
																		data-[disabled]:cursor-not-allowed
																		data-[disabled]:line-through
																		hover:bg-ocean-200
																		focus:outline-none focus:ring-2 focus:ring-ocean-500
																		{excludeWeekends && isWeekend(day) ? 'text-ocean-400 italic' : ''}"
																>
																	{day.day}
																</div>
															</td>
														{/each}
													</tr>
												{/each}
											</tbody>
										</table>
									</div>
								{/each}
							</div>

							<!-- Selection Footer -->
							<div class="flex items-center justify-between mt-4 pt-4 border-t border-ocean-100">
								<p class="text-sm text-ocean-600">
									{#if estimatedDays > 0}
										<span class="font-medium">{estimatedDays}</span> business {estimatedDays === 1 ? 'day' : 'days'}
										{#if excludeWeekends}
											<span class="text-ocean-400">(weekends excluded)</span>
										{/if}
									{:else}
										Select start and end dates
									{/if}
								</p>
							</div>
						</div>
					{/if}

					<!-- Reason Textarea -->
					<div>
						<label for="reason" class="block text-sm font-semibold text-ocean-800 mb-1.5">
							Reason (optional)
						</label>
						<textarea
							id="reason"
							bind:value={reason}
							class="w-full px-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 placeholder-ocean-500/50 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500 hover:border-ocean-500 resize-none"
							rows="3"
							placeholder="Family vacation, personal time, etc."
						></textarea>
					</div>

					<!-- Balance Info -->
					<div class="bg-ocean-500/10 rounded-xl p-4 space-y-2">
						<div class="flex justify-between text-sm">
							<span class="text-ocean-700">Current Balance:</span>
							<span class="font-medium text-ocean-900">{auth.user?.vacationBalance ?? 0} days</span>
						</div>
						{#if estimatedDays > 0}
							<div class="flex justify-between text-sm">
								<span class="text-ocean-700">Days Requested:</span>
								<span class="font-medium text-ocean-900">{estimatedDays} days</span>
							</div>
							<div class="flex justify-between text-sm">
								<span class="text-ocean-700">Remaining After:</span>
								<span
									class="font-medium {remainingBalance < 0 ? 'text-error' : 'text-ocean-900'}"
								>
									{remainingBalance} days
								</span>
							</div>
							{#if excludeWeekends}
								<p class="text-xs text-ocean-500 mt-1">* Weekends are excluded per company policy</p>
							{/if}
						{/if}
					</div>

					<!-- Actions -->
					<div class="flex gap-3 pt-2">
						<Button
							type="button"
							variant="outline"
							class="flex-1"
							onclick={() => dialogOpen.set(false)}
						>
							Cancel
						</Button>
						<Button type="submit" variant="primary" class="flex-1" loading={isSubmitting}>
							Submit Request
						</Button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}

<style>
	/* Hide second month on mobile */
	@media (max-width: 639px) {
		.hidden-on-mobile {
			display: none;
		}
	}
</style>
