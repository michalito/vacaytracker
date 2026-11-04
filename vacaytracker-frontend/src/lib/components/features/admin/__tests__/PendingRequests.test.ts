import { render, screen, waitFor } from '@testing-library/svelte';
import { userEvent } from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import type { VacationRequest } from '#lib/types/index.js';

vi.mock('#lib/api/admin.js', () => ({
	adminApi: { reviewRequest: vi.fn() }
}));

import { adminApi } from '#lib/api/admin.js';
import PendingRequests from '../PendingRequests.svelte';

function makeRequest(id: string, userName: string): VacationRequest {
	return {
		id,
		userId: `user-${id}`,
		userName,
		startDate: '2026-11-02',
		endDate: '2026-11-03',
		totalDays: 2,
		status: 'pending',
		createdAt: '2026-10-01T00:00:00Z',
		updatedAt: '2026-10-01T00:00:00Z'
	};
}

const requests = [makeRequest('r1', 'Alice Smith'), makeRequest('r2', 'Bob Jones')];

function reasonField() {
	return screen.findByLabelText<HTMLTextAreaElement>('Reason (optional)');
}

describe('PendingRequests', () => {
	beforeEach(() => {
		vi.mocked(adminApi.reviewRequest).mockResolvedValue(undefined as never);
	});

	it("does not carry one rejection reason into the next employee's rejection", async () => {
		const user = userEvent.setup();
		const onUpdate = vi.fn();
		render(PendingRequests, { props: { requests, onUpdate } });

		await user.click(screen.getAllByRole('button', { name: 'Reject' })[0]);
		await user.type(await reasonField(), 'Team offsite');
		await user.click(screen.getByRole('button', { name: 'Reject Request' }));

		expect(adminApi.reviewRequest).toHaveBeenCalledWith('r1', 'rejected', 'Team offsite');
		await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
		expect(onUpdate).toHaveBeenCalled();

		await user.click(screen.getAllByRole('button', { name: 'Reject' })[1]);
		expect((await reasonField()).value).toBe('');
	});

	it('clears a typed reason when the dialog is cancelled and reopened', async () => {
		const user = userEvent.setup();
		render(PendingRequests, { props: { requests, onUpdate: vi.fn() } });

		await user.click(screen.getAllByRole('button', { name: 'Reject' })[0]);
		await user.type(await reasonField(), 'Draft');
		await user.click(screen.getByRole('button', { name: 'Cancel' }));
		await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());

		await user.click(screen.getAllByRole('button', { name: 'Reject' })[0]);
		expect((await reasonField()).value).toBe('');
	});

	it('sends no reason when only whitespace was typed', async () => {
		const user = userEvent.setup();
		render(PendingRequests, { props: { requests, onUpdate: vi.fn() } });

		await user.click(screen.getAllByRole('button', { name: 'Reject' })[0]);
		await user.type(await reasonField(), '   ');
		await user.click(screen.getByRole('button', { name: 'Reject Request' }));

		expect(adminApi.reviewRequest).toHaveBeenCalledWith('r1', 'rejected', undefined);
	});

	it('approves the selected request and closes the dialog', async () => {
		const user = userEvent.setup();
		const onUpdate = vi.fn();
		render(PendingRequests, { props: { requests, onUpdate } });

		await user.click(screen.getAllByRole('button', { name: 'Approve' })[1]);
		expect(await screen.findByRole('dialog', { name: 'Approve Request' })).toHaveTextContent(
			'Bob Jones'
		);
		await user.click(screen.getByRole('button', { name: 'Approve Request' }));

		expect(adminApi.reviewRequest).toHaveBeenCalledWith('r2', 'approved');
		await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
		expect(onUpdate).toHaveBeenCalled();
	});

	it('keeps the dialog open when the review fails', async () => {
		vi.mocked(adminApi.reviewRequest).mockRejectedValueOnce(new Error('boom'));
		const user = userEvent.setup();
		const onUpdate = vi.fn();
		render(PendingRequests, { props: { requests, onUpdate } });

		await user.click(screen.getAllByRole('button', { name: 'Approve' })[0]);
		await user.click(await screen.findByRole('button', { name: 'Approve Request' }));

		expect(screen.getByRole('dialog', { name: 'Approve Request' })).toBeInTheDocument();
		expect(onUpdate).not.toHaveBeenCalled();
	});
});
