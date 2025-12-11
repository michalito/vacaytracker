<script lang="ts">
	import Skeleton from './Skeleton.svelte';

	interface Props {
		count?: number;
		variant?: 'simple' | 'withActions' | 'compact';
	}

	let { count = 3, variant = 'simple' }: Props = $props();
</script>

<div class="space-y-{variant === 'compact' ? '3' : '4'} py-{variant === 'compact' ? '2' : '4'}">
	{#each Array(count) as _}
		{#if variant === 'compact'}
			<div class="flex items-center gap-3 p-2">
				<Skeleton variant="avatar" />
				<div class="flex-1 space-y-2">
					<Skeleton variant="text" width="30%" />
					<Skeleton variant="text" width="50%" />
				</div>
			</div>
		{:else if variant === 'withActions'}
			<div class="flex items-center gap-4 p-4 border border-sand-200 rounded-lg">
				<Skeleton variant="avatar" />
				<div class="flex-1 space-y-2">
					<Skeleton variant="text" width="35%" />
					<Skeleton variant="text" width="55%" />
				</div>
				<div class="flex gap-2">
					<Skeleton variant="text" width="70px" height="32px" />
					<Skeleton variant="text" width="70px" height="32px" />
				</div>
			</div>
		{:else}
			<div class="flex items-center gap-4 p-4 border border-sand-200 rounded-lg">
				<Skeleton variant="avatar" />
				<div class="flex-1 space-y-2">
					<Skeleton variant="text" width="40%" />
					<Skeleton variant="text" width="60%" />
				</div>
				<Skeleton variant="text" width="80px" />
			</div>
		{/if}
	{/each}
</div>
