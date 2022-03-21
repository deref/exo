<script lang="ts">
  import { derived, writable } from 'svelte/store';
  import Layout from '../components/Layout.svelte';
  import TwoColumn from '../components/TwoColumn.svelte';
  import LogPanel from '../components/LogPanel.svelte';
  import type { StreamStore, Event } from '../components/LogPanel.svelte';
  import ComponentsPanel from '../components/ComponentsPanel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import { query, mutation } from '../lib/graphql';

  export let params = { workspace: '' };
  const workspaceId = params.workspace;

  const q = query(
    `#graphql
    query ($workspaceId: String!) {
      workspace: workspaceById(id: $workspaceId) {
        id
        stack {
          id
          displayName
          components {
            id
            name
            reconciling
            running
          }
        }
      }
    }`,
    {
      variables: {
        workspaceId,
      },
      pollInterval: 5000, // XXX Use a subscription?
    },
  );
  $: workspace = $q.data?.workspace;
  $: stack = workspace?.stack && {
    ...workspace.stack,
    detailsUrl: `/workspaces/${encodeURIComponent(workspace.id)}/details`,
    newComponentUrl: `#/workspaces/${encodeURIComponent(
      workspace.id,
    )}/new-component`,
    components: workspace.stack.components.map((component) => ({
      ...component,
      url: `/workspaces/${encodeURIComponent(
        workspace!.id,
      )}/components/${encodeURIComponent(component.id)}`,
      editUrl: `/workspaces/${encodeURIComponent(
        workspace!.id,
      )}/components/${encodeURIComponent(component.id)}/edit`,
    })),
  };

  const destroyStackMutation = mutation(
    `#graphql
    mutation ($id: String!) {
      destroyStack(ref: $id) {
        __typename
      }
    }`,
    {
      variables: {
        id: workspaceId,
      },
    },
  );
  const destroyStack = async () => {
    await destroyStackMutation();
  };

  const setComponentRun = null as any; // XXX
  const setLogsVisible = null as any; // XXX

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

<Layout loading={$q.loading} error={$q.error}>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  {#if stack}
    <TwoColumn>
      <!-- XXX loading & error -->
      <ComponentsPanel
        slot="left"
        {stack}
        {destroyStack}
        {setLogsVisible}
        setRun={setComponentRun}
        {disposeComponent}
      />
      <LogPanel slot="right" {stream} />
    </TwoColumn>
  {:else}
    <!-- TODO: style me & add some UI to initialize a stack in this workspace. -->
    No current stack.
  {/if}
</Layout>
