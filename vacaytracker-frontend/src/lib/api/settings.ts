import { request } from './client';

export interface PublicSettings {
	defaultVacationDays: number;
	vacationResetMonth: number;
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
