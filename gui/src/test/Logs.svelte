<script lang="ts">
  import FormattedLogMessage from '../components/logs/FormattedLogMessage.svelte';
  import Layout from '../components/Layout.svelte';
</script>

<Layout>
  <section>
    <h1>Logs</h1>
    <h3>Linkify</h3>
    <pre>
      <FormattedLogMessage
        message={'This message uses a   string with a https://deref.io URL in it and formats  an anchor tag with the link.'}
      />
    </pre>
  </section>
</Layout>

<style>
  section {
    padding: 60px;
  }
</style>
