---
purpose: "Propose component contract starter set for Milestone 3"
created_date: "2026-05-10"
owner: "platform"
status: "proposal"
related_docs:
  - "docs/KANBAN-MILESTONE-3.md"
  - "docs/COMPONENT-CONTRACTS.md"
---

# Component Contract Starter Set Proposal

## Rationale

This proposal identifies **17 commonly used web components** suitable for initial component contracts. Selection criteria:
- **Popularity:** Used frequently across content-driven websites
- **Genericity:** Reusable across different projects/domains
- **Integration:** Fit naturally into layout/section/block model
- **SEO Value:** Clear SEO metadata strategy per component type

## Component Categories

### Layout Components (6)

Pure structural containers—no text content, no SEO metadata.

1. **Container (Wrapper)**
   - Purpose: Generic layout wrapper with responsive padding/width constraints
   - Use: All pages need a root container
   - Genericity: Highest (100% reuse)
   - Model: layout
   - Props: maxWidth, padding, alignItems, backgroundColor
   - SEO: None

2. **Grid Layout**
   - Purpose: Multi-column responsive grid
   - Use: Product grids, testimonial grids, team galleries
   - Genericity: High (90% reuse)
   - Model: layout
   - Props: columns, gap, alignItems, justifyContent, minColumnWidth
   - SEO: None

3. **Flex Stack**
   - Purpose: Flexible row/column stacking
   - Use: Navigation, card layouts, sidebar arrangements
   - Genericity: High (85% reuse)
   - Model: layout
   - Props: direction (row|column), spacing, alignItems, justifyContent, wrap
   - SEO: None

4. **Sidebar Layout**
   - Purpose: Two-column layout with sidebar
   - Use: Blog with sidebar, dashboard layouts
   - Genericity: Medium-High (75% reuse)
   - Model: layout
   - Props: sidebarWidth, gap, sidebarPosition (left|right)
   - SEO: None

5. **Header/Navbar Layout**
   - Purpose: Fixed or sticky navigation bar container
   - Use: Site header on every page
   - Genericity: High (80% reuse)
   - Model: layout
   - Props: sticky, backgroundColor, height, padding
   - SEO: None (nav semantics handled by parent section)

6. **Footer Layout**
   - Purpose: Multi-column footer grid
   - Use: Footer on every page
   - Genericity: Medium-High (70% reuse)
   - Model: layout
   - Props: columns, gap, backgroundColor, linkColor
   - SEO: None (links/text in child blocks)

---

### Section Components (7)

Page-level content sections with SEO metadata support.

7. **Hero Section**
   - Purpose: Large banner with headline, subheading, CTA, background image
   - Use: Page top, landing pages
   - Genericity: High (95% reuse)
   - Model: section
   - Common props: headline, subheading, backgroundImage, backgroundVideo, ctaText, ctaUrl, minHeight, overlay
   - **SEO Metadata:**
     - `title` → `<meta property="og:title">` / schema.org headline
     - `subtitle` → `<meta property="og:description">` snippet
     - `backgroundImage` → `<meta property="og:image">`
     - `url` → canonical structure
   - Example: Landing page hero, product page header

8. **Features/Benefits Section**
   - Purpose: Grid of feature cards with icons, titles, descriptions
   - Use: Product pages, service pages
   - Genericity: High (90% reuse)
   - Model: section
   - Common props: title, subtitle, features (array of {icon, title, description}), layout (grid|list)
   - **SEO Metadata:**
     - `title` → H2 heading (section hierarchy)
     - `features[].title` → schema.org Feature type
     - `features[].description` → snippet enrichment

9. **Testimonials Section**
   - Purpose: Carousel or grid of customer testimonials
   - Use: Trust/social proof on any page
   - Genericity: High (85% reuse)
   - Model: section
   - Common props: title, testimonials (array of {quote, author, role, image}), layout, autoPlay
   - **SEO Metadata:**
     - `title` → H2 heading
     - `testimonials[].quote` → schema.org Review type (ratings, author)

10. **FAQ Section**
    - Purpose: Accordion or expandable Q&A pairs
    - Use: Product pages, support pages
    - Genericity: High (80% reuse)
    - Model: section
    - Common props: title, questions (array of {q, a}), accordion (true|false)
    - **SEO Metadata:**
      - `title` → H2 heading
      - `questions[].q`, `questions[].a` → schema.org FAQPage (indexed by Google)
      - Auto-generate `<ld+json>` FAQPage schema

11. **Pricing Table Section**
    - Purpose: Multi-tier pricing comparison
    - Use: Pricing pages
    - Genericity: Medium-High (75% reuse)
    - Model: section
    - Common props: title, plans (array of {name, price, features, ctaText, highlighted}), currency, billingCycle
    - **SEO Metadata:**
      - `title` → H2 heading
      - `plans[].price` → schema.org Product / Offer
      - `plans[].name` → product identifier

12. **Newsletter Signup Section**
    - Purpose: Email subscription form
    - Use: Sidebar, footer, mid-page CTA
    - Genericity: High (85% reuse)
    - Model: section
    - Common props: title, subtitle, placeholder, ctaText, formAction, disclosureText
    - **SEO Metadata:**
      - `title`, `subtitle` → engagement signals
      - Form endpoint → no direct SEO (privacy-focused)

