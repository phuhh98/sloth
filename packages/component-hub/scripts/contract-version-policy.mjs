import { createHash } from "node:crypto";
import { readdir, readFile } from "node:fs/promises";
import path from "node:path";
import { execFile } from "node:child_process";
import { promisify } from "node:util";

import {
  readReleaseLedger,
  validateReleaseLedger,
} from "./contract-release-ledger.mjs";

const execFileAsync = promisify(execFile);

export const MIN_DEPRECATION_MONTHS = 6;
const RELEASE_VERSION_PATTERN =
  /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$/;
const RESERVED_SOURCE_DIRS = new Set(["components", "releases"]);

export function isReleaseVersion(value) {
  return RELEASE_VERSION_PATTERN.test(value);
}

export function semverCompareDescending(a, b) {
  const pa = a.split(".").map((part) => Number.parseInt(part, 10));
  const pb = b.split(".").map((part) => Number.parseInt(part, 10));

  for (let index = 0; index < 3; index += 1) {
    const va = Number.isNaN(pa[index]) ? 0 : pa[index];
    const vb = Number.isNaN(pb[index]) ? 0 : pb[index];

    if (va !== vb) {
      return vb - va;
    }
  }

  return 0;
}

export function computeContentHash(rawContents) {
  return createHash("sha256").update(rawContents).digest("hex");
}

export function addMonths(date, months) {
  const nextDate = new Date(date);
  nextDate.setUTCMonth(nextDate.getUTCMonth() + months);
  return nextDate;
}

export function stripMutableManifestFields(manifest) {
  const { deprecatedAt, ...stableFields } = manifest;
  return stableFields;
}

export function validateReleaseSet(entries, referenceDate = new Date()) {
  if (entries.length === 0) {
    return [];
  }

  const latestVersion = [...entries].sort((left, right) =>
    semverCompareDescending(left.version, right.version),
  )[0].version;
  const minimumDeprecationDate = addMonths(
    referenceDate,
    MIN_DEPRECATION_MONTHS,
  );
  const errors = [];

  for (const entry of entries) {
    if (entry.version === latestVersion) {
      continue;
    }

    if (typeof entry.manifest.deprecatedAt !== "string") {
      errors.push(
        `Contract release ${entry.version} must declare deprecatedAt because a newer release exists (${latestVersion}).`,
      );
      continue;
    }

    const deprecatedAt = new Date(entry.manifest.deprecatedAt);
    if (Number.isNaN(deprecatedAt.getTime())) {
      errors.push(
        `Contract release ${entry.version} has an invalid deprecatedAt timestamp.`,
      );
      continue;
    }

    if (deprecatedAt < minimumDeprecationDate) {
      errors.push(
        `Contract release ${entry.version} must keep deprecatedAt at least ${MIN_DEPRECATION_MONTHS} months in the future.`,
      );
    }
  }

  return errors;
}

async function gitShow(rootDir, ref, relativePath) {
  if (!ref) {
    return null;
  }

  try {
    const { stdout } = await execFileAsync(
      "git",
      ["show", `${ref}:${relativePath}`],
      { cwd: rootDir },
    );
    return stdout;
  } catch (error) {
    return null;
  }
}

