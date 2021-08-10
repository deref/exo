<script lang="ts">
  import * as router from 'svelte-spa-router';
  import IconButton from './IconButton.svelte';
  import VersionInfo from './VersionInfo.svelte';

  import Feedback from './mono/feedback.svelte';

  const goHome = () => {
    router.push('#/');
  };
</script>

<main>
  <nav>
    <header>
      <div class="a logo" on:click={goHome}>
        <img src="/deref-rounded-icon.png" alt="Deref" height="24px" />
      </div>
    </header>
    <slot name="navbar" />
    <footer>
      <IconButton
        tooltip="Give feedback on GitHub"
        on:click={() => {
          window.location.href = 'https://github.com/deref/exo/discussions';
        }}><Feedback /></IconButton
      >
      <div class:devmode={import.meta.env.MODE === 'development'}>
        <VersionInfo />
      </div>
    </footer>
  </nav>
  <div>
    <slot />
  </div>
</main>

<style>
  :global(body) {
    overscroll-behavior: none;
  }

  main {
    display: grid;
    grid-template-columns: max-content 1fr;
    gap: 1px;
    height: 100vh;
    overflow: hidden;
    background: #cccccc;
  }

  nav {
    position: relative;
    width: 48px;
    height: 100vh;
    z-index: 3;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
    background: #dddddd;
    gap: 1px;
  }

  div {
    position: relative;
    z-index: 2;
    overflow: hidden;
  }

  header {
    padding: 8px 0;
    padding-top: 12px;
    width: 100%;
    background: #e9e9e9;
  }

  header,
  .logo {
    display: flex;
    align-items: center;
    flex-direction: column;
    gap: 8px;
  }

  footer {
    display: flex;
    flex-direction: column;
    align-items: center;
    font-size: 12px;
    width: 100%;
    color: #555555;
    background: #e9e9e9;
    overflow: hidden;
    padding-top: 8px;
    gap: 8px;
  }

  div.devmode {
    background: #007bd4;
    color: #ffffff;
  }
</style>
