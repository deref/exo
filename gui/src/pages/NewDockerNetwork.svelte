<script lang="ts">
  import CodeBlock from '../components/CodeBlock.svelte';
  import NewDockerComponent from './NewDockerComponent.svelte';
  import type { Params } from './NewDockerComponent.svelte';

  export let params: Params;

  const code = `name: mynetwork`;
  const href =
    'https://github.com/compose-spec/compose-spec/blob/master/spec.md#Networks-top-level-element';
</script>

<NewDockerComponent {params} componentType="network" displayType="Network">
  <p>
    See <a {href}>compose network reference</a>.
  </p>
  <CodeBlock {code} />
</NewDockerComponent>
