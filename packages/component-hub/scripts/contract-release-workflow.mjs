import { cp, mkdir, readdir, readFile, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";

import {
  MIN_DEPRECATION_MONTHS,
  addMonths,
  computeContentHash,
  isReleaseVersion,
} from "./contract-version-policy.mjs";
import {
  readReleaseLedger,
  writeReleaseLedger,
  validateReleaseLedger,
} from "./contract-release-ledger.mjs";
import { discoverComponentSources } from "./contract-source-discovery.mjs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const defaultRootDir = path.resolve(__dirname, "..");
const SEMVER_PATTERN =
  /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$/;

function isSemver(value) {
  return typeof value === "string" && SEMVER_PATTERN.test(value);
}

function parseArgs(argv) {
  const [command, ...rest] = argv;
  const flags = {};

  for (let i = 0; i < rest.length; i += 1) {
    const token = rest[i];
    if (!token.startsWith("--")) {
      continue;
    }

    const key = token.slice(2);
    const value = rest[i + 1];
    if (!value || value.startsWith("--")) {
      flags[key] = "true";
      continue;
    }

    flags[key] = value;
    i += 1;
  }

  return { command, flags };
}

async function readJson(filePath) {
  const raw = await readFile(filePath, "utf8");
  return JSON.parse(raw);
}

async function writeJson(filePath, payload) {
  await writeFile(filePath, `${JSON.stringify(payload, null, 2)}\n`, "utf8");
}

function buildDeprecationTimestamp(referenceDate = new Date()) {
  const minDate = addMonths(referenceDate, MIN_DEPRECATION_MONTHS);
  const paddedDate = new Date(minDate);
  paddedDate.setUTCDate(paddedDate.getUTCDate() + 1);
  return paddedDate.toISOString();
}

async function listComponentNames(componentsDir) {
  return (await readdir(componentsDir, { withFileTypes: true }))
    .filter((entry) => entry.isDirectory())
    .map((entry) => entry.name)
    .sort((a, b) => a.localeCompare(b));
}

export async function syncReleaseManifest({
  rootDir = defaultRootDir,
  releaseVersion,
}) {
  if (!isSemver(releaseVersion)) {
    throw new Error(`Invalid release version: ${releaseVersion}`);
  }

  const releaseDir = path.join(rootDir, "src", "contracts", releaseVersion);
  const componentsDir = path.join(releaseDir, "components");
  const manifestPath = path.join(releaseDir, "manifest.json");
  const manifest = await readJson(manifestPath);
  const componentNames = await listComponentNames(componentsDir);

  const components = {};
  let schemaVersion = manifest.schemaVersion;

  for (const componentName of componentNames) {
    const contractPath = path.join(
      componentsDir,
      componentName,
      "contract.json",
    );
    const contractRaw = await readFile(contractPath, "utf8");
    const contract = JSON.parse(contractRaw);

    if (contract.version !== releaseVersion) {
      throw new Error(
        `Contract ${componentName} has version ${contract.version} but release is ${releaseVersion}`,
      );
    }

    schemaVersion = schemaVersion ?? contract.schemaVersion;
    if (schemaVersion !== contract.schemaVersion) {
      throw new Error(
        `Contract ${componentName} schemaVersion mismatch: expected ${schemaVersion}, got ${contract.schemaVersion}`,
      );
    }

    components[componentName] = {
      contractPath: `./components/${componentName}/contract.json`,
      contentHash: computeContentHash(contractRaw),
    };
  }

  const nextManifest = {
    version: releaseVersion,
    schemaVersion: schemaVersion ?? "0.0.1",
    ...(typeof manifest.deprecatedAt === "string"
      ? { deprecatedAt: manifest.deprecatedAt }
      : {}),
    components,
  };

  await writeJson(manifestPath, nextManifest);
  return nextManifest;
}

export async function createReleaseFrom({
  rootDir = defaultRootDir,
  fromVersion,
  toVersion,
  deprecateFrom = true,
  referenceDate = new Date(),
}) {
  if (!isSemver(fromVersion) || !isSemver(toVersion)) {
    throw new Error(
      `Invalid semver values: from=${fromVersion} to=${toVersion}`,
    );
  }

  const sourceReleaseDir = path.join(rootDir, "src", "contracts", fromVersion);
  const targetReleaseDir = path.join(rootDir, "src", "contracts", toVersion);
  await mkdir(path.dirname(targetReleaseDir), { recursive: true });
  await cp(sourceReleaseDir, targetReleaseDir, { recursive: true });

  const componentNames = await listComponentNames(
    path.join(targetReleaseDir, "components"),
  );

  for (const componentName of componentNames) {
    const contractPath = path.join(
      targetReleaseDir,
      "components",
      componentName,
      "contract.json",
    );
    const contract = await readJson(contractPath);
    contract.version = toVersion;
    await writeJson(contractPath, contract);
  }

  const fromManifestPath = path.join(sourceReleaseDir, "manifest.json");
  if (deprecateFrom) {
    const fromManifest = await readJson(fromManifestPath);
    if (typeof fromManifest.deprecatedAt !== "string") {
      fromManifest.deprecatedAt = buildDeprecationTimestamp(referenceDate);
      await writeJson(fromManifestPath, fromManifest);
    }
  }

  await syncReleaseManifest({ rootDir, releaseVersion: toVersion });
}

// New ledger-driven workflow functions

export async function createReleaseLedgerEntry({
  rootDir = defaultRootDir,
  releaseVersion,
  sourceGitRef = null,
  referenceDate = new Date(),
}) {
  if (!isSemver(releaseVersion)) {
    throw new Error(`Invalid release version: ${releaseVersion}`);
  }

  // Discover component sources
  const sources = await discoverComponentSources({ rootDir });

  const components = {};
  let schemaVersion = null;

  for (const source of sources) {
    const { componentName, contract, contractRaw } = source;
    const contentHash = computeContentHash(contractRaw);

    if (contract.version !== releaseVersion) {
      throw new Error(
        `Component ${componentName} has version ${contract.version} but release is ${releaseVersion}`,
      );
    }

    schemaVersion = schemaVersion ?? contract.schemaVersion;
    if (schemaVersion !== contract.schemaVersion) {
      throw new Error(
        `Component ${componentName} schemaVersion mismatch: expected ${schemaVersion}, got ${contract.schemaVersion}`,
      );
    }

    components[componentName] = {
      contractPath: `components/${componentName}/contract.json`,
      contentHash,
    };
  }

  // Read current ledger and add entry
  const ledger = await readReleaseLedger({ rootDir });

  // Check for duplicate version
  const existingIndex = ledger.releases.findIndex(
    (r) => r.version === releaseVersion,
  );
  if (existingIndex >= 0) {
    throw new Error(`Release ${releaseVersion} already exists in ledger`);
  }

  const newRelease = {
    version: releaseVersion,
    schemaVersion: schemaVersion ?? "0.0.1",
    components,
    createdAt: new Date().toISOString(),
    ...(sourceGitRef ? { sourceGitRef } : {}),
  };

  ledger.releases.push(newRelease);
  ledger.releases.sort((a, b) => b.version.localeCompare(a.version));

  await writeReleaseLedger({ rootDir, ledger });

  return newRelease;
}

export async function verifyReleaseIntegrity({
  rootDir = defaultRootDir,
  releaseVersion,
}) {
  if (!isSemver(releaseVersion)) {
    throw new Error(`Invalid release version: ${releaseVersion}`);
  }

  const ledger = await readReleaseLedger({ rootDir });
  const releaseEntry = ledger.releases.find(
    (r) => r.version === releaseVersion,
  );

  if (!releaseEntry) {
    throw new Error(`Release ${releaseVersion} not found in ledger`);
  }

  const errors = [];
  const { components } = releaseEntry;

  for (const [componentName, componentInfo] of Object.entries(components)) {
    const sourceDir = path.join(rootDir, "src", "contracts");
    const contractPath = path.join(
      sourceDir,
      "components",
      componentName,
      "contract.json",
    );

    try {
      const contractRaw = await readFile(contractPath, "utf8");
      const contract = JSON.parse(contractRaw);
      const actualHash = computeContentHash(contractRaw);

      if (actualHash !== componentInfo.contentHash) {
        errors.push(
          `Component ${componentName}: hash mismatch (expected ${componentInfo.contentHash}, got ${actualHash})`,
        );
      }

      if (contract.version !== releaseVersion) {
        errors.push(
          `Component ${componentName}: version mismatch (expected ${releaseVersion}, got ${contract.version})`,
        );
      }
    } catch (error) {
      errors.push(
        `Component ${componentName}: could not verify (${error.message})`,
      );
    }
  }

  if (errors.length > 0) {
    throw new Error(
      `Release ${releaseVersion} verification failed:\n${errors.join("\n")}`,
    );
  }

  return {
    version: releaseVersion,
    status: "verified",
    componentCount: Object.keys(components).length,
  };
}

export async function materializeRelease({
  rootDir = defaultRootDir,
  releaseVersion,
}) {
  if (!isSemver(releaseVersion)) {
    throw new Error(`Invalid release version: ${releaseVersion}`);
  }

  const ledger = await readReleaseLedger({ rootDir });
  const releaseEntry = ledger.releases.find(
    (r) => r.version === releaseVersion,
  );

  if (!releaseEntry) {
    throw new Error(`Release ${releaseVersion} not found in ledger`);
  }

  const outputDir = path.join(
    rootDir,
    "dist",
    "registry",
    "contracts",
    releaseVersion,
  );
  await mkdir(path.dirname(outputDir), { recursive: true });

  // Copy component contracts to dist
  const sourceDir = path.join(rootDir, "src", "contracts");
  const componentsOutDir = path.join(outputDir, "components");
  await mkdir(componentsOutDir, { recursive: true });

  for (const componentName of Object.keys(releaseEntry.components)) {
    const sourceComponentDir = path.join(
      sourceDir,
      "components",
      componentName,
    );
    const targetComponentDir = path.join(componentsOutDir, componentName);
    await mkdir(path.dirname(targetComponentDir), { recursive: true });
    await cp(sourceComponentDir, targetComponentDir, { recursive: true });
  }

  // Generate manifest for release
  const manifest = {
    version: releaseEntry.version,
    schemaVersion: releaseEntry.schemaVersion,
    components: releaseEntry.components,
  };

  await writeJson(path.join(outputDir, "manifest.json"), manifest);

  return {
    version: releaseVersion,
    outputDir,
    componentCount: Object.keys(releaseEntry.components).length,
  };
}

function printUsage() {
  console.log("Usage:");
  console.log(
    "  node contract-release-workflow.mjs release create --version <x.y.z> [--git-ref <ref>]",
  );
  console.log(
    "  node contract-release-workflow.mjs release verify --version <x.y.z>",
  );
  console.log(
    "  node contract-release-workflow.mjs release materialize --version <x.y.z>",
  );
  console.log("");
  console.log("Legacy commands (for backward compatibility):");
  console.log("  node contract-release-workflow.mjs sync --release <x.y.z>");
  console.log(
    "  node contract-release-workflow.mjs create --from <x.y.z> --to <x.y.z> [--deprecate-from true|false]",
  );
}

export async function main(argv = process.argv.slice(2)) {
  const normalizedArgv = argv[0] === "--" ? argv.slice(1) : argv;
  const { command, flags } = parseArgs(normalizedArgv);

  if (!command || command === "help" || command === "--help") {
    printUsage();
    return;
  }

  // New release workflow commands
  if (command === "release") {
    const subcommand = normalizedArgv[1];
    const subcommandFlags = parseArgs(normalizedArgv.slice(2)).flags;

    if (!subcommand || subcommand === "help") {
      printUsage();
      return;
    }

    if (subcommand === "create") {
      const releaseVersion = subcommandFlags.version;
      if (!releaseVersion) {
        throw new Error("Missing --version for release create command.");
      }

      const result = await createReleaseLedgerEntry({
        releaseVersion,
        sourceGitRef: subcommandFlags["git-ref"] || null,
      });
      console.log(
        `Created release ${result.version} in ledger with ${Object.keys(result.components).length} components.`,
      );
      return;
    }

    if (subcommand === "verify") {
      const releaseVersion = subcommandFlags.version;
      if (!releaseVersion) {
        throw new Error("Missing --version for release verify command.");
      }

      const result = await verifyReleaseIntegrity({ releaseVersion });
      console.log(
        `✓ Release ${result.version} verified: ${result.componentCount} components match ledger.`,
      );
      return;
    }

    if (subcommand === "materialize") {
      const releaseVersion = subcommandFlags.version;
      if (!releaseVersion) {
        throw new Error("Missing --version for release materialize command.");
      }

      const result = await materializeRelease({ releaseVersion });
      console.log(
        `✓ Materialized release ${result.version} to ${result.outputDir} (${result.componentCount} components).`,
      );
      return;
    }

    throw new Error(`Unknown release subcommand: ${subcommand}`);
  }

  // Legacy commands (backward compatibility)
  if (command === "sync") {
    const releaseVersion = flags.release;
    if (!releaseVersion) {
      throw new Error("Missing --release for sync command.");
    }
    await syncReleaseManifest({ releaseVersion });
    console.log(`Synced manifest hashes for release ${releaseVersion}.`);
    return;
  }

  if (command === "create") {
    const fromVersion = flags.from;
    const toVersion = flags.to;
    if (!fromVersion || !toVersion) {
      throw new Error("Missing --from or --to for create command.");
    }

    await createReleaseFrom({
      fromVersion,
      toVersion,
      deprecateFrom: flags["deprecate-from"] !== "false",
    });
    console.log(`Created release ${toVersion} from ${fromVersion}.`);
    return;
  }

  throw new Error(`Unknown command: ${command}`);
}

if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  main().catch((error) => {
    console.error(error.message || error);
    process.exitCode = 1;
  });
}
