<script lang="ts">
	import { toast } from '$lib/stores/toast.svelte';
	import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-svelte';

	const icons = {
		success: CheckCircle,
		error: AlertCircle,
		warning: AlertTriangle,
		info: Info
	} as const;

	// Semantic color classes matching app.css theme
	const styles = {
		success: {
			container: 'bg-success/10 border-success',
			text: 'text-success',
			icon: 'text-success'
		},
		error: {
			container: 'bg-error/10 border-error',
			text: 'text-error',
			icon: 'text-error'
		},
		warning: {
			container: 'bg-warning/10 border-warning',
			text: 'text-warning',
			icon: 'text-warning'
		},
		info: {
			container: 'bg-info/10 border-info',
			text: 'text-info',
			icon: 'text-info'
		}
	} as const;
</script>

<!-- Toast container: bottom center, stacked in reverse order -->
<div
	class="fixed bottom-4 left-1/2 -translate-x-1/2 z-[9999] flex flex-col-reverse gap-2 w-full max-w-md px-4 pointer-events-none"
	role="region"
	aria-label="Notifications"
>
	{#each toast.toasts as t (t.id)}
		{@const Icon = icons[t.type]}
		{@const style = styles[t.type]}
		<div
			role="alert"
			aria-live="polite"
			class="flex items-start gap-3 p-4 rounded-lg border-l-4 shadow-lg bg-white pointer-events-auto animate-toast-in {style.container}"
		>
			<Icon class="w-5 h-5 flex-shrink-0 mt-0.5 {style.icon}" aria-hidden="true" />
			<div class="flex-1 min-w-0">
				<p class="font-semibold text-ocean-800">
					{t.title}
				</p>
				{#if t.description}
					<p class="text-sm mt-1 text-ocean-600">
						{t.description}
					</p>
				{/if}
			</div>
			<button
				type="button"
				onclick={() => toast.dismiss(t.id)}
				class="flex-shrink-0 p-1 rounded-md transition-colors cursor-pointer text-ocean-400 hover:text-ocean-600 hover:bg-ocean-50"
				aria-label="Dismiss notification"
			>
				<X class="w-4 h-4" />
			</button>
		</div>
	{/each}
</div>
