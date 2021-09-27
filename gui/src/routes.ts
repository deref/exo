import Home from './pages/Home.svelte';
import NewProject from './pages/NewProject.svelte';
import NewProjectConfig from './pages/NewProjectConfig.svelte';
import Workspace from './pages/Workspace.svelte';
import WorkspaceInfo from './pages/WorkspaceInfo.svelte';
import WorkspaceComponents from './pages/WorkspaceComponents.svelte';
import WorkspaceVariables from './pages/WorkspaceVariables.svelte';
import WorkspaceStorage from './pages/WorkspaceStorage.svelte';
import WorkspaceNetworking from './pages/WorkspaceNetworking.svelte';
import WorkspaceNewComponent from './pages/WorkspaceNewComponent.svelte';
import NewProcess from './pages/NewProcess.svelte';
import NewDockerContainer from './pages/NewDockerContainer.svelte';
import NewDockerNetwork from './pages/NewDockerNetwork.svelte';
import NewDockerVolume from './pages/NewDockerVolume.svelte';
import Process from './pages/Process.svelte';
import Preferences from './pages/Preferences.svelte';
import Jobs from './pages/Jobs.svelte';
import Job from './pages/Job.svelte';
import NotFound from './pages/NotFound.svelte';

export default {
  '/': Home,
  '/new-project': NewProject,
  '/new-project/:starter': NewProjectConfig,
  '/workspaces/:workspace': Workspace,
  '/workspaces/:workspace/info': WorkspaceInfo,
  '/workspaces/:workspace/components': WorkspaceComponents,
  '/workspaces/:workspace/variables': WorkspaceVariables,
  '/workspaces/:workspace/storage': WorkspaceStorage,
  '/workspaces/:workspace/networking': WorkspaceNetworking,
  '/workspaces/:workspace/new-component': WorkspaceNewComponent,
  '/workspaces/:workspace/new-process': NewProcess,
  '/workspaces/:workspace/new-container': NewDockerContainer,
  '/workspaces/:workspace/new-volume': NewDockerVolume,
  '/workspaces/:workspace/new-network': NewDockerNetwork,
  '/workspaces/:workspace/processes/:process': Process,
  '/preferences': Preferences,
  '/jobs': Jobs,
  '/jobs/:job': Job,
  '*': NotFound,
};
