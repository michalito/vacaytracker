<script lang="ts">
	import type { User, Role } from '$lib/types';
	import { admin } from '$lib/stores/admin.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import { X, UserPlus, Edit } from 'lucide-svelte';

	interface Props {
		open: boolean;
		user?: User | null;
		onClose: () => void;
	}

	let { open = $bindable(false), user = null, onClose }: Props = $props();

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let role = $state<Role>('employee');
	let vacationBalance = $state(25);
	let startDate = $state('');
	let isSubmitting = $state(false);
	let errors = $state<Record<string, string>>({});

	const isEditing = $derived(!!user);

	// Sync form with user when editing
	$effect(() => {
		if (user) {
			name = user.name;
			email = user.email;
			password = '';
			role = user.role;
			vacationBalance = user.vacationBalance;
			startDate = user.startDate ? user.startDate.split('T')[0] : '';
		} else {
			resetForm();
		}
	});

	function resetForm() {
		name = '';
		email = '';
		password = '';
		role = 'employee';
		vacationBalance = 25;
		startDate = '';
		errors = {};
	}

	function handleClose() {
		open = false;
		resetForm();
		onClose();
	}

	function validate(): boolean {
		errors = {};

		if (!name.trim()) {
			errors.name = 'Name is required';
		}

		if (!email.trim()) {
			errors.email = 'Email is required';
		} else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
			errors.email = 'Invalid email format';
		}

		if (!isEditing && !password) {
			errors.password = 'Password is required for new users';
		} else if (password && password.length < 8) {
			errors.password = 'Password must be at least 8 characters';
		}

		return Object.keys(errors).length === 0;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();

		if (!validate()) return;

		isSubmitting = true;

		try {
			if (isEditing && user) {
				await admin.updateUser(user.id, {
					name,
					email,
					role,
					vacationBalance,
					startDate: startDate || undefined
				});
				toast.success('User updated');
			} else {
				await admin.createUser({
					name,
					email,
					password,
					role,
					vacationBalance,
					startDate: startDate || undefined
				});
				toast.success('User created');
			}
			handleClose();
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to save user');
		} finally {
			isSubmitting = false;
		}
	}
</script>

{#if open}
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<div class="fixed inset-0 z-50">
		<div class="fixed inset-0 bg-black/50 backdrop-blur-sm" onclick={handleClose}></div>
		<div class="fixed inset-0 flex items-center justify-center p-4">
			<div
				class="bg-white rounded-xl shadow-xl w-full max-w-md"
				onclick={(e) => e.stopPropagation()}
			>
				<!-- Header -->
				<div class="flex items-center justify-between p-4 border-b border-sand-200">
					<div class="flex items-center gap-2">
						{#if isEditing}
							<Edit class="w-5 h-5 text-ocean-500" />
							<h2 class="text-lg font-semibold text-ocean-800">Edit User</h2>
						{:else}
							<UserPlus class="w-5 h-5 text-ocean-500" />
							<h2 class="text-lg font-semibold text-ocean-800">Create User</h2>
						{/if}
					</div>
					<button type="button" onclick={handleClose} class="text-ocean-400 hover:text-ocean-600">
						<X class="w-5 h-5" />
					</button>
				</div>

				<!-- Body -->
				<form onsubmit={handleSubmit} class="p-4 space-y-4">
					<Input label="Name" placeholder="John Doe" bind:value={name} error={errors.name} required />

					<Input
						type="email"
						label="Email"
						placeholder="john@company.com"
						bind:value={email}
						error={errors.email}
						required
					/>

					{#if !isEditing}
						<Input
							type="password"
							label="Password"
							placeholder="Enter password"
							bind:value={password}
							error={errors.password}
							required
						/>
					{/if}

					<div>
						<label for="role" class="block text-sm font-medium text-ocean-700 mb-1">Role</label>
						<select
							id="role"
							bind:value={role}
							class="w-full px-3 py-2 rounded-md border border-sand-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent"
						>
							<option value="employee">Employee</option>
							<option value="admin">Admin</option>
						</select>
					</div>

					<div class="grid grid-cols-2 gap-4">
						<div>
							<label for="vacationBalance" class="block text-sm font-medium text-ocean-700 mb-1">
								Vacation Balance
							</label>
							<input
								id="vacationBalance"
								type="number"
								bind:value={vacationBalance}
								min="0"
								max="365"
								class="w-full px-3 py-2 rounded-md border border-sand-300 focus:ring-2 focus:ring-ocean-500 focus:border-transparent"
							/>
						</div>
						<Input type="date" label="Start Date" bind:value={startDate} />
					</div>

					<!-- Actions -->
					<div class="flex gap-3 pt-2">
						<Button type="button" variant="outline" class="flex-1" onclick={handleClose}>
							Cancel
						</Button>
						<Button type="submit" variant="primary" class="flex-1" loading={isSubmitting}>
							{isEditing ? 'Save Changes' : 'Create User'}
						</Button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}
