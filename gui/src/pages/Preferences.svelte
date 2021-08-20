<script lang="ts">
  import Panel from '../components/Panel.svelte';
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';

  let theme = 'Auto';

  const setTheme = (s: string) => (theme = s);
</script>

<Layout>
  <Panel title="Preferences" --panel-padding="2rem">
    <div class="center-form">
      <!-- User preferences here -->
      <h2>Theme &amp; GUI</h2>
      <div class="button-row">
        <Button
          on:click={() => {
            setTheme('Auto');
          }}
          inset={theme === 'Auto'}
        >
          Auto
        </Button>
        <Button
          on:click={() => {
            setTheme('Light');
          }}
          inset={theme === 'Light'}
        >
          Light
        </Button>
        <Button
          on:click={() => {
            setTheme('Dark');
          }}
          inset={theme === 'Dark'}
        >
          Dark
        </Button>
        <Button
          on:click={() => {
            setTheme('Black');
          }}
          inset={theme === 'Black'}
        >
          Black
        </Button>
      </div>
      <hr />
    </div>
  </Panel>
</Layout>

<style>
  .center-form {
    max-width: 640px;
    margin: 0 auto;
  }

  .button-row {
    display: grid;
    grid-auto-flow: column;
    gap: 12px;
  }

  h2 {
    font-size: 20px;
    font-weight: 500;
    margin: 0;
    margin-bottom: 24px;
  }

  hr {
    border: none;
    display: block;
    height: 1px;
    background: var(--grey-c-color);
    margin: 32px 0;
  }
</style>
