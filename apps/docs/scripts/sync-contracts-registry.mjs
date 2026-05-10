import { createHash } from "node:crypto";
import { cp, mkdir, readFile, readdir, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

import {
  computeNextRegistryState,
  REGISTRY_FORMAT_VERSION,
} from "./registry-state.mjs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const docsDir = path.resolve(__dirname, "..");
const sourceRoot = path.resolve(
  docsDir,
  "..",
  "..",
  "packages",
  "contracts",
  "dist",
  "registry",
);
const targetRoot = path.join(docsDir, "static", "registry");

function getStateFilePath(targetRootPath) {
  return path.join(targetRootPath, "state.json");
}

async function collectFileEntries(rootDir, currentDir = rootDir) {
  const entries = await readdir(currentDir, { withFileTypes: true });
  const files = [];

  for (const entry of entries) {
    const absolutePath = path.join(currentDir, entry.name);

    if (entry.isDirectory()) {
      files.push(...(await collectFileEntries(rootDir, absolutePath)));
      continue;
    }

    const relativePath = path.relative(rootDir, absolutePath);
    if (relativePath === "state.json") {
      continue;
    }

    files.push(relativePath);
  }

  return files.sort((a, b) => a.localeCompare(b));
}

async function computeDirectoryHash(rootDir) {
  const hash = createHash("sha256");
  const files = await collectFileEntries(rootDir);

  for (const relativePath of files) {
    const absolutePath = path.join(rootDir, relativePath);
    const fileContents = await readFile(absolutePath);
    hash.update(relativePath);
    hash.update("\n");
    hash.update(fileContents);
    hash.update("\n");
  }

  return hash.digest("hex");
}

async function readPreviousStateAt(statePath) {
  try {
    const currentState = await readFile(statePath, "utf8");
    return JSON.parse(currentState);
  } catch {
    return undefined;
  }
}

export async function syncRegistry({
  sourceRootPath = sourceRoot,
  targetRootPath = targetRoot,
  now = new Date().toISOString(),
} = {}) {
  const stateFilePath = getStateFilePath(targetRootPath);
  const previousState = await readPreviousStateAt(stateFilePath);

  // Keep existing versioned artifacts; only merge new output and refresh mutable indexes/state.
  await mkdir(targetRootPath, { recursive: true });
  await mkdir(path.join(targetRootPath, "contracts"), { recursive: true });

  await cp(
    path.join(sourceRootPath, "contracts"),
    path.join(targetRootPath, "contracts"),
    {
      recursive: true,
    },
  );

  await mkdir(path.join(targetRootPath, "themes"), { recursive: true });

  await writeFile(
    path.join(targetRootPath, "themes", "index.json"),
    `${JSON.stringify({ registryFormatVersion: REGISTRY_FORMAT_VERSION, items: [] }, null, 2)}\n`,
    "utf8",
  );

  await writeFile(
    path.join(targetRootPath, "index.json"),
    `${JSON.stringify(
      {
        registryFormatVersion: REGISTRY_FORMAT_VERSION,
        contractsIndex: "/sloth/registry/contracts/index.json",
        themesIndex: "/sloth/registry/themes/index.json",
      },
      null,
      2,
    )}\n`,
    "utf8",
  );

  const nextHash = await computeDirectoryHash(targetRootPath);
  const nextState = computeNextRegistryState(previousState, nextHash, now);

  await writeFile(
    stateFilePath,
    `${JSON.stringify(nextState, null, 2)}\n`,
    "utf8",
  );

  return nextState;
}

if (process.argv[1] && process.argv[1] === __filename) {
  try {
    await syncRegistry();
  } catch (error) {
    console.error(error);
    process.exitCode = 1;
  }
}
