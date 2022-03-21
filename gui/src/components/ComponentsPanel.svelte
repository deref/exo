<script lang="ts" context="module">
  import type { Component } from './processes/ComponentStack.svelte';

  export type Stack = {
    id: string;
    displayName: string;
    components: Component[];
    detailsUrl: string;
    newComponentUrl: string;
  };

  export type { Component };
</script>

<script lang="ts">
  import Icon from './Icon.svelte';
  import Panel from './Panel.svelte';
  import IfEnabled from './IfEnabled.svelte';
  import IconButton from './IconButton.svelte';
  import ContextMenu from './ContextMenu.svelte';
  import MenuItem from './MenuItem.svelte';
  import ComponentStack from './processes/ComponentStack.svelte';
  import * as router from 'svelte-spa-router';

  import { modal } from '../lib/modal';
  import { bind } from '../components/modal/Modal.svelte';
  import ModalDialogue from '../components/modal/ModalDialogue.svelte';

  export let destroyStack: () => Promise<void>;
  export let setLogsVisible: (id: string, value: boolean) => Promise<void>;
  export let setRun: (id: string, value: boolean) => Promise<void>;
  export let disposeComponent: (id: string) => Promise<void>;

  export let stack: Stack;

  const confirmDestroy = () => {
    modal.set(
      bind(ModalDialogue, {
        h3: 'Delete stack?',
        bodyText: `Are you sure you want to delete the ${stack.displayName} stack?`,
        danger: true,
        actionLabel: 'Yes, delete',
        onOkay: async () => {
          await destroyStack(); // TODO: try/catch.
          router.push('/');
        },
      }),
    );
  };
</script>

<Panel title={stack.displayName} backRoute="/" --panel-padding="0 1rem">
  <div class="actions" slot="actions">
    <span>Logs</span>
    <div class="menu" tabindex="0">
      <IconButton
        glyph="Ellipsis"
        tooltip="Stack actions..."
        on:click={() => {}}
      />

      <ContextMenu title={stack.displayName}>
        <MenuItem glyph="Details" href={stack.detailsUrl}>
          View details
        </MenuItem>
        <MenuItem glyph="Add" href={stack.newComponentUrl}>
          Add component
        </MenuItem>
        <MenuItem glyph="Delete" danger on:click={() => confirmDestroy()}>
          Destroy stack
        </MenuItem>
      </ContextMenu>
    </div>
  </div>

  <section>
    <button
      id="add-component"
      on:click={() => {
        router.push(stack.newComponentUrl);
      }}
    >
      <Icon glyph="Add" /> Add component
    </button>
    <ComponentStack
      components={stack.components}
      dispose={disposeComponent}
      {setLogsVisible}
      {setRun}
    />
    <IfEnabled feature="export procfile" />
  </section>
</Panel>

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
