<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import Avatar from '$lib/components/ui/Avatar.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import { LogOut } from 'lucide-svelte';

	function handleLogout() {
		auth.logout();
		goto('/');
	}
</script>

<header class="bg-white shadow-sm border-b border-sand-200 px-6 py-4">
	<div class="flex items-center justify-between">
		<h1 class="text-xl font-semibold text-ocean-800">Admin Panel</h1>

		<div class="flex items-center gap-4">
			<div class="flex items-center gap-2">
				<Avatar name={auth.user?.name} size="sm" />
				<span class="text-sm font-medium text-ocean-700">
					{auth.user?.name}
				</span>
			</div>
			<Button variant="ghost" size="sm" onclick={handleLogout}>
				<LogOut class="w-4 h-4" />
			</Button>
		</div>
	</div>
</header>
