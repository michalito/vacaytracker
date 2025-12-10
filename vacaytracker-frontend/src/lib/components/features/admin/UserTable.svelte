<script lang="ts">
	import type { User } from '$lib/types';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { Edit, Trash2 } from 'lucide-svelte';

	interface Props {
		users: User[];
		onEdit: (user: User) => void;
		onDelete: (user: User) => void;
	}

	let { users, onEdit, onDelete }: Props = $props();

	function formatDate(dateStr: string | undefined): string {
		if (!dateStr) return '-';
		return new Date(dateStr).toLocaleDateString('en-GB', {
			day: 'numeric',
			month: 'short',
			year: 'numeric'
		});
	}
</script>

<div class="overflow-x-auto">
	<table class="w-full">
		<thead>
			<tr class="border-b border-sand-200">
				<th class="text-left py-3 px-4 text-sm font-medium text-ocean-600">User</th>
				<th class="text-left py-3 px-4 text-sm font-medium text-ocean-600">Role</th>
				<th class="text-left py-3 px-4 text-sm font-medium text-ocean-600">Balance</th>
				<th class="text-left py-3 px-4 text-sm font-medium text-ocean-600">Start Date</th>
				<th class="text-right py-3 px-4 text-sm font-medium text-ocean-600">Actions</th>
			</tr>
		</thead>
		<tbody>
			{#each users as user (user.id)}
				<tr class="border-b border-sand-100 hover:bg-sand-50 transition-colors">
					<td class="py-3 px-4">
						<div class="flex items-center gap-3">
							<Avatar name={user.name} size="sm" />
							<div>
								<p class="font-medium text-ocean-800">{user.name}</p>
								<p class="text-sm text-ocean-500">{user.email}</p>
							</div>
						</div>
					</td>
					<td class="py-3 px-4">
						<Badge variant={user.role === 'admin' ? 'approved' : 'pending'}>
							{user.role}
						</Badge>
					</td>
					<td class="py-3 px-4">
						<span class="font-medium text-ocean-800">{user.vacationBalance}</span>
						<span class="text-ocean-500"> days</span>
					</td>
					<td class="py-3 px-4 text-ocean-600">
						{formatDate(user.startDate)}
					</td>
					<td class="py-3 px-4">
						<div class="flex items-center justify-end gap-2">
							<Button variant="ghost" size="sm" onclick={() => onEdit(user)}>
								<Edit class="w-4 h-4" />
							</Button>
							<Button variant="ghost" size="sm" onclick={() => onDelete(user)}>
								<Trash2 class="w-4 h-4 text-error" />
							</Button>
						</div>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
