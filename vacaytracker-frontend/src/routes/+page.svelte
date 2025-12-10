<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth.svelte';
	import { toast } from '$lib/stores/toast.svelte';
	import Input from '$lib/components/ui/Input.svelte';

	let email = $state('');
	let password = $state('');
	let isSubmitting = $state(false);
	let errors = $state<{ email?: string; password?: string }>({});
	let hasError = $state(false);

	// Redirect if already authenticated
	$effect(() => {
		if (auth.isAuthenticated) {
			goto(auth.isAdmin ? '/admin' : '/employee');
		}
	});

	function validate(): boolean {
		errors = {};

		if (!email) {
			errors.email = 'Email is required';
		} else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
			errors.email = 'Invalid email format';
		}

		if (!password) {
			errors.password = 'Password is required';
		} else if (password.length < 6) {
			errors.password = 'Password must be at least 6 characters';
		}

		const isValid = Object.keys(errors).length === 0;
		if (!isValid) {
			hasError = true;
			setTimeout(() => (hasError = false), 500);
		}
		return isValid;
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();

		if (!validate()) return;

		isSubmitting = true;

		try {
			await auth.login(email, password);
			toast.success('Welcome back!');
			goto(auth.isAdmin ? '/admin' : '/employee');
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Login failed');
			hasError = true;
			setTimeout(() => (hasError = false), 500);
		} finally {
			isSubmitting = false;
		}
	}
</script>

<svelte:head>
	<title>Login - VacayTracker</title>
</svelte:head>

<div
	class="min-h-screen flex items-center justify-center bg-gradient-to-br from-ocean-400 via-ocean-500 to-ocean-600 p-4 relative overflow-hidden"
>
	<!-- Decorative floating elements -->
	<div class="absolute inset-0 overflow-hidden pointer-events-none">
		<div
			class="absolute -top-20 -left-20 w-64 h-64 bg-ocean-300/20 rounded-full blur-3xl animate-float"
		></div>
		<div
			class="absolute -top-10 -right-10 w-48 h-48 bg-coral-400/10 rounded-full blur-2xl animate-float"
			style="animation-delay: -2s;"
		></div>
		<div
			class="absolute bottom-20 left-1/4 w-32 h-32 bg-ocean-200/20 rounded-full blur-2xl animate-float"
			style="animation-delay: -4s;"
		></div>
	</div>

	<!-- Animated waves at bottom -->
	<div class="absolute bottom-0 left-0 right-0 overflow-hidden">
		<svg
			viewBox="0 0 1440 320"
			class="w-[200%] h-32 text-ocean-700/30 animate-wave"
			preserveAspectRatio="none"
		>
			<path
				fill="currentColor"
				d="M0,224L48,213.3C96,203,192,181,288,181.3C384,181,480,203,576,213.3C672,224,768,224,864,208C960,192,1056,160,1152,165.3C1248,171,1344,213,1392,234.7L1440,256L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"
			></path>
		</svg>
		<svg
			viewBox="0 0 1440 320"
			class="absolute bottom-0 w-[200%] h-24 text-ocean-800/20 animate-wave"
			style="animation-delay: -4s;"
			preserveAspectRatio="none"
		>
			<path
				fill="currentColor"
				d="M0,288L48,272C96,256,192,224,288,213.3C384,203,480,213,576,229.3C672,245,768,267,864,261.3C960,256,1056,224,1152,213.3C1248,203,1344,213,1392,218.7L1440,224L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"
			></path>
		</svg>
	</div>

	<!-- Login Card -->
	<div
		class="w-full max-w-md relative z-10 animate-slide-up {hasError ? 'animate-shake' : ''}"
		style="animation-duration: 0.4s;"
	>
		<div
			class="bg-white/95 backdrop-blur-sm rounded-2xl shadow-xl shadow-ocean-900/10 border border-white/20 overflow-hidden"
		>
			<!-- Header -->
			<div class="px-8 pt-8 pb-4 text-center">
				<div class="flex justify-center mb-3">
					<img
						src="/logo.png"
						alt="VacayTracker"
						class="w-36 h-36 drop-shadow-lg"
					/>
				</div>
				<p class="text-ocean-600">Sign in to manage your vacation</p>
			</div>

			<!-- Form -->
			<form onsubmit={handleSubmit} class="px-8 pb-8 space-y-5">
				<Input
					type="email"
					label="Email"
					placeholder="you@company.com"
					bind:value={email}
					error={errors.email}
					required
				/>

				<Input
					type="password"
					label="Password"
					placeholder="Enter your password"
					bind:value={password}
					error={errors.password}
					showPasswordToggle
					required
				/>

				<button
					type="submit"
					disabled={isSubmitting}
					class="btn-wave w-full py-3 px-4 text-white font-semibold rounded-lg focus:outline-none focus:ring-2 focus:ring-ocean-300 focus:ring-offset-2 disabled:opacity-60 disabled:cursor-not-allowed disabled:hover:transform-none"
				>
					{#if isSubmitting}
						<span class="inline-flex items-center gap-2">
							<svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Signing in...
						</span>
					{:else}
						Sign In
					{/if}
				</button>
			</form>

			<!-- Footer -->
			<div class="px-8 py-4 bg-ocean-50/50 border-t border-ocean-100">
				<p class="text-center text-sm text-ocean-600">
					Contact your administrator if you need access
				</p>
			</div>
		</div>
	</div>
</div>
