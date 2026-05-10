export function createRendererRegistry(entries) {
  const normalizedEntries = entries ?? {};
  const map = new Map(Object.entries(normalizedEntries));

  function get(rendererKey) {
    return map.get(rendererKey);
  }

  function has(rendererKey) {
    return map.has(rendererKey);
  }

  function register(rendererKey, renderer) {
    map.set(rendererKey, renderer);
  }

  function render(node, context = {}) {
    const renderer = get(node?.rendererKey);
    if (!renderer) {
      throw new Error(
        `Missing renderer for key: ${node?.rendererKey ?? "<undefined>"}`,
      );
    }
    return renderer(node?.props ?? {}, context);
  }

  return {
    get,
    has,
    register,
    render,
    keys: () => [...map.keys()].sort(),
  };
}
