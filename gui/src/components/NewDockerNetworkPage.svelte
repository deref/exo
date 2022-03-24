<script lang="ts">
  import CodeBlock from '../components/CodeBlock.svelte';
  import NewDockerComponentPage from './NewDockerComponentPage.svelte';
  import type { Params } from './NewDockerComponentPage.svelte';

  export let params: Params;

  const code = `name: mynetwork`;
  const href =
    'https://github.com/compose-spec/compose-spec/blob/master/spec.md#Networks-top-level-element';
</script>

<NewDockerComponentPage {params} componentType="network" displayType="Network">
  <p>See <a {href}>compose network reference</a>.</p>
  <CodeBlock {code} />
</NewDockerComponentPage>
