import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import rehypeBaseLinks from './plugins/base-links.mjs';

// Static output only: `astro build` emits a plain `dist/` folder that GitHub
// Pages (or any web server) can serve. Full-text search is Pagefind — a
// static index generated at build time under /_pagefind/. No database, no
// server runtime.
//
// llms.txt is generated after the build by scripts/generate-llms-txt.mjs
// (`npm run build` chains it) — a plain index of every page for AI assistants.
// base is '/' so internal links are simple absolute paths that behave
// identically in dev, preview and production. If you ever deploy to a
// project GitHub Pages URL (username.github.io/<repo>/), set
// base: '/images/' again and restore `site` accordingly.
export default defineConfig({
  site: 'https://websummoner.github.io',
  base: '/images/',

  // Hand-written root-relative links in Markdown are not base-aware on their
  // own; this rewrites them so the base stays a single setting.
  markdown: {
    rehypePlugins: [[rehypeBaseLinks, { base: '/images/' }]],
  },

  integrations: [
    starlight({
      title: 'Browser Images',
      description:
        'Browser images for WebSummoner — Chrome, Firefox, Edge, Opera, Brave, Yandex and WebKit, built from this repository.',
      favicon: '/img/favicon.png',
      head: [
        { tag: 'meta', attrs: { property: 'og:image', content: 'https://websummoner.github.io/images/img/og-image.jpg' } },
        { tag: 'meta', attrs: { name: 'twitter:card', content: 'summary_large_image' } },
        { tag: 'meta', attrs: { name: 'twitter:image', content: 'https://websummoner.github.io/images/img/og-image.jpg' } },
      ],
      customCss: ['./src/styles/custom.css'],
      social: [
        {
          icon: 'github',
          label: 'Source code',
          href: 'https://github.com/WebSummoner/images',
        },
        {
          icon: 'seti:docker',
          label: 'Docker image',
          href: 'https://hub.docker.com/u/websummoner',
        },
      ],
      sidebar: [
        {
          label: 'Getting started',
          items: [{ slug: 'quick-start' }],
        },
        {
          label: 'Guides',
          items: [{ autogenerate: { directory: 'guides' } }],
        },
        {
          label: 'Reference',
          items: [{ autogenerate: { directory: 'reference' } }],
        },
      ],
    }),
  ],
});
