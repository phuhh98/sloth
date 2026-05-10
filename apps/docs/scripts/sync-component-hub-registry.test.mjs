import test from "node:test";
import assert from "node:assert/strict";
import { mkdtemp, mkdir, readFile, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";

import {
  computeNextRegistryState,
  REGISTRY_FORMAT_VERSION,
} from "./registry-state.mjs";
import { syncRegistry } from "./sync-component-hub-registry.mjs";

test("starts revisioning at 1 for new registry state", () => {
  const nextState = computeNextRegistryState(
    undefined,
    "hash-a",
    "2026-05-10T00:00:00.000Z",
  );

  assert.deepEqual(nextState, {
    registryFormatVersion: REGISTRY_FORMAT_VERSION,
    revision: 1,
    contentHash: "hash-a",
    updatedAt: "2026-05-10T00:00:00.000Z",
  });
});

test("preserves revision when content hash is unchanged", () => {
  const nextState = computeNextRegistryState(
    {
      registryFormatVersion: REGISTRY_FORMAT_VERSION,
      revision: 3,
      contentHash: "hash-a",
      updatedAt: "2026-05-09T00:00:00.000Z",
    },
    "hash-a",
    "2026-05-10T00:00:00.000Z",
  );

  assert.equal(nextState.revision, 3);
  assert.equal(nextState.contentHash, "hash-a");
  assert.equal(nextState.updatedAt, "2026-05-09T00:00:00.000Z");
});

test("increments revision when content hash changes", () => {
  const nextState = computeNextRegistryState(
    {
      registryFormatVersion: REGISTRY_FORMAT_VERSION,
      revision: 3,
      contentHash: "hash-a",
      updatedAt: "2026-05-09T00:00:00.000Z",
    },
    "hash-b",
    "2026-05-10T00:00:00.000Z",
  );

  assert.equal(nextState.revision, 4);
  assert.equal(nextState.contentHash, "hash-b");
});

test("syncRegistry preserves existing historical release folders", async () => {
  const tempRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-docs-sync-"));
  const sourceRoot = path.join(tempRoot, "source-registry");
  const targetRoot = path.join(tempRoot, "target-registry");

  await mkdir(path.join(sourceRoot, "contracts", "1.0.0"), { recursive: true });
  await writeFile(
    path.join(sourceRoot, "contracts", "1.0.0", "manifest.json"),
    JSON.stringify({ version: "1.0.0", components: {} }),
    "utf8",
  );
  await writeFile(
    path.join(sourceRoot, "contracts", "index.json"),
    JSON.stringify({
      registryFormatVersion: "1",
      items: [{ version: "1.0.0" }],
    }),
    "utf8",
  );

  await mkdir(path.join(targetRoot, "contracts", "0.9.0"), { recursive: true });
  await writeFile(
    path.join(targetRoot, "contracts", "0.9.0", "manifest.json"),
    JSON.stringify({ version: "0.9.0", components: { legacy: true } }),
    "utf8",
  );

  await syncRegistry({
    sourceRootPath: sourceRoot,
    targetRootPath: targetRoot,
    now: "2026-05-10T10:00:00.000Z",
  });

  const preservedLegacyManifest = JSON.parse(
    await readFile(
      path.join(targetRoot, "contracts", "0.9.0", "manifest.json"),
      "utf8",
    ),
  );
  assert.equal(preservedLegacyManifest.version, "0.9.0");

  const copiedNewManifest = JSON.parse(
    await readFile(
      path.join(targetRoot, "contracts", "1.0.0", "manifest.json"),
      "utf8",
    ),
  );
  assert.equal(copiedNewManifest.version, "1.0.0");
});

test("syncRegistry keeps revision when resulting content does not change", async () => {
  const tempRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-docs-sync-"));
  const sourceRoot = path.join(tempRoot, "source-registry");
  const targetRoot = path.join(tempRoot, "target-registry");

  await mkdir(path.join(sourceRoot, "contracts", "1.0.0"), { recursive: true });
  await writeFile(
    path.join(sourceRoot, "contracts", "1.0.0", "manifest.json"),
    JSON.stringify({ version: "1.0.0", components: {} }),
    "utf8",
  );
  await writeFile(
    path.join(sourceRoot, "contracts", "index.json"),
    JSON.stringify({
      registryFormatVersion: "1",
      items: [{ version: "1.0.0" }],
    }),
    "utf8",
  );

  const firstState = await syncRegistry({
    sourceRootPath: sourceRoot,
    targetRootPath: targetRoot,
    now: "2026-05-10T10:00:00.000Z",
  });
  const secondState = await syncRegistry({
    sourceRootPath: sourceRoot,
    targetRootPath: targetRoot,
    now: "2026-05-10T11:00:00.000Z",
  });

  assert.equal(firstState.revision, 1);
  assert.equal(secondState.revision, 1);
  assert.equal(secondState.updatedAt, "2026-05-10T10:00:00.000Z");
});
