export const REGISTRY_FORMAT_VERSION = '1';

export function computeNextRegistryState(previousState, nextHash, updatedAt) {
  const previousHash = previousState?.contentHash;
  const previousRevision = Number(previousState?.revision ?? 0);
  const changed = previousHash !== nextHash;
  const revision = changed ? previousRevision + 1 : previousRevision;

  return {
    registryFormatVersion: REGISTRY_FORMAT_VERSION,
    revision,
    contentHash: nextHash,
    updatedAt: changed ? updatedAt : previousState?.updatedAt ?? updatedAt,
  };
}
