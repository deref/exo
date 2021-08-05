<script lang="ts">
  import type { SvelteConstructor } from '../lib/svelte';
  import ErrorLabel from './ErrorLabel.svelte';
  import Spinner from './mono/spinner.svelte';

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
  {/if}
{:catch ex}
  <ErrorLabel value={ex} />
{/await}

<style>
  table {
    border: none;
    border-collapse: collapse;
    border-radius: 0.25em;
    box-shadow: 0 0.33px 0 1px hsla(0, 0%, 100%, 0.15),
      0 6px 9px -4px rgba(0, 0, 0, 0.2), 0 0.4px 0 0.8px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }

  td,
  th {
    text-align: left;
    padding: 5px 15px;
  }

  th:nth-child(2n) {
    background: #f9f9f9;
  }

  tbody tr:nth-child(2n + 1) {
    background: #eeeeee;
  }

  tbody tr:nth-child(2n) {
    background: #ffffff;
  }

  tbody tr:nth-child(2n + 1) td:nth-child(2n) {
    background: #e7e7e7;
  }

  tbody tr:nth-child(2n) td:nth-child(2n) {
    background: #f9f9f9;
  }

  tr:first-child > * {
    padding-top: 10px;
  }

  tr:last-child > * {
    padding-bottom: 10px;
  }
</style>
