import test from "node:test";
import assert from "node:assert/strict";
import { mkdtemp, rm } from "node:fs/promises";
import os from "node:os";
import path from "node:path";

import {
  RELEASE_LEDGER_FORMAT_VERSION,
  readReleaseLedger,
  validateReleaseLedger,
  writeReleaseLedger,
} from "./contract-release-ledger.mjs";

test("readReleaseLedger returns empty ledger when file is missing", async () => {
  const tmpRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-ledger-read-"));

  try {
    const ledger = await readReleaseLedger({ rootDir: tmpRoot });
    assert.equal(ledger.ledgerFormatVersion, RELEASE_LEDGER_FORMAT_VERSION);
    assert.deepEqual(ledger.releases, []);
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});

test("writeReleaseLedger persists and validates valid release ledger", async () => {
  const tmpRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-ledger-write-"));

  try {
    const ledger = {
      ledgerFormatVersion: RELEASE_LEDGER_FORMAT_VERSION,
      releases: [
        {
          version: "0.0.1",
          schemaVersion: "0.0.1",
          createdAt: "2026-05-10T00:00:00.000Z",
          sourceGitRef: null,
          components: {
            "hero-banner": {
              contractPath: "./components/hero-banner/contract.json",
              contentHash: "hash",
            },
          },
        },
      ],
    };

    await writeReleaseLedger({ rootDir: tmpRoot, ledger });
    const nextLedger = await readReleaseLedger({ rootDir: tmpRoot });
    assert.deepEqual(nextLedger, ledger);
    assert.deepEqual(validateReleaseLedger(nextLedger), []);
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});

test("validateReleaseLedger rejects duplicate versions", () => {
  const errors = validateReleaseLedger({
    ledgerFormatVersion: RELEASE_LEDGER_FORMAT_VERSION,
    releases: [
      {
        version: "0.0.1",
        createdAt: "2026-05-10T00:00:00.000Z",
        components: {},
      },
      {
        version: "0.0.1",
        createdAt: "2026-05-10T00:00:00.000Z",
        components: {},
      },
    ],
  });

  assert.equal(errors.length, 1);
  assert.match(errors[0], /Duplicate release version/);
});
