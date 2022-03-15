<script lang="ts" context="module">
  export type Component = {
    id: string;
    name: string;
    reconciling: boolean;
    running: boolean;
    logsVisible: boolean;
    url: string;
    editUrl: string;
  };
</script>

<script lang="ts">
  import * as router from 'svelte-spa-router';
  import IconButton from '../IconButton.svelte';
  import ContextMenu from '../ContextMenu.svelte';
  import MenuItem from '../MenuItem.svelte';
  import CheckboxButton from '../CheckboxButton.svelte';
  import ComponentControls from './ComponentControls.svelte';
  import { logStyleFromHash } from '../../lib/color';

  const { link } = router;

  export let component: Component;
  const { name, reconciling, running, logsVisible } = component;

  export let setRun: (value: boolean) => void;
  export let setLogsVisible: (value: boolean) => void;
  export let dispose: () => void;
</script>

<div class="card" style={logStyleFromHash(name)}>
  <div>
    <ComponentControls {setRun} statusPending={reconciling} {running} />
  </div>

  <div>
    <a class="component-name" use:link href={component.url}>
      {name}
    </a>
  </div>

  <div class="checkbox">
    <CheckboxButton
      tooltip={logsVisible ? 'Hide logs' : 'Show logs'}
      on:click={() => {
        setLogsVisible(!logsVisible);
      }}
      active={logsVisible}
    />
  </div>

  <div class="actions" tabindex="0">
    <IconButton glyph="Ellipsis" />

    <ContextMenu title={name}>
      <MenuItem glyph="Details" href={component.url}>View details</MenuItem>
      <MenuItem glyph="Edit" href={component.editUrl}>Edit component</MenuItem>
      <MenuItem
        glyph="Logs"
        on:click={() => {
          setLogsVisible(!logsVisible);
        }}
      >
        Toggle logs visibility
      </MenuItem>
      <MenuItem
        glyph="Delete"
        danger
        on:click={() => {
          dispose();
          setLogsVisible(false);
        }}
      >
        Delete component
      </MenuItem>
    </ContextMenu>
  </div>
</div>

<style>
  .card {
    box-shadow: var(--card-shadow);
    display: grid;
    grid-template-columns: max-content auto max-content max-content;
    align-items: center;
    padding: 4px;
    margin: 0px -4px;
    border-radius: 4px;
    margin-bottom: 8px;
    border-left: 2px solid var(--log-color);
  }

  .card:hover {
    box-shadow: var(--card-hover-shadow);
  }

  .card > * {
    font-size: inherit;
    font-weight: inherit;
    align-items: center;
    justify-content: center;
  }

  .card .checkbox {
    margin-right: 18px;
  }

  .card > *:nth-child(2) {
    text-align: left;
  }

  .card > *:not(:nth-child(2)) {
    white-space: nowrap;
  }

  .card:not(:hover):not(:focus-within) .actions {
    opacity: 0.333;
  }

  .component-name {
    display: inline-block;
    text-decoration: none;
    margin: 0;
    margin-left: 6px;
    margin-right: 12px;
    line-height: 1;
    font-size: 16px;
    font-weight: 550;
    padding: 4px 7px;
    border-radius: 3px;
    color: var(--log-color);
    outline: none;
  }

  .component-name:hover,
  .component-name:focus {
    color: var(--log-hover-color);
    background: var(--log-bg-hover-color);
  }

  .actions {
    outline: none;
    position: relative;
  }

  .actions:focus :global(nav),
  .actions:focus-within :global(nav) {
    display: block;
  }
</style>
