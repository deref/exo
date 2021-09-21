<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';

  export let params = { starter: '' };

  const { starter } = params;

  let name = starter;
  let error = 'ee';

  let workingDirectory = '~';

  const setWorkingDirectory = (dir: string) => {
    workingDirectory = dir;
  };

  const describeWorkingDirectory = (dir: string) => ({
    children: directories[dir],
    parent: dir
      .split('/')
      .filter((_, j) => j !== dir.split('/').length - 1)
      .join('/'),
  });

  const directories = {
    '~': ['john', 'os', 'bin'],
    '~/john': ['dev', 'downloads', 'music', 'docs'],
    '~/john/dev': [
      'myorg',
      'calculator-app',
      'exo',
      'exo-website',
      'rust-wallet',
      'todo-app',
    ],
    '~/john/dev/myorg': [],
  };
</script>

<Layout>
  <CenterFormPanel title={`New project: ${starter}`} backRoute="#/new-project">
    <form on:submit|preventDefault={async () => {}}>
      <h1>New project</h1>
      <label for="name">Name:</label>
      <Textbox
        bind:value={name}
        name="name"
        id="name"
        placeholder={starter}
        --input-width="100%"
      />
      <div style="height:32px" />
      <label for="root">Root:</label>
      <h2><span>{workingDirectory}<span>/{name}</span></span></h2>
      <button
        on:click={() =>
          setWorkingDirectory(
            describeWorkingDirectory(workingDirectory).parent,
          )}>..</button
      >
      {#each describeWorkingDirectory(workingDirectory).children as child}
        <button
          on:click={() => setWorkingDirectory(workingDirectory + '/' + child)}
          >{child}</button
        >
      {/each}
      <SubmitButton>Create project</SubmitButton>
    </form>
    <ErrorLabel value={error} />
  </CenterFormPanel>
</Layout>

<style>
  h2 span span {
    color: var(--grey-9-color);
  }

  button {
    background: var(--button-background);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 3px;
    padding: 6px 12px;
    position: relative;
    display: grid;
    width: 100%;
    grid-template-columns: max-content 2fr;
    align-items: center;
    gap: 12px;
    margin-top: 4px;
  }

  button > b {
    max-width: 6em;
  }

  button > * {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  button:hover {
    background: var(--button-hover-background);
    box-shadow: var(--button-hover-shadow);
  }

  button:active {
    background: var(--button-active-background);
    box-shadow: var(--button-active-shadow);
  }
</style>
