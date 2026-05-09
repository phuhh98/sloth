import {
  mkdirSync,
  cpSync,
  writeFileSync,
  readFileSync,
  existsSync,
  rmSync,
} from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const packageRoot = path.resolve(__dirname, "..");
const distRoot = path.join(packageRoot, "dist", "publish-packages");
const binRoot = path.join(packageRoot, "dist", "bin");

const config = JSON.parse(
  readFileSync(path.join(packageRoot, "distribution.config.json"), "utf8"),
);

rmSync(distRoot, { recursive: true, force: true });
mkdirSync(distRoot, { recursive: true });

for (const target of config.platforms) {
  const packageName = `${config.scope}/${config.cliPackageName}-${target.os}-${target.arch}`;
  const folderName = `${config.cliPackageName}-${target.os}-${target.arch}`;
  const packageDir = path.join(distRoot, folderName);
  const sourceBinDir = path.join(binRoot, `${target.os}-${target.arch}`);
  const exe = target.os === "windows" ? "sloth.exe" : "sloth";
  const sourceBinary = path.join(sourceBinDir, exe);

  if (!existsSync(sourceBinary)) {
    throw new Error(`missing binary ${sourceBinary}; run build:go first`);
  }

  mkdirSync(packageDir, { recursive: true });
  mkdirSync(path.join(packageDir, "bin"), { recursive: true });

  cpSync(sourceBinary, path.join(packageDir, "bin", exe));

  const packageJson = {
    name: packageName,
    version: config.version,
    private: true,
    os: [target.os],
    cpu: [target.arch],
    files: ["bin"],
    bin: {
      sloth: `bin/${exe}`,
    },
  };

  writeFileSync(
    path.join(packageDir, "package.json"),
    JSON.stringify(packageJson, null, 2) + "\n",
    "utf8",
  );
}

console.log(`Generated publish packages in ${distRoot}`);
