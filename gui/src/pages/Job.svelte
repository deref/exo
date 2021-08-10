<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import JobTree from '../components/JobTree.svelte';
  import { api } from '../lib/api';

  export let params = { job: '' };

  const jobId = params.job;

  const tasks = api.kernel.describeTasks({ jobIds: [jobId] });
</script>

<Layout showBackButton backButtonRoute="#/jobs">
  <!-- TODO: common loader component -->
  {#await tasks}
    Loading...
  {:then tasks}
    <JobTree {jobId} {tasks} />
  {:catch}
    Error
  {/await}
</Layout>
