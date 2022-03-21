<script lang="ts">
  import Spinner from './Spinner.svelte';
  import ErrorLabel from './ErrorLabel.svelte';
  import IconButton from './IconButton.svelte';
  import CheckeredTable from './CheckeredTable.svelte';
  import type { IconGlyph } from './Icon.svelte';
  import type { SvelteConstructor } from '../lib/svelte';

  type Item = $$Generic;

  interface Column<T> {
    title: string;
    component: SvelteConstructor<{ value?: T }>;
    getValue: (item: Item) => T;
  }

  interface Action {
    tooltip: string;
    glyph: IconGlyph;
    execute(component: any): any;
  }

  export let load: () => Promise<Item[]>;

  export let columns: Column<any>[];
  export let actions: Action[] | undefined;

  let componentsPromise = load();
</script>

{#await componentsPromise}
  <Spinner />
{:then components}
  {#if components.length === 0}
    <div>No records</div>
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
      {#each components as component}
        <tr>
          {#each columns as column}
            <td>
              <svelte:component
                this={column.component}
                value={column.getValue(component)}
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
                      action.execute(component);
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
{:catch ex}
  <ErrorLabel value={ex} />
{/await}

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
