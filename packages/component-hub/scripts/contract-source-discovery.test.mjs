import test from "node:test";
import assert from "node:assert/strict";
import { mkdtemp, mkdir, rm, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";

import { discoverComponentSources } from "./contract-source-discovery.mjs";

test("discoverComponentSources returns empty list when source directory is missing", async () => {
  const tmpRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-source-scan-"));

  try {
    const discovered = await discoverComponentSources({ rootDir: tmpRoot });
    assert.deepEqual(discovered, []);
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});

test("discoverComponentSources returns sorted component contracts", async () => {
  const tmpRoot = await mkdtemp(path.join(os.tmpdir(), "sloth-source-scan-"));

  try {
    const componentsRoot = path.join(tmpRoot, "src", "contracts", "components");
    await mkdir(path.join(componentsRoot, "hero-banner"), { recursive: true });
    await mkdir(path.join(componentsRoot, "article-teaser"), {
      recursive: true,
    });

    await writeFile(
      path.join(componentsRoot, "hero-banner", "contract.json"),
      JSON.stringify({ name: "hero-banner", version: "0.0.2" }),
      "utf8",
    );
    await writeFile(
      path.join(componentsRoot, "article-teaser", "contract.json"),
      JSON.stringify({ name: "article-teaser", version: "0.0.2" }),
      "utf8",
    );

    const discovered = await discoverComponentSources({ rootDir: tmpRoot });

    assert.equal(discovered.length, 2);
    assert.equal(discovered[0].componentName, "article-teaser");
    assert.equal(discovered[1].componentName, "hero-banner");
    assert.equal(discovered[0].contract.name, "article-teaser");
    assert.equal(discovered[1].contract.version, "0.0.2");
  } finally {
    await rm(tmpRoot, { recursive: true, force: true });
  }
});