13. **CTA (Call-to-Action) Section**
    - Purpose: Simple banner with text and button(s)
    - Use: Mid-page conversions, exit-intent alternatives
    - Genericity: High (80% reuse)
    - Model: section
    - Common props: headline, description, ctaText, ctaUrl, alignment, backgroundColor
    - **SEO Metadata:**
      - `headline` → H2/H3 context
      - `ctaUrl` → internal link structure signal

---

### Block Components (5)

Reusable content blocks—atomic units with optional SEO metadata.

14. **Text Block (Rich Text)**
    - Purpose: Paragraph, formatted text with markdown/HTML
    - Use: Article content, descriptions, body copy
    - Genericity: Highest (100% reuse)
    - Model: block
    - Common props: content (markdown|html), fontSize, lineHeight, color, maxWidth
    - **SEO Metadata:**
      - `content` → primary text for indexing
      - Auto-parse for heading hierarchy, links, bold/italic signals

15. **Card Block**
    - Purpose: Self-contained content card (image, title, description, CTA)
    - Use: Team members, projects, product listings
    - Genericity: High (90% reuse)
    - Model: block
    - Common props: image, title, description, ctaText, ctaUrl, tags, metadata
    - **SEO Metadata:**
      - `title` → schema.org CreativeWork name
      - `description` → snippet
      - `image` → og:image for social sharing
      - `tags` → keyword signals

16. **Stat Block**
    - Purpose: Key metric display (number + label)
    - Use: KPI sections, team size, growth metrics
    - Genericity: High (80% reuse)
    - Model: block
    - Common props: value, unit, label, prefix, suffix, icon, trend
    - **SEO Metadata:**
      - `value`, `label` → structured data (schema.org Statistic)
      - No direct ranking impact (informational)

17. **Article/Post Block**
    - Purpose: Full-length article or blog post with metadata
    - Use: Blog posts, news articles, case studies, press releases
    - Genericity: High (85% reuse)
    - Model: block
    - Common props: title, body (markdown|html), featuredImage, author, publishedDate, readingTime, tags, categories, excerpt, canonicalUrl
    - **SEO Metadata:**
      - `title` → H1 (primary heading), schema.org Article headline
      - `body` → primary indexable content
      - `excerpt` → `<meta name="description">`
      - `featuredImage` → `<meta property="og:image">` for social sharing
      - `author` → schema.org author (Person or Organization)
      - `publishedDate` → schema.org datePublished, `<meta property="article:published_time">`
      - `tags` → schema.org keywords, internal linking structure
      - `canonicalUrl` → duplicate prevention
      - Auto-generate `<ld+json>` Article schema with all metadata
    - Variant: Can be used standalone or within an Article Section wrapper

---

## Summary Table

| Component | Category | Model | Genericity | SEO Value | Priority |
|-----------|----------|-------|-----------|-----------|----------|
| Container | Layout | layout | 100% | None | P0 |
| Grid Layout | Layout | layout | 90% | None | P0 |
| Flex Stack | Layout | layout | 85% | None | P0 |
| Hero | Section | section | 95% | High (og:, schema) | P0 |
| Features | Section | section | 90% | Medium (schema) | P1 |
| Testimonials | Section | section | 85% | Medium (review schema) | P1 |
| FAQ | Section | section | 80% | High (FAQPage schema) | P1 |
| Pricing | Section | section | 75% | Medium (offer schema) | P1 |
| Newsletter | Section | section | 85% | Low | P2 |
| CTA | Section | section | 80% | Low | P2 |
| Text Block | Block | block | 100% | High (indexing) | P0 |
| Card | Block | block | 90% | Medium (structured) | P1 |
| Article/Post | Block | block | 85% | High (article schema) | **P0** |
| Stat Block | Block | block | 80% | Low | P2 |
| Sidebar Layout | Layout | layout | 75% | None | P2 |
| Header Layout | Layout | layout | 80% | None | P2 |
| Footer Layout | Layout | layout | 70% | None | P2 |

---

## SEO Metadata Strategy Per Type

### Layout Components
- **No metadata:** Purely structural, no SEO relevance
- Focus: Responsive design, accessibility, rendering performance

### Section Components
- **Open Graph (og:):** Used when section is shared socially
  - `og:title` → section headline
  - `og:description` → section subtitle/summary
  - `og:image` → primary section image
- **Schema.org Markup:** Structured data for search engines
  - Hero → Article or NewsArticle (headline, image, datePublished)
  - Features → Collection with Feature items
  - Testimonials → AggregateRating or Review (author, rating, text)
  - FAQ → FAQPage (auto-generated by renderer)
  - Pricing → Product / Offer (name, price, availability)
- **Meta Tags:** Dynamically added to `<head>` per section on page

### Block Components
- **Embedded in parent section metadata:** Card block contributes to section schema
- **Link signals:** Internal links within text blocks → link structure
- **Keyword density:** Text content analyzed for relevance

---

## Recommended Phase-In

**Phase 1 (MVP):** 7 components
- Container, Grid, Flex Stack, Hero, Text Block, Card, Article/Post

**Phase 2:** +5 components
- Features, Testimonials, FAQ, CTA, Stat Block

**Phase 3:** +5 components
- Pricing, Newsletter, Sidebar, Header, Footer

---

## Next Steps

1. **Validate with design team:** Confirm layout/section/block model alignment
2. **Define JSON schemas:** Create contract schema for each component
3. **Generate starter pack:** Build React implementations for Phase 1
4. **Test SEO output:** Verify schema.org markup renders correctly
5. **Document recipes:** Show common composition patterns (e.g., hero + features + cta)
