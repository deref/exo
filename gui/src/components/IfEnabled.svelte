<script lang="ts">
  const features = new Set(['storybook feature']);

  export let feature: string;
</script>

{#if features.has(feature)}
  <slot />
{/if}
