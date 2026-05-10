import test from "node:test";
import assert from "node:assert/strict";
import { mkdir, writeFile, rm } from "node:fs/promises";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { randomUUID } from "node:crypto";

import {
  validateLedgerRelease,
  computeContentHash,
} from "./contract-version-policy.mjs";

test("validateLedgerRelease rejects invalid release object", async () => {
  const errors = await validateLedgerRelease({
    rootDir: tmpdir(),
    release: null,
  });
  assert.ok(errors.some((e) => e.includes("Release entry must be an object")));
});

test("validateLedgerRelease rejects invalid version format", async () => {
  const errors = await validateLedgerRelease({
    rootDir: tmpdir(),
    release: {
      version: "invalid-version",
      schemaVersion: "1",
      components: {},
    },
  });
  assert.ok(errors.some((e) => e.includes("Invalid release version")));
});

test("validateLedgerRelease detects missing component contractPath", async () => {
  const errors = await validateLedgerRelease({
    rootDir: tmpdir(),
    release: {
      version: "1.0.0",
      schemaVersion: "1",
      components: {
        myComponent: {
          contentHash: "abc123",
        },
      },
    },
  });
  assert.ok(
    errors.some(
      (e) => e.includes("myComponent") && e.includes("missing contractPath"),
    ),
  );
});

test("validateLedgerRelease detects missing component contentHash", async () => {
  const errors = await validateLedgerRelease({
    rootDir: tmpdir(),
    release: {
      version: "1.0.0",
      schemaVersion: "1",
      components: {
        myComponent: {
          contractPath: "components/myComponent/contract.json",
        },
      },
    },
  });
  assert.ok(
    errors.some(
      (e) => e.includes("myComponent") && e.includes("missing contentHash"),
    ),
  );
});

test("validateLedgerRelease validates contract against ledger entry successfully", async () => {
  const testDir = await mkdir(join(tmpdir(), `contract-test-${randomUUID()}`), {
    recursive: true,
  });

  try {
    const contractDir = join(
      testDir,
      "src",
      "contracts",
      "components",
      "button",
    );
    await mkdir(contractDir, { recursive: true });

    const contractContent = JSON.stringify({
      name: "button",
      version: "1.0.0",
      schemaVersion: "1",
      description: "A button component",
    });
    await writeFile(
      join(contractDir, "contract.json"),
      contractContent,
      "utf8",
    );

    const contentHash = computeContentHash(contractContent);

    const errors = await validateLedgerRelease({
      rootDir: testDir,
      release: {
        version: "1.0.0",
        schemaVersion: "1",
        components: {
          button: {
            contractPath: "components/button/contract.json",
            contentHash,
          },
        },
      },
    });

    assert.equal(errors.length, 0);
  } finally {
    await rm(testDir, { recursive: true, force: true });
  }
});

test("validateLedgerRelease detects content hash mismatch", async () => {
  const testDir = join(tmpdir(), `contract-test-${randomUUID()}`);
  await mkdir(testDir, { recursive: true });

  try {
    const contractDir = join(
      testDir,
      "src",
      "contracts",
      "components",
      "button",
    );
    await mkdir(contractDir, { recursive: true });

    const contractContent = JSON.stringify({
      name: "button",
      version: "1.0.0",
      schemaVersion: "1",
      description: "A button component",
    });
    await writeFile(
      join(contractDir, "contract.json"),
      contractContent,
      "utf8",
    );

    const errors = await validateLedgerRelease({
      rootDir: testDir,
      release: {
        version: "1.0.0",
        schemaVersion: "1",
        components: {
          button: {
            contractPath: "components/button/contract.json",
            contentHash: "wrong-hash",
          },
        },
      },
    });

    assert.ok(errors.some((e) => e.includes("contentHash mismatch")));
  } finally {
    await rm(testDir, { recursive: true, force: true });
  }
});

test("validateLedgerRelease detects contract version mismatch", async () => {
  const testDir = join(tmpdir(), `contract-test-${randomUUID()}`);
  await mkdir(testDir, { recursive: true });

  try {
    const contractDir = join(
      testDir,
      "src",
      "contracts",
      "components",
      "button",
    );
    await mkdir(contractDir, { recursive: true });

    const contractContent = JSON.stringify({
      name: "button",
      version: "1.0.0",
      schemaVersion: "1",
      description: "A button component",
    });
    await writeFile(
      join(contractDir, "contract.json"),
      contractContent,
      "utf8",
    );

    const contentHash = computeContentHash(contractContent);

    const errors = await validateLedgerRelease({
      rootDir: testDir,
      release: {
        version: "2.0.0",
        schemaVersion: "1",
        components: {
          button: {
            contractPath: "components/button/contract.json",
            contentHash,
          },
        },
      },
    });

    assert.ok(errors.some((e) => e.includes("Contract version mismatch")));
  } finally {
    await rm(testDir, { recursive: true, force: true });
  }
});

test("validateLedgerRelease detects missing contract file", async () => {
  const testDir = tmpdir();

  const errors = await validateLedgerRelease({
    rootDir: testDir,
    release: {
      version: "1.0.0",
      schemaVersion: "1",
      components: {
        button: {
          contractPath: "components/button/contract.json",
          contentHash: "abc123",
        },
      },
    },
  });

  assert.ok(errors.some((e) => e.includes("Could not read contract")));
});
