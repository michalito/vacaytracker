<script lang="ts">
	import type { VacationRequest } from '$lib/types';
	import { vacation } from '$lib/stores/vacation.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import { formatDateMedium, getDaysUntilStart } from '$lib/utils/date';
	import StatusBadge from '$lib/components/ui/StatusBadge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Calendar, Trash2, AlertCircle, Hourglass } from 'lucide-svelte';

	// ===================
	// Types
	// ===================
	type CardVariant = 'pending' | 'upcoming' | 'past';

	// ===================
	// Props
	// ===================
	interface Props {
		request: VacationRequest;
		variant: CardVariant;
	}

	const { request, variant }: Props = $props();

	// ===================
	// State
	// ===================
	let isCancelling = $state(false);

	// ===================
	// Styling Configuration
	// ===================
	const BORDER_COLORS: Record<CardVariant, string> = {
		pending: 'border-l-coral-400',
		upcoming: 'border-l-ocean-400',
		past: 'border-l-sand-300'
	};

	const CARD_BACKGROUNDS: Record<CardVariant, string> = {
		pending: 'bg-white',
		upcoming: 'bg-white',
		past: 'bg-sand-50/50'
	};

	// ===================
	// Derived Values
	// ===================
	const isPast = $derived(variant === 'past');
	const isRejected = $derived(request.status === 'rejected');
	const isPending = $derived(request.status === 'pending');

	// Map variant + request status to StatusBadge variant
	const statusBadgeVariant = $derived.by<'pending' | 'upcoming' | 'completed' | 'rejected'>(() => {
		if (request.status === 'rejected') return 'rejected';
		if (request.status === 'pending') return 'pending';
		if (variant === 'upcoming') return 'upcoming';
		return 'completed';
	});

	// Countdown text for upcoming requests
	const countdown = $derived.by<string | null>(() => {
		if (variant !== 'upcoming') return null;

		const days = getDaysUntilStart(request.startDate);
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
	});

	// Date visualization dots (max 14)
	const dateDots = $derived.by(() => {
		const count = Math.min(request.totalDays, 14);
		return Array.from({ length: count }, (_, i) => ({
			opacity: isPast || isRejected ? 0.5 : 0.4 + (i / count) * 0.6
		}));
	});

	const dotColor = $derived(isRejected ? 'bg-red-200' : isPast ? 'bg-sand-300' : 'bg-ocean-400');

	// ===================
	// Event Handlers
	// ===================
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

<article
	class="relative rounded-xl border-l-4 {BORDER_COLORS[variant]} {CARD_BACKGROUNDS[variant]} p-4 shadow-sm transition-all duration-200 hover:shadow-md hover:-translate-y-0.5"
	class:opacity-75={isPast}
>
	<!-- Header: Date range and status -->
	<header class="flex items-start justify-between gap-4">
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
				<span
					class="inline-flex items-center gap-1.5 px-2 py-1 text-xs font-medium rounded-md bg-ocean-100 text-ocean-700 border border-ocean-200"
				>
					<Hourglass class="w-3.5 h-3.5 text-ocean-500" />
					{countdown}
				</span>
			{/if}
			<StatusBadge status={statusBadgeVariant} size="sm" />
		</div>
	</header>

	<!-- Date range visualization -->
	<div class="mt-3 flex items-center gap-1" aria-hidden="true">
		{#each dateDots as dot}
			<div class="w-2 h-2 rounded-full {dotColor}" style="opacity: {dot.opacity}"></div>
		{/each}
		{#if request.totalDays > 14}
			<span class="text-xs text-ocean-400 ml-1">+{request.totalDays - 14} more</span>
		{/if}
	</div>

	<!-- Reason -->
	{#if request.reason}
		<p class="mt-3 text-sm text-ocean-600" class:italic={isPast}>
			"{request.reason}"
		</p>
	{/if}

	<!-- Rejection reason -->
	{#if isRejected && request.rejectionReason}
		<div class="mt-3 flex items-start gap-2 p-2 bg-red-50 rounded-lg">
			<AlertCircle class="w-4 h-4 text-error mt-0.5 shrink-0" />
			<p class="text-xs text-red-700">
				<span class="font-medium">Rejected:</span>
				{request.rejectionReason}
			</p>
		</div>
	{/if}

	<!-- Actions -->
	{#if isPending}
		<footer class="mt-4 flex justify-end">
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
		</footer>
	{/if}
</article>
