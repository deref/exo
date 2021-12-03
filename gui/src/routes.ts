import Home from './pages/Home.svelte';
import NewProject from './pages/NewProject.svelte';
import NewProjectConfig from './pages/NewProjectConfig.svelte';
import Workspace from './pages/Workspace.svelte';
import WorkspaceDetails from './pages/WorkspaceDetails.svelte';
import WorkspaceComponents from './pages/WorkspaceComponents.svelte';
import WorkspaceVariables from './pages/WorkspaceVariables.svelte';
import WorkspaceStorage from './pages/WorkspaceStorage.svelte';
import WorkspaceNetworking from './pages/WorkspaceNetworking.svelte';
import WorkspaceNewComponent from './pages/WorkspaceNewComponent.svelte';
import WorkspaceAddVault from './pages/WorkspaceAddVault.svelte';
import NewProcess from './pages/NewProcess.svelte';
import NewDockerContainer from './pages/NewDockerContainer.svelte';
import NewDockerNetwork from './pages/NewDockerNetwork.svelte';
import NewDockerVolume from './pages/NewDockerVolume.svelte';
import EditComponent from './pages/EditComponent.svelte';
import ComponentDetails from './pages/ComponentDetails.svelte';
import ApiGatewayDetails from './pages/ApiGatewayDetails.svelte';
import Preferences from './pages/Preferences.svelte';
import Jobs from './pages/Jobs.svelte';
import Job from './pages/Job.svelte';
import AuthEsv from './pages/AuthEsv.svelte';
import NotFound from './pages/NotFound.svelte';

export default {
  '/': Home,
  '/new-project': NewProject,
  '/new-project/:starter': NewProjectConfig,
  '/workspaces/:workspace': Workspace,
  '/workspaces/:workspace/details': WorkspaceDetails,
  '/workspaces/:workspace/components': WorkspaceComponents,
  '/workspaces/:workspace/variables': WorkspaceVariables,
  '/workspaces/:workspace/storage': WorkspaceStorage,
  '/workspaces/:workspace/networking': WorkspaceNetworking,
  '/workspaces/:workspace/new-component': WorkspaceNewComponent,
  '/workspaces/:workspace/new-process': NewProcess,
  '/workspaces/:workspace/new-container': NewDockerContainer,
  '/workspaces/:workspace/new-volume': NewDockerVolume,
  '/workspaces/:workspace/new-network': NewDockerNetwork,
  '/workspaces/:workspace/add-vault': WorkspaceAddVault,
  '/workspaces/:workspace/components/:component/edit': EditComponent,
  '/workspaces/:workspace/components/:component': ComponentDetails,
  '/workspaces/:workspace/api-gateways/:component': ApiGatewayDetails,
  '/preferences': Preferences,
  '/jobs': Jobs,
  '/jobs/:job': Job,
  '/auth-esv': AuthEsv,
  '*': NotFound,
};
