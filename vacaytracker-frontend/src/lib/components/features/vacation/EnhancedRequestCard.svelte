<script lang="ts">
	import type { VacationRequest } from '$lib/types';
	import { vacation } from '$lib/stores/vacation.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { formatDateMedium, getDaysUntilStart } from '$lib/utils/date';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Calendar, Trash2, AlertCircle } from 'lucide-svelte';

	type CardVariant = 'pending' | 'upcoming' | 'past';

	interface Props {
		request: VacationRequest;
		variant: CardVariant;
	}

	let { request, variant }: Props = $props();

	let isCancelling = $state(false);

	// Border colors based on variant
	const borderColors: Record<CardVariant, string> = {
		pending: 'border-l-coral-400',
		upcoming: 'border-l-ocean-400',
		past: 'border-l-sand-400'
	};

	// Card background styling based on variant
	const cardStyles: Record<CardVariant, string> = {
		pending: 'bg-white',
		upcoming: 'bg-white',
		past: 'bg-sand-50/50'
	};

	// Get countdown text for upcoming requests
	function getCountdown(startDate: string): string | null {
		const days = getDaysUntilStart(startDate);

		if (days < 0) return null;
		if (days === 0) return 'Starts today';
		if (days === 1) return 'Starts tomorrow';
		if (days <= 7) return `in ${days} days`;
		if (days <= 30) {
			const weeks = Math.ceil(days / 7);
			return `in ${weeks} week${weeks > 1 ? 's' : ''}`;
		}
		const months = Math.ceil(days / 30);
		return `in ${months} month${months > 1 ? 's' : ''}`;
	}

	// Generate dots for date range visualization (max 14 dots)
	function getDateDots(totalDays: number): { filled: boolean }[] {
		const maxDots = Math.min(totalDays, 14);
		const dots: { filled: boolean }[] = [];

		for (let i = 0; i < maxDots; i++) {
			dots.push({ filled: true });
		}

		return dots;
	}

	const countdown = $derived(variant === 'upcoming' ? getCountdown(request.startDate) : null);
	const dateDots = $derived(getDateDots(request.totalDays));
	const isPast = $derived(variant === 'past');
	const isRejected = $derived(request.status === 'rejected');

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

<div
	class="relative rounded-xl border-l-4 {borderColors[variant]} {cardStyles[variant]} p-4 shadow-sm hover:shadow-md transition-all duration-200 hover:-translate-y-0.5 {isPast ? 'opacity-75' : ''}"
>
	<!-- Header: Date range and status badge -->
	<div class="flex items-start justify-between gap-4">
		<div class="flex items-center gap-3">
			<div class="p-2 bg-ocean-100/80 rounded-lg">
				<Calendar class="w-5 h-5 text-ocean-600" />
			</div>
			<div>
				<p class="font-semibold text-ocean-800">
					{formatDateMedium(request.startDate)} - {formatDateMedium(request.endDate)}
				</p>
				<p class="text-sm text-ocean-500">
					{request.totalDays} day{request.totalDays !== 1 ? 's' : ''}
				</p>
			</div>
		</div>

		<div class="flex items-center gap-2">
			{#if countdown}
				<span class="px-2.5 py-1 text-xs font-medium rounded-md bg-ocean-100 text-ocean-700">
					{countdown}
				</span>
			{/if}
			<Badge variant={request.status} size="sm">
				{request.status}
			</Badge>
		</div>
	</div>

	<!-- Date range visualization -->
	<div class="mt-3 flex items-center gap-1">
		{#each dateDots as dot, i}
			<div
				class="w-2 h-2 rounded-full {isRejected ? 'bg-red-200' : isPast ? 'bg-sand-300' : 'bg-ocean-400'}"
				style="opacity: {isPast || isRejected ? 0.5 : 0.4 + (i / dateDots.length) * 0.6}"
			></div>
		{/each}
		{#if request.totalDays > 14}
			<span class="text-xs text-ocean-400 ml-1">+{request.totalDays - 14} more</span>
		{/if}
	</div>

	<!-- Reason -->
	{#if request.reason}
		<p class="mt-3 text-sm text-ocean-600 {isPast ? 'italic' : ''}">
			"{request.reason}"
		</p>
	{/if}

	<!-- Rejection reason (for rejected requests) -->
	{#if isRejected && request.rejectionReason}
		<div class="mt-3 flex items-start gap-2 p-2 bg-red-50 rounded-lg">
			<AlertCircle class="w-4 h-4 text-error mt-0.5 shrink-0" />
			<p class="text-xs text-red-700">
				<span class="font-medium">Rejected:</span>
				{request.rejectionReason}
			</p>
		</div>
	{/if}

	<!-- Actions (cancel button for pending) -->
	{#if request.status === 'pending'}
		<div class="mt-4 flex justify-end">
			<Button
				variant="ghost"
				size="sm"
				onclick={handleCancel}
				loading={isCancelling}
				class="text-error hover:bg-red-50"
			>
				<Trash2 class="w-4 h-4 mr-1.5" />
				Cancel Request
			</Button>
		</div>
	{/if}
</div>
