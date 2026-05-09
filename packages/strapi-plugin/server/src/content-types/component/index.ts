const schema = {
  kind: 'collectionType',
  collectionName: 'sloth_components',
  info: {
    singularName: 'component',
    pluralName: 'components',
    displayName: 'Component',
    description: 'Component materialized from contract ingest',
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
    componentKind: {
      type: 'enumeration',
      enum: ['CTA', 'Carousel', 'HeroSection', 'AsideLayout', 'Header', 'Footer', 'Custom'],
      required: true,
    },
    contractName: {
      type: 'string',
      required: true,
    },
    contractVersion: {
      type: 'string',
      required: true,
    },
    schemaVersion: {
      type: 'string',
      required: true,
    },
    contractPayload: {
      type: 'json',
      required: true,
    },
  },
};

export default {
  schema,
};
