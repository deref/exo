<script lang="ts">
  import { api } from '../lib/api';
  import { workspaceId } from '../lib/workspaces/store'

  let root: string = '';

  // XXX Janky parameter handling.
  const match = window.location.search.match(/\?root=(.*)/);
  if (match) {
    root = decodeURIComponent(match[1]);
  }
</script>

<main>
  <form on:submit|preventDefault={async e => {
    const id = await api.kernel.createWorkspace(root)
    workspaceId.set(id);
    history.pushState(null, null, `/workspaces/${encodeURIComponent(id)}`);
  }}>
    <label>
      Root:
      <input bind:value={root} name="root"/>
    </label>
    <button type="submit">Create Workspace</button>
  </form>
</main>

<style>
  main {
    padding: 20px;
  }
</style>
