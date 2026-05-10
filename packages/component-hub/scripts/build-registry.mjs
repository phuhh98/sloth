import { cp, mkdir, readdir, readFile, rm, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import {
  isReleaseVersion,
  validateContracts,
} from "./contract-version-policy.mjs";
import {
  readReleaseLedger,
  validateReleaseLedger,
} from "./contract-release-ledger.mjs";
import { materializeRelease } from "./contract-release-workflow.mjs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, "..");
const sourceDir = path.join(rootDir, "src", "contracts");
const outputDir = path.join(rootDir, "dist", "registry", "contracts");

function semverCompareDescending(a, b) {
  const pa = a.split(".").map((part) => Number.parseInt(part, 10));
  const pb = b.split(".").map((part) => Number.parseInt(part, 10));
  for (let i = 0; i < 3; i += 1) {
    const va = Number.isNaN(pa[i]) ? 0 : pa[i];
    const vb = Number.isNaN(pb[i]) ? 0 : pb[i];
    if (va !== vb) {
      return vb - va;
    }
  }
  return 0;
}

async function buildRegistry() {
  const validationErrors = await validateContracts({
    rootDir,
    compareRef: undefined,
    enforceGitImmutability: false,
  });

  if (validationErrors.length > 0) {
    throw new Error(validationErrors.join("\n"));
  }

  await rm(path.join(rootDir, "dist", "registry"), {
    recursive: true,
    force: true,
  });
  await mkdir(outputDir, { recursive: true });

  // Try to read ledger and materialize releases from it
  const ledger = await readReleaseLedger({ rootDir });
  const ledgerValidationErrors = validateReleaseLedger(ledger);

  const generatedVersions = new Set();
  const indexItems = [];

  // Materialize releases from ledger (primary source)
  for (const release of ledger.releases) {
    try {
      await materializeRelease({
        rootDir,
        releaseVersion: release.version,
      });
      generatedVersions.add(release.version);
    } catch (error) {
      // Fall back to legacy folder-based approach if materialization fails
      console.warn(
        `Warning: could not materialize release ${release.version} from ledger: ${error.message}`,
      );
    }
  }

  // Process materialized and legacy releases to build index
  const releaseVersions = (await readdir(sourceDir, { withFileTypes: true }))
    .filter((entry) => entry.isDirectory())
    .filter((entry) => isReleaseVersion(entry.name))
    .map((entry) => entry.name);

  // Also check dist output for materialized releases
  const distReleaseVersions = (
    await readdir(outputDir, { withFileTypes: true })
  )
    .filter((entry) => entry.isDirectory())
    .filter((entry) => isReleaseVersion(entry.name))
    .map((entry) => entry.name);

  const allVersions = new Set([...releaseVersions, ...distReleaseVersions]);

  for (const releaseVersion of Array.from(allVersions).sort(
    semverCompareDescending,
  )) {
    const releaseOutputDir = path.join(outputDir, releaseVersion);
    const manifestPath = path.join(releaseOutputDir, "manifest.json");

    // Try to read from already-materialized output
    let manifest;
    try {
      const manifestRaw = await readFile(manifestPath, "utf8");
      manifest = JSON.parse(manifestRaw);
    } catch {
      // Fall back to source folder if not yet materialized
      const releaseSourceDir = path.join(sourceDir, releaseVersion);
      try {
        const sourceManifestPath = path.join(releaseSourceDir, "manifest.json");
        const manifestRaw = await readFile(sourceManifestPath, "utf8");
        manifest = JSON.parse(manifestRaw);

        // Copy source folder to output if not materialized from ledger
        await mkdir(releaseOutputDir, { recursive: true });
        await cp(releaseSourceDir, releaseOutputDir, { recursive: true });
      } catch (error) {
        console.warn(
          `Warning: could not read manifest for release ${releaseVersion}`,
        );
        continue;
      }
    }

    const components = Object.keys(manifest.components ?? {}).sort((a, b) =>
      a.localeCompare(b),
    );

    indexItems.push({
      version: releaseVersion,
      manifest: `/sloth/registry/contracts/${releaseVersion}/manifest.json`,
      components,
      deprecatedAt: manifest.deprecatedAt ?? null,
    });
  }

  await writeFile(
    path.join(outputDir, "index.json"),
    `${JSON.stringify({ registryFormatVersion: "1", items: indexItems }, null, 2)}\n`,
    "utf8",
  );

  await writeFile(
    path.join(rootDir, "dist", "registry", "index.json"),
    `${JSON.stringify(
      {
        registryFormatVersion: "1",
        contractsIndex: "/sloth/registry/contracts/index.json",
        themesIndex: "/sloth/registry/themes/index.json",
      },
      null,
      2,
    )}\n`,
    "utf8",
  );
}

try {
  await buildRegistry();
} catch (error) {
  console.error(error);
  process.exitCode = 1;
}
