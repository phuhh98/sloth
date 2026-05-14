import type { Core } from '@strapi/strapi';
import { PERMISSION_ACTIONS, PLUGIN_ID } from './lib/constant/plugin';

const bootstrap = async ({ strapi }: { strapi: Core.Strapi }) => {
  // bootstrap phase
  // await strapi.service('admin::permission').actionProvider.registerMany([
  //   {
  //     section: 'plugins',
  //     displayName: 'Read',
  //     uid: PERMISSION_ACTIONS.READ,
  //     pluginName: PLUGIN_ID,
  //   },
  //   {
  //     section: 'plugins',
  //     displayName: 'Settings',
  //     uid: PERMISSION_ACTIONS.SETTINGS,
  //     pluginName: PLUGIN_ID,
  //   },
  // ]);
};

export default bootstrap;
