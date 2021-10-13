<script lang="ts">
  import { api } from '../lib/api';

  export let raw: string;
</script>

<span title={raw}>
  {#await api.kernel.getUserHomeDir()}
    {raw}
  {:then homeDir}
    {raw.substr(0, homeDir.length) === homeDir
      ? '~' + raw.substr(homeDir.length)
      : raw}
  {/await}
</span>
