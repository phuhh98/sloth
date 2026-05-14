import { getTranslation } from './utils/getTranslation';
import { PLUGIN_ID } from './pluginId';
import { Initializer } from './components/Initializer';
import { PluginIcon } from './components/PluginIcon';
import { Code } from '@strapi/icons';
import ContractSchemaInput from './components/ContractSchemaInput';
import { size } from 'lodash';
import type { StrapiApp } from '@strapi/strapi/admin';

export default {
  register(app: StrapiApp) {
    app.addMenuLink({
      to: `plugins/${PLUGIN_ID}`,
      icon: PluginIcon,
      intlLabel: {
        id: `${PLUGIN_ID}.plugin.name`,
        defaultMessage: PLUGIN_ID,
      },
      permissions: [],
      Component: async () => {
        const module = await import('./pages/App');

        return {
          default: module.App,
        };
      },
    });

    app.registerPlugin({
      id: PLUGIN_ID,
      initializer: Initializer,
      isReady: false,
      name: PLUGIN_ID,
    });

    // register sloth-component-contract custom field
    app.customFields.register({
      name: 'sloth-component-contract-schema',
      pluginId: PLUGIN_ID, // the custom field is created by a color-picker plugin
      type: 'json', // the color will be stored as json string in the database
      intlLabel: {
        id: `${PLUGIN_ID}.contract-schema-input.label`,
        defaultMessage: 'Sloth Component Contract',
      },
      intlDescription: {
        id: `${PLUGIN_ID}.contract-schema-input.description`,
        defaultMessage: 'Define the contract of your component with Sloth',
      },
      icon: Code, // don't forget to create/import your icon component
      components: {
        Input: async () =>
          import('./components/ContractSchemaInput').then((module) => ({
            default: module.default,
          })), // the React component used to edit the value of the custom field
      },
      // options: {
      //   // declare options here
      //   base: [
      //     /*
      //     Declare settings to be added to the "Base settings" section
      //     of the field in the Content-Type Builder
      //   */
      //     {
      //       sectionTitle: {
      //         // Add a "Format" settings section
      //         id: 'color-picker.color.section.format',
      //         defaultMessage: 'Format',
      //       },
      //       items: [
      //         // Add settings items to the section
      //         {
      //           /*
      //           Add a "Color format" dropdown
      //           to choose between 2 different format options
      //           for the color value: hexadecimal or RGBA
      //         */
      //           intlLabel: {
      //             id: 'color-picker.color.format.label',
      //             defaultMessage: 'Color format',
      //           },
      //           name: 'enum',
      //           type: 'select',
      //           options: [
      //             // List all available "Color format" options
      //             {
      //               key: 'hex',
      //               defaultValue: 'hex',
      //               value: 'hex',
      //               metadatas: {
      //                 intlLabel: {
      //                   id: 'color-picker.color.format.hex',
      //                   defaultMessage: 'Hexadecimal',
      //                 },
      //               },
      //             },
      //             {
      //               key: 'rgba',
      //               value: 'rgba',
      //               metadatas: {
      //                 intlLabel: {
      //                   id: 'color-picker.color.format.rgba',
      //                   defaultMessage: 'RGBA',
      //                 },
      //               },
      //             },
      //           ],
      //         },
      //       ],
      //     },
      //   ],
      // },
    });
  },

  async registerTrads({ locales }: { locales: string[] }) {
    return Promise.all(
      locales.map(async (locale) => {
        try {
          const { default: data } = await import(`./translations/${locale}.json`);

          return { data, locale };
        } catch {
          return { data: {}, locale };
        }
      })
    );
  },
};
