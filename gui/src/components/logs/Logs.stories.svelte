<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import Logs from './Logs.svelte';
</script>

<Meta title="Logs" component={Logs} />

<Template let:args>
  <Logs {...args} />
</Template>

<Story
  name="Logs"
  args={{
    events: [
      {
        id: 'one',
        timestamp: '2021-09-24T21:15:08.970Z',
        stream: 'some-program',
        message: 'Hello there',
      },
      {
        id: 'two',
        timestamp: '2021-09-24T21:15:09.970Z',
        stream: 'another-app',
        message: 'How are you?',
      },
    ],
  }}
/>
