<script lang="ts">
  import Panel from '../components/Panel.svelte';
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import IconButton from '../components/IconButton.svelte';
  import Textbox from '../components/Textbox.svelte';
  import ResetSVG from '../components/mono/ResetSVG.svelte';

  let theme = 'Auto';

  const setTheme = (s: string) => {
    theme = s;

    switch (s) {
      case 'Auto':
        document.body.classList.remove('light', 'dark', 'black');
        document.body.classList.add('auto');
      case 'Light':
        document.body.classList.remove('auto', 'dark', 'black');
        document.body.classList.add('light');
      case 'Dark':
        document.body.classList.remove('light', 'auto', 'black');
        document.body.classList.add('dark');
      case 'Black':
        document.body.classList.remove('light', 'dark', 'auto');
        document.body.classList.add('black');
    }
  };
</script>

<Layout>
  <Panel title="Preferences" --panel-padding="2rem">
    <div slot="actions">
      <IconButton
        tooltip="Reset to defaults"
        on:click={() => {
          /* Reset preferences */
        }}
      >
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
      </div>
      <hr />
      <div class="group">
        <div class="group-header">
          <h2>Typography</h2>
          <div class="add-new">
            <span>+ Add new</span>
            <ul class="add-new-menu">
              <li>code_ligatures</li>
              <li>font_code</li>
              <li>font_sans</li>
            </ul>
          </div>
        </div>
        <div class="input-text">
          <h3 title="code_ligatures">code_ligatures</h3>
          <Textbox value="false" --input-width="100%" />
        </div>
      </div>
      <hr />
      <div class="group">
        <div class="group-header">
          <h2>Telemetry</h2>
          <div class="add-new">
            <span>+ Add new</span>
            <ul class="add-new-menu">
              <li>disable_telemetry</li>
            </ul>
          </div>
        </div>
        <div class="input-text">
          <h3 title="disable_telemetry">disable_telemetry</h3>
          <Textbox value="true" --input-width="100%" />
        </div>
      </div>
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

  .add-new {
    position: relative;
  }

  .add-new > span {
    padding: 8px 16px;
    margin: -8px -16px;
    cursor: pointer;
  }

  .add-new:not(:hover) .add-new-menu {
    display: none;
  }

  .add-new:not(:hover) > span {
    color: var(--grey-7-color);
  }

  .add-new-menu {
    position: absolute;
    top: 100%;
    right: -16px;
    margin: 0;
    padding: 8px 0;
    border-radius: 6px;
    background: var(--primary-bg-color);
    box-shadow: var(--button-hover-shadow);
    list-style: none;
  }

  .add-new-menu li {
    padding: 6px 16px;
    font-family: var(--font-mono);
    color: var(--grey-5-color);
    font-size: 15px;
    font-weight: 450;
    cursor: pointer;
  }

  .add-new-menu li:hover {
    background: var(--grey-e-color);
    color: var(--strong-color);
  }

  .add-new-menu li:active {
    background: var(--grey-d-color);
    color: var(--strong-color);
  }

  .group h3 {
    font-family: var(--font-mono);
    font-size: 15px;
    font-weight: 450;
    color: var(--grey-5-color);
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .input-text {
    display: grid;
    grid-template-columns: 1fr 2fr;
    gap: 40px;
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
