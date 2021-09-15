<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Button from '../components/Button.svelte';
  import Textbox from '../components/Textbox.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api, isClientError } from '../lib/api';
  import * as qs from 'qs';
  import * as router from 'svelte-spa-router';
  import { querystring } from 'svelte-spa-router';

  const query = qs.parse($querystring);

  let root = typeof query.root === 'string' ? query.root : '';
  let error: Error | null = null;
</script>

<Layout>
  <CenterFormPanel title="New project" backRoute="/">
    <form
      on:submit|preventDefault={async () => {
        error = null;
        let workspaceId;
        try {
          workspaceId = await api.kernel.createWorkspace(root);
          router.push(`/workspaces/${encodeURIComponent(workspaceId)}`);
        } catch (ex) {
          if (!isClientError(ex)) {
            throw ex;
          }
          error = ex;
        }
        // XXX Hack to address lack of GUI for applying procfiles, etc.
        try {
          await api.workspace(workspaceId).apply();
        } catch (ex) {
          // Swallow error.
          console.error(ex);
        }
      }}
    >
      <h1>New project</h1>
      <p>Select your project root directory to create a new workspace.</p>
      <label for="root">Root:</label>
      <Textbox
        bind:value={root}
        name="root"
        id="root"
        placeholder="/home/user/path/to/project"
        --input-width="100%"
      />
      <div class="buttons">
        <Button type="submit">Create project</Button>
      </div>
    </form>
    <ErrorLabel value={error} />
  </CenterFormPanel>
</Layout>

<style>
  .buttons {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    margin: 24px 0;
  }
</style>
