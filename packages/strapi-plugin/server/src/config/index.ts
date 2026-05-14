export default {
  /**
   * Default configuration. for sloth
   * Strapi uses deep merge to merge default with user configurations - user take precedence.
   * https://docs.strapi.io/cms/plugins-development/server-configuration
   */
  default: {},

  /**
   * Validate the merged configuration object against a schema.
   * We could defined a configuration schema for sloth inside this server
   * then load it via ajv and validate with ajv for the merged configuration object.
   * Or simpler with zod?
   */
  validator() {},
};
