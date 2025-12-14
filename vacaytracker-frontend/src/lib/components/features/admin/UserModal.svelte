<script lang="ts">
	import { createDialog, melt } from '@melt-ui/svelte';
	import type { User, Role } from '$lib/types';
	import { admin } from '$lib/stores/admin.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Select from '$lib/components/ui/Select.svelte';
	import { X, UserPlus, Edit } from 'lucide-svelte';

	const roleOptions = [
		{ value: 'employee', label: 'Employee' },
		{ value: 'admin', label: 'Admin' }
	];

	interface Props {
		open: boolean;
		user?: User | null;
		onClose: () => void;
	}

	let { open = $bindable(false), user = null, onClose }: Props = $props();

	// Create Melt-UI dialog with controlled state
	const {
		elements: { overlay, content, title, close, portalled },
		states: { open: dialogOpen }
	} = createDialog({
		forceVisible: true,
		onOpenChange: ({ next }) => {
			open = next;
			if (!next) {
				resetForm();
				onClose();
			}
			return next;
		}
	});

	// Sync external open prop with dialog state
	$effect(() => {
		dialogOpen.set(open);
	});

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
			dialogOpen.set(false);
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Failed to save user');
		} finally {
			isSubmitting = false;
		}
	}
</script>

{#if open}
	<div use:melt={$portalled}>
		<!-- Overlay -->
		<div
			use:melt={$overlay}
			class="fixed inset-0 z-50 bg-ocean-900/50 backdrop-blur-sm transition-opacity duration-200
				data-[state=open]:opacity-100 data-[state=closed]:opacity-0"
		></div>

		<!-- Content Container -->
		<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
			<div
				use:melt={$content}
				class="bg-white/95 backdrop-blur-md rounded-2xl shadow-xl border border-white/30 w-full max-w-md
					transition-all duration-200 data-[state=open]:opacity-100 data-[state=open]:scale-100
					data-[state=closed]:opacity-0 data-[state=closed]:scale-95"
			>
				<!-- Header -->
				<div class="flex items-center justify-between p-4 border-b border-ocean-100/50">
					<div class="flex items-center gap-2">
						{#if isEditing}
							<Edit class="w-5 h-5 text-ocean-500" />
							<h2 use:melt={$title} class="text-lg font-semibold text-ocean-800">Edit User</h2>
						{:else}
							<UserPlus class="w-5 h-5 text-ocean-500" />
							<h2 use:melt={$title} class="text-lg font-semibold text-ocean-800">Create User</h2>
						{/if}
					</div>
					<button
						use:melt={$close}
						type="button"
						aria-label="Close dialog"
						class="p-1 rounded-lg text-ocean-400 hover:text-ocean-600 hover:bg-ocean-500/10 transition-all duration-200 cursor-pointer"
					>
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

					<Select
						label="Role"
						bind:value={role}
						options={roleOptions}
					/>

					<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
						<div>
							<label for="vacationBalance" class="block text-sm font-semibold text-ocean-800 mb-1.5">
								Vacation Balance
							</label>
							<input
								id="vacationBalance"
								type="number"
								bind:value={vacationBalance}
								min="0"
								max="365"
								class="w-full px-4 py-2.5 rounded-lg border-2 border-ocean-500/40 bg-white text-ocean-900 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ocean-500/30 focus:border-ocean-500 hover:border-ocean-500"
							/>
						</div>
						<Input type="date" label="Start Date" bind:value={startDate} />
					</div>

					<!-- Actions -->
					<div class="flex gap-3 pt-2">
						<Button type="button" variant="outline" class="flex-1" onclick={() => dialogOpen.set(false)}>
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
