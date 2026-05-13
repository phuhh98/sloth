import type { Core } from '@strapi/strapi';
import slothContract from '@sloth/contracts';

type GenericRecord = Record<string, unknown>;

type ContractInput = {
  name?: string;
  label?: string;
  componentKind?: string;
  version?: string;
  schemaVersion?: string;
  contractPayload?: GenericRecord;
};

const PLUGIN_NAME = 'sloth-strapi-plugin';
const COMPONENT_UID = 'plugin::sloth-strapi-plugin.component';
const PAGE_UID = 'plugin::sloth-strapi-plugin.page';
const DEFAULT_SCHEMA_VERSION = '0.0.1';
const DEFAULT_PLUGIN_VERSION = '0.0.0';
const SCHEMA_INSPECTION_PATH = '/sloth/inspection/contract-schema';

const COMPONENT_CONTRACT_SCHEMA = slothContract.schemas.componentContract;

const service = ({ strapi }: { strapi: Core.Strapi }) => ({
  async getPluginStatus() {
    const docs = (await strapi.documents(COMPONENT_UID).findMany({
      fields: ['name', 'label', 'componentKind', 'contractVersion', 'schemaVersion'],
      limit: 500,
      status: 'published',
    })) as Array<GenericRecord>;

    return {
      pluginName: PLUGIN_NAME,
      pluginVersion: process.env.npm_package_version ?? DEFAULT_PLUGIN_VERSION,
      compatibleSchemaVersions: [DEFAULT_SCHEMA_VERSION],
      components: docs.map((doc) => ({
        name: doc.name,
        label: doc.label,
        componentKind: doc.componentKind,
        contractVersion: doc.contractVersion,
        schemaVersion: doc.schemaVersion,
      })),
      totalComponents: docs.length,
    };
  },

  getContractSchema(schemaVersion?: string, inline?: boolean) {
    const resolvedVersion = schemaVersion ?? DEFAULT_SCHEMA_VERSION;

    if (inline) {
      return {
        schemaVersion: resolvedVersion,
        schema: COMPONENT_CONTRACT_SCHEMA,
      };
    }

    return {
      schemaVersion: resolvedVersion,
      schemaUrl: `${SCHEMA_INSPECTION_PATH}?schemaVersion=${encodeURIComponent(resolvedVersion)}&inline=true`,
    };
  },

  validateMinimalContractShape(contract: ContractInput) {
    const requiredFields: Array<keyof ContractInput> = [
      'name',
      'label',
      'version',
      'schemaVersion',
    ];

    const missing = requiredFields.filter((field) => {
      const value = contract[field];
      return typeof value !== 'string' || value.trim().length === 0;
    });

    return {
      valid: missing.length === 0,
      missing,
    };
  },

  async ingestContracts(contracts: ContractInput[]) {
    const created: string[] = [];
    const updated: string[] = [];
    const failed: Array<{ name: string; reason: string }> = [];

    for (const contract of contracts) {
      const contractName = contract.name ?? 'unknown';
      const shape = this.validateMinimalContractShape(contract);

      if (!shape.valid) {
        failed.push({
          name: contractName,
          reason: `missing_required_fields:${shape.missing.join(',')}`,
        });
        continue;
      }

      const normalizedName = contract.name;

      if (!normalizedName) {
        failed.push({
          name: contractName,
          reason: 'missing_required_fields:name',
        });
        continue;
      }

      const existing = (await strapi.documents(COMPONENT_UID).findMany({
        filters: { name: { $eq: normalizedName } },
        limit: 1,
        status: 'draft',
      })) as Array<GenericRecord>;

      const data = {
        name: normalizedName,
        label: contract.label,
        componentKind: contract.componentKind ?? 'Custom',
        contractName: normalizedName,
        contractVersion: contract.version,
        schemaVersion: contract.schemaVersion,
        contractPayload: contract.contractPayload ?? contract,
      };

      try {
        if (existing.length > 0 && typeof existing[0].documentId === 'string') {
          await strapi.documents(COMPONENT_UID).update({
            documentId: existing[0].documentId,
            data: data as any,
            status: 'draft',
          });
          updated.push(normalizedName);
        } else {
          await strapi.documents(COMPONENT_UID).create({
            data: data as any,
            status: 'draft',
          });
          created.push(normalizedName);
        }
      } catch (error) {
        failed.push({
          name: contractName,
          reason: error instanceof Error ? error.message : 'ingest_failed',
        });
      }
    }

    return {
      created,
      updated,
      failed,
      totalReceived: contracts.length,
    };
  },

  async getPageDelivery(pageDocumentId: string) {
    const page = (await strapi.documents(PAGE_UID).findOne({
      documentId: pageDocumentId,
      status: 'published',
    })) as GenericRecord | null;

    return page;
  },
});

export default service;
