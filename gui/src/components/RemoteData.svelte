<script lang="ts">
  import type { RequestLifecycle } from '../lib/api';

  type Item = $$Generic;
  interface $$Slots {
    pending: {};
    success: { data: Item };
    refetching: { data: Item };
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
{:else if data.stage === 'success'}
  <slot name="success" data={data.data}>Missing success slot!</slot>
{:else if data.stage === 'refetching'}
  <slot name="refetching" data={data.data}>
    <slot name="success" data={data.data}>
      Missing refetching and success slots!
    </slot>
  </slot>
{:else if data.stage === 'error'}
  <slot name="error" error={data.message}>
    Error: {data.message}.
  </slot>
{/if}
