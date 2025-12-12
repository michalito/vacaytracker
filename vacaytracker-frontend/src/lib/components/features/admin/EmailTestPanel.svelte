<script lang="ts">
	import { createDialog, melt } from '@melt-ui/svelte';
	import type { EmailTemplateType, EmailPreviewResponse } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Mail, Send, UserPlus, Clock, CheckCircle, XCircle, Bell, Eye, X, Newspaper } from 'lucide-svelte';
	import { adminApi } from '$lib/api/admin';
	import { toast } from '$lib/stores/toast.svelte';

	interface EmailTemplate {
		id: EmailTemplateType;
		name: string;
		description: string;
		icon: typeof Mail;
		color: string;
	}

	const templates: EmailTemplate[] = [
		{
			id: 'welcome',
			name: 'Welcome Email',
			description: 'Sent to new users with login credentials',
			icon: UserPlus,
			color: 'text-ocean-500'
		},
		{
			id: 'request_submitted',
			name: 'Request Submitted',
			description: 'Confirmation when vacation is requested',
			icon: Clock,
			color: 'text-warning'
		},
		{
			id: 'request_approved',
			name: 'Request Approved',
			description: 'Notification when request is approved',
			icon: CheckCircle,
			color: 'text-success'
		},
		{
			id: 'request_rejected',
			name: 'Request Rejected',
			description: 'Notification when request is rejected',
			icon: XCircle,
			color: 'text-error'
		},
		{
			id: 'admin_notification',
			name: 'Admin Notification',
			description: 'Alert to admins for new pending requests',
			icon: Bell,
			color: 'text-purple-500'
		},
		{
			id: 'newsletter',
			name: 'Newsletter',
			description: 'Monthly summary sent to all users',
			icon: Newspaper,
			color: 'text-ocean-600'
		}
	];

	let sendingTemplate = $state<EmailTemplateType | null>(null);
	let previewingTemplate = $state<EmailTemplateType | null>(null);
	let previewData = $state<EmailPreviewResponse | null>(null);

	// Preview Dialog
	const {
		elements: { overlay, content, title, close },
		states: { open: previewOpen }
	} = createDialog({
		forceVisible: true
	});

	async function sendTestEmail(template: EmailTemplateType) {
		sendingTemplate = template;
		try {
			const result = await adminApi.sendTestEmail(template);
			toast.success(result.message);
		} catch (error) {
			toast.error(`Failed to send test email: ${error instanceof Error ? error.message : 'Unknown error'}`);
		} finally {
			sendingTemplate = null;
		}
	}

	async function previewEmail(template: EmailTemplateType) {
		previewingTemplate = template;
		try {
			previewData = await adminApi.previewEmail(template);
			previewOpen.set(true);
		} catch (error) {
			toast.error(`Failed to load preview: ${error instanceof Error ? error.message : 'Unknown error'}`);
		} finally {
			previewingTemplate = null;
		}
	}

	function closePreview() {
		previewOpen.set(false);
		previewData = null;
	}
</script>

<Card>
	{#snippet header()}
		<div class="flex items-center gap-2">
			<Mail class="w-5 h-5 text-ocean-500" />
			<h2 class="text-lg font-semibold text-ocean-700">Email Templates</h2>
		</div>
	{/snippet}

	<div class="space-y-4">
		<p class="text-sm text-ocean-600">
			Preview or send test emails to yourself to see how each template looks.
		</p>

		<div class="grid gap-3">
			{#each templates as template}
				{@const Icon = template.icon}
				<div class="flex items-center justify-between p-3 bg-sand-50 rounded-lg border border-sand-200">
					<div class="flex items-center gap-3">
						<div class="p-2 bg-white rounded-md shadow-sm">
							<Icon class="w-5 h-5 {template.color}" />
						</div>
						<div>
							<p class="font-medium text-ocean-800">{template.name}</p>
							<p class="text-xs text-ocean-500">{template.description}</p>
						</div>
					</div>
					<div class="flex gap-2">
						<Button
							variant="outline"
							size="sm"
							onclick={() => previewEmail(template.id)}
							loading={previewingTemplate === template.id}
							disabled={sendingTemplate !== null || previewingTemplate !== null}
						>
							<Eye class="w-4 h-4 mr-1" />
							Preview
						</Button>
						<Button
							variant="primary"
							size="sm"
							onclick={() => sendTestEmail(template.id)}
							loading={sendingTemplate === template.id}
							disabled={sendingTemplate !== null || previewingTemplate !== null}
						>
							<Send class="w-4 h-4 mr-1" />
							Send
						</Button>
					</div>
				</div>
			{/each}
		</div>

		<p class="text-xs text-ocean-500 pt-2 border-t border-sand-200">
			<strong>Note:</strong> Test emails use sample data and are sent to your admin email address.
		</p>
	</div>
</Card>

<!-- Preview Dialog -->
{#if $previewOpen && previewData}
	<div use:melt={$overlay} class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm"></div>
	<div
		use:melt={$content}
		class="fixed left-1/2 top-1/2 z-50 -translate-x-1/2 -translate-y-1/2 bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30 w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col"
	>
		<!-- Header -->
		<div class="flex items-center justify-between p-4 border-b border-ocean-100/50">
			<div>
				<h3 use:melt={$title} class="text-lg font-semibold text-ocean-800">{previewData.template} Preview</h3>
				<p class="text-sm text-ocean-600">
					<strong>Subject:</strong> {previewData.subject}
				</p>
			</div>
			<button
				use:melt={$close}
				class="p-2 rounded-lg hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
			>
				<X class="w-5 h-5 text-ocean-600" />
			</button>
		</div>

		<!-- Content -->
		<div class="flex-1 overflow-auto p-4 bg-ocean-50/30">
			<div class="bg-white rounded-xl shadow-sm p-4">
				<!-- eslint-disable-next-line svelte/no-at-html-tags -->
				{@html previewData.htmlBody}
			</div>
		</div>

		<!-- Footer -->
		<div class="flex justify-end gap-2 p-4 border-t border-ocean-100/50">
			<Button variant="outline" onclick={closePreview}>
				Close
			</Button>
			<Button
				onclick={() => {
					const templateId = templates.find(t => t.name === previewData?.template)?.id;
					if (templateId) {
						closePreview();
						sendTestEmail(templateId);
					}
				}}
				loading={sendingTemplate !== null}
			>
				<Send class="w-4 h-4 mr-1" />
				Send Test Email
			</Button>
		</div>
	</div>
{/if}
