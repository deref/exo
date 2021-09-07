<script lang="ts">
  import NewDockerComponent from './NewDockerComponent.svelte';
  import type { Params } from './NewDockerComponent.svelte';

  export let params: Params;

  const specExample = `# An empty YAML is a valid volume.
# See https://github.com/compose-spec/compose-spec/blob/master/spec.md#volumes-top-level-element
`;
</script>

<NewDockerComponent
  {params}
  componentType="volume"
  displayType="Volume"
  {specExample}
/>
