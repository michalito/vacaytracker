import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import CardWrapper from './CardWrapper.svelte';

describe('Card', () => {
	it('renders children content', () => {
		render(CardWrapper, { bodyText: 'Hello World' });
		expect(screen.getByText('Hello World')).toBeInTheDocument();
	});

	it('renders header when provided', () => {
		render(CardWrapper, { bodyText: 'Body', showHeader: true, headerText: 'My Header' });
		expect(screen.getByText('My Header')).toBeInTheDocument();
	});

	it('renders footer when provided', () => {
		render(CardWrapper, { bodyText: 'Body', showFooter: true, footerText: 'My Footer' });
		expect(screen.getByText('My Footer')).toBeInTheDocument();
	});

	it('renders all sections together', () => {
		render(CardWrapper, {
			bodyText: 'Body Content',
			showHeader: true,
			headerText: 'Top Section',
			showFooter: true,
			footerText: 'Bottom Section'
		});
		expect(screen.getByText('Top Section')).toBeInTheDocument();
		expect(screen.getByText('Body Content')).toBeInTheDocument();
		expect(screen.getByText('Bottom Section')).toBeInTheDocument();
	});

	it('does not render header or footer by default', () => {
		render(CardWrapper, { bodyText: 'Just body' });
		expect(screen.getByText('Just body')).toBeInTheDocument();
		expect(screen.queryByText('Header')).not.toBeInTheDocument();
		expect(screen.queryByText('Footer')).not.toBeInTheDocument();
	});
});
