import { VStack } from '@chakra-ui/react';
import NavbarItem from './NavbarItem';
import React from 'react';

export default function Navbar() {
  return (
    <VStack>
      <NavbarItem href="/">Home</NavbarItem>
      <NavbarItem href="/compute">Compute</NavbarItem>
      <NavbarItem href="/logs">Logs</NavbarItem>
      <NavbarItem href="/apis">APIs</NavbarItem>
    </VStack>
  );
}
