import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('./client', () => ({
	request: vi.fn()
}));

import { request } from './client';
import { settingsApi } from './settings';
import type { PublicSettings } from './settings';

describe('settingsApi', () => {
	beforeEach(() => {
		vi.mocked(request).mockReset();
	});

	describe('getPublic', () => {
		it('calls request with /settings/public and returns settings', async () => {
			const settings: PublicSettings = {
				defaultVacationDays: 25,
				vacationResetMonth: 1
			};
			vi.mocked(request).mockResolvedValue(settings);

			const result = await settingsApi.getPublic();

			expect(request).toHaveBeenCalledWith('/settings/public');
			expect(result).toEqual(settings);
		});
	});
});
