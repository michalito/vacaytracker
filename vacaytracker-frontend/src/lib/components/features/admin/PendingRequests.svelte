<script lang="ts">
	import type { VacationRequest } from '$lib/types';
	import { adminApi } from '$lib/api/admin';
	import { toast } from '$lib/stores/toast.svelte';
	import { formatDateShort } from '$lib/utils/date';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import { Check, X, Clock, XCircle, CheckCircle, Calendar, Umbrella } from 'lucide-svelte';

	interface Props {
		requests: VacationRequest[];
		onUpdate: () => void;
	}

	let { requests, onUpdate }: Props = $props();

	let processingId = $state<string | null>(null);
	let rejectionReason = $state('');
	let rejectingRequest = $state<VacationRequest | null>(null);
	let approvingRequest = $state<VacationRequest | null>(null);

	let approveOpen = $state(false);
	let rejectOpen = $state(false);

	function openApproveModal(request: VacationRequest) {
		approvingRequest = request;
		approveOpen = true;
	}

	function openRejectModal(request: VacationRequest) {
		// Each rejection starts with an empty reason; it is emailed to the employee
		rejectionReason = '';
		rejectingRequest = request;
		rejectOpen = true;
	}

	async function confirmApprove() {
		if (!approvingRequest) return;

		processingId = approvingRequest.id;
		try {
			await adminApi.reviewRequest(approvingRequest.id, 'approved');
			toast.success('Request approved');
			approveOpen = false;
			onUpdate();
		} catch (error) {
			toast.error('Failed to approve request');
		} finally {
			processingId = null;
		}
	}

	async function confirmReject() {
		if (!rejectingRequest) return;

		processingId = rejectingRequest.id;
		try {
			await adminApi.reviewRequest(
				rejectingRequest.id,
				'rejected',
				rejectionReason.trim() || undefined
			);
			toast.success('Request rejected');
			rejectOpen = false;
			onUpdate();
		} catch (error) {
			toast.error('Failed to reject request');
		} finally {
			processingId = null;
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
						onclick={() => openRejectModal(request)}
						disabled={processingId !== null}
					>
						<X class="w-4 h-4 mr-1" />
						Reject
					</Button>
					<Button
						variant="primary"
						size="sm"
						onclick={() => openApproveModal(request)}
						disabled={processingId !== null}
					>
						<Check class="w-4 h-4 mr-1" />
						Approve
					</Button>
				</div>
			</div>
		{/each}
	</div>
{/if}

<!-- Approval Dialog -->
<Modal
	bind:open={approveOpen}
	title="Approve Request"
	icon={CheckCircle}
	iconClass="text-green-600"
	closeLabel="Close approval dialog"
	onClosed={() => (approvingRequest = null)}
>
	{#if approvingRequest}
		<div class="space-y-4 p-4">
			<!-- Request Context -->
			<div class="flex items-center gap-3 rounded-xl border border-green-100 bg-green-50 p-3">
				<Avatar name={approvingRequest.userName} size="md" />
				<div class="min-w-0 flex-1">
					<p class="truncate font-medium text-ocean-800">{approvingRequest.userName}</p>
					<p class="text-sm text-ocean-500">
						{formatDateShort(approvingRequest.startDate)} – {formatDateShort(
							approvingRequest.endDate
						)}
						<span class="text-ocean-400">({approvingRequest.totalDays} days)</span>
					</p>
				</div>
			</div>

			<!-- Request Details -->
			<div class="space-y-3 rounded-xl bg-ocean-50/50 p-4">
				<div class="flex items-center gap-3">
					<Calendar class="h-4 w-4 text-ocean-400" aria-hidden="true" />
					<div class="flex-1">
						<p class="text-xs text-ocean-500">Duration</p>
						<p class="text-sm font-medium text-ocean-800">
							{approvingRequest.totalDays} day{approvingRequest.totalDays !== 1 ? 's' : ''}
						</p>
					</div>
				</div>
				{#if approvingRequest.reason}
					<div class="flex items-start gap-3">
						<Umbrella class="mt-0.5 h-4 w-4 text-ocean-400" aria-hidden="true" />
						<div class="flex-1">
							<p class="text-xs text-ocean-500">Reason</p>
							<p class="text-sm text-ocean-700">{approvingRequest.reason}</p>
						</div>
					</div>
				{/if}
			</div>

			<p class="text-sm text-ocean-600">
				This will deduct <span class="font-semibold">{approvingRequest.totalDays} days</span>
				from {approvingRequest.userName ? `${approvingRequest.userName.split(' ')[0]}'s` : 'their'} vacation
				balance.
			</p>

			<div class="flex gap-3 pt-2">
				<Button variant="outline" class="flex-1" onclick={() => (approveOpen = false)}>
					Cancel
				</Button>
				<Button
					variant="success"
					class="flex-1"
					onclick={confirmApprove}
					loading={processingId !== null}
				>
					<CheckCircle class="mr-1.5 h-4 w-4" aria-hidden="true" />
					Approve Request
				</Button>
			</div>
		</div>
	{/if}
</Modal>

<!-- Rejection Dialog -->
<Modal
	bind:open={rejectOpen}
	title="Reject Request"
	icon={XCircle}
	iconClass="text-error"
	closeLabel="Close rejection dialog"
	onClosed={() => (rejectingRequest = null)}
>
	{#if rejectingRequest}
		<div class="space-y-4 p-4">
			<!-- Request Context -->
			<div class="flex items-center gap-3 rounded-xl bg-sand-100 p-3">
				<Avatar name={rejectingRequest.userName} size="md" />
				<div class="min-w-0 flex-1">
					<p class="truncate font-medium text-ocean-800">{rejectingRequest.userName}</p>
					<p class="text-sm text-ocean-500">
						{formatDateShort(rejectingRequest.startDate)} – {formatDateShort(
							rejectingRequest.endDate
						)}
						<span class="text-ocean-400">({rejectingRequest.totalDays} days)</span>
					</p>
				</div>
			</div>

			<div>
				<label for="rejection-reason" class="mb-1.5 block text-sm font-semibold text-ocean-800">
					Reason (optional)
				</label>
				<textarea
					id="rejection-reason"
					bind:value={rejectionReason}
					rows="3"
					placeholder="Let the crew member know why their request was rejected..."
					aria-describedby="rejection-reason-hint"
					class="w-full resize-none rounded-lg border-2 border-ocean-500/40 bg-white px-4 py-2.5 text-ocean-900
						placeholder-ocean-500/50 transition-all duration-200 hover:border-ocean-500
						focus:border-ocean-500 focus:outline-none focus:ring-2 focus:ring-ocean-500/30"
				></textarea>
				<p id="rejection-reason-hint" class="mt-1.5 text-xs text-ocean-400">
					This message will be included in the notification email.
				</p>
			</div>

			<div class="flex gap-3 pt-2">
				<Button variant="outline" class="flex-1" onclick={() => (rejectOpen = false)}>
					Cancel
				</Button>
				<Button
					variant="danger"
					class="flex-1"
					onclick={confirmReject}
					loading={processingId !== null}
				>
					<XCircle class="mr-1.5 h-4 w-4" aria-hidden="true" />
					Reject Request
				</Button>
			</div>
		</div>
	{/if}
</Modal>
