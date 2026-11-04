<script lang="ts">
	import type { VacationRequest } from '#lib/types/index.js';
	import RequestCard from './RequestCard.svelte';

	interface Props {
		requests: VacationRequest[];
	}

	let { requests }: Props = $props();
</script>

<div class="divide-y divide-sand-200">
	{#each requests as request (request.id)}
		<RequestCard {request} />
	{/each}
</div>
