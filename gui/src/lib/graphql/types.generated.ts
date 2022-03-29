// GENERATED FILE. DO NOT EDIT.

import type { Instant } from './scalars';

export type Fragment_EnvironmentVariables_Data = {
  __typename: 'Environment';
  variables: {
    __typename: 'Variable';
    name: string;
    source: string;
    value: string;
  }[];
};
export type Fragment_EnvironmentVariables_Variables = {};

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
  '#graphql\n    query ($componentId: String!) {\n      component: componentById(id: $componentId) {\n        id\n        name\n        running\n        asProcess {\n          cpuPercent\n          residentBytes\n          started\n          ports\n          environment {\n            ...EnvironmentVariables\n          }\n        }\n        environment {\n          ...EnvironmentVariables\n        }\n      }\n    }\n    \n    # TODO: Lift to EnvironmentTable.\n    fragment EnvironmentVariables on Environment {\n      variables {\n        name\n        value\n        source\n      }\n    }\n    ': {
    data: {
      __typename: 'Query';
      component: {
        __typename: 'Component';
        asProcess: {
          __typename: 'ProcessComponent';
          cpuPercent: number | null;
          environment:
            | ({
                __typename: 'Environment';
              } & Fragment_EnvironmentVariables_Data)
            | null;
          ports: number[] | null;
          residentBytes: number | null;
          started: Instant | null;
        } | null;
        environment: {
          __typename: 'Environment';
        } & Fragment_EnvironmentVariables_Data;
        id: string;
        name: string;
        running: boolean;
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
  '#graphql\n    { __typename }\n  ': {
    data: { __typename: 'Query' };
    variables: {};
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
      destroyStack: { __typename: 'Reconciliation' };
    };
    variables: { id: string };
  };
  '#graphql\n    mutation ($id: String!) {\n      destroyComponent(ref: $id) {\n        __typename\n      }\n    }': {
    data: {
      __typename: 'Mutation';
      destroyComponent: { __typename: 'Reconciliation' };
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
  '#graphql\n    query ($workspaceId: String!) {\n      workspace: workspaceById(id: $workspaceId) {\n        id\n        stack {\n          networks {\n            type\n            name\n            componentId\n          }\n        }\n      }\n    }': {
    data: {
      __typename: 'Query';
      workspace: {
        __typename: 'Workspace';
        id: string;
        stack: {
          __typename: 'Stack';
          networks: {
            __typename: 'NetworkComponent';
            componentId: string;
            name: string;
            type: string;
          }[];
        } | null;
      } | null;
    };
    variables: { workspaceId: string };
  };
  '#graphql\n    query ($workspaceId: String!) {\n      workspace: workspaceById(id: $workspaceId) {\n        id\n        stack {\n          stores {\n            type\n            name\n            componentId\n            sizeMiB\n          }\n        }\n      }\n    }': {
    data: {
      __typename: 'Query';
      workspace: {
        __typename: 'Workspace';
        id: string;
        stack: {
          __typename: 'Stack';
          stores: {
            __typename: 'StoreComponent';
            componentId: string;
            name: string;
            sizeMiB: number | null;
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
  '#graphql\n    mutation ($id: String!) {\n      forgetVault(id: $id) {\n        __typename\n      }\n    }': {
    data: {
      __typename: 'Mutation';
      forgetVault: { __typename: 'Void' } | null;
    };
    variables: { id: string };
  };
};
