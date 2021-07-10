import { Box, Flex } from '@chakra-ui/react';
import React from 'react';
import Navbar from './Navbar';

interface Props {
  children: React.ReactNode;
}

export default function Layout(props) {
  return (
    <Flex>
      <Navbar />
      <Box>{props.children}</Box>
    </Flex>
  );
}
