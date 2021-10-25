module.exports = {
  stories: [
    '../src/**/*.stories.mdx',
    '../src/**/*.stories.@(js|jsx|ts|tsx|svelte)',
  ],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/addon-svelte-csf',
    'storybook-dark-mode',
  ],
  svelteOptions: {
    preprocess: require('svelte-preprocess')(),
  },
  babel: (options) => {
    return {
      ...options,
      plugins: [
        ...(options.plugins || []),
        'babel-plugin-transform-vite-meta-env',
      ],
    };
  },
  webpackFinal: async (config) => {
    config.infrastructureLogging = {
      ...config.infrastructureLogging,
      level: 'error',
    };
    return config;
  },
};
