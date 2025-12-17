<script lang="ts">
	import type { VacationRequest, VacationStatus } from '$lib/types';
	import { vacation } from '$lib/stores/vacation.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { formatDateMedium } from '$lib/utils/date';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Calendar, Trash2 } from 'lucide-svelte';

	interface Props {
		request: VacationRequest;
	}

	let { request }: Props = $props();

	let isCancelling = $state(false);

	// Map VacationStatus to Badge variant (semantic traffic light)
	const statusVariantMap: Record<VacationStatus, 'success' | 'warning' | 'error'> = {
		pending: 'warning',
		approved: 'success',
		rejected: 'error'
	};

	const badgeVariant = $derived(statusVariantMap[request.status]);

	async function handleCancel() {
		if (!confirm('Are you sure you want to cancel this request?')) return;

		isCancelling = true;
		try {
			await vacation.cancelRequest(request.id);
			toast.success('Request cancelled');
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to cancel');
		} finally {
			isCancelling = false;
		}
	}
</script>

<div class="flex items-center justify-between py-4">
	<div class="flex items-center gap-4">
		<div class="p-2 bg-ocean-100 rounded-lg">
			<Calendar class="w-5 h-5 text-ocean-600" />
		</div>
		<div>
			<p class="font-medium text-ocean-800">
				{formatDateMedium(request.startDate)} - {formatDateMedium(request.endDate)}
			</p>
			<p class="text-sm text-ocean-500">
				{request.totalDays} day{request.totalDays !== 1 ? 's' : ''}
				{#if request.reason}
					&middot; {request.reason}
				{/if}
			</p>
			{#if request.status === 'rejected' && request.rejectionReason}
				<p class="text-xs text-error mt-1">Reason: {request.rejectionReason}</p>
			{/if}
		</div>
	</div>

	<div class="flex items-center gap-3">
		<Badge variant={badgeVariant}>
			{request.status}
		</Badge>

		{#if request.status === 'pending'}
			<Button
				variant="ghost"
				size="sm"
				onclick={handleCancel}
				loading={isCancelling}
				aria-label="Cancel request"
				title="Cancel request"
			>
				<Trash2 class="w-4 h-4 text-error" />
			</Button>
		{/if}
	</div>
</div>
