import { mkdir, readFile, writeFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, "..");

function parseArgs(argv) {
  const flags = {};
  for (let i = 0; i < argv.length; i += 1) {
    const token = argv[i];
    if (token === "--") {
      continue;
    }
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

async function readJson(filePath) {
  const raw = await readFile(filePath, "utf8");
  return JSON.parse(raw);
}

export async function buildOCIReleasePayload({
  contractsRoot,
  version,
  outputPath,
}) {
  if (!version || typeof version !== "string") {
    throw new Error("version is required");
  }

  const releaseRoot = path.join(contractsRoot, version);
  const manifestPath = path.join(releaseRoot, "manifest.json");
  const manifest = await readJson(manifestPath);

  const contracts = [];
  for (const [name, metadata] of Object.entries(manifest.components ?? {})) {
    const contractPath = path.join(releaseRoot, metadata.contractPath);
    const payload = await readJson(contractPath);
    contracts.push({
      name,
      label: payload.label,
      version: payload.version || manifest.version,
      schemaVersion: payload.schemaVersion || manifest.schemaVersion,
      contentHash: metadata.contentHash,
      payload,
    });
  }

  contracts.sort((left, right) => left.name.localeCompare(right.name));

  const releasePayload = {
    version: manifest.version || version,
    schemaVersion: manifest.schemaVersion || "",
    contracts,
  };

  if (outputPath) {
    await mkdir(path.dirname(outputPath), { recursive: true });
    await writeFile(outputPath, `${JSON.stringify(releasePayload, null, 2)}\n`, "utf8");
  }

  return releasePayload;
}

if (process.argv[1] && path.resolve(process.argv[1]) === __filename) {
  const flags = parseArgs(process.argv.slice(2));
  const version = flags.version;
  const outputPath = flags.out;

  if (version === "") {
    console.error("Missing required --version argument");
    process.exitCode = 1;
  } else {
    const contractsRoot = path.join(rootDir, "dist", "registry", "contracts");
    try {
      await buildOCIReleasePayload({ contractsRoot, version, outputPath });
    } catch (error) {
      console.error(error.message);
      process.exitCode = 1;
    }
  }
}
