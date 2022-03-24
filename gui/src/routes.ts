import HomePage from './components/HomePage.svelte';
import NewProjectPage from './components/NewProjectPage.svelte';
import NewProjectConfigPage from './components/NewProjectConfigPage.svelte';
import WorkspacePage from './components/WorkspacePage.svelte';
import WorkspaceDetailsPage from './components/WorkspaceDetailsPage.svelte';
import WorkspaceComponentsPage from './components/WorkspaceComponentsPage.svelte';
import WorkspaceVariablesPage from './components/WorkspaceVariablesPage.svelte';
import WorkspaceStoragePage from './components/WorkspaceStoragePage.svelte';
import WorkspaceNetworkingPage from './components/WorkspaceNetworkingPage.svelte';
import WorkspaceNewComponentPage from './components/WorkspaceNewComponentPage.svelte';
import WorkspaceAddVaultPage from './components/WorkspaceAddVaultPage.svelte';
import NewProcessPage from './components/NewProcessPage.svelte';
import NewDockerContainerPage from './components/NewDockerContainerPage.svelte';
import NewDockerNetworkPage from './components/NewDockerNetworkPage.svelte';
import NewDockerVolumePage from './components/NewDockerVolumePage.svelte';
import EditComponentPage from './components/EditComponentPage.svelte';
import ComponentDetailsPage from './components/ComponentDetailsPage.svelte';
import PreferencesPage from './components/PreferencesPage.svelte';
import JobsPage from './components/JobsPage.svelte';
import JobPage from './components/JobPage.svelte';
import AuthEsvPage from './components/AuthEsvPage.svelte';
import NotFoundPage from './components/NotFoundPage.svelte';

export default {
  '/': HomePage,
  '/new-project': NewProjectPage,
  '/new-project/:starter': NewProjectConfigPage,
  '/workspaces/:workspace': WorkspacePage,
  '/workspaces/:workspace/details': WorkspaceDetailsPage,
  '/workspaces/:workspace/components': WorkspaceComponentsPage,
  '/workspaces/:workspace/variables': WorkspaceVariablesPage,
  '/workspaces/:workspace/storage': WorkspaceStoragePage,
  '/workspaces/:workspace/networking': WorkspaceNetworkingPage,
  '/workspaces/:workspace/new-component': WorkspaceNewComponentPage,
  '/workspaces/:workspace/new-process': NewProcessPage,
  '/workspaces/:workspace/new-container': NewDockerContainerPage,
  '/workspaces/:workspace/new-volume': NewDockerVolumePage,
  '/workspaces/:workspace/new-network': NewDockerNetworkPage,
  '/workspaces/:workspace/add-vault': WorkspaceAddVaultPage,
  '/workspaces/:workspace/components/:component/edit': EditComponentPage,
  '/workspaces/:workspace/components/:component': ComponentDetailsPage,
  '/preferences': PreferencesPage,
  '/jobs': JobsPage,
  '/jobs/:job': JobPage,
  '/auth-esv': AuthEsvPage,
  '*': NotFoundPage,
};
