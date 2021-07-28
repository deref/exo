<script lang="ts">
  import IconButton from './IconButton.svelte';
  import * as router from 'svelte-spa-router';
  import VersionInfo from './VersionInfo.svelte';

  import Feedback from './mono/feedback.svelte';
  import GoBack from './mono/leftarrow.svelte';

  export let showBackButton: boolean = false;
  export let backButtonRoute: string = '#/';

  const goHome = () => {
    router.push('#/');
  };

  const goBack = () => {
    router.push(backButtonRoute);
  };
</script>

<main>
  <header>
    <div class="logo">
      <div class="a logo" on:click={goHome}>
        <img src="/deref-rounded-icon.png" alt="Deref" height="24px" />
        <h1>exo</h1>
      </div>
      {#if showBackButton}
        <IconButton tooltip="Go back" on:click={goBack}><GoBack /></IconButton>
        <span>Go back</span>
      {/if}
    </div>
    <div class="logo">
      <span>Feedback?</span>
      <IconButton
        tooltip="Give feedback on GitHub"
        on:click={() => {
          window.location.href = 'https://github.com/deref/exo/discussions';
        }}><Feedback /></IconButton
      >
    </div>
  </header>
  <div>
    <slot />
  </div>
  <footer>
    <VersionInfo />
  </footer>
</main>

<style>
  :global(body) {
    overscroll-behavior: none;
  }

  main {
    display: grid;
    grid-auto-flow: row;
    grid-auto-rows: 48px 1fr 28px;
    height: 100vh;
    overflow: hidden;
  }

  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 12px;
    box-shadow: 0px 6px 9px -6px #00000022, 0px 0.25px 0px 1px #00000022;
  }

  header .logo {
    display: flex;
    align-items: center;
    gap: 6px;
  }

  h1 {
    font-size: 20px;
    font-weight: 550;
    margin: 0;
    margin-top: -3px;
  }

  div {
    position: relative;
    overflow-y: auto;
    overflow-x: hidden;
  }

  footer {
    display: flex;
    flex-direction: row-reverse;
    padding-right: 30px;
  }
</style>
