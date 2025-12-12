<script lang="ts">
	const currentYear = new Date().getFullYear();
</script>

<footer class="fixed bottom-0 left-0 right-0 z-40">
	<!-- Wave Container -->
	<div class="relative h-32 overflow-hidden">
		<!-- First Wave Layer -->
		<svg
			viewBox="0 0 1440 320"
			class="absolute bottom-0 w-[200%] h-24 text-ocean-400/30 animate-wave"
			preserveAspectRatio="none"
		>
			<path
				fill="currentColor"
				d="M0,224L48,213.3C96,203,192,181,288,181.3C384,181,480,203,576,213.3C672,224,768,224,864,208C960,192,1056,160,1152,165.3C1248,171,1344,213,1392,234.7L1440,256L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"
			></path>
		</svg>

		<!-- Second Wave Layer -->
		<svg
			viewBox="0 0 1440 320"
			class="absolute bottom-0 w-[200%] h-20 text-ocean-500/40 animate-wave"
			style="animation-delay: -3s;"
			preserveAspectRatio="none"
		>
			<path
				fill="currentColor"
				d="M0,288L48,272C96,256,192,224,288,213.3C384,203,480,213,576,229.3C672,245,768,267,864,261.3C960,256,1056,224,1152,213.3C1248,203,1344,213,1392,218.7L1440,224L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"
			></path>
		</svg>

		<!-- Third Wave Layer (darkest, on top) -->
		<svg
			viewBox="0 0 1440 320"
			class="absolute bottom-0 w-[200%] h-16 text-ocean-600/50 animate-wave"
			style="animation-delay: -6s;"
			preserveAspectRatio="none"
		>
			<path
				fill="currentColor"
				d="M0,256L48,240C96,224,192,192,288,186.7C384,181,480,203,576,218.7C672,235,768,245,864,234.7C960,224,1056,192,1152,181.3C1248,171,1344,181,1392,186.7L1440,192L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"
			></path>
		</svg>
	</div>

	<!-- Footer Content -->
	<div class="bg-ocean-700 py-4">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex flex-col sm:flex-row items-center justify-between gap-2">
				<div class="flex items-center gap-2">
					<img src="/logo.png" alt="VacayTracker" class="w-6 h-6 opacity-80" />
					<span class="text-ocean-200 text-sm font-medium">VacayTracker</span>
				</div>
				<p class="text-ocean-300 text-sm">
					&copy; {currentYear} All rights reserved
				</p>
			</div>
		</div>
	</div>
</footer>
