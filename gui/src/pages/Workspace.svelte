<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import TwoColumn from '../components/TwoColumn.svelte';
  import LogPanel from '../components/LogPanel.svelte';
  import ComponentsPanel from '../components/ComponentsPanel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import { query, mutation } from '../lib/graphql';

  export let params = { workspace: '' };

  $: workspaceId = params.workspace;

  const q = query(
    `#graphql
    query ($workspaceId: String!) {
      workspace: workspaceById(id: $workspaceId) {
        __typename #XXX
      }
    }`,
    {
      variables: {
        workspaceId,
      },
      pollInterval: 5000, // XXX Use a subscription.
    },
  );
  $: workspace = $q.data?.workspace;

  const destroyWorkspace = mutation(
    `#graphql
    mutation ($workspaceId: String!) {
      destroyWorkspace(ref: $id) {
        __typename
      }
    }`,
    {
      variables: {
        id: workspaceId,
      },
    },
  );

  const disposeComponent = mutation(
    `#graphql
    mutation ($componentId: String!) {
      disposeComponent(ref: $id) {
        __typename
      }
    }`,
  );
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <TwoColumn>
    <!-- XXX loading & error -->
    {#if workspace}
      <ComponentsPanel
        slot="left"
        {workspace}
        {destroyWorkspace}
        {disposeComponent}
      />
      <LogPanel slot="right" {workspace} />
    {/if}
  </TwoColumn>
</Layout>
