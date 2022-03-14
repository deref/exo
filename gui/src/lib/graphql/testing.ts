import type { ApolloClient } from '@apollo/client';
import { setClient } from 'svelte-apollo';

export const initMockClient = () => {
  setClient({} as ApolloClient<unknown>);
};
