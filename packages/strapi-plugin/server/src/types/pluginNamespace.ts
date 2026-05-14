import { Schema } from '@strapi/strapi';

declare namespace Plugin {
  interface ComponentContractSchemaFieldAttribute {
    attributes: Record<string, Schema.ContentType['attributes'][string] | CustomFieldAttribute>;
  }

  interface CustomFieldAttribute {
    type: 'customField';
    customField: `plugin::${string}.${string}`;
    options?: any;
  }

  export enum CUSTOM_FIELD_NAME {
    SLOTH_COMPONENT_CONTRACT_SCHEMA = 'sloth-component-contract-schema',
  }

  export type ContentTypeSchemaDefinition = Omit<
    Schema.ContentType,
    'uid' | 'modelType' | 'modelName' | 'globalId' | 'attributes'
  > &
    ComponentContractSchemaFieldAttribute;
}

export { Plugin };
