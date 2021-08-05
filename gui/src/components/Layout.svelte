<script lang="ts">
  import * as router from 'svelte-spa-router';
  import IconButton from './IconButton.svelte';
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
  <footer class:devmode="{import.meta.env.MODE === 'development'}">
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
    grid-auto-rows: max-content 1fr max-content;
    gap: 1px;
    height: 100vh;
    overflow: hidden;
    background: #cccccc;
  }

  header {
    position: relative;
    height: 40px;
    z-index: 3;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 13px;
    background: #dddddd;
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
    z-index: 2;
    overflow: hidden;
  }

  footer {
    display: flex;
    flex-direction: row-reverse;
    padding: 4px 6px;
    height: 20px;
    align-items: center;
    background: #eeeeee;
  }

  footer.devmode {
    background: rgb(21,3,33);
    background: linear-gradient(158deg, rgba(21,3,33,0.8403069846102503) 5%, rgba(127,85,183,1) 43%, rgba(144,218,245,1) 100%);
  }
</style>
