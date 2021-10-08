<script lang="ts">
  import type { RequestLifecycle } from '../lib/api';

  type Item = $$Generic;
  interface $$Slots {
    pending: {};
    success: { data: Item; loading: boolean };
    error: { error: string };
    default: {
      data: Item;
      error: string;
    };
  }
  interface $$Props {
    data: RequestLifecycle<Item>;
  }

  export let data: RequestLifecycle<Item>;
</script>

{#if data.stage === 'pending'}
  <slot name="pending">Loading...</slot>
{:else if data.stage === 'success' || data.stage === 'refetching'}
  <slot name="success" data={data.data} loading={data.stage == 'refetching'}>
    Missing success slot!
  </slot>
{:else if data.stage === 'error'}
  <slot name="error" error={data.message}>
    Error: {data.message}.
  </slot>
{/if}
