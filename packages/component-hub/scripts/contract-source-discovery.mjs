import { readdir, readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const defaultRootDir = path.resolve(__dirname, "..");

export function getComponentSourceDir(rootDir = defaultRootDir) {
  return path.join(rootDir, "src", "contracts", "components");
}

export async function discoverComponentSources({
  rootDir = defaultRootDir,
} = {}) {
  const sourceDir = getComponentSourceDir(rootDir);

  let componentDirs;
  try {
    componentDirs = await readdir(sourceDir, { withFileTypes: true });
  } catch (error) {
    if (
      typeof error === "object" &&
      error !== null &&
      "code" in error &&
      error.code === "ENOENT"
    ) {
      return [];
    }

    throw error;
  }

  const components = componentDirs
    .filter((entry) => entry.isDirectory())
    .map((entry) => entry.name)
    .sort((left, right) => left.localeCompare(right));

  const discovered = [];

  for (const componentName of components) {
    const contractPath = path.join(sourceDir, componentName, "contract.json");
    const contractRaw = await readFile(contractPath, "utf8");
    const contract = JSON.parse(contractRaw);

    discovered.push({
      componentName,
      relativePath: path.join(
        "src",
        "contracts",
        "components",
        componentName,
        "contract.json",
      ),
      contractPath,
      contractRaw,
      contract,
    });
  }

  return discovered;
}
