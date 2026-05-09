import { createHash } from "node:crypto";
import { cp, mkdir, readFile, readdir, rm, writeFile } from "node:fs/promises";
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
  "component-hub",
  "dist",
  "registry",
);
const targetRoot = path.join(docsDir, "static", "registry");
const stateFilePath = path.join(targetRoot, "state.json");

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

  return files.sort();
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

async function readPreviousState() {
  try {
    const currentState = await readFile(stateFilePath, "utf8");
    return JSON.parse(currentState);
  } catch (error) {
    return undefined;
  }
}

async function syncRegistry() {
  const previousState = await readPreviousState();

  await rm(targetRoot, { recursive: true, force: true });
  await mkdir(targetRoot, { recursive: true });

  await cp(
    path.join(sourceRoot, "contracts"),
    path.join(targetRoot, "contracts"),
    {
      recursive: true,
    },
  );

  await mkdir(path.join(targetRoot, "themes"), { recursive: true });
  await mkdir(path.join(targetRoot, "packs"), { recursive: true });

  await writeFile(
    path.join(targetRoot, "themes", "index.json"),
    `${JSON.stringify({ registryFormatVersion: REGISTRY_FORMAT_VERSION, items: [] }, null, 2)}\n`,
    "utf8",
  );
  await writeFile(
    path.join(targetRoot, "packs", "index.json"),
    `${JSON.stringify({ registryFormatVersion: REGISTRY_FORMAT_VERSION, items: [] }, null, 2)}\n`,
    "utf8",
  );

  await writeFile(
    path.join(targetRoot, "index.json"),
    `${JSON.stringify(
      {
        registryFormatVersion: REGISTRY_FORMAT_VERSION,
        contractsIndex: "/sloth/registry/contracts/index.json",
        themesIndex: "/sloth/registry/themes/index.json",
        packsIndex: "/sloth/registry/packs/index.json",
      },
      null,
      2,
    )}\n`,
    "utf8",
  );

  const nextHash = await computeDirectoryHash(targetRoot);
  const nextState = computeNextRegistryState(
    previousState,
    nextHash,
    new Date().toISOString(),
  );

  await writeFile(
    stateFilePath,
    `${JSON.stringify(nextState, null, 2)}\n`,
    "utf8",
  );
}

syncRegistry().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
