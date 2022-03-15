// GENERATED FILE. DO NOT EDIT.

export type QueryTypes = {
  '#graphql\n    subscription {\n      system: systemChange {\n        version {\n          installed\n          managed\n          upgrade\n        }\n      }\n    }\n  ': {
    data: {
      __typename: 'Subscription';
      system: {
        __typename: 'System';
        version: {
          __typename: 'VersionInfo';
          installed: string;
          managed: boolean;
          upgrade: string | null;
        };
      };
    };
    variables: {};
  };
  '#graphql\n    {\n      workspaces: allWorkspaces {\n        id\n        root\n        displayName\n      }\n    }\n  ': {
    data: {
      __typename: 'Query';
      workspaces: {
        __typename: 'Workspace';
        displayName: string;
        id: string;
        root: string;
      }[];
    };
    variables: {};
  };
  '#graphql\n    query ($workspaceId: String!) {\n      workspace: workspaceById(id: $workspaceId) {\n        components {\n          id\n          name\n        }\n      }\n    }': {
    data: {
      __typename: 'Query';
      workspace: {
        __typename: 'Workspace';
        components: { __typename: 'Component'; id: string; name: string }[];
      } | null;
    };
    variables: { workspaceId: string };
  };
  '#graphql\n    mutation ($id: String!) {\n      destroyWorkspace(ref: $id) {\n        __typename\n      }\n    }': {
    data: {
      __typename: 'Mutation';
      destroyWorkspace: {
        __typename: 'Void';
        __typename: string | null;
      } | null;
    };
    variables: { id: string };
  };
  '#graphql\n    mutation ($id: String!) {\n      disposeComponent(ref: $id) {\n        __typename\n      }\n    }': {
    data: {
      __typename: 'Mutation';
      disposeComponent: {
        __typename: 'Reconciliation';
        __typename: string | null;
      };
    };
    variables: { id: string };
  };
};