import { ChakraProvider } from '@chakra-ui/react';
import React, { useMemo } from 'react';
import routes from '../routes';
import { RouteRenderer } from '../hooks/useRouter';
import { ApolloProvider } from '@apollo/client';
import { newApolloClient } from '../lib/graphql';

export default function App() {
  const client = useMemo(newApolloClient, []);
  return (
    <ChakraProvider>
      <ApolloProvider client={client}>
        <RouteRenderer routes={routes} />
      </ApolloProvider>
    </ChakraProvider>
  );
}
