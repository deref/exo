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
    stack: {
      id: 'stack1',
      detailsUrl: '/#stacks/1',
      newComponentUrl: '/#stacks/1/new-component',
      displayName: 'example',
      components: [],
    },
  }}
/>

<Story
  name="Non-empty"
  args={{
    stack: {
      id: 'stack2',
      displayName: 'sample',
      detailsUrl: '/#stack2',
      newComponentUrl: '/#stacks/2/new-component',
      components: [
        {
          id: 'backend',
          name: 'Backend',
          url: '/#components/backend',
          editUrl: '/#components/backend/edit',
          ...componentDefaults,
        },
        {
          id: 'frontend',
          name: 'Frontend',
          url: '/#components/frontend',
          editUrl: '/#components/frontend/edit',
          ...componentDefaults,
        },
      ],
    },
  }}
/>
