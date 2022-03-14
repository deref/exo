import '../src/style/main.css';
import { initMockClient as initMockGraphqlClient } from '../src/lib/graphql/testing';

export const parameters = {
  actions: { argTypesRegex: '^on[A-Z].*' },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
  darkMode: {
    darkClass: 'dark',
    lightClass: 'light',
    stylePreview: true,
  },
};

export const decorators = [
  () => {
    initMockGraphqlClient();
  },
];
