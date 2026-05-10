import { cp, mkdir, readdir, readFile, rm, writeFile } from "node:fs/promises";
import { createHash } from "node:crypto";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, "..");
const sourceDir = path.join(rootDir, "src", "contracts", "components");
const outputDir = path.join(rootDir, "dist", "registry", "contracts");

function computeContentHash(rawContents) {
  return createHash("sha256").update(rawContents).digest("hex");
}

async function readComponentContracts(componentsDir) {
  const entries = await readdir(componentsDir, { withFileTypes: true });
  const componentNames = entries
    .filter((entry) => entry.isDirectory())
    .map((entry) => entry.name)
    .sort((a, b) => a.localeCompare(b));

  if (componentNames.length === 0) {
    throw new Error(`No contracts found under ${componentsDir}`);
  }

  const contracts = [];
  for (const componentName of componentNames) {
    const contractPath = path.join(componentsDir, componentName, "contract.json");
    const contractRaw = await readFile(contractPath, "utf8");
    const contract = JSON.parse(contractRaw);

    contracts.push({
      componentName,
      contract,
      contentHash: computeContentHash(contractRaw),
    });
  }

  return contracts;
}

function inferReleaseMetadata(contracts) {
  const versions = new Set(contracts.map((item) => item.contract.version));
  const schemaVersions = new Set(
    contracts.map((item) => item.contract.schemaVersion),
  );

  if (versions.size !== 1) {
    throw new Error(
      `Contracts must share one release version. Found: ${[...versions].join(", ")}`,
    );
  }
  if (schemaVersions.size !== 1) {
    throw new Error(
      `Contracts must share one schemaVersion. Found: ${[
        ...schemaVersions,
      ].join(", ")}`,
    );
  }

  return {
    releaseVersion: [...versions][0],
    schemaVersion: [...schemaVersions][0],
  };
}

export async function buildRegistry({
  sourceComponentsDir = sourceDir,
  outputContractsDir = outputDir,
  registryRootDir = path.join(rootDir, "dist", "registry"),
} = {}) {
  const contracts = await readComponentContracts(sourceComponentsDir);
  const { releaseVersion, schemaVersion } = inferReleaseMetadata(contracts);

  await rm(registryRootDir, { recursive: true, force: true });
  const releaseOutputDir = path.join(outputContractsDir, releaseVersion);
  const componentsOutputDir = path.join(releaseOutputDir, "components");
  await mkdir(componentsOutputDir, { recursive: true });

  const manifestComponents = {};
  for (const { componentName, contentHash } of contracts) {
    const sourceComponentDir = path.join(sourceComponentsDir, componentName);
    const targetComponentDir = path.join(componentsOutputDir, componentName);
    await cp(sourceComponentDir, targetComponentDir, { recursive: true });

    manifestComponents[componentName] = {
      contractPath: `./components/${componentName}/contract.json`,
      contentHash,
    };
  }

  const releaseManifest = {
    version: releaseVersion,
    schemaVersion,
    components: manifestComponents,
  };

  await writeFile(
    path.join(releaseOutputDir, "manifest.json"),
    `${JSON.stringify(releaseManifest, null, 2)}\n`,
    "utf8",
  );

  const indexItems = [
    {
      version: releaseVersion,
      manifest: `/sloth/registry/contracts/${releaseVersion}/manifest.json`,
      components: Object.keys(manifestComponents).sort((a, b) =>
        a.localeCompare(b),
      ),
      deprecatedAt: null,
    },
  ];

  await writeFile(
    path.join(outputContractsDir, "index.json"),
    `${JSON.stringify({ registryFormatVersion: "1", items: indexItems }, null, 2)}\n`,
    "utf8",
  );

  await writeFile(
    path.join(registryRootDir, "index.json"),
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

  return { releaseVersion, schemaVersion, componentCount: contracts.length };
}

try {
  await buildRegistry();
} catch (error) {
  console.error(error);
  process.exitCode = 1;
}
