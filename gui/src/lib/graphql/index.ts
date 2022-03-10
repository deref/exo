import type { SubscriptionOptions, WatchQueryOptions } from '@apollo/client';
import type { DocumentNode } from 'graphql';
import type { ReadableQuery, ReadableResult } from 'svelte-apollo';
import * as sa from 'svelte-apollo';
import gql from 'graphql-tag';
import type {
  QueryToData,
  QueryToVariables,
} from '~/src/lib/graphql/types.generated';

export type FragmentOptions = {
  fragments?: string[];
};

export type QueryOptions<Q extends string> = Omit<
  WatchQueryOptions<QueryToVariables<Q>, QueryToData<Q>>,
  'query'
> &
  FragmentOptions;

// Converts the query and associated fragments to a DocumentNode.
// Omit 'fragments' from options.
const makeArgs = <Options>(
  query: string,
  options: Options & FragmentOptions,
): [DocumentNode, Options] => {
  // Abuse the gql template literal tag for its caching behavior.
  const { fragments: optionalFragments, ...otherOptions } = options;
  const fragments = optionalFragments ?? [];
  const template = [
    query,
    Array(fragments.length).fill(''),
  ] as any as TemplateStringsArray;
  return [gql(template, fragments), otherOptions as Options];
};

export const query = <Q extends string>(
  query: Q,
  options: QueryOptions<Q> = {},
) => sa.query(...makeArgs(query, options)) as ReadableQuery<QueryToData<Q>>;

export type SubscribeOptions<Q extends string> = Omit<
  SubscriptionOptions<QueryToVariables<Q>>,
  'query'
> &
  FragmentOptions;

export const subscribe = <Q extends string>(
  query: Q,
  options: SubscribeOptions<Q> = {},
) =>
  sa.subscribe(...makeArgs(query, options)) as ReadableResult<QueryToData<Q>>;

// TODO: mutation.
// TODO: restore?

export { initClient } from './client';