export async function validateLedgerRelease({
  rootDir,
  release,
  referenceDate = new Date(),
}) {
  const errors = [];

  if (!release || typeof release !== "object") {
    return ["Release entry must be an object."];
  }

  const { version, schemaVersion, components = {} } = release;

  if (!isReleaseVersion(version)) {
    errors.push(`Invalid release version in ledger: ${version}.`);
    return errors;
  }

  const sourceDir = path.join(rootDir, "src", "contracts");

  for (const [componentName, componentInfo] of Object.entries(components)) {
    if (!componentInfo || typeof componentInfo !== "object") {
      errors.push(
        `Release ${version} component ${componentName} entry must be an object.`,
      );
      continue;
    }

    const { contractPath: relativeContractPath, contentHash } = componentInfo;

    if (typeof relativeContractPath !== "string") {
      errors.push(
        `Release ${version} component ${componentName} is missing contractPath.`,
      );
      continue;
    }

    if (typeof contentHash !== "string") {
      errors.push(
        `Release ${version} component ${componentName} is missing contentHash.`,
      );
      continue;
    }

    const actualContractPath = path.join(
      sourceDir,
      "components",
      componentName,
      "contract.json",
    );
    const versionSpecificContractPath = path.join(
      sourceDir,
      version,
      "components",
      componentName,
      "contract.json",
    );

    let contractRaw;
    try {
      try {
        contractRaw = await readFile(actualContractPath, "utf8");
      } catch {
        // Fall back to version-specific path during migration
        contractRaw = await readFile(versionSpecificContractPath, "utf8");
      }
    } catch (error) {
      errors.push(
        `Could not read contract for ${componentName} in release ${version}: ${error.message}.`,
      );
      continue;
    }

    try {
      const contract = JSON.parse(contractRaw);
      const actualHash = computeContentHash(contractRaw);

      if (contract.name !== componentName) {
        errors.push(
          `Contract name mismatch for ${componentName}: contract.name is ${contract.name} but ledger references ${componentName}.`,
        );
      }

      if (contract.version !== version) {
        errors.push(
          `Contract version mismatch for ${componentName}: contract.version is ${contract.version} but ledger release is ${version}.`,
        );
      }

      if (schemaVersion && contract.schemaVersion !== schemaVersion) {
        errors.push(
          `Contract schemaVersion mismatch for ${componentName}: contract.schemaVersion is ${contract.schemaVersion} but ledger release is ${schemaVersion}.`,
        );
      }

      if (actualHash !== contentHash) {
        errors.push(
          `Contract contentHash mismatch for ${componentName}: actual is ${actualHash} but ledger records ${contentHash}.`,
        );
      }
    } catch (error) {
      errors.push(
        `Could not parse or validate contract for ${componentName}: ${error.message}.`,
      );
    }
  }

  return errors;
}

