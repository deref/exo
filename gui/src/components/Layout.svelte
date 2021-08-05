<script lang="ts">
  import * as router from 'svelte-spa-router';
  import IconButton from './IconButton.svelte';
  import VersionInfo from './VersionInfo.svelte';
  import LeftNav from './LeftNav.svelte';

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
    <LeftNav workspaceId="xbe1fc9j61kkk3d46dtjmkepnc" active="Dashboard" />
    <section>
      <slot />
    </section>
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
    display: grid;
    grid-template-columns: max-content auto;
    gap: 1px;
  }

  div > section {
    position: relative;
    overflow-y: auto;
    overflow-x: hidden;
    background: #ffffff;
  }

  footer {
    display: flex;
    flex-direction: row-reverse;
    padding: 4px 6px;
    height: 20px;
    align-items: center;
    background: #eeeeee;
  }
</style>
