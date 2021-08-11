<script lang="ts">
  import CheckSVG from './mono/check.svelte';

  export let active: boolean = false;
  export let tooltip: string = '';
</script>

<button
  on:click
  class={active ? 'active' : ''}
  title={tooltip !== '' ? tooltip : undefined}
>
  {#if active}
    <div>
      <CheckSVG />
    </div>
  {/if}
</button>

<style>
  button {
    font-family: inherit;
    font-size: 0;
    font-weight: 450;
    width: 22px;
    height: 22px;
    padding: 0;
    color: #777777;
    background-color: #77777700;
    border-radius: 4px;
    border: 1px solid #77777777;
    outline: none;
    margin-left: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  button > div {
    width: 32px;
    height: 30px;
    padding-left: 4px;
    padding-bottom: 2px;
  }

  button > div > :global(svg) {
    width: 28px;
    height: 28px;
  }

  button.active {
    border-color: #008ac577;
  }

  button :global(svg),
  button :global(svg *) {
    fill: currentColor;
  }

  button.active :global(svg),
  button.active :global(svg *) {
    fill: #008ac5;
  }

  button:hover {
    background-color: #77777711;
    color: #444444;
  }

  button.active:hover {
    background-color: #55ccff22;
  }

  button:focus {
    background-color: #77777722;
  }

  button.active:focus {
    background-color: #55ccff44;
  }
</style>
