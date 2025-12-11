<script lang="ts">
	import Skeleton from './Skeleton.svelte';

	interface Props {
		count?: number;
	}

	let { count = 4 }: Props = $props();
</script>

{#each Array(count) as _}
	<div class="bg-white rounded-lg shadow-sm p-6 border border-sand-200">
		<div class="flex items-center justify-between">
			<div class="space-y-2 flex-1">
				<Skeleton variant="text" width="60%" />
				<Skeleton variant="text" width="40%" height="2rem" />
			</div>
			<Skeleton variant="circle" width="48px" height="48px" />
		</div>
	</div>
{/each}
