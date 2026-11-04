<script lang="ts">
	import { Pagination } from 'bits-ui';
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';

	interface Props {
		page: number;
		totalPages: number;
		onPageChange: (page: number) => void;
	}

	let { page, totalPages, onPageChange }: Props = $props();
</script>

{#if totalPages > 1}
	<Pagination.Root
		count={totalPages}
		perPage={1}
		{page}
		onPageChange={(p) => onPageChange(p)}
		siblingCount={2}
	>
		{#snippet children({ pages })}
			<nav class="flex items-center justify-center gap-1">
				<Pagination.PrevButton
					class="p-2 rounded-lg text-ocean-600 hover:bg-sand-100 transition-colors cursor-pointer
						disabled:opacity-50 disabled:cursor-not-allowed"
				>
					<ChevronLeft class="w-4 h-4" />
				</Pagination.PrevButton>

				{#each pages as p (p.key)}
					{#if p.type === 'ellipsis'}
						<span class="px-2 text-ocean-400">...</span>
					{:else}
						<Pagination.Page
							page={p}
							class="w-8 h-8 rounded-md text-sm font-medium transition-colors cursor-pointer
								text-ocean-700 hover:bg-sand-100
								data-[selected]:bg-ocean-500 data-[selected]:text-white data-[selected]:hover:bg-ocean-600"
						>
							{p.value}
						</Pagination.Page>
					{/if}
				{/each}

				<Pagination.NextButton
					class="p-2 rounded-lg text-ocean-600 hover:bg-sand-100 transition-colors cursor-pointer
						disabled:opacity-50 disabled:cursor-not-allowed"
				>
					<ChevronRight class="w-4 h-4" />
				</Pagination.NextButton>
			</nav>
		{/snippet}
	</Pagination.Root>
{/if}
