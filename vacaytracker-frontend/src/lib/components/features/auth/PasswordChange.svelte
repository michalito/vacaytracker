<script lang="ts">
	import { authApi } from '$lib/api/auth';
	import { toast } from '$lib/stores/toast.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import { Lock } from 'lucide-svelte';

	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let isSubmitting = $state(false);
	let errors = $state<{
		currentPassword?: string;
		newPassword?: string;
		confirmPassword?: string;
	}>({});

	function validate(): boolean {
		errors = {};

		if (!currentPassword) {
			errors.currentPassword = 'Current password is required';
		}

		if (!newPassword) {
			errors.newPassword = 'New password is required';
		} else if (newPassword.length < 8) {
			errors.newPassword = 'Password must be at least 8 characters';
		}

		if (!confirmPassword) {
			errors.confirmPassword = 'Please confirm your password';
		} else if (newPassword !== confirmPassword) {
			errors.confirmPassword = 'Passwords do not match';
		}

		return Object.keys(errors).length === 0;
	}

	function resetForm() {
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
		errors = {};
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();

		if (!validate()) return;

		isSubmitting = true;

		try {
			await authApi.changePassword(currentPassword, newPassword);
			toast.success('Password changed successfully');
			resetForm();
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to change password');
		} finally {
			isSubmitting = false;
		}
	}
</script>

<Card>
	{#snippet header()}
		<div class="flex items-center gap-2">
			<Lock class="w-5 h-5 text-ocean-500" />
			<h2 class="text-lg font-semibold text-ocean-700">Change Password</h2>
		</div>
	{/snippet}

	<form onsubmit={handleSubmit} class="space-y-4">
		<Input
			type="password"
			label="Current Password"
			placeholder="Enter your current password"
			bind:value={currentPassword}
			error={errors.currentPassword}
			required
		/>

		<Input
			type="password"
			label="New Password"
			placeholder="Enter your new password"
			bind:value={newPassword}
			error={errors.newPassword}
			required
		/>

		<Input
			type="password"
			label="Confirm New Password"
			placeholder="Confirm your new password"
			bind:value={confirmPassword}
			error={errors.confirmPassword}
			required
		/>

		<div class="pt-2">
			<Button type="submit" variant="primary" loading={isSubmitting}>Change Password</Button>
		</div>
	</form>
</Card>
