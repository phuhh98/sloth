import {
  mkdirSync,
  writeFileSync,
  existsSync,
  chmodSync,
  readFileSync,
} from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { spawnSync } from "node:child_process";
import crypto from "node:crypto";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const packageRoot = path.resolve(__dirname, "..");

const configPath = path.join(packageRoot, "distribution.config.json");
const distDir = path.join(packageRoot, "dist", "bin");
const mainPackage = "./cmd/sloth";

const config = JSON.parse(readFileSync(configPath, "utf8"));
mkdirSync(distDir, { recursive: true });

const checksums = [];

for (const target of config.platforms) {
  const outDir = path.join(distDir, `${target.os}-${target.arch}`);
  mkdirSync(outDir, { recursive: true });

  const exe = target.os === "windows" ? "sloth.exe" : "sloth";
  const outPath = path.join(outDir, exe);

  const env = {
    ...process.env,
    CGO_ENABLED: "0",
    GOOS: target.os === "windows" ? "windows" : target.os,
    GOARCH: target.arch === "amd64" ? "amd64" : "arm64",
  };

  const result = spawnSync(
    "go",
    ["build", "-trimpath", "-o", outPath, mainPackage],
    {
      cwd: packageRoot,
      env,
      stdio: "inherit",
    },
  );

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }

  if (target.os !== "windows") {
    chmodSync(outPath, 0o755);
  }

  const content = readFileSync(outPath);
  const hash = crypto.createHash("sha256").update(content).digest("hex");
  checksums.push(`${hash}  ${path.relative(packageRoot, outPath)}`);
}

const checksumPath = path.join(packageRoot, "dist", "checksums.txt");
writeFileSync(checksumPath, checksums.join("\n") + "\n", "utf8");

if (!existsSync(checksumPath)) {
  console.error("failed to generate checksums");
  process.exit(1);
}

console.log(
  `Built ${checksums.length} binaries and checksums at ${checksumPath}`,
);
