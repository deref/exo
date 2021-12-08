<script lang="ts">
  import { getContext } from 'svelte';
  import Button from '../Button.svelte';

  export let onCancel = () => {};
  export let onOkay = () => {};
  export let danger = false;

  export let h3: string | undefined;
  export let bodyText: string | undefined;

  export let cancelLabel = 'Cancel';
  export let actionLabel = 'Okay';

  const { close } = getContext('simple-modal');

  function _onCancel() {
    onCancel();
    close();
  }

  function _onOkay() {
    onOkay();
    close();
  }
</script>

{#if h3}
  <h3>{h3}</h3>
{/if}

{#if bodyText}
  <p>{bodyText}</p>
{/if}

<div class="buttons">
  <Button on:click={_onCancel}>{cancelLabel}</Button>
  <Button on:click={_onOkay} {danger}>{actionLabel}</Button>
</div>

<style>
  .buttons {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 12px;
    margin-top: 2rem;
  }

  p {
    white-space: pre-wrap;
  }
</style>
