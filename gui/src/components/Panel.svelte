<script lang="ts">
  import IconButton from './IconButton.svelte';
  import BackSVG from './mono/leftarrow.svelte';
  import * as router from 'svelte-spa-router';

  export let title: string = '';
  export let backRoute: string = '';
</script>

<div class="panel">
  <header>
    <div class="header-title">
      {#if backRoute}
        <IconButton
          tooltip="Go back"
          on:click={() => {
            router.push(backRoute ?? '');
          }}
        >
          <BackSVG />
        </IconButton>
      {/if}

      {#if title}
        <h1>{title}</h1>
      {/if}
    </div>

    <div class="header-actions">
      <slot name="actions" />
    </div>
  </header>
  <section>
    <slot />
  </section>
  {#if $$slots.bottom}
    <slot name="bottom" />
  {/if}
</div>

<style>
  section,
  header,
  .panel {
    position: relative;
  }

  .panel {
    height: 100vh;
    overflow: hidden;
    background: var(--panel-bg-color);
    display: grid;
    grid-auto-flow: row;
    grid-template-rows: max-content auto max-content;
  }

  section {
    overflow-y: var(--panel-overflow-y);
    overflow-x: var(--panel-overflow-x);
    padding: var(--panel-padding);
  }

  header {
    border-bottom: 1px solid var(--layout-bg-color);
    overflow: hidden;
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 32px;
  }

  h1 {
    margin: 0;
    padding: 0;
    font-size: 16px;
    font-weight: 550;
  }

  h1:first-child {
    margin-left: 8px;
  }

  .header-title,
  .header-actions {
    display: flex;
    align-items: center;
  }
</style>
