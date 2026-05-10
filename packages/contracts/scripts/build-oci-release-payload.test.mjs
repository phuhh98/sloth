import test from "node:test";
import assert from "node:assert/strict";
import { mkdir, writeFile, readFile, mkdtemp } from "node:fs/promises";
import os from "node:os";
import path from "node:path";

import { buildOCIReleasePayload } from "./build-oci-release-payload.mjs";

async function writeJson(filePath, payload) {
  await mkdir(path.dirname(filePath), { recursive: true });
  await writeFile(filePath, `${JSON.stringify(payload, null, 2)}\n`, "utf8");
}

test("buildOCIReleasePayload materializes release payload from manifest/contracts", async () => {
  const tmpRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-oci-payload-"));
  const contractsRoot = path.join(tmpRoot, "dist", "registry", "contracts");

  const releaseRoot = path.join(contractsRoot, "0.0.7");
  await writeJson(path.join(releaseRoot, "manifest.json"), {
    version: "0.0.7",
    schemaVersion: "0.0.1",
    components: {
      "hero-banner": {
        contractPath: "./components/hero-banner/contract.json",
        contentHash: "hash-hero",
      },
      "faq-list": {
        contractPath: "./components/faq-list/contract.json",
        contentHash: "hash-faq",
      },
    },
  });

  await writeJson(path.join(releaseRoot, "components", "hero-banner", "contract.json"), {
    name: "hero-banner",
    label: "Hero Banner",
    version: "0.0.7",
    schemaVersion: "0.0.1",
    kind: "section",
  });

  await writeJson(path.join(releaseRoot, "components", "faq-list", "contract.json"), {
    name: "faq-list",
    label: "FAQ List",
    version: "0.0.7",
    schemaVersion: "0.0.1",
    kind: "section",
  });

  const outPath = path.join(tmpRoot, "release.json");
  const payload = await buildOCIReleasePayload({
    contractsRoot,
    version: "0.0.7",
    outputPath: outPath,
  });

  assert.equal(payload.version, "0.0.7");
  assert.equal(payload.schemaVersion, "0.0.1");
  assert.equal(payload.contracts.length, 2);
  assert.deepEqual(
    payload.contracts.map((entry) => entry.name),
    ["faq-list", "hero-banner"],
  );

  const written = JSON.parse(await readFile(outPath, "utf8"));
  assert.equal(written.contracts[0].contentHash, "hash-faq");
  assert.equal(written.contracts[1].contentHash, "hash-hero");
});
