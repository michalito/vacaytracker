import { describe, it, expect } from 'vitest';
import { USER_COLOR_PALETTE, getUserColorIndex, getUserColor } from '$lib/utils/colors';
import type { UserColor } from '$lib/utils/colors';

describe('USER_COLOR_PALETTE', () => {
	it('has exactly 8 entries', () => {
		expect(USER_COLOR_PALETTE).toHaveLength(8);
	});

	it('each entry has background, text, and combined properties', () => {
		for (const color of USER_COLOR_PALETTE) {
			expect(color).toHaveProperty('background');
			expect(color).toHaveProperty('text');
			expect(color).toHaveProperty('combined');
			expect(typeof color.background).toBe('string');
			expect(typeof color.text).toBe('string');
			expect(typeof color.combined).toBe('string');
		}
	});

	it('each combined is the concatenation of background and text', () => {
		for (const color of USER_COLOR_PALETTE) {
			expect(color.combined).toBe(`${color.background} ${color.text}`);
		}
	});
});

describe('getUserColorIndex', () => {
	it('returns the same value for the same userId (deterministic)', () => {
		const index1 = getUserColorIndex('user-abc-123');
		const index2 = getUserColorIndex('user-abc-123');
		expect(index1).toBe(index2);
	});

	it('returns a value within palette bounds (0 to 7)', () => {
		const testIds = [
			'user-1',
			'user-2',
			'abc',
			'some-long-user-id-with-many-characters',
			'',
			'12345',
			'a',
			'zzz'
		];
		for (const id of testIds) {
			const index = getUserColorIndex(id);
			expect(index).toBeGreaterThanOrEqual(0);
			expect(index).toBeLessThan(USER_COLOR_PALETTE.length);
		}
	});

	it('returns an integer', () => {
		const index = getUserColorIndex('user-test');
		expect(Number.isInteger(index)).toBe(true);
	});

	it('different userIds can produce different indices', () => {
		const indices = new Set<number>();
		const testIds = [
			'alice',
			'bob',
			'charlie',
			'diana',
			'eve',
			'frank',
			'grace',
			'heidi',
			'ivan',
			'judy'
		];
		for (const id of testIds) {
			indices.add(getUserColorIndex(id));
		}
		// With 10 different inputs and 8 possible outputs, we expect at least 2 distinct values
		expect(indices.size).toBeGreaterThan(1);
	});
});

describe('getUserColor', () => {
	it('returns a valid UserColor object', () => {
		const color: UserColor = getUserColor('user-test');
		expect(color).toHaveProperty('background');
		expect(color).toHaveProperty('text');
		expect(color).toHaveProperty('combined');
	});

	it('returns a color from the palette', () => {
		const color = getUserColor('user-test');
		expect(USER_COLOR_PALETTE).toContainEqual(color);
	});

	it('is consistent for the same userId', () => {
		const color1 = getUserColor('user-abc-123');
		const color2 = getUserColor('user-abc-123');
		expect(color1).toEqual(color2);
	});

	it('returns the color at the index from getUserColorIndex', () => {
		const userId = 'user-xyz-789';
		const index = getUserColorIndex(userId);
		const color = getUserColor(userId);
		expect(color).toEqual(USER_COLOR_PALETTE[index]);
	});

	it('different userIds can produce different colors', () => {
		const colors = new Set<string>();
		const testIds = [
			'alice',
			'bob',
			'charlie',
			'diana',
			'eve',
			'frank',
			'grace',
			'heidi',
			'ivan',
			'judy'
		];
		for (const id of testIds) {
			colors.add(getUserColor(id).combined);
		}
		expect(colors.size).toBeGreaterThan(1);
	});
});
