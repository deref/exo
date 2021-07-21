<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Button from '../components/Button.svelte';
  import Textbox from '../components/Textbox.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import { api, isClientError } from '../lib/api';
  import * as qs from 'qs'
  import * as router from 'svelte-spa-router'
  import { querystring } from 'svelte-spa-router'

  const query = qs.parse($querystring);

  let root = typeof query.root === 'string' ? query.root : '';
  let error: Error | null = null;
</script>

<Layout>
<section>
  <form on:submit|preventDefault={async () => {
    error = null;
    try {
      const id = await api.kernel.createWorkspace(root)
      router.push(`/workspaces/${encodeURIComponent(id)}`);
    } catch (ex) {
      if (!isClientError(ex)) {
        throw ex;
      }
      error = ex;
    }
  }}>
    <label for="root">Root:</label>
    <Textbox bind:value={root} name="root" id="root"/>
    <Button type="submit">Create Workspace</Button>
  </form>
  <ErrorLabel value={error}/>
</section>
</Layout>

<style>
  section {
    display: flex;
    flex-direction: column;
    height: 100%;
    width: 100%;
    align-items: center;
    justify-content: center;
  }

  form {
    padding: 20px;
    display: flex;
    flex-direction: row;
    gap: 20px;
    align-items: center;
  }
</style>
