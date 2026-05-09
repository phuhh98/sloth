import test from "node:test";
import assert from "node:assert/strict";
import { mkdtemp, mkdir, rm, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";

import {
  MIN_DEPRECATION_MONTHS,
  computeContentHash,
  stripMutableManifestFields,
  validateContracts,
  validateReleaseSet,
} from "./contract-version-policy.mjs";

test("computes deterministic content hash", () => {
  assert.equal(
    computeContentHash('{"name":"hero-banner"}'),
    computeContentHash('{"name":"hero-banner"}'),
  );
});

test("allows only deprecatedAt as mutable manifest field", () => {
  assert.deepEqual(
    stripMutableManifestFields({
      name: "hero-banner",
      version: "1.0.0",
      deprecatedAt: "2026-11-10T00:00:00.000Z",
    }),
    { name: "hero-banner", version: "1.0.0" },
  );
});

test("requires deprecatedAt on non-latest releases", () => {
  const errors = validateReleaseSet(
    [
      { version: "1.1.0", manifest: {} },
      { version: "1.0.0", manifest: {} },
    ],
    new Date("2026-05-10T00:00:00.000Z"),
  );

  assert.equal(errors.length, 1);
  assert.match(errors[0], /must declare deprecatedAt/);
});

test("requires at least six months deprecation window", () => {
  const errors = validateReleaseSet(
    [
      { version: "1.1.0", manifest: {} },
      {
        version: "1.0.0",
        manifest: { deprecatedAt: "2026-08-10T00:00:00.000Z" },
      },
    ],
    new Date("2026-05-10T00:00:00.000Z"),
  );

  assert.equal(errors.length, 1);
  assert.match(errors[0], new RegExp(`${MIN_DEPRECATION_MONTHS} months`));
});

test("accepts non-latest releases with sufficient deprecation window", () => {
  const errors = validateReleaseSet(
    [
      { version: "1.1.0", manifest: {} },
      {
        version: "1.0.0",
        manifest: { deprecatedAt: "2026-11-10T00:00:00.000Z" },
      },
    ],
    new Date("2026-05-10T00:00:00.000Z"),
  );

  assert.deepEqual(errors, []);
});

test("invalid release folder error uses latest valid release as example", async () => {
  const tmpRoot = await mkdtemp(
    path.join(os.tmpdir(), "sloth-contract-policy-"),
  );

  try {
    await mkdir(path.join(tmpRoot, "src", "contracts", "1.2.3", "components"), {
      recursive: true,
    });
    await mkdir(path.join(tmpRoot, "src", "contracts", "hero-banner"), {
      recursive: true,
    });

    await writeFile(
      path.join(tmpRoot, "src", "contracts", "1.2.3", "manifest.json"),
      JSON.stringify({
        version: "1.2.3",
        schemaVersion: "0.0.1",
        components: {},
      }),
      "utf8",
    );

    const errors = await validateContracts({
      rootDir: tmpRoot,
      compareRef: undefined,
      enforceGitImmutability: false,
    });

    assert.equal(errors.length, 1);
    assert.match(errors[0], /src\/contracts\/hero-banner/);
    assert.match(
      errors[0],
      /Example: src\/contracts\/1\.2\.3\/components\/<component-name>\/contract\.json/,
    );
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});

test("invalid release folder error falls back to 0.0.1 when no valid release exists", async () => {
  const tmpRoot = await mkdtemp(
    path.join(os.tmpdir(), "sloth-contract-policy-"),
  );

  try {
    await mkdir(path.join(tmpRoot, "src", "contracts", "hero-banner"), {
      recursive: true,
    });

    const errors = await validateContracts({
      rootDir: tmpRoot,
      compareRef: undefined,
      enforceGitImmutability: false,
    });

    assert.equal(errors.length, 1);
    assert.match(
      errors[0],
      /Example: src\/contracts\/0\.0\.1\/components\/<component-name>\/contract\.json/,
    );
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});
