<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import ComponentStack from './ComponentStack.svelte';

  const componentDefaults = {
    logsVisible: true,
    reconciling: false,
    running: true,
  };

  const showPath = (id: string) => `/#show-${id}`;
  const editPath = (id: string) => `/#edit-${id}`;
</script>

<Meta title="ComponentStack" component={ComponentStack} />

<Template let:args>
  <ComponentStack {...args} />
</Template>

<Story
  name="Default"
  args={{
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
        running: false,
      },
      {
        id: 'helper',
        name: 'helper',
        ...componentDefaults,
        reconciling: true,
        logsVisible: false,
      },
    ],
    showPath,
    editPath,
  }}
/>
