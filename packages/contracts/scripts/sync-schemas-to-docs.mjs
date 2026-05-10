import { cp, mkdir, readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const packageRoot = path.resolve(__dirname, "..");
const repoRoot = path.resolve(packageRoot, "..", "..");
const defaultSourceRoot = path.join(packageRoot, "src", "schemas");
const defaultTargetRoot = path.join(
  repoRoot,
  "apps",
  "docs",
  "static",
  "schemas",
);
const defaultArtifact = "component-contract";

function parseArgs(argv) {
  const flags = {};
  for (let i = 0; i < argv.length; i += 1) {
    const token = argv[i];
    if (!token.startsWith("--")) {
      continue;
    }
    const key = token.slice(2);
    const next = argv[i + 1];
    if (!next || next.startsWith("--")) {
      flags[key] = "true";
      continue;
    }
    flags[key] = next;
    i += 1;
  }
  return flags;
}

async function assertVersionMatches({ sourceRoot, targetRoot, version }) {
  const relativePath = path.join(version, "schema.json");
  const sourcePath = path.join(sourceRoot, relativePath);
  const targetPath = path.join(targetRoot, relativePath);

  const source = JSON.parse(await readFile(sourcePath, "utf8"));
  const target = JSON.parse(await readFile(targetPath, "utf8"));

  if (JSON.stringify(source) !== JSON.stringify(target)) {
    throw new Error(
      `Docs schema mismatch for ${version}. Run sync:schemas:docs to update ${targetPath}.`,
    );
  }
}

function resolveArtifactRoots({ sourceRoot, targetRoot, artifact }) {
  return {
    sourceArtifactRoot: path.join(sourceRoot, artifact),
    targetArtifactRoot: path.join(targetRoot, artifact),
  };
}

export async function syncSchemasToDocs({
  sourceRoot = defaultSourceRoot,
  targetRoot = defaultTargetRoot,
  version,
  artifact = defaultArtifact,
  check = false,
} = {}) {
  const { sourceArtifactRoot, targetArtifactRoot } = resolveArtifactRoots({
    sourceRoot,
    targetRoot,
    artifact,
  });

  if (!version || version === "all") {
    if (check) {
      throw new Error("--check requires --version <x.y.z>");
    }
    await mkdir(targetRoot, { recursive: true });
    await cp(sourceRoot, targetRoot, { recursive: true });
    return;
  }

  const sourceVersionDir = path.join(sourceArtifactRoot, version);
  const targetVersionDir = path.join(targetArtifactRoot, version);

  if (check) {
    await assertVersionMatches({
      sourceRoot: sourceArtifactRoot,
      targetRoot: targetArtifactRoot,
      version,
    });
    return;
  }

  await mkdir(targetVersionDir, { recursive: true });
  await cp(sourceVersionDir, targetVersionDir, { recursive: true });
}

if (process.argv[1] && path.resolve(process.argv[1]) === __filename) {
  const flags = parseArgs(process.argv.slice(2));
  const version = flags.version;
  const artifact = flags.artifact || defaultArtifact;
  const check = flags.check === "true";

  try {
    await syncSchemasToDocs({ version, artifact, check });
  } catch (error) {
    console.error(error.message || error);
    process.exitCode = 1;
  }
}
