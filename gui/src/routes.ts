import Home from './pages/Home.svelte';
import Workspace from './pages/Workspace.svelte';
import NewWorkspace from './pages/NewWorkspace.svelte';
import NewProcess from './pages/NewProcess.svelte';
import NotFound from './pages/NotFound.svelte';
import Process from './pages/Process.svelte';
import WorkspaceComponents from './pages/WorkspaceComponents.svelte';
import WorkspaceStorage from './pages/WorkspaceStorage.svelte';
import WorkspaceNetworking from './pages/WorkspaceNetworking.svelte';

export default {
  '/': Home,
  '/new-workspace': NewWorkspace,
  '/workspaces/:workspace': Workspace,
  '/workspaces/:workspace/components': WorkspaceComponents,
  '/workspaces/:workspace/storage': WorkspaceStorage,
  '/workspaces/:workspace/networking': WorkspaceNetworking,
  '/workspaces/:workspace/new-process': NewProcess,
  '/workspaces/:workspace/processes/:process': Process,
  '*': NotFound,
};
