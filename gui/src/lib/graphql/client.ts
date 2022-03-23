import { ApolloClient, InMemoryCache, from, ApolloLink } from '@apollo/client';
import { split, HttpLink } from '@apollo/client';
import { getMainDefinition, Observable } from '@apollo/client/utilities';
import { ServerSentEventsLink } from '@graphql-sse/apollo-client';
import { setClient } from 'svelte-apollo';
import { derived, writable } from 'svelte/store';

export const connected = writable(true);
export const authenticated = writable(true);

// NOTE [ONLINE]: Since the server is always expected on localhost, the network
// should be extremely reliable. Assume that any network error means the server
// is not running. This also means that we can treat offline as a global, fatal
// condition, with captive offline mode UI.
export const online = derived(
  [connected, authenticated],
  ([$connected, $authenticated]) => $connected && $authenticated,
);

export const initClient = () => {
  const apiPort = parseInt(import.meta.env.VITE_API_PORT as string, 10);
  const endpoint = `http://localhost:${apiPort}/_exo/graphql`;

  const httpLink = new HttpLink({
    uri: endpoint,
  });

  const sseLink = new ServerSentEventsLink({
    graphQlSubscriptionUrl: endpoint,
  });

  const splitLink = split(
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

  const errorLink = new ApolloLink((operation, forward) => {
    return new Observable((observer) => {
      const sub = forward(operation).subscribe({
        next: (result) => {
          connected.set(true);
          authenticated.set(true);
          observer.next(result);
        },
        complete: () => {
          connected.set(true);
          authenticated.set(true);
          observer.complete();
        },
        error: (err) => {
          // Extract HTTP status code form messaged produced here:
          // https://github.com/apollographql/apollo-client/blob/c57547e63482641b29a0a93841d382799695b1cf/src/link/http/parseAndCheckHttpResponse.ts#L35
          const status = /Received status code (\d+)/.exec(err.message)?.[1];
          if (status === '401' || status === '403') {
            authenticated.set(false);
          } else {
            const isGraphqlErr = !!(err as any)?.result?.errors;
            if (!isGraphqlErr) {
              connected.set(false);
            }
          }
          observer.error(err);
        },
      });
      return () => {
        sub.unsubscribe();
      };
    });
  });

  const client = new ApolloClient({
    cache: new InMemoryCache(),
    link: from([errorLink, splitLink]),
  });

  setClient(client);
};
