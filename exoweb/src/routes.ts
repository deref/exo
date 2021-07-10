import { Route } from './lib/routing';
import { route as homeRoute } from './pages/home';
import { route as computeRoute } from './pages/compute';

const routes: Route<any>[] = [homeRoute, computeRoute];

export default routes;
