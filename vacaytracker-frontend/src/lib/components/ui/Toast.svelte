<script lang="ts">
	import { toast } from '$lib/stores/toast.svelte';
	import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-svelte';
	import { clsx } from 'clsx';

	const icons = {
		success: CheckCircle,
		error: AlertCircle,
		warning: AlertTriangle,
		info: Info
	};

	const styles = {
		success: 'bg-green-50 border-green-500 text-green-800',
		error: 'bg-red-50 border-red-500 text-red-800',
		warning: 'bg-yellow-50 border-yellow-500 text-yellow-800',
		info: 'bg-blue-50 border-blue-500 text-blue-800'
	};
</script>

<div class="fixed top-4 right-4 z-50 flex flex-col gap-2 w-80">
	{#each toast.toasts as t (t.id)}
		{@const Icon = icons[t.type]}
		<div
			class={clsx('flex items-start gap-3 p-4 rounded-lg border-l-4 shadow-lg', styles[t.type])}
			role="alert"
		>
			<Icon class="w-5 h-5 flex-shrink-0" />
			<p class="flex-1 text-sm">{t.message}</p>
			<button onclick={() => toast.dismiss(t.id)} class="flex-shrink-0 hover:opacity-70">
				<X class="w-4 h-4" />
			</button>
		</div>
	{/each}
</div>
