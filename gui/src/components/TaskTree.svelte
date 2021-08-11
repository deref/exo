<script lang="ts" context="module">
  export type Status = 'pending' | 'running' | 'success' | 'failure';

  export interface TaskNode {
    id: string;
    name: string;
    status: Status;
    children: TaskNode[];
  }
</script>

<script lang="ts">
  export let value: TaskNode;
</script>

<div class="container">
  <div class="row">
    <div class="status">{value.status}</div>
    <div class="message">{value.name}</div>
  </div>
  {#if value.children.length > 0}
    <ul>
      {#each value.children as child (child.id)}
        <li><svelte:self value={child} /></li>
      {/each}
    </ul>
  {/if}
</div>

<style>
  .row {
    display: flex;
    flex-direction: row;
  }

  .status {
    padding-right: 16px;
  }
</style>
