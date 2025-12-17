<script lang="ts">
	type Theme = 'healthy' | 'warning' | 'critical';

	interface Segment {
		value: number;
		label: string;
	}

	interface Props {
		segments: Segment[]; // [used, upcoming, available] - order matters for shading
		theme?: Theme;
		size?: number;
		strokeWidth?: number;
		centerLabel?: string;
		centerSubLabel?: string;
		class?: string;
	}

	let {
		segments,
		theme = 'healthy',
		size = 140,
		strokeWidth = 10,
		centerLabel = '',
		centerSubLabel = '',
		class: className = ''
	}: Props = $props();

	// Spa/Resort palette: [available (dynamic), upcoming (caribbean), used (stone)]
	// Segment order matches legend: Available → Upcoming → Used (clockwise from 12 o'clock)
	const themeShades: Record<Theme, [string, string, string]> = {
		healthy: ['text-mint-400', 'text-caribbean-400', 'text-stone-300'],      // >50% - Soft Mint
		warning: ['text-sunshine-400', 'text-caribbean-400', 'text-stone-300'],  // 20-50% - Sunshine
		critical: ['text-salmon-400', 'text-caribbean-400', 'text-stone-300']    // <20% - Soft Coral
	};

	const radius = $derived((size - strokeWidth) / 2);
	const circumference = $derived(2 * Math.PI * radius);
	const total = $derived(segments.reduce((sum, s) => sum + s.value, 0));
	const shades = $derived(themeShades[theme]);

	// Calculate segment arcs with offsets
	const segmentArcs = $derived.by(() => {
		if (total === 0) return [];

		const arcs: Array<{
			dashArray: string;
			dashOffset: number;
			shade: string;
			label: string;
		}> = [];

		let cumulativeOffset = 0;

		segments.forEach((segment, index) => {
			if (segment.value === 0) return;

			const segmentLength = (segment.value / total) * circumference;

			arcs.push({
				dashArray: `${segmentLength} ${circumference - segmentLength}`,
				dashOffset: -cumulativeOffset,
				shade: shades[index] || shades[0],
				label: segment.label
			});

			cumulativeOffset += segmentLength;
		});

		return arcs;
	});

	// Center label color based on theme (matches "Available" segment)
	const centerColor = $derived(
		theme === 'healthy'
			? 'text-mint-500'
			: theme === 'warning'
				? 'text-sunshine-500'
				: 'text-salmon-500'
	);
</script>

<div class="relative inline-flex items-center justify-center {className}">
	<svg width={size} height={size} class="-rotate-90">
		<!-- Background circle -->
		<circle
			cx={size / 2}
			cy={size / 2}
			r={radius}
			fill="none"
			stroke="currentColor"
			stroke-width={strokeWidth}
			class="text-sand-200"
		/>

		<!-- Segment circles (drawn in reverse order so first segment is on top) -->
		{#each segmentArcs.toReversed() as arc}
			<circle
				cx={size / 2}
				cy={size / 2}
				r={radius}
				fill="none"
				stroke="currentColor"
				stroke-width={strokeWidth}
				stroke-dasharray={arc.dashArray}
				stroke-dashoffset={arc.dashOffset}
				class={arc.shade}
				style="transition: stroke-dasharray 0.5s ease, stroke-dashoffset 0.5s ease, stroke 0.3s ease"
			/>
		{/each}
	</svg>

	<!-- Center content -->
	<div class="absolute inset-0 flex flex-col items-center justify-center">
		{#if centerLabel}
			<span class="text-2xl font-bold {centerColor}">{centerLabel}</span>
		{/if}
		{#if centerSubLabel}
			<span class="text-xs text-ocean-500">{centerSubLabel}</span>
		{/if}
	</div>
</div>
