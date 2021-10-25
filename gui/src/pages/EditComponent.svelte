<script lang="ts">
  import Icon from '../components/Icon.svelte';
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import Textbox from '../components/Textbox.svelte';
  import EditAs from '../components/form/EditAs.svelte';
  import CodeBlock from '../components/CodeBlock.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import TextEditor from '../components/TextEditor.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import type { IconGlyph } from '../components/Icon.svelte';
  import type { ComponentDescription } from '../lib/api';
  import { api, isClientError } from '../lib/api';
  import * as router from 'svelte-spa-router';

  export let params = { workspace: '', component: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceComponentsRoute = `/workspaces/${encodeURIComponent(
    workspaceId,
  )}/components`;

  const componentId = params.component;

  let error: Error | null = null;

  let mode = 'raw';
  const editorModes = [
    {
      id: 'raw',
      name: 'Raw Spec',
    },
  ];

  let name: string = '';
  let spec: string = '';

  const getComponent = async () => {
    const component = (
      await workspace.describeComponents({ refs: [componentId] })
    )[0];

    name = component.name;
    spec = component.spec;

    return component;
  };

  const pageTitle = (component: ComponentDescription) =>
    `Edit ${component.type} “${component.name}”`;

  interface ComponentType {
    glyph: IconGlyph;
    example: string;
  }

  // TODO: Extract these to a registry.
  const componentTypes: Record<string, ComponentType> = {
    process: {
      glyph: 'Layers',
      example: `{
  "directory": "/",
  "environment": {},
  "program": "npm",
  "arguments": ["run", "dev"]
}`,
    },
    container: {
      glyph: 'LogoDocker',
      example: `image: postgres:9.4
environment:
    POSTGRES_USER: "postgres"
    POSTGRES_PASSWORD: "postgres"`,
    },
    network: {
      glyph: 'LogoDocker',
      example: `name: mynetwork`,
    },
    volume: {
      glyph: 'LogoDocker',
      example: `# An empty YAML is a valid volume.`,
    },
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
          <Icon glyph={componentTypes[component.type].glyph} />
          {pageTitle(component)}
        </h1>
        <form
          on:submit|preventDefault={async () => {
            try {
              await workspace.updateComponent(component.id, name, spec);

              router.push(workspaceComponentsRoute);
            } catch (ex) {
              if (!isClientError(ex)) {
                throw ex;
              }
              error = ex;
            }
          }}
        >
          <div>
            <label for="name">Name:</label>
            <Textbox
              id="name"
              name="name"
              bind:value={name}
              --input-width="100%"
            />
          </div>

          <EditAs bind:mode {editorModes} />
          <div>
            <label for="spec">Spec:</label>
            <TextEditor id="spec" bind:value={spec} language="yaml" />
          </div>
          <details>
            <summary>Show/hide example</summary>
            <CodeBlock code={componentTypes[component.type].example} />
          </details>
          <SubmitButton>Save changes</SubmitButton>
          <div style="margin: 24px 0;">
            <ErrorLabel value={error} />
          </div>
        </form>
      </CenterFormPanel>
    {/if}
  {/await}
</Layout>
