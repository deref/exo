<script lang="ts">
  import NewDockerComponent from './NewDockerComponent.svelte';
  import type { Params } from './NewDockerComponent.svelte';

  export let params: Params;

  const specExample = `
# See https://github.com/compose-spec/compose-spec/blob/master/spec.md#services-top-level-element
image: postgres:9.4
environment:
    POSTGRES_USER: "postgres"
    POSTGRES_PASSWORD: "postgres"`;
</script>

<NewDockerComponent
  {params}
  componentType="container"
  displayType="Container"
  {specExample}
/>
