<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { authApi } from '$lib/api/auth';
	import { toast } from '$lib/stores/toast.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Mail } from 'lucide-svelte';

	let vacationUpdates = $state(auth.user?.emailPreferences?.vacationUpdates ?? true);
	let weeklyDigest = $state(auth.user?.emailPreferences?.weeklyDigest ?? true);
	let teamNotifications = $state(auth.user?.emailPreferences?.teamNotifications ?? true);
	let isSubmitting = $state(false);

	// Sync with auth store changes
	$effect(() => {
		if (auth.user?.emailPreferences) {
			vacationUpdates = auth.user.emailPreferences.vacationUpdates;
			weeklyDigest = auth.user.emailPreferences.weeklyDigest;
			teamNotifications = auth.user.emailPreferences.teamNotifications;
		}
	});

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		isSubmitting = true;

		try {
			const result = await authApi.updateEmailPreferences({
				vacationUpdates,
				weeklyDigest,
				teamNotifications
			});

			// Update local auth store
			auth.updateUser({
				emailPreferences: result.emailPreferences
			});

			toast.success('Email preferences updated');
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to update preferences');
		} finally {
			isSubmitting = false;
		}
	}
</script>

<Card>
	{#snippet header()}
		<div class="flex items-center gap-2">
			<Mail class="w-5 h-5 text-ocean-500" />
			<h2 class="text-lg font-semibold text-ocean-700">Email Preferences</h2>
		</div>
	{/snippet}

	<form onsubmit={handleSubmit} class="space-y-4">
		<p class="text-sm text-ocean-600 mb-4">
			Choose which email notifications you'd like to receive.
		</p>

		<label class="flex items-center gap-3 cursor-pointer">
			<input
				type="checkbox"
				bind:checked={vacationUpdates}
				class="w-5 h-5 rounded border-sand-300 text-ocean-500 focus:ring-ocean-500"
			/>
			<div>
				<p class="font-medium text-ocean-800">Vacation Updates</p>
				<p class="text-sm text-ocean-500">
					Receive notifications when your requests are approved or rejected
				</p>
			</div>
		</label>

		<label class="flex items-center gap-3 cursor-pointer">
			<input
				type="checkbox"
				bind:checked={weeklyDigest}
				class="w-5 h-5 rounded border-sand-300 text-ocean-500 focus:ring-ocean-500"
			/>
			<div>
				<p class="font-medium text-ocean-800">Weekly Digest</p>
				<p class="text-sm text-ocean-500">Get a weekly summary of team vacation schedules</p>
			</div>
		</label>

		<label class="flex items-center gap-3 cursor-pointer">
			<input
				type="checkbox"
				bind:checked={teamNotifications}
				class="w-5 h-5 rounded border-sand-300 text-ocean-500 focus:ring-ocean-500"
			/>
			<div>
				<p class="font-medium text-ocean-800">Team Notifications</p>
				<p class="text-sm text-ocean-500">
					Be notified when teammates have upcoming vacations
				</p>
			</div>
		</label>

		<div class="pt-2">
			<Button type="submit" variant="primary" loading={isSubmitting}>Save Preferences</Button>
		</div>
	</form>
</Card>
