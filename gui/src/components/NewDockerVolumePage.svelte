<script lang="ts">
  import CodeBlock from '../components/CodeBlock.svelte';
  import NewDockerComponentPage from './NewDockerComponentPage.svelte';
  import type { Params } from './NewDockerComponentPage.svelte';

  export let params: Params;

  const code = `# An empty YAML is a valid volume.`;
  const href =
    'https://github.com/compose-spec/compose-spec/blob/master/spec.md#volumes-top-level-element';
</script>

<NewDockerComponentPage {params} componentType="volume" displayType="Volume">
  <p>See <a {href}>compose volume reference</a>.</p>
  <CodeBlock {code} />
</NewDockerComponentPage>
