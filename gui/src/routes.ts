import Home from './pages/Home.svelte';
import Workspace from './pages/Workspace.svelte';
import NewWorkspace from './pages/NewWorkspace.svelte';
import NewProcess from './pages/NewProcess.svelte';
import NotFound from './pages/NotFound.svelte';
import Process from './pages/Process.svelte';
import Logs from './test/Logs.svelte';

export default {
  '/': Home,
  '/new-workspace': NewWorkspace,
  '/workspaces/:workspace': Workspace,
  '/workspaces/:workspace/new-process': NewProcess,
  '/workspaces/:workspace/processes/:process': Process,
  '/test/logs': Logs,
  '*': NotFound,
};
