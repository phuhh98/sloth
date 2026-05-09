import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";

import { validateContracts } from "./contract-version-policy.mjs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const rootDir = path.resolve(__dirname, "..");

export function readCompareRef(argv) {
  const compareRefIndex = argv.findIndex(
    (argument) => argument === "--compare-ref",
  );
  if (compareRefIndex >= 0 && argv[compareRefIndex + 1]) {
    return argv[compareRefIndex + 1];
  }

  return "HEAD";
}

export async function main() {
  const compareRef = readCompareRef(process.argv.slice(2));
  const errors = await validateContracts({
    rootDir,
    compareRef,
    enforceGitImmutability: true,
  });

  if (errors.length === 0) {
    return;
  }

  console.error("Component contract version policy failed:");
  for (const error of errors) {
    console.error(`- ${error}`);
  }
  console.error(
    "Action required: create a new release folder for changed contracts and mark older releases deprecated for at least 6 months.",
  );
  process.exitCode = 1;
}

if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
}
