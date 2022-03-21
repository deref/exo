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
  '#graphql\n    query ($componentId: String!) {\n      component: componentById(id: $componentId) {\n        id\n        name\n        asProcess {\n          running\n          cpuPercentage\n        }\n        environment {\n          variables {\n            name\n            value\n            source\n          }\n        }\n      }\n    }': {
    data: {
      __typename: 'Query';
      component: {
        __typename: 'Component';
        asProcess: unknown;
        environment: unknown;
        id: string;
        name: string;
      } | null;
    };
    variables: { componentId: string };
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
  '#graphql\n    {\n      templates: allTemplates {\n        name\n        displayName\n        iconGlyph\n        url\n      }\n    }': {
    data: {
      __typename: 'Query';
      templates: {
        __typename: 'Template';
        displayName: string;
        iconGlyph: string;
        name: string;
        url: string;
      }[];
    };
    variables: {};
  };
  '#graphql\n    query ($directoryPath: String) {\n      fileSystem {\n        homePath\n        file: fileOrHome(path: $directoryPath) {\n          path\n          parentPath\n          children {\n            name\n            path\n            isDirectory\n          }\n        }\n      }\n    }': {
    data: {
      __typename: 'Query';
      fileSystem: {
        __typename: 'FileSystem';
        file: {
          __typename: 'File';
          children: {
            __typename: 'File';
            isDirectory: boolean;
            name: string;
            path: string;
          }[];
          parentPath: string | null;
          path: string;
        } | null;
        homePath: string;
      };
    };
    variables: { directoryPath: string | null };
  };
  '#graphql\n    query ($workspaceId: String!) {\n      workspace: workspaceById(id: $workspaceId) {\n        id\n        stack {\n          id\n          displayName\n          components {\n            id\n            name\n            reconciling\n            running\n          }\n        }\n      }\n    }': {
    data: {
      __typename: 'Query';
      workspace: {
        __typename: 'Workspace';
        id: string;
        stack: {
          __typename: 'Stack';
          components: {
            __typename: 'Component';
            id: string;
            name: string;
            reconciling: boolean;
            running: boolean;
          }[];
          displayName: string;
          id: string;
        } | null;
      } | null;
    };
    variables: { workspaceId: string };
  };
  '#graphql\n    mutation ($id: String!) {\n      destroyStack(ref: $id) {\n        __typename\n      }\n    }': {
    data: {
      __typename: 'Mutation';
      destroyStack: { __typename: 'Reconciliation'; __typename: string | null };
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
  '#graphql\n    query ($workspaceId: String!) {\n      workspace: workspaceById(id: $workspaceId) {\n        id\n        stack {\n          components {\n            id\n            name\n            type\n          }\n        }\n      }\n    }': {
    data: {
      __typename: 'Query';
      workspace: {
        __typename: 'Workspace';
        id: string;
        stack: {
          __typename: 'Stack';
          components: {
            __typename: 'Component';
            id: string;
            name: string;
            type: string;
          }[];
        } | null;
      } | null;
    };
    variables: { workspaceId: string };
  };
  '#graphql\n    query ($workspaceId: String!) {\n      workspace: workspaceById(id: $workspaceId) {\n        id\n        stack {\n          vaults {\n            id\n            name\n            url\n            connected\n            authenticated\n          }\n          environment {\n            # TODO: Use a fragment.\n            variables {\n              name\n              value\n              source\n            }\n          }\n        }\n      }\n    }': {
    data: {
      __typename: 'Query';
      workspace: {
        __typename: 'Workspace';
        id: string;
        stack: {
          __typename: 'Stack';
          environment: {
            __typename: 'Environment';
            variables: {
              __typename: 'Variable';
              name: string;
              source: string;
              value: string;
            }[];
          };
          vaults: {
            __typename: 'Vault';
            authenticated: boolean;
            connected: boolean;
            id: string;
            name: string;
            url: string;
          }[];
        } | null;
      } | null;
    };
    variables: { workspaceId: string };
  };
};
