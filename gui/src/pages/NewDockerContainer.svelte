<script lang="ts">
  import CodeBlock from '../components/CodeBlock.svelte';
  import NewDockerComponent from './NewDockerComponent.svelte';
  import type { Params } from './NewDockerComponent.svelte';

  export let params: Params;

  const specExample = `image: postgres:9.4
environment:
    POSTGRES_USER: "postgres"
    POSTGRES_PASSWORD: "postgres"`;
  const href =
    'https://github.com/compose-spec/compose-spec/blob/master/spec.md#services-top-level-element';
</script>

<NewDockerComponent {params} componentType="container" displayType="Container">
  <p>
    See <a {href}>services volume reference</a>.
  </p>
  <CodeBlock>
    {specExample}
  </CodeBlock>
</NewDockerComponent>
