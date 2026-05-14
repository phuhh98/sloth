import { Plugin } from '../../types/pluginNamespace';
import { snakeCase, startCase } from 'lodash';
import pluralize from 'pluralize';
import { SLOTH_CONTENT_TYPE_KEY } from '../../lib/constant/contentType';
import { PLUGIN_ID } from '../../lib/constant/plugin';

const schema: Plugin.ContentTypeSchemaDefinition = {
  kind: 'collectionType',
  collectionName: snakeCase(SLOTH_CONTENT_TYPE_KEY.COMPONENT_CONTRACT),
  info: {
    singularName: SLOTH_CONTENT_TYPE_KEY.COMPONENT_CONTRACT,
    pluralName: pluralize.plural(SLOTH_CONTENT_TYPE_KEY.COMPONENT_CONTRACT),
    displayName: startCase(SLOTH_CONTENT_TYPE_KEY.COMPONENT_CONTRACT),
  },
  options: {
    comment: '',
  },
  attributes: {
    name: {
      type: 'string',
    },
    description: {
      type: 'string',
    },
    contract: {
      type: 'customField',
      customField: `plugin::${PLUGIN_ID}.sloth-component-contract-schema`,
    },
  },
};

export default {
  schema,
};
