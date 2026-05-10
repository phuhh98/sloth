import test from "node:test";
import assert from "node:assert/strict";

import { createRendererRegistry } from "../src/react/runtime/renderer-mapping.mjs";

test("renderer registry renders node by rendererKey", () => {
  const registry = createRendererRegistry({
    "hero-banner": (props) => `<h1>${props.headline}</h1>`,
  });

  const rendered = registry.render({
    rendererKey: "hero-banner",
    props: { headline: "Hello" },
  });

  assert.equal(rendered, "<h1>Hello</h1>");
});

test("renderer registry throws when renderer key is missing", () => {
  const registry = createRendererRegistry({});

  assert.throws(
    () => registry.render({ rendererKey: "missing", props: {} }),
    /Missing renderer/,
  );
});
