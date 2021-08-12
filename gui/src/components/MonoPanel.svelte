<script lang="ts">
  import IconButton from './IconButton.svelte';
  import BackSVG from './mono/leftarrow.svelte';
  import * as router from 'svelte-spa-router';

  export let title: string | undefined = undefined;
  export let backRoute: string | undefined = undefined;
</script>

{#if $$slots.actions || title !== undefined || backRoute !== undefined}
  <header>
    <div class="header-title">
      {#if backRoute !== undefined}
        <IconButton
          tooltip="Go back"
          on:click={() => {
            router.push(backRoute ?? '');
          }}
        >
          <BackSVG />
        </IconButton>
      {/if}

      {#if title !== undefined}
        <h1>{title}</h1>
      {/if}
    </div>

    <div class="header-actions">
      <slot name="actions" />
    </div>
  </header>
{/if}
<section>
  <slot />
</section>

<style>
  section {
    position: relative;
    overflow-y: scroll;
    overflow-x: hidden;
    background: var(--primary-bg-color);
    padding: var(--panel-padding);
    height: 100vh;
  }

  header {
    background: var(--primary-bg-color);
    border-bottom: 1px solid var(--layout-bg-color);
    overflow: hidden;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 32px;
    padding: 0 8px;
  }

  h1 {
    margin: 0;
    padding: 0;
    font-size: 16px;
    font-weight: 550;
  }

  .header-title,
  .header-actions {
    display: flex;
    align-items: center;
  }
</style>
