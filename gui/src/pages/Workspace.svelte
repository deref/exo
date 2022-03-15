<script lang="ts">
  import { derived, writable } from 'svelte/store';
  import Layout from '../components/Layout.svelte';
  import TwoColumn from '../components/TwoColumn.svelte';
  import LogPanel from '../components/LogPanel.svelte';
  import type { StreamStore, Event } from './LogPanel.svelte';
  import ComponentsPanel from '../components/ComponentsPanel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import { query, mutation } from '../lib/graphql';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;

  const q = query(
    `#graphql
    query ($workspaceId: String!) {
      workspace: workspaceById(id: $workspaceId) {
        components {
          id
          name
        }
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

  const destroyWorkspaceMutation = mutation(
    `#graphql
    mutation ($id: String!) {
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
  const destroyWorkspace = async () => {
    await destroyWorkspaceMutation();
  };

  const setComponentRun = null as any; // XXX

  const disposeComponentMutation = mutation(
    `#graphql
    mutation ($id: String!) {
      disposeComponent(ref: $id) {
        __typename
      }
    }`,
  );
  const disposeComponent = async (id: string) => {
    await disposeComponentMutation({ variables: { id } });
  };

  const events = writable<Event[]>([]); // XXX append from query. cache prior results.
  const filterString = writable(''); // XXX Use this as a graphql variable.

  const stream: StreamStore = {
    ...derived([events, filterString], ([events, filterString]) => ({
      events: events.filter(
        (event) => event.message.indexOf(filterString) >= 0,
      ),
      filterString,
    })),
    clearEvents: () => {
      events.set([]);
    },
    setFilterString: filterString.set,
  };
</script>

<Layout loading={$q.loading}>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <p>
    error: {$q.error}
  </p>
  <p>
    data: {$q.data}
  </p>
  {#if workspace}
    <TwoColumn>
      <!-- XXX loading & error -->
      <ComponentsPanel
        slot="left"
        {workspace}
        {destroyWorkspace}
        {setComponentRun}
        {disposeComponent}
      />
      <LogPanel slot="right" {stream} />
    </TwoColumn>
  {/if}
</Layout>
