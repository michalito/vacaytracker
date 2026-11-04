<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { clsx } from 'clsx';
	import { X } from 'lucide-svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		/** Whether the modal is open. Bindable. */
		open: boolean;
		title: string;
		/** Optional icon shown before the title. */
		icon?: typeof X;
		iconClass?: string;
		/** Supporting text under the title; announced as the dialog description. */
		description?: string | Snippet;
		size?: 'md' | 'xl';
		/** Accessible label for the icon-only close button. */
		closeLabel?: string;
		/**
		 * Runs once the modal has fully closed, however it was closed
		 * (close button, Escape, overlay click, or `open = false` from code).
		 */
		onClosed?: () => void;
		class?: string;
		children: Snippet;
		/** Optional actions pinned below the scrollable body. */
		footer?: Snippet;
	}

	let {
		open = $bindable(false),
		title,
		icon: Icon,
		iconClass = 'text-ocean-500',
		description,
		size = 'md',
		closeLabel = 'Close dialog',
		onClosed,
		class: className = '',
		children,
		footer
	}: Props = $props();

	const sizeClasses = {
		md: 'max-w-md',
		xl: 'max-w-4xl'
	};
</script>

<Dialog.Root
	bind:open
	onOpenChangeComplete={(isOpen) => {
		if (!isOpen) onClosed?.();
	}}
>
	<Dialog.Portal>
		<Dialog.Overlay class="presence-fade fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm" />

		<Dialog.Content
			class={clsx(
				'presence-pop fixed left-1/2 top-1/2 z-50 -translate-x-1/2 -translate-y-1/2',
				'flex w-[calc(100%-2rem)] max-h-[calc(100dvh-2rem)] flex-col overflow-hidden',
				'rounded-2xl border border-white/40 bg-white/95 shadow-xl backdrop-blur-md focus:outline-none',
				sizeClasses[size],
				className
			)}
		>
			<!-- Header (pr-12 leaves room for the close button) -->
			<div class="shrink-0 border-b border-ocean-100/50 p-4 pr-12">
				<div class="flex items-center gap-2">
					{#if Icon}
						<Icon class={clsx('h-5 w-5 shrink-0', iconClass)} aria-hidden="true" />
					{/if}
					<Dialog.Title class="text-lg font-semibold text-ocean-800">{title}</Dialog.Title>
				</div>
				{#if description}
					<Dialog.Description class="mt-1 text-sm text-ocean-600">
						{#if typeof description === 'string'}
							{description}
						{:else}
							{@render description()}
						{/if}
					</Dialog.Description>
				{/if}
			</div>

			<!-- Body -->
			<div class="min-h-0 flex-1 overflow-y-auto">
				{@render children()}
			</div>

			{#if footer}
				<div class="flex shrink-0 justify-end gap-2 border-t border-ocean-100/50 p-4">
					{@render footer()}
				</div>
			{/if}

			<!-- Close button is last in DOM order so initial focus lands on the first field -->
			<Dialog.Close
				aria-label={closeLabel}
				class="absolute right-3 top-3 cursor-pointer rounded-lg p-1 text-ocean-400 transition-all duration-200
					hover:bg-ocean-500/10 hover:text-ocean-600
					focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ocean-500"
			>
				<X class="h-5 w-5" aria-hidden="true" />
			</Dialog.Close>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
