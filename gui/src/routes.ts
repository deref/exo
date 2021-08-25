import Home from './pages/Home.svelte';
import NewWorkspace from './pages/NewWorkspace.svelte';
import Workspace from './pages/Workspace.svelte';
import WorkspaceComponents from './pages/WorkspaceComponents.svelte';
import WorkspaceStorage from './pages/WorkspaceStorage.svelte';
import WorkspaceNetworking from './pages/WorkspaceNetworking.svelte';
import NewProcess from './pages/NewProcess.svelte';
import Process from './pages/Process.svelte';
import Preferences from './pages/Preferences.svelte';
import Jobs from './pages/Jobs.svelte';
import Job from './pages/Job.svelte';
import NotFound from './pages/NotFound.svelte';

export default {
  '/': Home,
  '/new-workspace': NewWorkspace,
  '/workspaces/:workspace': Workspace,
  '/workspaces/:workspace/components': WorkspaceComponents,
  '/workspaces/:workspace/storage': WorkspaceStorage,
  '/workspaces/:workspace/networking': WorkspaceNetworking,
  '/workspaces/:workspace/new-process': NewProcess,
  '/workspaces/:workspace/processes/:process': Process,
  '/preferences': Preferences,
  '/jobs': Jobs,
  '/jobs/:job': Job,
  '*': NotFound,
};
