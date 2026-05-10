import test from "node:test";
import assert from "node:assert/strict";
import { mkdtemp, mkdir, readFile, rm, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";

import {
  createReleaseFrom,
  syncReleaseManifest,
  createReleaseLedgerEntry,
  verifyReleaseIntegrity,
  materializeRelease,
} from "./contract-release-workflow.mjs";

async function writeJson(filePath, payload) {
  await writeFile(filePath, `${JSON.stringify(payload, null, 2)}\n`, "utf8");
}

test("syncReleaseManifest recomputes component hashes", async () => {
  const tmpRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-release-sync-"));

  try {
    const releaseDir = path.join(tmpRoot, "src", "contracts", "0.0.1");
    await mkdir(path.join(releaseDir, "components", "hero-banner"), {
      recursive: true,
    });

    await writeJson(path.join(releaseDir, "manifest.json"), {
      version: "0.0.1",
      schemaVersion: "0.0.1",
      components: {
        "hero-banner": {
          contractPath: "./components/hero-banner/contract.json",
          contentHash: "stale",
        },
      },
    });

    await writeJson(
      path.join(releaseDir, "components", "hero-banner", "contract.json"),
      {
        name: "hero-banner",
        label: "Hero Banner",
        kind: "section",
        version: "0.0.1",
        schemaVersion: "0.0.1",
        dataset: [{ key: "headline", label: "Headline", type: "string" }],
        renderMeta: { rendererKey: "hero-banner" },
      },
    );

    const nextManifest = await syncReleaseManifest({
      rootDir: tmpRoot,
      releaseVersion: "0.0.1",
    });

    assert.equal(nextManifest.version, "0.0.1");
    assert.equal(nextManifest.schemaVersion, "0.0.1");
    assert.notEqual(
      nextManifest.components["hero-banner"].contentHash,
      "stale",
    );
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});

test("createReleaseFrom clones release and updates contract versions", async () => {
  const tmpRoot = await mkdtemp(
    path.join(os.tmpdir(), "sloth-release-create-"),
  );

  try {
    const sourceDir = path.join(tmpRoot, "src", "contracts", "0.0.1");
    await mkdir(path.join(sourceDir, "components", "hero-banner"), {
      recursive: true,
    });

    await writeJson(path.join(sourceDir, "manifest.json"), {
      version: "0.0.1",
      schemaVersion: "0.0.1",
      components: {
        "hero-banner": {
          contractPath: "./components/hero-banner/contract.json",
          contentHash: "hash",
        },
      },
    });

    await writeJson(
      path.join(sourceDir, "components", "hero-banner", "contract.json"),
      {
        name: "hero-banner",
        label: "Hero Banner",
        kind: "section",
        version: "0.0.1",
        schemaVersion: "0.0.1",
        dataset: [{ key: "headline", label: "Headline", type: "string" }],
        renderMeta: { rendererKey: "hero-banner" },
      },
    );

    await createReleaseFrom({
      rootDir: tmpRoot,
      fromVersion: "0.0.1",
      toVersion: "0.0.2",
      deprecateFrom: true,
      referenceDate: new Date("2026-05-10T00:00:00.000Z"),
    });

    const nextContractRaw = await readFile(
      path.join(
        tmpRoot,
        "src",
        "contracts",
        "0.0.2",
        "components",
        "hero-banner",
        "contract.json",
      ),
      "utf8",
    );
    const nextContract = JSON.parse(nextContractRaw);
    assert.equal(nextContract.version, "0.0.2");

    const fromManifestRaw = await readFile(
      path.join(tmpRoot, "src", "contracts", "0.0.1", "manifest.json"),
      "utf8",
    );
    const fromManifest = JSON.parse(fromManifestRaw);
    assert.equal(typeof fromManifest.deprecatedAt, "string");
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});

test("createReleaseLedgerEntry adds component sources to ledger", async () => {
  const tmpRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-ledger-create-"));

  try {
    // Create ledger file
    const releasesDir = path.join(tmpRoot, "src", "contracts", "releases");
    await mkdir(releasesDir, { recursive: true });
    await writeJson(path.join(releasesDir, "index.json"), {
      ledgerFormatVersion: "1",
      releases: [],
    });

    // Create component source
    const componentDir = path.join(
      tmpRoot,
      "src",
      "contracts",
      "components",
      "button",
    );
    await mkdir(componentDir, { recursive: true });
    await writeJson(path.join(componentDir, "contract.json"), {
      name: "button",
      version: "1.0.0",
      schemaVersion: "1",
      description: "A button component",
    });

    const result = await createReleaseLedgerEntry({
      rootDir: tmpRoot,
      releaseVersion: "1.0.0",
    });

    assert.equal(result.version, "1.0.0");
    assert.equal(Object.keys(result.components).length, 1);
    assert.ok(result.components.button);
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});

test("verifyReleaseIntegrity checks ledger against source contracts", async () => {
  const tmpRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-verify-"));

  try {
    // Create ledger with entry
    const releasesDir = path.join(tmpRoot, "src", "contracts", "releases");
    await mkdir(releasesDir, { recursive: true });

    // Create a component first
    const componentDir = path.join(
      tmpRoot,
      "src",
      "contracts",
      "components",
      "button",
    );
    await mkdir(componentDir, { recursive: true });
    const contractContent = JSON.stringify({
      name: "button",
      version: "1.0.0",
      schemaVersion: "1",
    });
    await writeFile(
      path.join(componentDir, "contract.json"),
      contractContent,
      "utf8",
    );

    // Create ledger with matching hash
    const { createHash } = await import("node:crypto");
    const hash = createHash("sha256").update(contractContent).digest("hex");

    await writeJson(path.join(releasesDir, "index.json"), {
      ledgerFormatVersion: "1",
      releases: [
        {
          version: "1.0.0",
          schemaVersion: "1",
          components: {
            button: {
              contractPath: "components/button/contract.json",
              contentHash: hash,
            },
          },
          createdAt: new Date().toISOString(),
        },
      ],
    });

    const result = await verifyReleaseIntegrity({
      rootDir: tmpRoot,
      releaseVersion: "1.0.0",
    });

    assert.equal(result.version, "1.0.0");
    assert.equal(result.status, "verified");
    assert.equal(result.componentCount, 1);
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});

test("materializeRelease copies contracts to dist", async () => {
  const tmpRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-materialize-"));

  try {
    // Create ledger with entry
    const releasesDir = path.join(tmpRoot, "src", "contracts", "releases");
    await mkdir(releasesDir, { recursive: true });

    // Create component source
    const componentDir = path.join(
      tmpRoot,
      "src",
      "contracts",
      "components",
      "button",
    );
    await mkdir(componentDir, { recursive: true });
    const contractContent = JSON.stringify({
      name: "button",
      version: "1.0.0",
      schemaVersion: "1",
    });
    await writeFile(
      path.join(componentDir, "contract.json"),
      contractContent,
      "utf8",
    );

    // Create ledger
    await writeJson(path.join(releasesDir, "index.json"), {
      ledgerFormatVersion: "1",
      releases: [
        {
          version: "1.0.0",
          schemaVersion: "1",
          components: {
            button: {
              contractPath: "components/button/contract.json",
              contentHash: "hash",
            },
          },
          createdAt: new Date().toISOString(),
        },
      ],
    });

    const result = await materializeRelease({
      rootDir: tmpRoot,
      releaseVersion: "1.0.0",
    });

    assert.equal(result.version, "1.0.0");
    assert.equal(result.componentCount, 1);

    const materializedContract = await readFile(
      path.join(result.outputDir, "components", "button", "contract.json"),
      "utf8",
    );
    assert.ok(materializedContract.includes("button"));
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});
