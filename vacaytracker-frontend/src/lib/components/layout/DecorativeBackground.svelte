<script lang="ts">
	interface Props {
		variant?: 'subtle' | 'standard' | 'prominent';
	}

	let { variant = 'standard' }: Props = $props();
</script>

<div class="fixed inset-0 overflow-hidden pointer-events-none z-0" aria-hidden="true">
	{#if variant === 'subtle'}
		<!-- Subtle variant - minimal decorations -->
		<div
			class="absolute -top-32 -right-32 w-64 h-64 bg-ocean-300/10 rounded-full blur-3xl animate-float"
		></div>
		<div
			class="absolute bottom-1/4 -left-16 w-48 h-48 bg-ocean-400/8 rounded-full blur-2xl animate-float"
			style="animation-delay: -3s;"
		></div>
	{:else if variant === 'standard'}
		<!-- Standard variant - balanced decorations -->
		<div
			class="absolute -top-20 -left-20 w-72 h-72 bg-ocean-300/15 rounded-full blur-3xl animate-float"
		></div>
		<div
			class="absolute -top-10 right-1/4 w-56 h-56 bg-coral-400/8 rounded-full blur-2xl animate-float"
			style="animation-delay: -2s;"
		></div>
		<div
			class="absolute bottom-1/3 -right-16 w-48 h-48 bg-ocean-400/10 rounded-full blur-2xl animate-float"
			style="animation-delay: -4s;"
		></div>
		<div
			class="absolute bottom-1/4 left-1/4 w-40 h-40 bg-ocean-200/12 rounded-full blur-2xl animate-float"
			style="animation-delay: -5s;"
		></div>
	{:else}
		<!-- Prominent variant - more visible decorations -->
		<div
			class="absolute -top-24 -left-24 w-80 h-80 bg-ocean-300/20 rounded-full blur-3xl animate-float"
		></div>
		<div
			class="absolute -top-12 right-1/3 w-64 h-64 bg-coral-400/12 rounded-full blur-2xl animate-float"
			style="animation-delay: -2s;"
		></div>
		<div
			class="absolute top-1/2 -right-20 w-56 h-56 bg-ocean-400/15 rounded-full blur-2xl animate-float"
			style="animation-delay: -3s;"
		></div>
		<div
			class="absolute bottom-1/4 left-1/3 w-48 h-48 bg-ocean-200/15 rounded-full blur-2xl animate-float"
			style="animation-delay: -4s;"
		></div>
		<div
			class="absolute bottom-1/3 -left-12 w-40 h-40 bg-coral-400/8 rounded-full blur-3xl animate-float"
			style="animation-delay: -5s;"
		></div>
	{/if}
</div>
