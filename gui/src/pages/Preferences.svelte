<script lang="ts">
  import Panel from '../components/Panel.svelte';
  import Layout from '../components/Layout.svelte';

  let theme = 'Auto';

  const setTheme = (s: string) => (theme = s);
</script>

<Layout>
  <Panel title="Preferences">
    <div class="center-form">
      <!-- User preferences here -->
      <h2>Theme</h2>
      <div class="button-row">
        <button
          on:click={() => {
            setTheme('Auto');
          }}
          class:active={theme === 'Auto'}>Auto</button
        >
        <button
          on:click={() => {
            setTheme('Light');
          }}
          class:active={theme === 'Light'}>Light</button
        >
        <button
          on:click={() => {
            setTheme('Dark');
          }}
          class:active={theme === 'Dark'}>Dark</button
        >
        <button
          on:click={() => {
            setTheme('Black');
          }}
          class:active={theme === 'Black'}>Black</button
        >
      </div>
    </div>
  </Panel>
</Layout>

<style>
  .center-form {
    max-width: 640px;
    margin: 0 auto;
  }

  button.active {
    background: blue;
  }
</style>
