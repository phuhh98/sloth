import type { Core } from '@strapi/strapi';
import { PLUGIN_ID } from './lib/constant/plugin';

/**
 * One of the lifecycle configuration
 * https://docs.strapi.io/cms/plugins-development/server-lifecycle#register
 */
const register = ({ strapi }: { strapi: Core.Strapi }) => {
  // register phase
  strapi.customFields.register({
    name: 'sloth-component-contract-schema',
    plugin: PLUGIN_ID,
    type: 'json',
    inputSize: {
      // optional
      default: 4,
      isResizable: true,
    },
  });
};

export default register;
