import { request } from './client';
import type { WeekendPolicy } from '#lib/types/index.js';

export interface PublicSettings {
	defaultVacationDays: number;
	vacationResetMonth: number;
	weekendPolicy: WeekendPolicy;
}

export const settingsApi = {
	/**
	 * Get public (non-sensitive) application settings
	 * Available to all authenticated users
	 */
	getPublic: (): Promise<PublicSettings> => {
		return request('/settings/public');
	}
};
