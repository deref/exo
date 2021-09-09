<script lang="ts">
  import CodeBlock from '../components/CodeBlock.svelte';
  import NewDockerComponent from './NewDockerComponent.svelte';
  import type { Params } from './NewDockerComponent.svelte';

  export let params: Params;

  const specExample = `name: mynetwork`;
  const href =
    'https://github.com/compose-spec/compose-spec/blob/master/spec.md#Networks-top-level-element';
</script>

<NewDockerComponent {params} componentType="network" displayType="Network">
  <div slot="help">
    <p>
      See <a {href}>compose network reference</a>.
    </p>
    <CodeBlock>
      {specExample}
    </CodeBlock>
  </div>
</NewDockerComponent>
