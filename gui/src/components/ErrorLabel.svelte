<script type="ts">
  import Warning from './mono/warning.svelte';

  export let error: Error | null = null;
</script>


{#if error}
<div>
  <Warning/>
  <span>{error.message}</span>
</div>
{/if}

<style>

div {
  display: flex;
  align-items: center;
  color: red;
}

span {
  display: block;
  padding-left: 8px;
}

</style>
