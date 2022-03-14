<script lang="ts">
  import Icon from './Icon.svelte';
  import Panel from './Panel.svelte';
  import IfEnabled from './IfEnabled.svelte';
  import IconButton from './IconButton.svelte';
  import RemoteData from './RemoteData.svelte';
  import ContextMenu from './ContextMenu.svelte';
  import MenuItem from './MenuItem.svelte';
  import ProcfileChecker from './processes/ProcfileChecker.svelte';
  import ProcessListTable from './processes/ProcessListTable.svelte';
  import * as router from 'svelte-spa-router';
  import { query, mutation } from '../lib/graphql';

  import { modal } from '../lib/modal';
  import { bind } from '../components/modal/Modal.svelte';
  import ModalDialogue from '../components/modal/ModalDialogue.svelte';

  export let workspaceId: string;

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
      destroyWorkspace(ref: $workspaceId) {
        __typename
      }
    }`,
    {
      variables: {
        workspaceId,
      },
    },
  );

  // TODO: Display out-of-sync manifests somehow.
  //let procfileExport: string | null = null;
  //async function checkProcfile() {
  //  const current = await workspace.readFile('Procfile');
  //  const computed = await workspace.exportProcfile();

  //  procfileExport = current && current === computed ? null : computed;
  //}

  const showWorkspaceDeleteModal = (displayName: string) => {
    modal.set(
      bind(ModalDialogue, {
        h3: 'Delete workspace?',
        bodyText: `Are you sure you want to delete ${displayName}?\nThis is irreversible, but will only delete the workspace in exo, not the files.`,
        danger: true,
        actionLabel: 'Yes, delete',
        onOkay: async () => {
          await destroyWorkspace({}); // TODO: try/catch.
          router.push('/');
        },
      }),
    );
  };
</script>

{#await workspace.describeSelf()}
  <Panel title="" backRoute="/" />
{:then description}
  <Panel title={description.displayName} backRoute="/" --panel-padding="0 1rem">
    <div class="actions" slot="actions">
      <span>Logs</span>
      <div class="menu" tabindex="0">
        <IconButton
          glyph="Ellipsis"
          tooltip="Workspace actions..."
          on:click={() => {}}
        />

        <ContextMenu title={description.displayName}>
          <MenuItem
            glyph="Details"
            href={`/workspaces/${encodeURIComponent(workspaceId)}/details`}
          >
            View details
          </MenuItem>
          <MenuItem
            glyph="Add"
            href={`#/workspaces/${encodeURIComponent(
              workspaceId,
            )}/new-component`}
          >
            Add component
          </MenuItem>
          <MenuItem
            glyph="Delete"
            danger
            on:click={() => showWorkspaceDeleteModal(description.displayName)}
          >
            Destroy workspace
          </MenuItem>
        </ContextMenu>
      </div>
    </div>

    <section>
      <button
        id="add-component"
        on:click={() => {
          router.push(
            `#/workspaces/${encodeURIComponent(workspaceId)}/new-component`,
          );
        }}
      >
        <Icon glyph="Add" /> Add component
      </button>
      <RemoteData data={processList} let:data let:error>
        <div slot="success">
          <ProcessListTable {data} {workspace} {workspaceId} />
        </div>

        <div slot="error">
          Error fetching process list: {error}
        </div>
      </RemoteData>
      <IfEnabled feature="export procfile">
        <ProcfileChecker
          {procfileExport}
          clickHandler={async () => {
            if (procfileExport == null) {
              return;
            }
            await workspace.writeFile('Procfile', procfileExport);
            checkProcfile();
          }}
        />
      </IfEnabled>
    </section>
  </Panel>
{/await}

<style>
  #add-component {
    background: none;
    font-size: 0.9em;
    color: var(--grey-5-color);
    border: 1px dashed var(--grey-c-color);
    width: calc(100% + 8px);
    display: flex;
    align-items: center;
    border-radius: 4px;
    gap: 6px;
    height: 40px;
    margin: 12px -4px;
    padding: 0 8px;
  }

  #add-component:hover,
  #add-component:focus,
  #add-component:focus-within {
    background: var(--grey-e-color);
    color: var(--strong-color);
  }

  #add-component :global(svg) {
    height: 16px;
  }

  #add-component :global(*) {
    fill: currentColor;
  }
  .actions {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-right: 16px;
  }
  .actions span {
    color: var(--grey-7-color);
  }

  .menu {
    outline: none;
    position: relative;
  }

  .menu:focus :global(nav),
  .menu:focus-within :global(nav) {
    display: block;
  }
</style>
