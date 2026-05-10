import { mkdir, readFile, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import {
  isReleaseVersion,
  semverCompareDescending,
} from "./contract-version-policy.mjs";

export const RELEASE_LEDGER_FORMAT_VERSION = "1";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const defaultRootDir = path.resolve(__dirname, "..");

export function getReleaseLedgerPath(rootDir = defaultRootDir) {
  return path.join(rootDir, "src", "contracts", "releases", "index.json");
}

export function createEmptyReleaseLedger() {
  return {
    ledgerFormatVersion: RELEASE_LEDGER_FORMAT_VERSION,
    releases: [],
  };
}

export async function readReleaseLedger({
  rootDir = defaultRootDir,
  allowMissing = true,
} = {}) {
  const ledgerPath = getReleaseLedgerPath(rootDir);

  try {
    const raw = await readFile(ledgerPath, "utf8");
    return JSON.parse(raw);
  } catch (error) {
    if (
      allowMissing &&
      typeof error === "object" &&
      error !== null &&
      "code" in error &&
      error.code === "ENOENT"
    ) {
      return createEmptyReleaseLedger();
    }

    throw error;
  }
}

export async function writeReleaseLedger({ rootDir = defaultRootDir, ledger }) {
  const ledgerPath = getReleaseLedgerPath(rootDir);
  await mkdir(path.dirname(ledgerPath), { recursive: true });
  await writeFile(ledgerPath, `${JSON.stringify(ledger, null, 2)}\n`, "utf8");
}

export function validateReleaseLedger(ledger) {
  const errors = [];

  if (!ledger || typeof ledger !== "object") {
    return ["Release ledger must be a JSON object."];
  }

  if (ledger.ledgerFormatVersion !== RELEASE_LEDGER_FORMAT_VERSION) {
    errors.push(
      `Release ledger format must be ${RELEASE_LEDGER_FORMAT_VERSION}.`,
    );
  }

  if (!Array.isArray(ledger.releases)) {
    errors.push("Release ledger must include a releases array.");
    return errors;
  }

  const seenVersions = new Set();

  for (const release of ledger.releases) {
    if (!release || typeof release !== "object") {
      errors.push("Each release entry must be an object.");
      continue;
    }

    if (!isReleaseVersion(release.version)) {
      errors.push(`Invalid release version in ledger: ${release.version}.`);
    }

    if (seenVersions.has(release.version)) {
      errors.push(`Duplicate release version in ledger: ${release.version}.`);
    }

    seenVersions.add(release.version);

    if (typeof release.createdAt !== "string") {
      errors.push(
        `Release ${release.version} must include createdAt as an ISO string.`,
      );
    }

    if (typeof release.components !== "object" || release.components === null) {
      errors.push(`Release ${release.version} must define components object.`);
      continue;
    }

    for (const [componentName, componentInfo] of Object.entries(
      release.components,
    )) {
      if (!componentInfo || typeof componentInfo !== "object") {
        errors.push(
          `Release ${release.version} component ${componentName} must be an object.`,
        );
        continue;
      }

      if (typeof componentInfo.contractPath !== "string") {
        errors.push(
          `Release ${release.version} component ${componentName} is missing contractPath.`,
        );
      }

      if (typeof componentInfo.contentHash !== "string") {
        errors.push(
          `Release ${release.version} component ${componentName} is missing contentHash.`,
        );
      }
    }
  }

  const sortedVersions = [...ledger.releases]
    .map((release) => release.version)
    .sort(semverCompareDescending);

  for (let index = 0; index < sortedVersions.length; index += 1) {
    if (ledger.releases[index]?.version === sortedVersions[index]) {
      continue;
    }

    errors.push(
      "Release ledger must keep releases sorted descending by semantic version.",
    );
    break;
  }

  return errors;
}
