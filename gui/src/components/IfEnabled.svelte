<script lang="ts">
  import { api } from '../lib/api';
  import { onMount } from 'svelte';

  export let feature: string;

  let enabled: boolean;
  onMount(async () => {
    enabled = await api.kernel.isEnabled(feature);
  });
</script>

{#if enabled}
  <slot />
{/if}
