import React, { useContext, useMemo, useState } from 'react';
import { Route } from '../lib/routing';
import { route as notFoundRoute } from '../pages/notFound';

interface Router {
  push(path: string): void;
}

const RouterContext = React.createContext<Router | null>(null);

interface Props {
  routes: Route<any>[];
}

export function RouteRenderer(props: Props) {
  const { routes } = props;

  const [epoch, setEpoch] = useState(0);
  const router = useMemo(
    (): Router => ({
      push(path: string) {
        window.history.pushState(undefined, '', path);
        // XXX handle loading, etc.
        setEpoch((epoch) => epoch + 1);
      },
    }),
    [window.location.href],
  );

  const route =
    routes.find(
      (route) => route.pattern === window.location.pathname, // XXX
    ) ?? notFoundRoute;
  const routeProps = {}; // XXX

  return (
    <RouterContext.Provider value={router}>
      {route.render(routeProps)}
    </RouterContext.Provider>
  );
}

export default function useRouter(): Router {
  const router = useContext(RouterContext);
  if (router == null) {
    throw Error('expected router in context');
  }
  return router;
}
