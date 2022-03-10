// NOT YET GENERATED, but will be!

export type QueryType<Data = unknown, Variables = unknown> = {
  data: Data;
  variables: Variables;
};

export type QueryTypes = {
  '#graphql\n    {\n      workspaces: allWorkspaces {\n        id\n        root\n        displayName\n      }\n    }\n  ': {
    data: {
      workspaces: {
        id: string;
        root: string;
        displayName: string;
      }[];
    };
    variables: {};
  };
  '#graphql\n    subscription {\n      system: systemChange {\n        version {\n          installed\n          managed\n          upgrade\n        }\n      }\n    }\n  ': {
    data: {
      system: {
        version: {
          installed: string;
          managed: boolean;
          upgrade: string;
        };
      };
    };
    variables: {};
  };
};

export type QueryToTypes<Q> = Q extends keyof QueryTypes
  ? QueryTypes[Q]
  : QueryType;

export type QueryToData<Q> = QueryToTypes<Q>['data'];
export type QueryToVariables<Q> = QueryToTypes<Q>['variables'];
