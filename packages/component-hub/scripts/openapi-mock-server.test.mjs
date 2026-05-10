import test from "node:test";
import assert from "node:assert/strict";

import { createMockServer } from "./openapi-mock-server.mjs";

test("mock server serves plugin status and ingest", async () => {
  const server = await createMockServer();
  await new Promise((resolve) => server.listen(0, resolve));

  try {
    const address = server.address();
    const baseUrl = `http://127.0.0.1:${address.port}`;

    const statusRes = await fetch(`${baseUrl}/sloth/inspection/plugin-status`);
    assert.equal(statusRes.status, 200);
    const statusPayload = await statusRes.json();
    assert.equal(statusPayload.totalComponents, 2);

    const ingestRes = await fetch(`${baseUrl}/sloth/contracts/ingest`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        contracts: [
          {
            name: "feature-grid",
            version: "0.1.0",
            schemaVersion: "0.0.1",
          },
        ],
      }),
    });

    assert.equal(ingestRes.status, 200);
    const ingestPayload = await ingestRes.json();
    assert.equal(ingestPayload.totalReceived, 1);
    assert.deepEqual(ingestPayload.created, ["feature-grid"]);
  } finally {
    await new Promise((resolve, reject) => {
      server.close((error) => {
        if (error) {
          reject(error);
          return;
        }
        resolve();
      });
    });
  }
});
