<script lang="ts">
	import { createDialog, melt } from '@melt-ui/svelte';
	import { vacation } from '$lib/stores/vacation.svelte';
	import { auth } from '$lib/stores/auth.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { toApiDateFormat } from '$lib/utils/date';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import { X, Palmtree } from 'lucide-svelte';

	interface Props {
		open: boolean;
	}

	let { open = $bindable(false) }: Props = $props();

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

	let startDate = $state('');
	let endDate = $state('');
	let reason = $state('');
	let isSubmitting = $state(false);
	let errors = $state<{ startDate?: string; endDate?: string }>({});

	function resetForm() {
		startDate = '';
		endDate = '';
		reason = '';
		errors = {};
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

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();

		if (!validate()) return;

		isSubmitting = true;

		try {
			await vacation.createRequest({
				startDate: toApiDateFormat(startDate),
				endDate: toApiDateFormat(endDate),
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

	// Calculate estimated days (simple version - doesn't account for weekends)
	const estimatedDays = $derived.by(() => {
		if (!startDate || !endDate) return 0;
		const start = new Date(startDate);
		const end = new Date(endDate);
		const diff = Math.ceil((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)) + 1;
		return diff > 0 ? diff : 0;
	});
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
						class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600 hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
					>
						<X class="w-5 h-5" />
					</button>
				</div>

				<!-- Body -->
				<form onsubmit={handleSubmit} class="p-4 space-y-4">
					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
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
					<div class="bg-ocean-500/10 rounded-xl p-4 space-y-1">
						<div class="flex justify-between text-sm">
							<span class="text-ocean-700">Current Balance:</span>
							<span class="font-medium text-ocean-900">{auth.user?.vacationBalance ?? 0} days</span>
						</div>
						{#if estimatedDays > 0}
							<div class="flex justify-between text-sm">
								<span class="text-ocean-700">Estimated Days:</span>
								<span class="font-medium text-ocean-900">{estimatedDays} days</span>
							</div>
							<p class="text-xs text-ocean-500 mt-1">
								* Final days calculated by admin based on weekend policy
							</p>
						{/if}
					</div>

					<!-- Actions -->
					<div class="flex gap-3 pt-2">
						<Button type="button" variant="outline" class="flex-1" onclick={() => dialogOpen.set(false)}>
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
