import type { Core } from '@strapi/strapi';

const controller = ({ strapi }: { strapi: Core.Strapi }) => ({
  async getPluginStatus(ctx: any) {
    ctx.body = await strapi.plugin('sloth-strapi-plugin').service('service').getPluginStatus();
  },

  getContractSchema(ctx: any) {
    const schemaVersion =
      typeof ctx.query.schemaVersion === 'string' ? ctx.query.schemaVersion : undefined;
    const inline =
      ctx.query.inline === 'true' || ctx.query.inline === true || ctx.query.inline === 1;

    ctx.body = strapi
      .plugin('sloth-strapi-plugin')
      .service('service')
      .getContractSchema(schemaVersion, inline);
  },

  async ingestContracts(ctx: any) {
    const payload = (ctx.request?.body ?? {}) as { contracts?: unknown };

    if (!Array.isArray(payload.contracts)) {
      ctx.status = 400;
      ctx.body = {
        error: 'invalid_payload',
        message: 'Expected body.contracts to be an array',
      };
      return;
    }

    const result = await strapi
      .plugin('sloth-strapi-plugin')
      .service('service')
      .ingestContracts(payload.contracts);

    ctx.body = result;
  },

  async getPageDelivery(ctx: any) {
    const pageDocumentId = ctx.params?.id;

    if (typeof pageDocumentId !== 'string' || pageDocumentId.trim().length === 0) {
      ctx.status = 400;
      ctx.body = {
        error: 'invalid_page_id',
        message: 'Expected page document id in path parameter',
      };
      return;
    }

    const page = await strapi
      .plugin('sloth-strapi-plugin')
      .service('service')
      .getPageDelivery(pageDocumentId);

    if (!page) {
      ctx.status = 404;
      ctx.body = {
        error: 'page_not_found',
      };
      return;
    }

    ctx.body = {
      page,
      note: 'Only first-level linked data should be consumed from this payload.',
    };
  },
});

export default controller;
