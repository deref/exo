import { Box } from '@chakra-ui/react';
import React from 'react';
import Link from './Link';

interface Props {
  href: string;
  children: React.ReactNode;
}

export default function NavbarItem({ href, children }: Props) {
  return (
    <Box>
      <Link href={href}>{children}</Link>
    </Box>
  );
}
