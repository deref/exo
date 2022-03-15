<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import LogScroller from './LogScroller.svelte';
</script>

<Meta title="LogScroller" component={LogScroller} />

<Template let:args>
  <div>
    <LogScroller {...args} />
  </div>
</Template>

<Story
  name="Default"
  args={{
    events: [
      {
        id: 'one',
        timestamp: '2021-09-24T21:15:08.970Z',
        sourceId: 'some-program',
        sourceName: 'some-program',
        message: 'Hello there',
      },
      {
        id: 'two',
        timestamp: '2021-09-24T21:15:09.970Z',
        sourceId: 'another-app',
        sourceName: 'another-app',
        message: 'How are you?',
      },
    ],
  }}
/>

<Story
  name="Scrollback"
  args={{
    events: (() => {
      const events = [];
      for (let i = 1000; i < 2000; i++) {
        events.push({
          id: `event-${i}`,
          timestamp: new Date(1632518109970 + i * 1000).toISOString(),
          stream: `component-${Math.floor(Math.sin(i) * 10000) % 5}`,
          message: `message ${i}`,
        });
      }
      return events;
    })(),
  }}
/>

<style>
  div {
    height: 200px;
    border: solid 1px #333;
  }
</style>
