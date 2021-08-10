<script lang="ts">
  import * as router from 'svelte-spa-router';
  import VersionInfo from './VersionInfo.svelte';
  import NavbarButton from './nav/NavbarButton.svelte';
  import FeedbackSVG from './mono/feedback.svelte';

  const goHome = () => {
    router.push('#/');
  };
</script>

<main>
  <nav>
    <header>
      <NavbarButton title="Home" on:click={goHome}>
        <img src="/deref-rounded-icon.png" alt="Deref" height="24px" />
      </NavbarButton>
    </header>
    <slot name="navbar" />
    <footer>
      <NavbarButton
        title="Give feedback on GitHub"
        on:click={() => {
          window.location.href = 'https://github.com/deref/exo/discussions';
        }}
      >
        <FeedbackSVG />
      </NavbarButton>
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
    background: var(--layout-bg-color);
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
    gap: 1px;
    background: var(--nav-bg-color);
  }

  div {
    position: relative;
    z-index: 2;
    overflow: hidden;
  }

  footer {
    display: flex;
    flex-direction: column;
    align-items: center;
    font-size: 12px;
    width: 100%;
    color: var(--grey-5-text-color);
    border-top: 1px solid var(--layout-bg-color);
    overflow: hidden;
  }

  div.devmode {
    background: var(--dev-mode-bg-color);
    color: var(--dev-mode-text-color);
  }
</style>
