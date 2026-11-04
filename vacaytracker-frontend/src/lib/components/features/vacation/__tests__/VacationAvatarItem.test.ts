import { render, screen } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect } from 'vitest';
import type { TeamVacation } from '$lib/types';
import TooltipHarness from '$lib/components/ui/__tests__/TooltipHarness.svelte';
import VacationAvatarItem from '../VacationAvatarItem.svelte';

const vacation: TeamVacation = {
	id: 'v1',
	userId: 'u1',
	userName: 'Alice Smith',
	startDate: '2099-06-01',
	endDate: '2099-06-05',
	totalDays: 5
};

const FOCUSABLE = 'button, a[href], input, select, textarea, [tabindex]:not([tabindex="-1"])';

function renderItem() {
	return render(TooltipHarness, { props: { component: VacationAvatarItem, props: { vacation } } });
}

describe('VacationAvatarItem', () => {
	it('exposes a single named, focusable trigger (no nested interactive elements)', () => {
		const { container } = renderItem();

		const focusables = container.querySelectorAll(FOCUSABLE);
		expect(focusables).toHaveLength(1);
		expect(focusables[0].tagName).toBe('BUTTON');
		expect(screen.getByRole('button', { name: "Alice Smith's vacation details" })).toBe(
			focusables[0]
		);
	});

	it('opens the details popover with a labelled close button', async () => {
		const user = userEvent.setup();
		renderItem();

		await user.click(screen.getByRole('button', { name: "Alice Smith's vacation details" }));

		expect(await screen.findByText('Vacation Details')).toBeInTheDocument();
		// jsdom never finishes floating-ui positioning, so the popover stays `visibility: hidden`
		// and role queries compute an empty name; match the aria-label directly instead
		expect(screen.getByLabelText('Close vacation details').tagName).toBe('BUTTON');
	});

	it('labels a future vacation as upcoming rather than "On Vacation"', async () => {
		const user = userEvent.setup();
		renderItem();

		await user.click(screen.getByRole('button', { name: "Alice Smith's vacation details" }));

		expect(await screen.findByText('Upcoming')).toBeInTheDocument();
		expect(screen.queryByText('On Vacation')).not.toBeInTheDocument();
	});
});
