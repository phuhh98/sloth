import test from "node:test";
import assert from "node:assert/strict";

import { readCompareRef } from "./validate-contract-versions.mjs";

test("readCompareRef returns provided compare ref", () => {
  const compareRef = readCompareRef(["--compare-ref", "origin/main"]);

  assert.equal(compareRef, "origin/main");
});

test("readCompareRef falls back to HEAD when flag is missing", () => {
  const compareRef = readCompareRef([]);

  assert.equal(compareRef, "HEAD");
});

test("readCompareRef falls back to HEAD when flag has no value", () => {
  const compareRef = readCompareRef(["--compare-ref"]);

  assert.equal(compareRef, "HEAD");
});
