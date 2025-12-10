<script lang="ts">
	import { vacation } from '$lib/stores/vacation.svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import { X, Palmtree } from 'lucide-svelte';

	interface Props {
		open: boolean;
	}

	let { open = $bindable(false) }: Props = $props();

	let startDate = $state('');
	let endDate = $state('');
	let reason = $state('');
	let isSubmitting = $state(false);
	let errors = $state<{ startDate?: string; endDate?: string }>({});

	// Get today's date in YYYY-MM-DD format for min attribute
	const today = new Date().toISOString().split('T')[0];

	function resetForm() {
		startDate = '';
		endDate = '';
		reason = '';
		errors = {};
	}

	function handleClose() {
		open = false;
		resetForm();
	}

	function validate(): boolean {
		errors = {};

		if (!startDate) {
			errors.startDate = 'Start date is required';
		}

		if (!endDate) {
			errors.endDate = 'End date is required';
		}

		if (startDate && endDate && new Date(endDate) < new Date(startDate)) {
			errors.endDate = 'End date must be after start date';
		}

		return Object.keys(errors).length === 0;
	}

	// Convert YYYY-MM-DD to DD/MM/YYYY for API
	function formatDateForApi(dateStr: string): string {
		const [year, month, day] = dateStr.split('-');
		return `${day}/${month}/${year}`;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();

		if (!validate()) return;

		isSubmitting = true;

		try {
			await vacation.createRequest({
				startDate: formatDateForApi(startDate),
				endDate: formatDateForApi(endDate),
				reason: reason || undefined
			});

			toast.success('Vacation request submitted!');
			handleClose();

			// Refresh user to get updated balance
			// Note: Balance is only deducted when approved, so we don't update here
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to submit request');
		} finally {
			isSubmitting = false;
		}
	}

	// Calculate estimated days (simple version - doesn't account for weekends)
	const estimatedDays = $derived(() => {
		if (!startDate || !endDate) return 0;
		const start = new Date(startDate);
		const end = new Date(endDate);
		const diff = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)) + 1;
		return diff > 0 ? diff : 0;
	});
</script>

{#if open}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<div class="fixed inset-0 z-50">
		<!-- Overlay -->
		<div class="fixed inset-0 bg-black/50 backdrop-blur-sm" onclick={handleClose}></div>

		<!-- Content -->
		<div class="fixed inset-0 flex items-center justify-center p-4">
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div class="bg-white rounded-xl shadow-xl w-full max-w-md" onclick={(e) => e.stopPropagation()}>
				<!-- Header -->
				<div class="flex items-center justify-between p-4 border-b border-sand-200">
					<div class="flex items-center gap-2">
						<Palmtree class="w-5 h-5 text-ocean-500" />
						<h2 class="text-lg font-semibold text-ocean-800">Request Vacation</h2>
					</div>
					<button type="button" onclick={handleClose} class="text-ocean-400 hover:text-ocean-600">
						<X class="w-5 h-5" />
					</button>
				</div>

				<!-- Body -->
				<form onsubmit={handleSubmit} class="p-4 space-y-4">
					<div class="grid grid-cols-2 gap-4">
						<Input
							type="date"
							label="Start Date"
							bind:value={startDate}
							error={errors.startDate}
							required
						/>
						<Input
							type="date"
							label="End Date"
							bind:value={endDate}
							error={errors.endDate}
							required
						/>
					</div>

					<div>
						<label for="reason" class="block text-sm font-medium text-ocean-700 mb-1">
							Reason (optional)
						</label>
						<textarea
							id="reason"
							bind:value={reason}
							class="w-full px-3 py-2 rounded-md border border-sand-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent resize-none"
							rows="3"
							placeholder="Family vacation, personal time, etc."
						></textarea>
					</div>

					<!-- Balance Info -->
					<div class="bg-ocean-50 rounded-lg p-3 space-y-1">
						<div class="flex justify-between text-sm">
							<span class="text-ocean-700">Current Balance:</span>
							<span class="font-medium text-ocean-900">{auth.user?.vacationBalance ?? 0} days</span>
						</div>
						{#if estimatedDays() > 0}
							<div class="flex justify-between text-sm">
								<span class="text-ocean-700">Estimated Days:</span>
								<span class="font-medium text-ocean-900">{estimatedDays()} days</span>
							</div>
							<p class="text-xs text-ocean-500 mt-1">
								* Final days calculated by admin based on weekend policy
							</p>
						{/if}
					</div>

					<!-- Actions -->
					<div class="flex gap-3 pt-2">
						<Button type="button" variant="outline" class="flex-1" onclick={handleClose}>
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
