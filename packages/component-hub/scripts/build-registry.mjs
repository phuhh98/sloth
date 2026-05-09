import { cp, mkdir, readdir, readFile, rm, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import { validateContracts } from "./contract-version-policy.mjs";

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

  const releaseVersions = (await readdir(sourceDir, { withFileTypes: true }))
    .filter((entry) => entry.isDirectory())
    .map((entry) => entry.name);

  const indexItems = [];

  for (const releaseVersion of releaseVersions.sort(semverCompareDescending)) {
    const releaseSourceDir = path.join(sourceDir, releaseVersion);
    const manifestPath = path.join(releaseSourceDir, "manifest.json");
    const manifestRaw = await readFile(manifestPath, "utf8");
    const manifest = JSON.parse(manifestRaw);
    const components = Object.keys(manifest.components ?? {}).sort();

    const releaseOutputDir = path.join(outputDir, releaseVersion);
    await mkdir(releaseOutputDir, { recursive: true });
    await cp(releaseSourceDir, releaseOutputDir, { recursive: true });

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
        packsIndex: "/sloth/registry/packs/index.json",
      },
      null,
      2,
    )}\n`,
    "utf8",
  );
}

buildRegistry().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
