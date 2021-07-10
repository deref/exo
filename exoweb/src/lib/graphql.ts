import {
  ApolloClient,
  InMemoryCache,
  ApolloLink,
  Observable,
  FetchResult,
} from '@apollo/client';
import { graphql, buildSchema, print } from 'graphql';

const schema = buildSchema(`
  type Query {
    hello: String
  }
`);

const root = { hello: () => 'Hello world!' };

const link = new ApolloLink(
  (operation, _forward) =>
    new Observable((observer) => {
      (async () => {
        try {
          const context = operation.getContext();
          const { extensions, variables } = operation;
          const query = print(operation.query);
          const response = await graphql(
            schema,
            query,
            root,
            context,
            variables,
            operation.operationName,
          );
          const result: FetchResult = {
            context,
            data: response.data,
            errors: response.errors,
            extensions,
          };
          observer.next(result);
          observer.complete();
        } catch (ex: unknown) {
          observer.error(ex);
        }
      })();
    }),
);

export const newApolloClient = () => {
  return new ApolloClient({
    link,
    cache: new InMemoryCache(),
  });
};
