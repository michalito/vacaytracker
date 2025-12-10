<script lang="ts">
	import { page } from '$app/stores';
	import Button from '$lib/components/ui/Button.svelte';
	import { Umbrella, Home, RefreshCw } from 'lucide-svelte';

	const errorMessages: Record<number, { title: string; message: string }> = {
		404: {
			title: 'Page Not Found',
			message: "Looks like this page sailed away. The page you're looking for doesn't exist or has been moved."
		},
		500: {
			title: 'Server Error',
			message: "Something went wrong on our end. Our crew is working to fix it. Please try again later."
		},
		403: {
			title: 'Access Denied',
			message: "You don't have permission to access this page. Please contact your administrator if you believe this is an error."
		}
	};

	const statusCode = $derived($page.status);
	const errorInfo = $derived(errorMessages[statusCode] ?? {
		title: 'Something Went Wrong',
		message: 'An unexpected error occurred. Please try again.'
	});
</script>

<svelte:head>
	<title>{errorInfo.title} - VacayTracker</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-b from-ocean-50 to-sand-100 flex items-center justify-center p-4">
	<div class="text-center max-w-md">
		<!-- Beach illustration -->
		<div class="mb-8 relative">
			<div class="w-32 h-32 mx-auto bg-ocean-100 rounded-full flex items-center justify-center">
				<Umbrella class="w-16 h-16 text-ocean-400" />
			</div>
			<div class="absolute -bottom-2 left-1/2 -translate-x-1/2 w-48 h-4 bg-sand-300 rounded-full blur-md"></div>
		</div>

		<!-- Error code -->
		<p class="text-6xl font-bold text-ocean-300 mb-2">{statusCode}</p>

		<!-- Error title -->
		<h1 class="text-2xl font-bold text-ocean-800 mb-3">{errorInfo.title}</h1>

		<!-- Error message -->
		<p class="text-ocean-600 mb-8">{errorInfo.message}</p>

		<!-- Actions -->
		<div class="flex flex-col sm:flex-row gap-3 justify-center">
			<Button onclick={() => history.back()} variant="outline">
				<RefreshCw class="w-4 h-4 mr-2" />
				Go Back
			</Button>
			<Button onclick={() => (window.location.href = '/')}>
				<Home class="w-4 h-4 mr-2" />
				Go Home
			</Button>
		</div>

		<!-- Additional help -->
		{#if $page.error?.message}
			<details class="mt-8 text-left bg-white/50 rounded-lg p-4">
				<summary class="cursor-pointer text-sm text-ocean-500 hover:text-ocean-700">
					Technical Details
				</summary>
				<pre class="mt-2 text-xs text-ocean-600 whitespace-pre-wrap break-words">{$page.error.message}</pre>
			</details>
		{/if}
	</div>
</div>
