<script lang="ts">
	import type { VacationRequest } from '$lib/types';
	import { adminApi } from '$lib/api/admin';
	import { toast } from '$lib/stores/toast.svelte';
	import { formatDateShort } from '$lib/utils/date';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import { Check, X, Clock } from 'lucide-svelte';

	interface Props {
		requests: VacationRequest[];
		onUpdate: () => void;
	}

	let { requests, onUpdate }: Props = $props();

	let processingId = $state<string | null>(null);
	let rejectionReason = $state('');
	let showRejectModal = $state(false);
	let rejectingId = $state<string | null>(null);

	async function approve(id: string) {
		processingId = id;
		try {
			await adminApi.reviewRequest(id, 'approved');
			toast.success('Request approved');
			onUpdate();
		} catch (error) {
			toast.error('Failed to approve request');
		} finally {
			processingId = null;
		}
	}

	function openRejectModal(id: string) {
		rejectingId = id;
		rejectionReason = '';
		showRejectModal = true;
	}

	async function confirmReject() {
		if (!rejectingId) return;

		processingId = rejectingId;
		try {
			await adminApi.reviewRequest(rejectingId, 'rejected', rejectionReason || undefined);
			toast.success('Request rejected');
			showRejectModal = false;
			onUpdate();
		} catch (error) {
			toast.error('Failed to reject request');
		} finally {
			processingId = null;
			rejectingId = null;
		}
	}
</script>

{#if requests.length === 0}
	<EmptyState icon={Clock} message="No pending requests" iconSize="lg" />
{:else}
	<div class="divide-y divide-sand-200">
		{#each requests as request (request.id)}
			<div class="flex items-center justify-between py-4">
				<div class="flex items-center gap-4">
					<Avatar name={request.userName} size="md" />
					<div>
						<p class="font-medium text-ocean-800">{request.userName}</p>
						<p class="text-sm text-ocean-500">
							{formatDateShort(request.startDate)} - {formatDateShort(request.endDate)}
							({request.totalDays} days)
						</p>
						{#if request.reason}
							<p class="text-xs text-ocean-400 mt-1">{request.reason}</p>
						{/if}
					</div>
				</div>

				<div class="flex items-center gap-2">
					<Button
						variant="outline"
						size="sm"
						onclick={() => openRejectModal(request.id)}
						disabled={processingId !== null}
					>
						<X class="w-4 h-4 mr-1" />
						Reject
					</Button>
					<Button
						variant="primary"
						size="sm"
						onclick={() => approve(request.id)}
						loading={processingId === request.id}
						disabled={processingId !== null && processingId !== request.id}
					>
						<Check class="w-4 h-4 mr-1" />
						Approve
					</Button>
				</div>
			</div>
		{/each}
	</div>
{/if}

<!-- Rejection Modal -->
{#if showRejectModal}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<div class="fixed inset-0 z-50">
		<div class="fixed inset-0 bg-black/50" onclick={() => (showRejectModal = false)}></div>
		<div class="fixed inset-0 flex items-center justify-center p-4">
			<div class="bg-white rounded-xl shadow-xl w-full max-w-md p-6" onclick={(e) => e.stopPropagation()}>
				<h3 class="text-lg font-semibold text-ocean-800 mb-4">Reject Request</h3>
				<div class="mb-4">
					<label for="rejectionReason" class="block text-sm font-medium text-ocean-700 mb-1">
						Reason (optional)
					</label>
					<textarea
						id="rejectionReason"
						bind:value={rejectionReason}
						class="w-full px-3 py-2 rounded-md border border-sand-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent resize-none"
						rows="3"
						placeholder="Enter reason for rejection..."
					></textarea>
				</div>
				<div class="flex gap-3">
					<Button variant="outline" class="flex-1" onclick={() => (showRejectModal = false)}>
						Cancel
					</Button>
					<Button
						variant="primary"
						class="flex-1 !bg-error hover:!bg-red-600"
						onclick={confirmReject}
						loading={processingId !== null}
					>
						Reject Request
					</Button>
				</div>
			</div>
		</div>
	</div>
{/if}
