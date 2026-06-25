import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: 'Go DAML SDK',
  description:
    'Go toolkit for building applications that talk to DAML / Canton ledgers over the gRPC Ledger API.',
  lang: 'en-US',

  // Build only the documentation, ignore stray markdown elsewhere in the repo.
  srcExclude: ['**/node_modules/**'],

  // Serve dev/README.md at /dev/ (VitePress does not treat README as an index).
  rewrites: {
    'dev/README.md': 'dev/index.md',
  },

  // Pretty URLs without trailing .html
  cleanUrls: true,

  // Last-updated timestamps from git
  lastUpdated: true,

  themeConfig: {
    nav: [
      {
        text: 'GitHub',
        link: 'https://github.com/noders-team/go-daml',
      },
    ],

    sidebar: {
      '/dev/': [
        {
          text: 'Developer Documentation',
          items: [
            { text: 'Overview', link: '/dev/' },
            { text: 'Getting Started', link: '/dev/getting-started' },
            { text: 'Code Generation', link: '/dev/code-generation' },
            { text: 'Minimal Examples', link: '/dev/minimal-examples' },
            { text: 'Compatibility', link: '/dev/compatibility' },
            { text: 'Troubleshooting', link: '/dev/troubleshooting' },
            { text: 'Triage & SLA', link: '/dev/triage-and-sla' },
          ],
        },
      ],
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/noders-team/go-daml' },
    ],

    search: {
      provider: 'local',
    },

    editLink: {
      pattern:
        'https://github.com/noders-team/go-daml/edit/main/docs/:path',
      text: 'Edit this page on GitHub',
    },
  },
})
