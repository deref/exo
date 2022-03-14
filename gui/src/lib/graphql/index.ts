import type {
  WatchQueryOptions as QueryOptionsBase,
  SubscriptionOptions as SubscriptionOptionsBase,
  MutationOptions as MutationOptionsBase,
} from '@apollo/client';
import type { DocumentNode } from 'graphql';
import type { ReadableQuery, ReadableResult } from 'svelte-apollo';
import * as sa from 'svelte-apollo';
import gql from 'graphql-tag';
import type { QueryTypes } from '~/src/lib/graphql/types.generated';

type QueryType<Data = unknown, Variables = unknown> = {
  data: Data;
  variables: Variables;
};

type QueryToTypes<Q> = Q extends keyof QueryTypes ? QueryTypes[Q] : QueryType;

type QueryToData<Q> = QueryToTypes<Q>['data'];
type QueryToVariables<Q> = QueryToTypes<Q>['variables'];

export type ExtraOptions = {
  // GraphQL fragment strings to be included.
  include?: string[];
};

// Converts the query and associated includes to a DocumentNode.
// Omit 'include' from options.
const makeArgs = <Options>(
  query: string,
  options: Options & ExtraOptions,
): [DocumentNode, NonNullable<Options>] => {
  // Abuse the gql template literal tag for its caching behavior.
  const { include: optionalFragments, ...otherOptions } = options;
  const include = optionalFragments ?? [];
  const template = [
    query,
    Array(include.length).fill(''),
  ] as any as TemplateStringsArray;
  return [gql(template, include), otherOptions as NonNullable<Options>];
};

export type QueryOptions<Q extends string> = Omit<
  QueryOptionsBase<QueryToVariables<Q>, QueryToData<Q>>,
  'query'
> &
  ExtraOptions;

export const query = <Q extends string>(
  query: Q,
  options: QueryOptions<Q> = {},
) => sa.query(...makeArgs(query, options)) as ReadableQuery<QueryToData<Q>>;

export type SubscribeOptions<Q extends string> = Omit<
  SubscriptionOptionsBase<QueryToVariables<Q>, QueryToData<Q>>,
  'query'
> &
  ExtraOptions;

export const subscribe = <Q extends string>(
  query: Q,
  options: SubscribeOptions<Q> = {},
) =>
  sa.subscribe(...makeArgs(query, options)) as ReadableResult<QueryToData<Q>>;

export type MutationOptions<Q extends string> = Omit<
  MutationOptionsBase<QueryToData<Q>, QueryToVariables<Q>>,
  'mutation'
> &
  ExtraOptions;

export type MutateOptions<Q extends string> = Omit<
  MutationOptionsBase<QueryToData<Q>, QueryToVariables<Q>>,
  'mutation'
>;

export const mutation = <Q extends string>(
  query: Q,
  options: MutationOptions<Q> = {},
): ((options: MutateOptions<Q> | undefined) => Promise<QueryToData<Q>>) => {
  const mutate = sa.mutation(...makeArgs(query, options));
  return (options: MutationOptions<Q> = {}) => mutate(options);
};

// TODO: Expose sa.restore?

export { initClient } from './client';
