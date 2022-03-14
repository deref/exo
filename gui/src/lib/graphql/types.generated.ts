// GENERATED FILE. DO NOT EDIT.

export type QueryTypes = {
  '#graphql\n    query ($workspaceId: String!) {\n      workspace: workspaceById(id: $workspaceId) {\n        __typename #XXX\n      }\n    }': {
    data: {
      __typename: 'Query';
      workspace: { __typename: 'Workspace'; __typename: string | null } | null;
    };
    variables: { workspaceId: string };
  };
  '#graphql\n    mutation ($workspaceId: String!) {\n      destroyWorkspace(ref: $workspaceId) {\n        __typename\n      }\n    }': {
    data: {
      __typename: 'Mutation';
      destroyWorkspace: {
        __typename: 'Void';
        __typename: string | null;
      } | null;
    };
    variables: { workspaceId: string };
  };
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
};
