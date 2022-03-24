<script lang="ts">
  import CodeBlock from '../components/CodeBlock.svelte';
  import NewDockerComponentPage from './NewDockerComponentPage.svelte';
  import type { Params } from './NewDockerComponentPage.svelte';

  export let params: Params;

  const code = `image: postgres:9.4
environment:
    POSTGRES_USER: "postgres"
    POSTGRES_PASSWORD: "postgres"`;
  const href =
    'https://github.com/compose-spec/compose-spec/blob/master/spec.md#services-top-level-element';
</script>

<NewDockerComponentPage
  {params}
  componentType="container"
  displayType="Container"
>
  <p>See <a {href}>services volume reference</a>.</p>
  <CodeBlock {code} />
</NewDockerComponentPage>
