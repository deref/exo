<script lang="ts">
  import Link from '../components/Link.svelte';
  import Code from '../components/Code.svelte';
  import Panel from '../components/Panel.svelte';
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import IconButton from '../components/IconButton.svelte';
  import ResetSVG from '../components/mono/ResetSVG.svelte';
  import { guiTheme, setGuiTheme } from '../lib/gui-theme';

  const resetAllPreferences = () => {
    guiTheme.auto();
    setGuiTheme($guiTheme);
  };
</script>

<Layout>
  <Panel title="Preferences" --panel-padding="2rem">
    <div slot="actions">
      <IconButton tooltip="Reset to defaults" on:click={resetAllPreferences}>
        <ResetSVG />
      </IconButton>
    </div>
    <div class="center-form">
      <div class="group">
        <div class="group-header">
          <h2>Theme &amp; GUI</h2>
        </div>
        <div class="button-row">
          <Button
            on:click={() => {
              guiTheme.auto();
              setGuiTheme($guiTheme);
            }}
            inset={$guiTheme === 'auto'}
          >
            Auto
          </Button>
          <Button
            on:click={() => {
              guiTheme.light();
              setGuiTheme($guiTheme);
            }}
            inset={$guiTheme === 'light'}
          >
            Light
          </Button>
          <Button
            on:click={() => {
              guiTheme.dark();
              setGuiTheme($guiTheme);
            }}
            inset={$guiTheme === 'dark'}
          >
            Dark
          </Button>
          <Button
            on:click={() => {
              guiTheme.black();
              setGuiTheme($guiTheme);
            }}
            inset={$guiTheme === 'black'}
          >
            Black
          </Button>
        </div>
      </div>
      <hr />
      <p>
        More preferences will be available here soon. In the meantime, please
        refer to the following help:
        <Link href="https://github.com/deref/exo#telemetry">
          disabling telemetry manually
        </Link>
        with the <Code>{'config.toml'}</Code> file.
      </p>
    </div>
  </Panel>
</Layout>

<style>
  .group-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

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
  }

  hr {
    border: none;
    display: block;
    height: 1px;
    background: var(--grey-c-color);
    margin: 32px 0;
  }
</style>
