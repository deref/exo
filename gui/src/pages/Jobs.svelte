<script lang="ts">
  import Link from '../components/Link.svelte';
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import { api } from '../lib/api';

  const tasks = api.kernel.describeTasks();
  const jobs = tasks.then((tasks) =>
    tasks.filter((task) => task.parentId === null),
  );
</script>

<Layout>
  <Panel title="Jobs" backRoute="/">
    <!-- TODO: common loader component -->
    {#await jobs}
      Loading...
    {:then jobs}
      {#if jobs.length === 0}
        <i>None</i>
      {:else}
        <table>
          <thead>
            <tr>
              <td>ID</td>
              <td>Status</td>
              <td>Message</td>
              <td>Created</td>
              <td>Finished</td>
            </tr>
          </thead>
          {#each jobs as job (job.id)}
            <tr>
              <td>
                <Link href={`#/jobs/${job.id}`}><pre>{job.id}</pre></Link>
              </td>
              <td>
                {job.status}
              </td>
              <td>
                {job.message}
              </td>
              <td>
                {job.created}
              </td>
              <td>
                {job.finished}
              </td>
            </tr>
          {/each}
        </table>
      {/if}
    {:catch}
      Error.
    {/await}
  </Panel>
</Layout>
