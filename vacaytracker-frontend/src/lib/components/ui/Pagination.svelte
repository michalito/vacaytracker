<script lang="ts">
	import Button from './Button.svelte';
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';

	interface Props {
		page: number;
		totalPages: number;
		onPageChange: (page: number) => void;
	}

	let { page, totalPages, onPageChange }: Props = $props();

	const visiblePages = $derived(() => {
		const pages: (number | 'ellipsis')[] = [];
		const delta = 2;

		for (let i = 1; i <= totalPages; i++) {
			if (i === 1 || i === totalPages || (i >= page - delta && i <= page + delta)) {
				pages.push(i);
			} else if (pages[pages.length - 1] !== 'ellipsis') {
				pages.push('ellipsis');
			}
		}

		return pages;
	});
</script>

{#if totalPages > 1}
	<div class="flex items-center justify-center gap-1">
		<Button variant="ghost" size="sm" onclick={() => onPageChange(page - 1)} disabled={page <= 1}>
			<ChevronLeft class="w-4 h-4" />
		</Button>

		{#each visiblePages() as p, i}
			{#if p === 'ellipsis'}
				<span class="px-2 text-ocean-400">...</span>
			{:else}
				<button
					type="button"
					onclick={() => onPageChange(p)}
					class="w-8 h-8 rounded-md text-sm font-medium transition-colors {p === page
						? 'bg-ocean-500 text-white'
						: 'text-ocean-700 hover:bg-sand-100'}"
				>
					{p}
				</button>
			{/if}
		{/each}

		<Button
			variant="ghost"
			size="sm"
			onclick={() => onPageChange(page + 1)}
			disabled={page >= totalPages}
		>
			<ChevronRight class="w-4 h-4" />
		</Button>
	</div>
{/if}
