<!-- Captive network error page. See note [ONLINE]. -->
<script lang="ts">
  import Code from '../components/Code.svelte';
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import { connected, authenticated, query } from '../lib/graphql';
  import Spinner from '../components/Spinner.svelte';

  // Ping server periodically, rely on the Apollo client's onError link to
  // reset the online state.
  query(
    `#graphql
    { __typename }
  `,
    { pollInterval: 1000 /* ms */ },
  );
</script>

<Layout loader={null}>
  <Panel>
    <div>
      {#if !$connected}
        <h3>
          <Code>exo</Code> server is offline
        </h3>
        <p>
          Run <Code>exo daemon</Code> in your terminal to start.
        </p>
      {:else if !$authenticated}
        <h3>Not authenticated</h3>
        <p>
          Run <Code>exo gui</Code> in your terminal to start.
        </p>
      {:else}
        <!-- A higher-up component prevents this case from occuring. -->
        <Spinner />
      {/if}
    </div>
  </Panel>
</Layout>

<style>
  div {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    font-size: 18px;
    height: 100%;
  }
</style>
