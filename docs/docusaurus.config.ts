import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';
import { PresetEntry } from 'redocusaurus';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'Prism Documentation',
  tagline: 'Documentation for Prism platform',
  favicon: 'img/favicon.ico',

  // Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
  future: {
    v4: true, // Improve compatibility with the upcoming Docusaurus v4
  },

   url: 'https://dan-sones.github.io',
  baseUrl: '/prism/',
  organizationName: 'Dan-Sones',
  projectName: 'prism',
  trailingSlash: false,

  onBrokenLinks: 'throw',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },


  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
        },
        blog: false, // Disable blog
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
    [
      'redocusaurus',
      {
        specs: [
          {
            spec: 'openapi/assignment-service.yaml',
            route: '/api/assignment-service',
          },
          {
            spec: 'openapi/admin-service/admin-service.yaml',
            route: '/api/admin-service',
          },
        ],
        theme: {
          primaryColor: '#1890ff',
        },
      } satisfies PresetEntry,
    ],
  ],

  themeConfig: {
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Prism Documentation',
      logo: {
        alt: 'Prism Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          to: '/docs/developers/getting-started',
          label: 'For Developers',
          position: 'left',
        },
        {
          to: '/docs/experiment-owners/getting-started',
          label: 'For Experiment Owners',
          position: 'left',
        },
      ],
    },
    footer: {
      style: 'dark',
      copyright: `Copyright © ${new Date().getFullYear()} Prism. Built with Docusaurus`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
