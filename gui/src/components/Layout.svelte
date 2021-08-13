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
    <div class="navbar-wrapper">
      <div>
        <slot name="navbar" />
      </div>
    </div>
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
    display: grid;
    grid-template-rows: auto 1fr max-content;
    grid-auto-flow: column;
    gap: 1px;
    background: var(--nav-bg-color);
  }

  .navbar-wrapper {
    width: 48px;
    overflow-y: scroll;
    direction: rtl;
  }

  nav:not(:hover) ::-webkit-scrollbar {
    display: none;
  }

  nav:hover .navbar-wrapper > div {
    margin-left: calc(-1 * var(--scrollbar-width));
  }

  .navbar-wrapper :global(*) {
    direction: ltr;
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
    font-size: 11px;
    width: 100%;
    color: var(--grey-5-color);
    overflow: hidden;
  }

  div.devmode {
    background: var(--dev-mode-bg-color);
    color: var(--dev-mode-color);
  }
</style>
