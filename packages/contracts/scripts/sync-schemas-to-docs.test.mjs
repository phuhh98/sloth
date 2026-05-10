import test from "node:test";
import assert from "node:assert/strict";
import { mkdtemp, mkdir, readFile, writeFile } from "node:fs/promises";
import os from "node:os";
import path from "node:path";

import { syncSchemasToDocs } from "./sync-schemas-to-docs.mjs";

test("syncSchemasToDocs copies all schema versions by default", async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), "sloth-schema-sync-"));
  const sourceRoot = path.join(root, "source", "schemas");
  const targetRoot = path.join(root, "target", "schemas");

  await mkdir(path.join(sourceRoot, "component-contract", "0.0.1"), {
    recursive: true,
  });
  await mkdir(path.join(sourceRoot, "component-contract", "0.0.2"), {
    recursive: true,
  });
  await mkdir(path.join(sourceRoot, "cli-config", "0.0.1"), {
    recursive: true,
  });
  await writeFile(
    path.join(sourceRoot, "component-contract", "0.0.1", "schema.json"),
    '{"v":"0.0.1"}\n',
    "utf8",
  );
  await writeFile(
    path.join(sourceRoot, "component-contract", "0.0.2", "schema.json"),
    '{"v":"0.0.2"}\n',
    "utf8",
  );
  await writeFile(
    path.join(sourceRoot, "cli-config", "0.0.1", "schema.json"),
    '{"v":"cli-0.0.1"}\n',
    "utf8",
  );

  await syncSchemasToDocs({ sourceRoot, targetRoot });

  const copiedA = await readFile(
    path.join(targetRoot, "component-contract", "0.0.1", "schema.json"),
    "utf8",
  );
  const copiedB = await readFile(
    path.join(targetRoot, "component-contract", "0.0.2", "schema.json"),
    "utf8",
  );
  const copiedCli = await readFile(
    path.join(targetRoot, "cli-config", "0.0.1", "schema.json"),
    "utf8",
  );
  assert.equal(copiedA, '{"v":"0.0.1"}\n');
  assert.equal(copiedB, '{"v":"0.0.2"}\n');
  assert.equal(copiedCli, '{"v":"cli-0.0.1"}\n');
});

test("syncSchemasToDocs copies only requested version", async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), "sloth-schema-sync-"));
  const sourceRoot = path.join(root, "source", "schemas");
  const targetRoot = path.join(root, "target", "schemas");

  await mkdir(path.join(sourceRoot, "cli-config", "0.0.3"), {
    recursive: true,
  });
  await writeFile(
    path.join(sourceRoot, "cli-config", "0.0.3", "schema.json"),
    '{"v":"0.0.3"}\n',
    "utf8",
  );

  await syncSchemasToDocs({
    sourceRoot,
    targetRoot,
    version: "0.0.3",
    artifact: "cli-config",
  });

  const copied = await readFile(
    path.join(targetRoot, "cli-config", "0.0.3", "schema.json"),
    "utf8",
  );
  assert.equal(copied, '{"v":"0.0.3"}\n');
});

test("syncSchemasToDocs --check validates parity and fails on mismatch", async () => {
  const root = await mkdtemp(path.join(os.tmpdir(), "sloth-schema-sync-"));
  const sourceRoot = path.join(root, "source", "schemas");
  const targetRoot = path.join(root, "target", "schemas");

  await mkdir(path.join(sourceRoot, "component-contract", "1.0.0"), {
    recursive: true,
  });
  await mkdir(path.join(targetRoot, "component-contract", "1.0.0"), {
    recursive: true,
  });
  await writeFile(
    path.join(sourceRoot, "component-contract", "1.0.0", "schema.json"),
    '{"v":"1.0.0"}\n',
    "utf8",
  );
  await writeFile(
    path.join(targetRoot, "component-contract", "1.0.0", "schema.json"),
    '{"v":"outdated"}\n',
    "utf8",
  );

  await assert.rejects(
    () =>
      syncSchemasToDocs({
        sourceRoot,
        targetRoot,
        version: "1.0.0",
        artifact: "component-contract",
        check: true,
      }),
    /Docs schema mismatch/,
  );

  await writeFile(
    path.join(targetRoot, "component-contract", "1.0.0", "schema.json"),
    '{"v":"1.0.0"}\n',
    "utf8",
  );
  await syncSchemasToDocs({
    sourceRoot,
    targetRoot,
    version: "1.0.0",
    artifact: "component-contract",
    check: true,
  });
});

test("syncSchemasToDocs --check requires explicit version", async () => {
  await assert.rejects(
    () => syncSchemasToDocs({ check: true }),
    /--check requires --version/,
  );
});
