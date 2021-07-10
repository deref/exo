import { ChakraProvider } from '@chakra-ui/react';
import React from 'react';
import routes from '../routes';
import { RouteRenderer } from '../hooks/useRouter';

export default function App() {
  return (
    <ChakraProvider>
      <RouteRenderer routes={routes} />
    </ChakraProvider>
  );
}
