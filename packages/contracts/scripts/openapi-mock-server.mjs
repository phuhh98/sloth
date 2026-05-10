import http from "node:http";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import { readFile } from "node:fs/promises";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const defaultSeedPath = path.join(
  __dirname,
  "..",
  "src",
  "mock",
  "seed-contracts.json",
);

function parseArgs(argv) {
  let port = 4010;
  for (let index = 0; index < argv.length; index += 1) {
    if (argv[index] === "--port" && argv[index + 1]) {
      port = Number.parseInt(argv[index + 1], 10);
    }
  }
  return { port };
}

function jsonResponse(res, status, payload) {
  const body = `${JSON.stringify(payload)}\n`;
  res.writeHead(status, {
    "Content-Type": "application/json",
    "Content-Length": Buffer.byteLength(body),
  });
  res.end(body);
}

function errorResponse(res, status, code, message) {
  return jsonResponse(res, status, { code, message });
}

function buildPluginStatus(seedContracts) {
  return {
    pluginName: "sloth-host-mock",
    pluginVersion: "0.1.0",
    compatibleSchemaVersions: ["0.0.1"],
    components: seedContracts.map((item) => ({
      name: item.name,
      version: item.version,
      schemaVersion: item.schemaVersion,
    })),
    totalComponents: seedContracts.length,
  };
}

function buildContractSchemaResponse(requestUrl) {
  const schemaVersion = requestUrl.searchParams.get("schemaVersion") || "0.0.1";
  const inline = requestUrl.searchParams.get("inline") === "true";

  return {
    schemaVersion,
    schemaUrl: `/sloth/inspection/contract-schema?schemaVersion=${schemaVersion}&inline=true`,
    ...(inline
      ? {
          schema: {
            $schema: "https://json-schema.org/draft/2020-12/schema",
            type: "object",
          },
        }
      : {}),
  };
}

function buildContractsPage(requestUrl, seedContracts) {
  const page = Number.parseInt(requestUrl.searchParams.get("page") ?? "1", 10);
  const pageSize = Number.parseInt(
    requestUrl.searchParams.get("pageSize") ?? "20",
    10,
  );
  const start = (page - 1) * pageSize;
  const items = seedContracts.slice(start, start + pageSize);

  return {
    items,
    pagination: {
      page,
      pageSize,
      totalItems: seedContracts.length,
      totalPages: Math.max(1, Math.ceil(seedContracts.length / pageSize)),
    },
  };
}

function findContractByPath(requestUrl, seedContracts) {
  const name = requestUrl.pathname.replace("/sloth/contracts/", "");
  return seedContracts.find((item) => item.name === name);
}

async function handleIngest(req, res, seedContracts) {
  const chunks = [];
  for await (const chunk of req) {
    chunks.push(chunk);
  }

  let payload = {};
  try {
    payload = JSON.parse(Buffer.concat(chunks).toString("utf8") || "{}");
  } catch {
    return errorResponse(res, 400, "BAD_REQUEST", "Invalid JSON payload");
  }

  const contracts = Array.isArray(payload.contracts) ? payload.contracts : [];
  if (contracts.length === 0) {
    return errorResponse(
      res,
      400,
      "BAD_REQUEST",
      "contracts array is required",
    );
  }

  const created = [];
  const updated = [];
  for (const item of contracts) {
    const existing = seedContracts.find(
      (contract) => contract.name === item.name,
    );
    if (existing) {
      updated.push(item.name);
      Object.assign(existing, item);
    } else {
      created.push(item.name);
      seedContracts.push(item);
    }
  }

  return jsonResponse(res, 200, {
    created,
    updated,
    failed: [],
    totalReceived: contracts.length,
  });
}

function handleGet(requestUrl, res, seedContracts) {
  if (requestUrl.pathname === "/healthz") {
    return jsonResponse(res, 200, { ok: true });
  }
  if (requestUrl.pathname === "/sloth/inspection/plugin-status") {
    return jsonResponse(res, 200, buildPluginStatus(seedContracts));
  }
  if (requestUrl.pathname === "/sloth/inspection/contract-schema") {
    return jsonResponse(res, 200, buildContractSchemaResponse(requestUrl));
  }
  if (requestUrl.pathname === "/sloth/contracts") {
    return jsonResponse(
      res,
      200,
      buildContractsPage(requestUrl, seedContracts),
    );
  }
  if (requestUrl.pathname.startsWith("/sloth/contracts/")) {
    const contract = findContractByPath(requestUrl, seedContracts);
    if (!contract) {
      return errorResponse(
        res,
        404,
        "NOT_FOUND",
        `Contract ${requestUrl.pathname.replace("/sloth/contracts/", "")} not found`,
      );
    }
    return jsonResponse(res, 200, { contract });
  }

  return errorResponse(
    res,
    404,
    "NOT_FOUND",
    `No route for GET ${requestUrl.pathname}`,
  );
}

export async function loadSeedContracts(seedPath = defaultSeedPath) {
  const raw = await readFile(seedPath, "utf8");
  return JSON.parse(raw);
}

export async function createMockServer({ seedPath } = {}) {
  const seedContracts = await loadSeedContracts(seedPath);

  const server = http.createServer(async (req, res) => {
    const requestUrl = new URL(req.url ?? "/", "http://localhost");

    if (req.method === "GET") {
      return handleGet(requestUrl, res, seedContracts);
    }
    if (
      req.method === "POST" &&
      requestUrl.pathname === "/sloth/contracts/ingest"
    ) {
      return handleIngest(req, res, seedContracts);
    }

    return errorResponse(
      res,
      404,
      "NOT_FOUND",
      `No route for ${req.method} ${requestUrl.pathname}`,
    );
  });

  return server;
}

export async function startMockServer({ port, seedPath } = {}) {
  const server = await createMockServer({ seedPath });
  await new Promise((resolve) => server.listen(port, resolve));
  const address = server.address();
  if (address && typeof address === "object") {
    process.stdout.write(`SLOTH_MOCK_SERVER_READY:${address.port}\n`);
  }
  return server;
}

if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  const args = parseArgs(process.argv.slice(2));
  try {
    await startMockServer({ port: args.port });
  } catch (error) {
    console.error(error);
    process.exitCode = 1;
  }
}
