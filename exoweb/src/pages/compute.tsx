import React from 'react';
import Layout from '../components/Layout';
import { makeRoute } from '../lib/routing';

function Main() {
  return <>compute resources go here</>;
}

export const route = makeRoute({
  pattern: '/compute',
  load: async () => ({}),
  render: (props) => (
    <Layout>
      <Main />
    </Layout>
  ),
});
