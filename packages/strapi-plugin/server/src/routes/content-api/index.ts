const createContentApiRoutes = () => ({
  type: 'content-api',
  routes: [
    {
      method: 'GET',
      path: '/inspection/plugin-status',
      handler: 'controller.getPluginStatus',
      config: {
        auth: false,
        policies: [],
      },
    },
    {
      method: 'GET',
      path: '/inspection/contract-schema',
      handler: 'controller.getContractSchema',
      config: {
        auth: false,
        policies: [],
      },
    },
    {
      method: 'POST',
      path: '/contracts/ingest',
      handler: 'controller.ingestContracts',
      config: {
        auth: false,
        policies: [],
      },
    },
    {
      method: 'GET',
      path: '/pages/:id/delivery',
      handler: 'controller.getPageDelivery',
      config: {
        auth: false,
        policies: [],
      },
    },
  ],
});

export default createContentApiRoutes;
