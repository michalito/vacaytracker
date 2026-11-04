import { render } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';

vi.mock('#lib/api/admin.js', () => ({
	adminApi: { listUsers: vi.fn(), getSettings: vi.fn() }
}));

import { adminApi } from '#lib/api/admin.js';
import { admin } from './admin.svelte';
import EffectHarness from './__tests__/EffectHarness.svelte';

// Kept apart from admin.test.ts: that suite's fake timers and module resets hide effect re-runs.
describe('admin store commands inside an $effect', () => {
	it('do not re-run the effect when a fetch refreshes the cache', async () => {
		vi.mocked(adminApi.listUsers).mockResolvedValue({
			users: [],
			pagination: { page: 1, limit: 10, total: 0, totalPages: 0 }
		});
		// Settings failing to load keeps the settings cache cold, as on a flaky network
		vi.mocked(adminApi.getSettings).mockRejectedValue(new Error('offline'));

		// One effect that loads settings (reads the cache) and a filtered user list (always hits
		// the API and refreshes the cache). If the cache read were tracked, every refresh would
		// re-run the effect and refetch forever.
		let runs = 0;
		render(EffectHarness, {
			props: {
				run: () => {
					runs++;
					// Circuit breaker: a regression must fail this test, not starve the event loop
					if (runs > 10) return;
					admin.fetchSettings();
					admin.fetchUsers({ search: 'alice' });
				}
			}
		});
		await new Promise((resolve) => setTimeout(resolve, 100));

		expect(runs).toBe(1);
		expect(adminApi.listUsers).toHaveBeenCalledTimes(1);
	});
});
