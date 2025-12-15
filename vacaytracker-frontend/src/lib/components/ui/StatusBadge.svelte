<script lang="ts">
	import { clsx } from 'clsx';
	import { Clock, Sun, Palmtree, CloudRain } from 'lucide-svelte';
	import type { VacationStatus } from '$lib/types';

	type StatusVariant = VacationStatus | 'upcoming' | 'completed';

	interface Props {
		status: StatusVariant;
		size?: 'sm' | 'md';
		showLabel?: boolean;
		class?: string;
	}

	let { status, size = 'sm', showLabel = true, class: className = '' }: Props = $props();

	// Status configuration using app's ocean/coral/sand palette
	const statusConfig: Record<
		StatusVariant,
		{
			icon: typeof Clock;
			label: string;
			bgClass: string;
			textClass: string;
			iconClass: string;
		}
	> = {
		pending: {
			icon: Clock,
			label: 'Awaiting',
			bgClass: 'bg-coral-400/15 border-coral-400/30',
			textClass: 'text-coral-600',
			iconClass: 'text-coral-500'
		},
		approved: {
			icon: Sun,
			label: 'Confirmed',
			bgClass: 'bg-ocean-100 border-ocean-200',
			textClass: 'text-ocean-700',
			iconClass: 'text-ocean-500'
		},
		upcoming: {
			icon: Sun,
			label: 'Confirmed',
			bgClass: 'bg-ocean-100 border-ocean-200',
			textClass: 'text-ocean-700',
			iconClass: 'text-ocean-500'
		},
		completed: {
			icon: Palmtree,
			label: 'Enjoyed',
			bgClass: 'bg-sand-200/60 border-sand-300',
			textClass: 'text-sand-500',
			iconClass: 'text-sand-400'
		},
		rejected: {
			icon: CloudRain,
			label: 'Declined',
			bgClass: 'bg-red-50 border-red-200/60',
			textClass: 'text-red-500',
			iconClass: 'text-red-400'
		}
	};

	const config = $derived(statusConfig[status]);

	const sizeStyles = {
		sm: {
			container: 'px-2 py-1 gap-1.5',
			icon: 'w-3.5 h-3.5',
			text: 'text-xs'
		},
		md: {
			container: 'px-2.5 py-1.5 gap-2',
			icon: 'w-4 h-4',
			text: 'text-sm'
		}
	};

	const sizeConfig = $derived(sizeStyles[size]);
</script>

<span
	class={clsx(
		'inline-flex items-center font-medium rounded-md border transition-colors',
		config.bgClass,
		sizeConfig.container,
		className
	)}
>
	<config.icon class={clsx(sizeConfig.icon, config.iconClass)} />
	{#if showLabel}
		<span class={clsx(sizeConfig.text, config.textClass)}>{config.label}</span>
	{/if}
</span>
