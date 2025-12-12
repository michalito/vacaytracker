<script lang="ts">
	import { createPagination, melt } from '@melt-ui/svelte';
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';

	interface Props {
		page: number;
		totalPages: number;
		onPageChange: (page: number) => void;
	}

	let { page, totalPages, onPageChange }: Props = $props();

	// Create Melt-UI pagination - note: count/perPage are reactive via derived
	const {
		elements: { root, pageTrigger, prevButton, nextButton },
		states: { pages, page: currentPage },
		options: { count }
	} = createPagination({
		count: 1,
		perPage: 1,
		defaultPage: 1,
		siblingCount: 2,
		onPageChange: ({ next }) => {
			onPageChange(next);
			return next;
		}
	});

	// Sync external totalPages prop with count option
	$effect(() => {
		count.set(totalPages);
	});

	// Sync external page prop with internal state
	$effect(() => {
		if (page !== $currentPage) {
			currentPage.set(page);
		}
	});
</script>

{#if totalPages > 1}
	<nav use:melt={$root} class="flex items-center justify-center gap-1">
		<button
			use:melt={$prevButton}
			class="p-2 rounded-lg text-ocean-600 hover:bg-sand-100 transition-colors cursor-pointer
				disabled:opacity-50 disabled:cursor-not-allowed"
		>
			<ChevronLeft class="w-4 h-4" />
		</button>

		{#each $pages as p}
			{#if p.type === 'ellipsis'}
				<span class="px-2 text-ocean-400">...</span>
			{:else}
				<button
					use:melt={$pageTrigger(p)}
					class="w-8 h-8 rounded-md text-sm font-medium transition-colors cursor-pointer
						text-ocean-700 hover:bg-sand-100
						data-[selected]:bg-ocean-500 data-[selected]:text-white data-[selected]:hover:bg-ocean-600"
				>
					{p.value}
				</button>
			{/if}
		{/each}

		<button
			use:melt={$nextButton}
			class="p-2 rounded-lg text-ocean-600 hover:bg-sand-100 transition-colors cursor-pointer
				disabled:opacity-50 disabled:cursor-not-allowed"
		>
			<ChevronRight class="w-4 h-4" />
		</button>
	</nav>
{/if}
