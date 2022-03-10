import { ApolloClient, InMemoryCache } from '@apollo/client';
import { split, HttpLink } from '@apollo/client';
import { getMainDefinition } from '@apollo/client/utilities';
import { ServerSentEventsLink } from '@graphql-sse/apollo-client';
import { setClient } from 'svelte-apollo';

export const initClient = () => {
  const apiPort = parseInt(import.meta.env.VITE_API_PORT as string, 10);
  const endpoint = `http://localhost:${apiPort}/_exo/graphql`;

  const httpLink = new HttpLink({
    uri: endpoint,
  });

  const sseLink = new ServerSentEventsLink({
    graphQlSubscriptionUrl: endpoint,
  });

  const link = split(
    ({ query }) => {
      const definition = getMainDefinition(query);
      return (
        definition.kind === 'OperationDefinition' &&
        definition.operation === 'subscription'
      );
    },
    sseLink,
    httpLink,
  );

  const client = new ApolloClient({
    cache: new InMemoryCache(),
    link: link,
  });

  setClient(client);
};
