<script lang="ts">
	import { createDialog, melt } from '@melt-ui/svelte';
	import type { VacationRequest } from '$lib/types';
	import { adminApi } from '$lib/api/admin';
	import { toast } from '$lib/stores/toast.svelte';
	import { formatDateShort } from '$lib/utils/date';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import EmptyState from '$lib/components/ui/EmptyState.svelte';
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

	// Approval Dialog
	const {
		elements: {
			overlay: approveOverlay,
			content: approveContent,
			title: approveTitle,
			description: approveDescription,
			close: approveClose,
			portalled: approvePortalled
		},
		states: { open: approveDialogOpen }
	} = createDialog({
		forceVisible: true,
		onOpenChange: ({ next }) => {
			if (!next) approvingRequest = null;
			return next;
		}
	});

	// Rejection Dialog
	const {
		elements: {
			overlay: rejectOverlay,
			content: rejectContent,
			title: rejectTitle,
			description: rejectDescription,
			close: rejectClose,
			portalled: rejectPortalled
		},
		states: { open: rejectDialogOpen }
	} = createDialog({
		forceVisible: true,
		onOpenChange: ({ next }) => {
			if (!next) {
				rejectionReason = '';
				rejectingRequest = null;
			}
			return next;
		}
	});

	function openApproveModal(request: VacationRequest) {
		approvingRequest = request;
		approveDialogOpen.set(true);
	}

	function openRejectModal(request: VacationRequest) {
		rejectingRequest = request;
		rejectDialogOpen.set(true);
	}

	async function confirmApprove() {
		if (!approvingRequest) return;

		processingId = approvingRequest.id;
		try {
			await adminApi.reviewRequest(approvingRequest.id, 'approved');
			toast.success('Request approved');
			approveDialogOpen.set(false);
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
			await adminApi.reviewRequest(rejectingRequest.id, 'rejected', rejectionReason || undefined);
			toast.success('Request rejected');
			rejectDialogOpen.set(false);
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
{#if $approveDialogOpen && approvingRequest}
	<div use:melt={$approvePortalled}>
		<div
			use:melt={$approveOverlay}
			class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm transition-opacity duration-200
				data-[state=open]:opacity-100 data-[state=closed]:opacity-0"
		></div>

		<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
			<div
				use:melt={$approveContent}
				class="bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30 w-full max-w-md
					transition-all duration-200 data-[state=open]:opacity-100 data-[state=open]:scale-100
					data-[state=closed]:opacity-0 data-[state=closed]:scale-95"
			>
				<!-- Header -->
				<div class="flex items-center justify-between p-4 border-b border-ocean-100/50">
					<div class="flex items-center gap-2">
						<CheckCircle class="w-5 h-5 text-green-600" />
						<h3 use:melt={$approveTitle} class="text-lg font-semibold text-ocean-800">
							Approve Request
						</h3>
					</div>
					<button
						use:melt={$approveClose}
						type="button"
						aria-label="Close dialog"
						class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600 hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
					>
						<X class="w-5 h-5" />
					</button>
				</div>

				<!-- Body -->
				<div class="p-4 space-y-4">
					<!-- Request Context -->
					<div
						use:melt={$approveDescription}
						class="flex items-center gap-3 p-3 bg-green-50 rounded-xl border border-green-100"
					>
						<Avatar name={approvingRequest.userName} size="md" />
						<div class="min-w-0 flex-1">
							<p class="font-medium text-ocean-800 truncate">{approvingRequest.userName}</p>
							<p class="text-sm text-ocean-500">
								{formatDateShort(approvingRequest.startDate)} – {formatDateShort(approvingRequest.endDate)}
								<span class="text-ocean-400">({approvingRequest.totalDays} days)</span>
							</p>
						</div>
					</div>

					<!-- Request Details -->
					<div class="bg-ocean-50/50 rounded-xl p-4 space-y-3">
						<div class="flex items-center gap-3">
							<Calendar class="w-4 h-4 text-ocean-400" />
							<div class="flex-1">
								<p class="text-xs text-ocean-500">Duration</p>
								<p class="text-sm font-medium text-ocean-800">
									{approvingRequest.totalDays} day{approvingRequest.totalDays !== 1 ? 's' : ''}
								</p>
							</div>
						</div>
						{#if approvingRequest.reason}
							<div class="flex items-start gap-3">
								<Umbrella class="w-4 h-4 text-ocean-400 mt-0.5" />
								<div class="flex-1">
									<p class="text-xs text-ocean-500">Reason</p>
									<p class="text-sm text-ocean-700">{approvingRequest.reason}</p>
								</div>
							</div>
						{/if}
					</div>

					<!-- Confirmation Message -->
					<p class="text-sm text-ocean-600">
						This will deduct <span class="font-semibold">{approvingRequest.totalDays} days</span>
						from {approvingRequest.userName?.split(' ')[0]}'s vacation balance.
					</p>

					<!-- Actions -->
					<div class="flex gap-3 pt-2">
						<Button variant="outline" class="flex-1" onclick={() => approveDialogOpen.set(false)}>
							Cancel
						</Button>
						<Button
							variant="success"
							class="flex-1"
							onclick={confirmApprove}
							loading={processingId !== null}
						>
							<CheckCircle class="w-4 h-4 mr-1.5" />
							Approve Request
						</Button>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- Rejection Dialog -->
{#if $rejectDialogOpen && rejectingRequest}
	<div use:melt={$rejectPortalled}>
		<div
			use:melt={$rejectOverlay}
			class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm transition-opacity duration-200
				data-[state=open]:opacity-100 data-[state=closed]:opacity-0"
		></div>

		<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
			<div
				use:melt={$rejectContent}
				class="bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30 w-full max-w-md
					transition-all duration-200 data-[state=open]:opacity-100 data-[state=open]:scale-100
					data-[state=closed]:opacity-0 data-[state=closed]:scale-95"
			>
				<!-- Header -->
				<div class="flex items-center justify-between p-4 border-b border-ocean-100/50">
					<div class="flex items-center gap-2">
						<XCircle class="w-5 h-5 text-error" />
						<h3 use:melt={$rejectTitle} class="text-lg font-semibold text-ocean-800">
							Reject Request
						</h3>
					</div>
					<button
						use:melt={$rejectClose}
						type="button"
						aria-label="Close dialog"
						class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600 hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
					>
						<X class="w-5 h-5" />
					</button>
				</div>

				<!-- Body -->
				<div class="p-4 space-y-4">
					<!-- Request Context -->
					<div
						use:melt={$rejectDescription}
						class="flex items-center gap-3 p-3 bg-sand-100 rounded-xl"
					>
						<Avatar name={rejectingRequest.userName} size="md" />
						<div class="min-w-0 flex-1">
							<p class="font-medium text-ocean-800 truncate">{rejectingRequest.userName}</p>
							<p class="text-sm text-ocean-500">
								{formatDateShort(rejectingRequest.startDate)} – {formatDateShort(rejectingRequest.endDate)}
								<span class="text-ocean-400">({rejectingRequest.totalDays} days)</span>
							</p>
						</div>
					</div>

					<!-- Reason Input -->
					<div>
						<label for="rejectionReason" class="block text-sm font-semibold text-ocean-800 mb-1.5">
							Reason (optional)
						</label>
						<textarea
							id="rejectionReason"
							bind:value={rejectionReason}
							class="w-full px-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 placeholder-ocean-500/50 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500 hover:border-ocean-500 resize-none"
							rows="3"
							placeholder="Let the crew member know why their request was rejected..."
						></textarea>
						<p class="mt-1.5 text-xs text-ocean-400">
							This message will be included in the notification email.
						</p>
					</div>

					<!-- Actions -->
					<div class="flex gap-3 pt-2">
						<Button variant="outline" class="flex-1" onclick={() => rejectDialogOpen.set(false)}>
							Cancel
						</Button>
						<Button
							variant="danger"
							class="flex-1"
							onclick={confirmReject}
							loading={processingId !== null}
						>
							<XCircle class="w-4 h-4 mr-1.5" />
							Reject Request
						</Button>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
