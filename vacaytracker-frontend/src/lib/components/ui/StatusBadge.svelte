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

	// Status configuration using semantic traffic light colors
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
			bgClass: 'bg-warning-light border-yellow-200',
			textClass: 'text-yellow-800',
			iconClass: 'text-yellow-600'
		},
		approved: {
			icon: Sun,
			label: 'Confirmed',
			bgClass: 'bg-success-light border-green-200',
			textClass: 'text-green-800',
			iconClass: 'text-green-600'
		},
		upcoming: {
			icon: Sun,
			label: 'Confirmed',
			bgClass: 'bg-success-light border-green-200',
			textClass: 'text-green-800',
			iconClass: 'text-green-600'
		},
		completed: {
			icon: Palmtree,
			label: 'Enjoyed',
			bgClass: 'bg-success-light border-green-200',
			textClass: 'text-green-800',
			iconClass: 'text-green-600'
		},
		rejected: {
			icon: CloudRain,
			label: 'Declined',
			bgClass: 'bg-error-light border-red-200',
			textClass: 'text-red-800',
			iconClass: 'text-red-600'
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
