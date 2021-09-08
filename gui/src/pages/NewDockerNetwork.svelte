<script lang="ts">
  import NewDockerComponent from './NewDockerComponent.svelte';
  import type { Params } from './NewDockerComponent.svelte';

  export let params: Params;

  const specExample = `# See https://github.com/compose-spec/compose-spec/blob/master/spec.md#Networks-top-level-element
name: mynetwork`;
</script>

<NewDockerComponent
  {params}
  componentType="network"
  displayType="Network"
  {specExample}
/>