export async function validateContracts({
  rootDir,
  compareRef,
  enforceGitImmutability = true,
  referenceDate = new Date(),
}) {
  const sourceDir = path.join(rootDir, "src", "contracts");
  const releaseDirs = (await readdir(sourceDir, { withFileTypes: true }))
    .filter((entry) => entry.isDirectory())
    .map((entry) => entry.name);
  const errors = [];
  const releaseEntries = [];

  // Read and validate release ledger if present
  const ledger = await readReleaseLedger({ rootDir });
  const ledgerErrors = validateReleaseLedger(ledger);
  if (ledgerErrors.length > 0) {
    errors.push("Release ledger validation failed:", ...ledgerErrors);
  }

  // Validate each ledger release against source contracts
  for (const release of ledger.releases) {
    const ledgerReleaseErrors = await validateLedgerRelease({
      rootDir,
      release,
      referenceDate,
    });
    if (ledgerReleaseErrors.length > 0) {
      errors.push(
        `Release ${release.version} ledger entry failed validation:`,
        ...ledgerReleaseErrors,
      );
    }
  }
  const validReleaseExample =
    [...releaseDirs]
      .filter((name) => isReleaseVersion(name))
      .sort(semverCompareDescending)[0] ?? "0.0.1";

  const invalidReleaseDirs = releaseDirs.filter(
    (name) => !isReleaseVersion(name) && !RESERVED_SOURCE_DIRS.has(name),
  );
  for (const folderName of invalidReleaseDirs) {
    errors.push(
      `Invalid contract release folder src/contracts/${folderName}. Use semver release folders. Example: src/contracts/${validReleaseExample}/components/<component-name>/contract.json.`,
    );
  }

  for (const releaseVersion of releaseDirs.filter((name) =>
    isReleaseVersion(name),
  )) {
    const releaseDir = path.join(sourceDir, releaseVersion);
    const manifestPath = path.join(releaseDir, "manifest.json");
    const manifestRelativePath = path.relative(rootDir, manifestPath);
    let manifestRaw;
    try {
      manifestRaw = await readFile(manifestPath, "utf8");
    } catch {
      errors.push(`Missing release manifest at ${manifestRelativePath}.`);
      continue;
    }
    const manifest = JSON.parse(manifestRaw);

    if (manifest.version !== releaseVersion) {
      errors.push(
        `Contract release version mismatch in ${manifestRelativePath}: folder is ${releaseVersion} but manifest.version is ${manifest.version}.`,
      );
    }

    const componentsDir = path.join(releaseDir, "components");
    let componentNames = [];
    try {
      componentNames = (await readdir(componentsDir, { withFileTypes: true }))
        .filter((entry) => entry.isDirectory())
        .map((entry) => entry.name);
    } catch {
      errors.push(
        `Missing components directory at ${path.relative(rootDir, componentsDir)}.`,
      );
      continue;
    }

    const manifestComponentNames = Object.keys(manifest.components ?? {});
    for (const componentName of manifestComponentNames) {
      if (!componentNames.includes(componentName)) {
        errors.push(
          `Contract release ${releaseVersion} manifest references missing component ${componentName}.`,
        );
      }
    }

    for (const componentName of componentNames) {
      const contractPath = path.join(
        componentsDir,
        componentName,
        "contract.json",
      );
      const contractRelativePath = path.relative(rootDir, contractPath);
      const contractRaw = await readFile(contractPath, "utf8");

      const contract = JSON.parse(contractRaw);
      const contentHash = computeContentHash(contractRaw);
      const manifestEntry = manifest.components?.[componentName];

      if (!manifestEntry) {
        errors.push(
          `Contract release ${releaseVersion} is missing manifest entry for component ${componentName}.`,
        );
        continue;
      }

      if (contract.version !== releaseVersion) {
        errors.push(
          `Contract version mismatch in ${contractRelativePath}: release folder is ${releaseVersion} but contract.version is ${contract.version}.`,
        );
      }

      if (contract.name !== componentName) {
        errors.push(
          `Contract name mismatch in ${contractRelativePath}: folder is ${componentName} but contract.name is ${contract.name}.`,
        );
      }

      if (
        manifestEntry.contractPath !==
        `./components/${componentName}/contract.json`
      ) {
        errors.push(
          `Manifest contractPath for ${componentName} in ${manifestRelativePath} must be ./components/${componentName}/contract.json.`,
        );
      }

      if (manifest.schemaVersion !== contract.schemaVersion) {
        errors.push(
          `Manifest schemaVersion mismatch in ${manifestRelativePath}: expected ${contract.schemaVersion}.`,
        );
      }

      if (manifestEntry.contentHash !== contentHash) {
        errors.push(
          `Manifest contentHash mismatch for ${componentName} in ${manifestRelativePath}. Expected ${contentHash}. If contracts changed, create a new release folder instead of editing ${releaseVersion} in place.`,
        );
      }

      if (enforceGitImmutability) {
        const previousContractRaw = await gitShow(
          rootDir,
          compareRef,
          contractRelativePath,
        );

        if (previousContractRaw && previousContractRaw !== contractRaw) {
          errors.push(
            `Existing contract ${contractRelativePath} changed after introduction. Create a new release folder and leave ${releaseVersion} immutable.`,
          );
        }
      }
    }

    if (enforceGitImmutability) {
      const previousManifestRaw = await gitShow(
        rootDir,
        compareRef,
        manifestRelativePath,
      );

      if (previousManifestRaw) {
        const previousManifest = JSON.parse(previousManifestRaw);
        const previousStableManifest = JSON.stringify(
          stripMutableManifestFields(previousManifest),
        );
        const currentStableManifest = JSON.stringify(
          stripMutableManifestFields(manifest),
        );

        if (previousStableManifest !== currentStableManifest) {
          errors.push(
            `Existing release manifest ${manifestRelativePath} changed in immutable fields. Only deprecatedAt may change on older releases; create a new release folder for contract changes.`,
          );
        }
      }
    }

    releaseEntries.push({ version: releaseVersion, manifest });
  }

  errors.push(...validateReleaseSet(releaseEntries, referenceDate));

  return errors;
}
