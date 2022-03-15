<script lang="ts">
  import { derived, writable } from 'svelte/store';
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import LogPanel from './LogPanel.svelte';
  import type { StreamStore, Event } from './LogPanel.svelte';

  const allEvents = writable<Event[]>([
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
  ]);
  const filterString = writable('');

  const stream: StreamStore = {
    ...derived([allEvents, filterString], ([events, filterString]) => ({
      events: events.filter(
        (event) => event.message.indexOf(filterString) >= 0,
      ),
      filterString,
    })),
    clearEvents: () => {
      allEvents.set([]);
    },
    setFilterString: filterString.set,
  };
</script>

<Meta title="LogPanel" component={LogPanel} />

<Template let:args>
  <LogPanel {...args} />
</Template>

<Story
  name="Default"
  args={{
    stream,
  }}
/>
