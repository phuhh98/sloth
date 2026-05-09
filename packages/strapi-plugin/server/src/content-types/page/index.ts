const schema = {
  kind: 'collectionType',
  collectionName: 'sloth_pages',
  info: {
    singularName: 'page',
    pluralName: 'pages',
    displayName: 'Page',
    description: 'Page template configuration',
  },
  options: {
    draftAndPublish: true,
  },
  pluginOptions: {},
  attributes: {
    name: {
      type: 'string',
      required: true,
      unique: true,
      minLength: 1,
    },
    label: {
      type: 'string',
      required: true,
      minLength: 1,
    },
    pageType: {
      type: 'enumeration',
      enum: ['static', 'dynamic'],
      required: true,
    },
    route: {
      type: 'string',
      required: true,
      unique: true,
    },
    dataset: {
      type: 'json',
      required: true,
    },
    puckConfig: {
      type: 'json',
      required: true,
    },
  },
};

export default {
  schema,
};
