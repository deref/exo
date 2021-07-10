import { gql, useQuery } from '@apollo/client';
import React from 'react';
import Layout from '../components/Layout';
import { makeRoute } from '../lib/routing';

const query = gql`
  {
    hello
  }
`;

function Main() {
  const { loading, error, data } = useQuery(query);
  return (
    <>
      RESUME WORK HERE, call an API to get list of compute resources.
      {JSON.stringify({ loading, error, data })}
    </>
  );
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
