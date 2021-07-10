import { Link as ChakraLink } from '@chakra-ui/react';
import React, { useCallback } from 'react';
import useRouter from '../hooks/useRouter';

type Props = Exclude<React.ComponentProps<typeof ChakraLink>, 'as'>;

export default function Link(props: Props) {
  const router = useRouter();
  const handleClick: React.MouseEventHandler<HTMLAnchorElement> = useCallback(
    (e) => {
      if (props.href?.startsWith('/')) {
        router.push(props.href);
        e.preventDefault();
        e.stopPropagation();
      }
      props.onClick?.(e);
    },
    [props.onClick],
  );
  return <ChakraLink {...props} onClick={handleClick} />;
}
