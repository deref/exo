<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import ComponentsPanel from './ComponentsPanel.svelte';

  const componentDefaults = {
    running: true,
    logsVisible: true,
  };
</script>

<Meta title="ComponentsPanel" component={ComponentsPanel} />

<Template let:args>
  <ComponentsPanel {...args} />
</Template>

<Story
  name="Empty"
  args={{
    workspace: {
      id: 'workspace1',
      displayName: 'example',
      components: [],
    },
  }}
/>

<Story
  name="Non-empty"
  args={{
    workspace: {
      id: 'workspace2',
      displayName: 'sample',
      components: [
        {
          id: 'backend',
          name: 'Backend',
          ...componentDefaults,
        },
        {
          id: 'frontend',
          name: 'Frontend',
          ...componentDefaults,
        },
      ],
    },
  }}
/>
