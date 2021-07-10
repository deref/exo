import React from 'react';
import Layout from '../components/Layout';
import { makeRoute } from '../lib/routing';

function Main() {
  return <>home</>;
}

export const route = makeRoute({
  pattern: '/',
  load: async () => ({}),
  render: (props) => (
    <Layout>
      <Main />
    </Layout>
  ),
});
