import { render, screen, waitFor } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

vi.mock('$lib/stores/vacation.svelte', () => ({
	vacation: { createRequest: vi.fn() }
}));

import RequestModal from '../RequestModal.svelte';

function reasonField() {
	return screen.findByLabelText<HTMLTextAreaElement>('Reason (optional)');
}

describe('RequestModal', () => {
	it('starts with an empty form when reopened after Cancel', async () => {
		const user = userEvent.setup();
		const { rerender } = render(RequestModal, { props: { open: true } });

		await user.type(await reasonField(), 'Old reason');
		await user.click(screen.getByRole('button', { name: 'Cancel' }));
		await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());

		await rerender({ open: false });
		await rerender({ open: true });
		expect((await reasonField()).value).toBe('');
	});

	it('requires a date range before submitting', async () => {
		const user = userEvent.setup();
		render(RequestModal, { props: { open: true } });

		await user.click(await screen.findByRole('button', { name: 'Submit Request' }));

		expect(await screen.findByRole('alert')).toHaveTextContent('Please select a date range');
	});

	it('labels its icon-only buttons', async () => {
		render(RequestModal, { props: { open: true } });

		expect(await screen.findByRole('dialog', { name: 'Request Vacation' })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Close request form' })).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Open calendar' })).toBeInTheDocument();
	});

	// Bits must render exactly the months it navigates; a CSS-hidden month would stay reachable
	// by keyboard on small screens
	it.each([
		[false, 1],
		[true, 2]
	])('wide screen = %s renders %i month grid(s)', async (wide, grids) => {
		vi.spyOn(window, 'matchMedia').mockImplementation(
			(query: string) =>
				({
					matches: wide,
					media: query,
					addEventListener: () => {},
					removeEventListener: () => {}
				}) as unknown as MediaQueryList
		);
		const user = userEvent.setup();
		render(RequestModal, { props: { open: true } });

		await user.click(await screen.findByRole('button', { name: 'Open calendar' }));

		await waitFor(() =>
			expect(document.querySelectorAll('[data-popover-content] table')).toHaveLength(grids)
		);
	});

	describe('business-day preview', () => {
		beforeEach(() => {
			// Fake only Date so user-event timers keep working; Monday 1 June 2026
			vi.useFakeTimers({ toFake: ['Date'] });
			vi.setSystemTime(new Date(2026, 5, 1, 9, 0));
		});

		afterEach(() => {
			vi.useRealTimers();
		});

		async function pickRange(start: string, end: string) {
			const user = userEvent.setup();
			await user.click(await screen.findByRole('button', { name: 'Open calendar' }));
			for (const value of [start, end]) {
				const day = await waitFor(() => {
					const el = document.querySelector<HTMLElement>(
						`.calendar-day[data-value="${value}"]:not([data-outside-month])`
					);
					if (!el) throw new Error(`day ${value} not rendered`);
					return el;
				});
				await user.click(day);
			}
		}

		function requestedDays() {
			const row = screen.getByText('Days Requested:').parentElement;
			return row?.textContent?.replace(/\s+/g, ' ').trim();
		}

		it('counts days with the company policy instead of a hard-coded weekend', async () => {
			// Only Fridays are non-working for this company
			render(RequestModal, {
				props: { open: true, weekendPolicy: { excludeWeekends: true, excludedDays: [5] } }
			});

			await pickRange('2026-06-08', '2026-06-14'); // Mon .. Sun

			await waitFor(() => expect(requestedDays()).toBe('Days Requested: 6 days'));
		});

		it('counts every day when the policy is off', async () => {
			render(RequestModal, {
				props: { open: true, weekendPolicy: { excludeWeekends: false, excludedDays: [0, 6] } }
			});

			await pickRange('2026-06-08', '2026-06-14');

			await waitFor(() => expect(requestedDays()).toBe('Days Requested: 7 days'));
		});
	});
});
