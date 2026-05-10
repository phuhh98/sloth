---
purpose: "Milestone 3 shortlist decisions for dynamic/static components and shared base contract + SEO slot model."
status: "active"
owner: "platform"
last_updated: "2026-05-10"
related_docs:
  - "docs/archive/KANBAN-MILESTONE-3.md"
  - "docs/IDEAS.md"
  - "docs/COMPONENT-CONTRACTS.md"
---

# Component Hub Base Contract and Shortlist

## Dynamic Content Shortlist (HUB-000)

Approved first-class dynamic contracts:

- content-feed
- related-carousel
- article-teaser
- author-bio
- breadcrumb-trail
- seo-head

SEO mapping baseline:

- seo-head contributes page-level metadata fields (`title`, `description`, `canonicalUrl`, `noIndex`).
- article-teaser contributes optional social image fallback metadata.
- breadcrumb-trail contributes structured-data friendly breadcrumb labels and URLs.

## Static Block Shortlist (HUB-008)

Approved static/reusable contracts:

- text-block
- card-grid
- stat-strip
- cta-banner
- feature-grid
- testimonials
- faq-list
- pricing-table
- layout-shell

## Shared Base Contract (HUB-009)

Shared runtime fields to be modeled consistently across contracts:

- identity: `name`, `version`, `schemaVersion`, `renderMeta.rendererKey`
- visibility: dataset field `visibility` with `visible|hidden` options
- style tokens: dataset fields `themeTone`, `spacingY`, `containerWidth`
- data source: optional dynamic dataset fields (`valueDropdown`) for linked content
- seo contribution boundary:
  - each contract can expose metadata contribution fields in dataset
  - page-level merge policy is host/runtime responsibility
  - `seo-head` has highest priority for direct page SEO values

## Notes

- Current generated release keeps schemaVersion `0.0.1`.
- New contract families should be introduced as new contract releases to preserve immutability.
