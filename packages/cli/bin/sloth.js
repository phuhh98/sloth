#!/usr/bin/env node

import { existsSync } from "node:fs";
import { spawnSync } from "node:child_process";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const packageRoot = path.resolve(__dirname, "..");

const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const archMap = {
  arm64: "arm64",
  x64: "amd64",
};

const osName = platformMap[process.platform];
const archName = archMap[process.arch];

if (!osName || !archName) {
  console.error(`Unsupported platform: ${process.platform}/${process.arch}`);
  process.exit(1);
}

const exeName = osName === "windows" ? "sloth.exe" : "sloth";
const binaryPath = path.join(
  packageRoot,
  "dist",
  "bin",
  `${osName}-${archName}`,
  exeName,
);

if (!existsSync(binaryPath)) {
  console.error(
    `Missing CLI binary at ${binaryPath}. Run: pnpm --filter @sloth/cli run build:go`,
  );
  process.exit(1);
}

const result = spawnSync(binaryPath, process.argv.slice(2), {
  stdio: "inherit",
});

if (result.error) {
  console.error(result.error.message);
  process.exit(1);
}

process.exit(result.status ?? 0);
