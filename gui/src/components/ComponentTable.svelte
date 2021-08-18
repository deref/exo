<script lang="ts">
  import type { SvelteConstructor } from '../lib/svelte';
  import CheckeredTableWrapper from './CheckeredTableWrapper.svelte';
  import ErrorLabel from './ErrorLabel.svelte';
  import Spinner from './Spinner.svelte';

  type Item = $$Generic;

  interface Column<T> {
    title: string;
    component: SvelteConstructor<{ value?: T }>;
    getValue: (item: Item) => T;
  }

  export let load: () => Promise<Item[]>;

  export let columns: Column<any>[];

  const components = load();
</script>

{#await components}
  <Spinner />
{:then components}
  {#if components.length === 0}
    <div>No records</div>
  {:else}
    <CheckeredTableWrapper>
      <table>
        <thead>
          <tr>
            {#each columns as column}
              <th>
                {column.title}
              </th>
            {/each}
          </tr>
        </thead>
        <tbody>
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
            </tr>
          {/each}
        </tbody>
      </table>
    </CheckeredTableWrapper>
  {/if}
{:catch ex}
  <ErrorLabel value={ex} />
{/await}
