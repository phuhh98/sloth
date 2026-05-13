import type { SlothComponentContract } from "@/types/generated/schema.js";
import articleTeaser from "./components/article-teaser.contract.json" with { type: "json" };
import authorBio from "./components/author-bio.contract.json" with { type: "json" };
import breadCrumbTrail from "./components/breadcrumb.contract.json" with { type: "json" };

const components = {
  articleTeaser: articleTeaser as unknown as SlothComponentContract,
  authorBio: authorBio as unknown as SlothComponentContract,
  breadCrumbTrail: breadCrumbTrail as unknown as SlothComponentContract,
};

const contracts = {
  components,
};

export { components, contracts };

export default contracts;
