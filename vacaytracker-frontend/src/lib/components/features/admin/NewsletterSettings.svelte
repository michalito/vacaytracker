<script lang="ts">
	import { createDialog, melt } from '@melt-ui/svelte';
	import type { NewsletterConfig } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Checkbox from '$lib/components/ui/Checkbox.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import { Mail, Send, Eye, X } from 'lucide-svelte';
	import { admin } from '$lib/stores/admin.svelte';
	import { toast } from '$lib/stores/toast.svelte';

	interface Props {
		config: NewsletterConfig;
		onChange: (config: NewsletterConfig) => void;
	}

	let { config, onChange }: Props = $props();
	let isSending = $state(false);
	let isLoadingPreview = $state(false);

	const frequencyOptions = [
		{ value: 'weekly', label: 'Weekly' },
		{ value: 'monthly', label: 'Monthly' }
	];

	// Preview Dialog
	const {
		elements: { trigger: previewTrigger, overlay, content, title, description, close },
		states: { open: previewOpen }
	} = createDialog({
		forceVisible: true
	});

	function toggleEnabled() {
		onChange({
			...config,
			enabled: !config.enabled
		});
	}

	function handleFrequencyChange(value: string) {
		onChange({
			...config,
			frequency: value as 'weekly' | 'monthly'
		});
	}

	function handleDayChange(e: Event) {
		const target = e.target as HTMLInputElement;
		onChange({
			...config,
			dayOfMonth: parseInt(target.value, 10)
		});
	}

	async function handleSendNow() {
		isSending = true;
		try {
			const result = await admin.sendNewsletter();
			toast.success(result.message);
		} catch (error) {
			toast.error('Failed to send newsletter');
		} finally {
			isSending = false;
		}
	}

	async function handlePreview() {
		isLoadingPreview = true;
		try {
			await admin.fetchNewsletterPreview();
			previewOpen.set(true);
		} catch (error) {
			toast.error('Failed to load preview');
		} finally {
			isLoadingPreview = false;
		}
	}

	function formatLastSent(lastSentAt: string | null | undefined): string {
		if (!lastSentAt) return 'Never';
		return new Date(lastSentAt).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<Card>
	{#snippet header()}
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-2">
				<Mail class="w-5 h-5 text-ocean-500" />
				<h2 class="text-lg font-semibold text-ocean-700">Newsletter Settings</h2>
			</div>
			{#if config.enabled}
				<div class="flex gap-2">
					<Button variant="outline" size="sm" onclick={handlePreview} loading={isLoadingPreview}>
						<Eye class="w-4 h-4 mr-1" />
						Preview
					</Button>
					<Button size="sm" onclick={handleSendNow} loading={isSending}>
						<Send class="w-4 h-4 mr-1" />
						Send Now
					</Button>
				</div>
			{/if}
		</div>
	{/snippet}

	<div class="space-y-4">
		<p class="text-sm text-ocean-600">
			Configure automatic email digests for your team.
		</p>

		<label class="flex items-center gap-3 cursor-pointer">
			<Checkbox checked={config.enabled} onchange={toggleEnabled} />
			<div>
				<p class="font-medium text-ocean-800">Enable newsletter</p>
				<p class="text-sm text-ocean-500">
					Send periodic vacation summaries to all users
				</p>
			</div>
		</label>

		{#if config.enabled}
			<div class="pt-4 border-t border-ocean-100/50 space-y-4">
				<Select
					label="Frequency"
					value={config.frequency}
					options={frequencyOptions}
					onchange={handleFrequencyChange}
				/>

				{#if config.frequency === 'monthly'}
					<div>
						<label for="dayOfMonth" class="block text-sm font-semibold text-ocean-800 mb-1.5">
							Day of Month
						</label>
						<input
							id="dayOfMonth"
							type="number"
							value={config.dayOfMonth}
							onchange={handleDayChange}
							min="1"
							max="28"
							class="w-full px-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500 hover:border-ocean-500"
						/>
						<p class="text-xs text-ocean-500 mt-1.5">
							Choose a day between 1-28 to ensure delivery in all months.
						</p>
					</div>
				{:else}
					<p class="text-sm text-ocean-600">
						Newsletter will be sent every Monday.
					</p>
				{/if}

				<div class="pt-4 border-t border-ocean-100/50">
					<p class="text-sm text-ocean-600">
						<strong>Last sent:</strong> {formatLastSent(config.lastSentAt)}
					</p>
				</div>
			</div>
		{/if}
	</div>
</Card>

<!-- Preview Dialog -->
{#if $previewOpen && admin.newsletterPreview}
	<div use:melt={$overlay} class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm"></div>
	<div
		use:melt={$content}
		class="fixed left-1/2 top-1/2 z-50 -translate-x-1/2 -translate-y-1/2 bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30 w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col"
	>
		<div class="flex items-center justify-between p-4 border-b border-ocean-100/50">
			<div>
				<h3 use:melt={$title} class="text-lg font-semibold text-ocean-800">Newsletter Preview</h3>
				<p use:melt={$description} class="text-sm text-ocean-600">
					Will be sent to {admin.newsletterPreview.recipientCount} recipient{admin.newsletterPreview.recipientCount !== 1 ? 's' : ''}
				</p>
			</div>
			<button
				use:melt={$close}
				class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600 hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
			>
				<X class="w-5 h-5" />
			</button>
		</div>
		<div class="p-4 border-b border-ocean-100/50 bg-ocean-50/30">
			<p class="text-sm text-ocean-700">
				<strong>Subject:</strong> {admin.newsletterPreview.subject}
			</p>
		</div>
		<div class="flex-1 overflow-auto p-4">
			<!-- eslint-disable-next-line svelte/no-at-html-tags -->
			{@html admin.newsletterPreview.htmlBody}
		</div>
		<div class="flex justify-end gap-2 p-4 border-t border-ocean-100/50">
			<Button variant="outline" onclick={() => previewOpen.set(false)}>
				Close
			</Button>
			<Button onclick={handleSendNow} loading={isSending}>
				<Send class="w-4 h-4 mr-1" />
				Send Now
			</Button>
		</div>
	</div>
{/if}
