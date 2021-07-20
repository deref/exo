<script lang="ts">
  import { api } from '../lib/api';
  import * as qs from 'qs'
  import * as router from 'svelte-spa-router'
  import { querystring } from 'svelte-spa-router'

  const query = qs.parse($querystring);

  let root = typeof query.root === 'string' ? query.root : '';
</script>

<main>
  <form on:submit|preventDefault={async () => {
    const id = await api.kernel.createWorkspace(root)
    router.push(`/workspaces/${encodeURIComponent(id)}`);
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
