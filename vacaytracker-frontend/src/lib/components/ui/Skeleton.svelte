<script lang="ts">
	import { clsx } from 'clsx';

	interface Props {
		variant?: 'text' | 'card' | 'avatar' | 'table-row' | 'circle';
		width?: string;
		height?: string;
		class?: string;
		count?: number;
	}

	let {
		variant = 'text',
		width,
		height,
		class: className = '',
		count = 1
	}: Props = $props();

	const baseStyles = 'animate-pulse bg-sand-200 rounded';

	const variantStyles = {
		text: 'h-4 rounded',
		card: 'h-24 rounded-lg',
		avatar: 'w-10 h-10 rounded-full',
		circle: 'rounded-full',
		'table-row': 'h-12 rounded'
	};

	const getStyles = (index: number) => {
		const style: Record<string, string> = {};

		if (width) {
			style.width = width;
		} else if (variant === 'text' && count > 1 && index === count - 1) {
			// Last text line is shorter for natural look
			style.width = '60%';
		}

		if (height) {
			style.height = height;
		}

		return style;
	};
</script>

{#each Array(count) as _, i}
	<div
		class={clsx(baseStyles, variantStyles[variant], className)}
		style={Object.entries(getStyles(i)).map(([k, v]) => `${k}: ${v}`).join('; ')}
	></div>
{/each}
