<script lang="ts">
  import IconButton from './IconButton.svelte';
  import CheckeredTable from './CheckeredTable.svelte';
  import type { IconGlyph } from './Icon.svelte';
  import type { SvelteConstructor } from '../lib/svelte';

  type Item = $$Generic;

  interface Column<T> {
    title: string;
    label: SvelteConstructor<{ value?: T }>;
    getValue: (item: Item) => T;
  }

  interface Action {
    tooltip: string;
    glyph: IconGlyph;
    execute: (item: Item) => void;
  }

  export let columns: Column<unknown>[];
  export let actions: Action[] | undefined;
  export let items: Item[];
  export let plural = 'records';
</script>

{#if items.length === 0}
  <div>No {plural}</div>
{:else}
  <CheckeredTable>
    <svelte:fragment slot="head">
      {#each columns as column}
        <th>
          {column.title}
        </th>
      {/each}
      {#if actions && actions.length > 0}
        <th />
      {/if}
    </svelte:fragment>
    {#each items as item}
      <tr>
        {#each columns as column}
          <td>
            <svelte:component
              this={column.label}
              value={column.getValue(item)}
            />
          </td>
        {/each}
        {#if actions && actions.length > 0}
          <td class="actions">
            <div>
              {#each actions as action}
                <IconButton
                  tooltip={action.tooltip}
                  glyph={action.glyph}
                  on:click={() => {
                    action.execute(item);
                  }}
                />
              {/each}
            </div>
          </td>
        {/if}
      </tr>
    {/each}
  </CheckeredTable>
{/if}

<style>
  td.actions {
    padding: 5px 8px;
  }

  td.actions div {
    display: flex;
    align-items: center;
  }

  td.actions :global(button) {
    width: 30px;
    height: 30px;
  }
</style>
