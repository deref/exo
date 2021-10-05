<script lang="ts">
  import Icon from '../components/Icon.svelte';
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import Textbox from '../components/Textbox.svelte';
  import EditAs from '../components/form/EditAs.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import TextEditor from '../components/TextEditor.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import type { ComponentDescription } from '../lib/api';
  import { api } from '../lib/api';

  export let params = { workspace: '', component: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceComponentsRoute = `/workspaces/${encodeURIComponent(
    workspaceId,
  )}/components`;

  const componentId = params.component;

  let error: Error | null = null;

  let mode = 'compose';
  const editorModes = [
    {
      id: 'compose',
      name: 'Compose YAML',
    },
    // {
    //   id: 'form',
    //   name: 'Form',
    // },
  ];

  let name: string = '';
  let spec: string = '';

  const getComponent = async () =>
    (await workspace.describeComponents({ refs: [componentId] }))[0];

  const pageTitle = (component: ComponentDescription) =>
    `Edit ${component.type} “${component.name}”`;

  const componentGlyph = (cType: string) => {
    switch (cType) {
      case 'process':
        return 'Layers';
      default:
        return 'LogoDocker';
    }
  };
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  {#await getComponent()}
    <CenterFormPanel title="Loading..." backRoute={workspaceComponentsRoute}>
      <Spinner />
    </CenterFormPanel>
  {:then component}
    {#if component !== undefined}
      <CenterFormPanel
        title={pageTitle(component)}
        backRoute={workspaceComponentsRoute}
      >
        <h1>
          <Icon glyph={componentGlyph(component.type)} />{pageTitle(component)}
        </h1>
        <form on:submit|preventDefault={async () => {}}>
          <div class="group">
            <label for="name">Name:</label>
            <Textbox
              id="name"
              name="name"
              value={component.name}
              --input-width="100%"
            />
          </div>

          <EditAs bind:mode {editorModes} />
          {#if mode === 'compose'}
            <div class="group">
              <label for="spec">Spec:</label>
              <TextEditor id="spec" bind:value={spec} language="yaml" />
            </div>
            <details>
              <summary>Show/hide example</summary>
              <slot />
            </details>
          {:else}
            <!-- GUI form edit mode -->
          {/if}
          <SubmitButton>Save changes</SubmitButton>
          <div style="margin: 24px 0;">
            <ErrorLabel value={error} />
          </div>
        </form>
      </CenterFormPanel>
    {/if}
  {/await}
</Layout>
