<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import IfEnabled from '../components/IfEnabled.svelte';
</script>

<Layout>
  <h1>Sandbox</h1>

  <div>
    <h2>Feature Flags</h2>
    <IfEnabled feature="sandbox test">
      <p>I should be enabled.</p>
    </IfEnabled>

    <IfEnabled feature="bogus">
      <p>I should NOT be enabled.</p>
    </IfEnabled>
  </div>
</Layout>
