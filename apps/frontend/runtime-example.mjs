import { createRendererRegistry } from "../../packages/component-hub/src/react/runtime/renderer-mapping.mjs";

const payload = {
  page: {
    id: "landing-home",
    components: [
      {
        id: "hero-1",
        rendererKey: "hero-banner",
        props: {
          headline: "Build faster with sloth starter packs",
          subheadline:
            "Map renderer keys to components with deterministic runtime behavior.",
        },
      },
      {
        id: "teaser-1",
        rendererKey: "article-teaser",
        props: {
          entryId: "article-42",
        },
      },
    ],
  },
  linkedContent: {
    "article-42": {
      title: "Release Notes",
      excerpt:
        "Milestone 3 introduces starter packs and runtime mapping examples.",
      href: "/blog/release-notes",
    },
  },
};

const registry = createRendererRegistry({
  "hero-banner": (props) =>
    `\n<section class="hero">\n  <h1>${props.headline}</h1>\n  <p>${props.subheadline ?? ""}</p>\n</section>\n`,
  "article-teaser": (props, context) => {
    const linked = context.linkedContent?.[props.entryId];
    if (!linked) {
      return '\n<section class="teaser"><p>Missing linked content</p></section>\n';
    }

    return `\n<section class="teaser">\n  <h2>${linked.title}</h2>\n  <p>${linked.excerpt}</p>\n  <a href="${linked.href}">Read more</a>\n</section>\n`;
  },
});

const html = payload.page.components
  .map((node) =>
    registry.render(node, { linkedContent: payload.linkedContent }),
  )
  .join("\n");

process.stdout.write(`${html}\n`);
