<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import ComponentCard from './ComponentCard.svelte';
</script>

<Meta title="ComponentCard" component={ComponentCard} />

<Template let:args>
  <ComponentCard {...args} />
</Template>

<Story
  name="Default"
  args={{
    detailsPath: '/#',
    editPath: '/#',
    setRun: () => {},
    setLogsVisible: () => {},
    deleteComponent: () => {},
    name: 'Component',
    logsVisible: true,
    reconciling: false,
    running: true,
  }}
/>
