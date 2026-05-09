import test from 'node:test';
import assert from 'node:assert/strict';

import { computeNextRegistryState, REGISTRY_FORMAT_VERSION } from './registry-state.mjs';

test('starts revisioning at 1 for new registry state', () => {
  const nextState = computeNextRegistryState(undefined, 'hash-a', '2026-05-10T00:00:00.000Z');

  assert.deepEqual(nextState, {
    registryFormatVersion: REGISTRY_FORMAT_VERSION,
    revision: 1,
    contentHash: 'hash-a',
    updatedAt: '2026-05-10T00:00:00.000Z',
  });
});

test('preserves revision when content hash is unchanged', () => {
  const nextState = computeNextRegistryState(
    {
      registryFormatVersion: REGISTRY_FORMAT_VERSION,
      revision: 3,
      contentHash: 'hash-a',
      updatedAt: '2026-05-09T00:00:00.000Z',
    },
    'hash-a',
    '2026-05-10T00:00:00.000Z',
  );

  assert.equal(nextState.revision, 3);
  assert.equal(nextState.contentHash, 'hash-a');
  assert.equal(nextState.updatedAt, '2026-05-09T00:00:00.000Z');
});

test('increments revision when content hash changes', () => {
  const nextState = computeNextRegistryState(
    {
      registryFormatVersion: REGISTRY_FORMAT_VERSION,
      revision: 3,
      contentHash: 'hash-a',
      updatedAt: '2026-05-09T00:00:00.000Z',
    },
    'hash-b',
    '2026-05-10T00:00:00.000Z',
  );

  assert.equal(nextState.revision, 4);
  assert.equal(nextState.contentHash, 'hash-b');
});
