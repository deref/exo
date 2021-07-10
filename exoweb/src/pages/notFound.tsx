import React from 'react';
import Layout from '../components/Layout';
import { makeRoute } from '../lib/routing';

function Main() {
  return <>not found</>;
}

export const route = makeRoute({
  pattern: '/404',
  load: async () => ({}),
  render: (props) => (
    <Layout>
      <Main />
    </Layout>
  ),
});
